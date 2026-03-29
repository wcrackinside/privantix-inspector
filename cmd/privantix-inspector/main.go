package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"privantix-source-inspector/config"
	"privantix-source-inspector/exporter"
	"privantix-source-inspector/inspector"
	"privantix-source-inspector/models"
	"privantix-source-inspector/scanner"
	"privantix-source-inspector/utils"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "scan" {
		fmt.Println("Usage: privantix-inspector scan --path <path> [--output <dir>] [--config <file>]")
		fmt.Println("  [--max-sample-rows N] [--workers N] [--extensions csv,xlsx,...]")
		fmt.Println("  [--output-name <name>] [--trusted-groups <file.json>]")
		fmt.Println("  [--created-since <date>] [--owners <user1,user2,...>]")
		fmt.Println("  [--stream] [--log error|basic|detail] [--hide-samples] [--recursive=true] [--recursive-archives] [--incremental]")
		fmt.Println("  [--security] [--checksum]")
		os.Exit(1)
	}

	fs := flag.NewFlagSet("scan", flag.ExitOnError)
	path := fs.String("path", "", "Repository path to analyze")
	output := fs.String("output", "./output", "Output directory")
	configPath := fs.String("config", "", "Optional YAML config file")
	maxSampleRows := fs.Int("max-sample-rows", 0, "Optional max sample rows override")
	workersFlag := fs.Int("workers", 0, "Number of parallel workers (0=use config)")
	extensions := fs.String("extensions", "", "Optional comma separated extensions override")
	outputName := fs.String("output-name", "", "Base name for output files (e.g. CONADI_20250306). If empty, uses analysis.json, files.csv, etc.")
	trustedGroupsPath := fs.String("trusted-groups", "", "JSON file with trusted groups for compliance classification (requires --security)")
	securityFlag := fs.Bool("security", false, "Extract extra file security ACL info via OS commands")
	checksumFlag := fs.Bool("checksum", false, "Calculate SHA256 checksum for each file")
	hideSamplesFlag := fs.Bool("hide-samples", false, "Omit sample_values from output (for data governance / sensitive data)")
	recursiveFlag := fs.Bool("recursive", true, "Recursive directory search")
	recursiveArchivesFlag := fs.Bool("recursive-archives", false, "Process nested archives (ZIP/7z/RAR inside archives)")
	incrementalFlag := fs.Bool("incremental", false, "Only analyze files modified since last run (uses output dir)")
	createdSinceStr := fs.String("created-since", "", "Only include files with creation date >= this date (RFC3339 or 2006-01-02)")
	ownersFilter := fs.String("owners", "", "Only include files owned by these users/groups (comma-separated; e.g. DOMAIN\\user,group1)")
	streamFlag := fs.Bool("stream", false, "Write results incrementally (reduces memory; outputs JSONL instead of JSON)")
	logLevel := fs.String("log", "basic", "Log level: error (minimal), basic (file being analyzed), detail (file + size, ext)")
	if err := fs.Parse(os.Args[2:]); err != nil {
		log.Fatal(err)
	}
	if *path == "" {
		log.Fatal("--path is required")
	}
	pathClean := filepath.Clean(*path)
	logLvl := parseLogLevel(*logLevel)

	cfg, err := config.Load(*configPath)
	if err != nil && *configPath != "" {
		log.Fatalf("failed to load config: %v", err)
	}
	if *maxSampleRows > 0 {
		cfg.MaxSampleRows = *maxSampleRows
	}
	if *workersFlag > 0 {
		cfg.Workers = *workersFlag
	}
	if strings.TrimSpace(*extensions) != "" {
		cfg.SupportedExtensions = utils.NormalizeExtensions(strings.Split(*extensions, ","))
	}
	if *securityFlag {
		cfg.Security = true
	}
	trustedSet, _ := utils.LoadTrustedGroups(*trustedGroupsPath)
	if *trustedGroupsPath != "" && len(trustedSet) > 0 {
		cfg.Security = true
		logIf(logLvl, logBasic, "loaded %d trusted principals from %s", len(trustedSet), *trustedGroupsPath)
	} else if *trustedGroupsPath != "" {
		logIf(logLvl, logBasic, "warning: no trusted groups loaded from %s (file empty or not found)", *trustedGroupsPath)
	}
	if *checksumFlag {
		cfg.Checksum = true
	}
	if *hideSamplesFlag {
		cfg.HideSampleValues = true
	}
	if *recursiveArchivesFlag {
		cfg.RecursiveArchives = true
	}
	if *incrementalFlag {
		cfg.Incremental = true
	}
	cfg.Recursive = *recursiveFlag
	streamMode := *streamFlag

	if len(cfg.SupportedExtensions) == 0 {
		cfg.SupportedExtensions = config.Default().SupportedExtensions
		logIf(logLvl, logBasic, "warning: no supported extensions configured, using defaults: %v", cfg.SupportedExtensions)
	}
	if cfg.Workers < 1 {
		cfg.Workers = 1
	}

	outputBaseName := strings.TrimSpace(*outputName)
	analysisFilename := "analysis.json"
	if outputBaseName != "" {
		analysisFilename = outputBaseName + ".json"
	}

	started := time.Now()
	logIf(logLvl, logBasic, "starting scan: %s", pathClean)
	discovered, scanErrors, err := scanner.Scan(pathClean, cfg.SupportedExtensions, cfg.Recursive)
	if err != nil {
		log.Fatalf("scan error: %v", err)
	}

	var createdSince time.Time
	if strings.TrimSpace(*createdSinceStr) != "" {
		for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05Z07:00", "2006-01-02"} {
			if t, err := time.Parse(layout, strings.TrimSpace(*createdSinceStr)); err == nil {
				createdSince = t
				logIf(logLvl, logBasic, "filter: created-since >= %s", createdSince.Format(time.RFC3339))
				break
			}
		}
		if createdSince.IsZero() {
			log.Fatalf("invalid --created-since value %q (use RFC3339 or 2006-01-02)", *createdSinceStr)
		}
	}

	ownersSet := make(map[string]struct{})
	if strings.TrimSpace(*ownersFilter) != "" {
		for _, o := range strings.Split(*ownersFilter, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				ownersSet[o] = struct{}{}
			}
		}
		logIf(logLvl, logBasic, "filter: owners %s", *ownersFilter)
	}

	if !createdSince.IsZero() || len(ownersSet) > 0 {
		filtered := discovered[:0]
		for _, f := range discovered {
			if !createdSince.IsZero() && f.Created.Before(createdSince) {
				continue
			}
			if len(ownersSet) > 0 {
				if _, ok := ownersSet[f.Owner]; !ok {
					continue
				}
			}
			filtered = append(filtered, f)
		}
		logIf(logLvl, logBasic, "filter: %d files after created-since/owners (of %d discovered)", len(filtered), len(discovered))
		discovered = filtered
	}

	if len(discovered) == 0 {
		logIf(logLvl, logBasic, "warning: 0 files found at %s (extensions: %v). Check path exists and contains .csv, .xlsx, .parquet, etc.", pathClean, cfg.SupportedExtensions)
	}
	for _, e := range scanErrors {
		logIf(logLvl, logError, "scan error: %s", e)
	}

	var previousResult *models.AnalysisResult
	useIncrementalByPath := false
	if cfg.Incremental {
		jsonPath := filepath.Join(*output, analysisFilename)
		prev, err := loadPreviousAnalysis(jsonPath)
		if err == nil {
			previousResult = prev
			logIf(logLvl, logBasic, "incremental: found previous run from %s, filtering modified files", prev.Run.CompletedAt.Format(time.RFC3339))
		} else {
			jsonlPath := filepath.Join(*output, strings.TrimSuffix(analysisFilename, ".json")+".jsonl")
			partial, errPartial := loadPartialFromJSONL(jsonlPath)
			if errPartial == nil && len(partial) > 0 {
				previousResult = &models.AnalysisResult{Files: partial}
				useIncrementalByPath = true
				logIf(logLvl, logBasic, "incremental: resuming from partial run (%d files already done), analyzing the rest", len(partial))
			}
		}
	}

	var toAnalyze []models.FileDiscovered
	if previousResult != nil {
		if useIncrementalByPath {
			donePaths := make(map[string]bool)
			for _, p := range previousResult.Files {
				donePaths[filepath.ToSlash(p.Path)] = true
			}
			for _, f := range discovered {
				if !donePaths[filepath.ToSlash(f.Path)] {
					toAnalyze = append(toAnalyze, f)
				}
			}
			logIf(logLvl, logBasic, "incremental: %d files left to analyze (of %d total, %d already done)", len(toAnalyze), len(discovered), len(previousResult.Files))
		} else {
			cutoff := previousResult.Run.CompletedAt
			for _, f := range discovered {
				if f.Modified.After(cutoff) {
					toAnalyze = append(toAnalyze, f)
				}
			}
			logIf(logLvl, logBasic, "incremental: %d files modified since last run (of %d total)", len(toAnalyze), len(discovered))
		}
	} else {
		toAnalyze = discovered
	}

	jobs := make(chan models.FileDiscovered)
	results := make(chan []models.FileProfile)
	var wg sync.WaitGroup
	var currentFileMu sync.Mutex
	var currentFilePath string
	for i := 0; i < cfg.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range jobs {
				currentFileMu.Lock()
				currentFilePath = file.Path
				currentFileMu.Unlock()
				if logLvl >= logBasic {
					if logLvl >= logDetail {
						log.Printf("Analyzing: %s (%d bytes, %s)", file.Path, file.Size, file.Ext)
					} else {
						log.Printf("Analyzing: %s", file.Path)
					}
				}
				results <- inspector.Analyze(file, cfg)
			}
		}()
	}
	go func() {
		for _, file := range toAnalyze {
			jobs <- file
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	failed := 0
	completed := 0
	total := len(toAnalyze)
	hasCompliance := len(trustedSet) > 0

	if streamMode {
		sw, err := exporter.NewStreamingWriter(*output, outputBaseName, hasCompliance)
		if err != nil {
			log.Fatalf("streaming writer: %v", err)
		}
		defer sw.Close()

		runMeta := models.RunMetadata{
			StartedAt:           started,
			Path:                pathClean,
			OutputDir:           *output,
			SupportedFiles:      len(discovered),
			MaxSampleRows:       cfg.MaxSampleRows,
			Workers:             cfg.Workers,
			SupportedExtensions: cfg.SupportedExtensions,
		}
		if err := sw.WriteRunStart(runMeta); err != nil {
			log.Fatalf("write run start: %v", err)
		}

		if total == 0 && previousResult != nil {
			for _, p := range previousResult.Files {
				applyProfileTransforms(&p, cfg.HideSampleValues, trustedSet)
				if err := sw.WriteProfile(p); err != nil {
					log.Fatalf("write profile: %v", err)
				}
			}
			completed = len(previousResult.Files)
		} else {
			var previousToWrite []models.FileProfile
			if cfg.Incremental && previousResult != nil && len(toAnalyze) > 0 {
				reAnalyzedSet := make(map[string]bool)
				for _, f := range toAnalyze {
					reAnalyzedSet[filepath.ToSlash(f.Path)] = true
				}
				for _, p := range previousResult.Files {
					pPath := filepath.ToSlash(p.Path)
					keep := true
					for base := range reAnalyzedSet {
						if pPath == base || strings.HasPrefix(pPath, base+"/") {
							keep = false
							break
						}
					}
					if keep {
						previousToWrite = append(previousToWrite, p)
					}
				}
			}
			prevCount := len(previousToWrite)
			for _, p := range previousToWrite {
				applyProfileTransforms(&p, cfg.HideSampleValues, trustedSet)
				if err := sw.WriteProfile(p); err != nil {
					log.Fatalf("write profile: %v", err)
				}
			}
			completed = prevCount

			done := make(chan bool)
			if logLvl >= logBasic {
				go func() {
					ticker := time.NewTicker(200 * time.Millisecond)
					defer ticker.Stop()
					for {
						select {
						case <-ticker.C:
							pct := float64(completed) / float64(max(1, total)) * 100
							line := fmt.Sprintf("[%6.2f%%] Analyzed %d/%d files...", pct, completed, total)
							if cfg.Workers == 1 {
								currentFileMu.Lock()
								cf := currentFilePath
								currentFileMu.Unlock()
								if cf != "" {
									base := filepath.Base(cf)
									if len(base) > 45 {
										base = base[:42] + "..."
									}
									line += " " + base
								}
							}
							fmt.Printf("\r\033[K%s", line)
						case <-done:
							fmt.Printf("\r\033[K[100.00%%] Analyzed %d/%d files. Done!\n", total, total)
							return
						}
					}
				}()
			}

			for batch := range results {
				for _, profile := range batch {
					if len(profile.Errors) > 0 {
						failed++
						logIf(logLvl, logError, "error analyzing %s: %s", profile.Path, strings.Join(profile.Errors, "; "))
					}
					applyProfileTransforms(&profile, cfg.HideSampleValues, trustedSet)
					if err := sw.WriteProfile(profile); err != nil {
						log.Fatalf("write profile: %v", err)
					}
					completed++
				}
			}
			if logLvl >= logBasic {
				done <- true
			}
		}

		streamCompletedAt := time.Now()
		if err := sw.WriteRunEnd(streamCompletedAt, completed, failed, scanErrors); err != nil {
			log.Fatalf("write run end: %v", err)
		}
		streamDuration := streamCompletedAt.Sub(started).Seconds()
		logIf(logLvl, logBasic, "analysis completed: %d files, %d failed, output=%s (streaming: JSONL, %.2fs)", completed, failed, *output, streamDuration)
		return
	}

	profiles := make([]models.FileProfile, 0, len(toAnalyze))

	if total == 0 && previousResult != nil {
		profiles = previousResult.Files
		logIf(logLvl, logBasic, "incremental: no files modified, reusing previous %d profiles", len(profiles))
	} else {
		done := make(chan bool)
		if logLvl >= logBasic {
			go func() {
				ticker := time.NewTicker(200 * time.Millisecond)
				defer ticker.Stop()
				for {
					select {
					case <-ticker.C:
						pct := float64(completed) / float64(max(1, total)) * 100
						line := fmt.Sprintf("[%6.2f%%] Analyzed %d/%d files...", pct, completed, total)
						if cfg.Workers == 1 {
							currentFileMu.Lock()
							cf := currentFilePath
							currentFileMu.Unlock()
							if cf != "" {
								base := filepath.Base(cf)
								if len(base) > 45 {
									base = base[:42] + "..."
								}
								line += " " + base
							}
						}
						fmt.Printf("\r\033[K%s", line)
					case <-done:
						fmt.Printf("\r\033[K[100.00%%] Analyzed %d/%d files. Done!\n", total, total)
						return
					}
				}
			}()
		}

		for batch := range results {
			for _, profile := range batch {
				if len(profile.Errors) > 0 {
					failed++
					logIf(logLvl, logError, "error analyzing %s: %s", profile.Path, strings.Join(profile.Errors, "; "))
				}
				profiles = append(profiles, profile)
			}
			completed++
		}
		if logLvl >= logBasic {
			done <- true
		}
	}

	if cfg.Incremental && previousResult != nil && len(toAnalyze) > 0 {
		profiles = mergeIncrementalProfiles(previousResult.Files, profiles, toAnalyze)
	}

	completedAt := time.Now()
	duration := completedAt.Sub(started).Seconds()
	result := models.AnalysisResult{
		Run: models.RunMetadata{
			StartedAt:           started,
			CompletedAt:         completedAt,
			DurationSeconds:     duration,
			Path:                pathClean,
			OutputDir:           *output,
			SupportedFiles:      len(discovered),
			AnalyzedFiles:       len(profiles),
			FailedFiles:         failed,
			MaxSampleRows:       cfg.MaxSampleRows,
			Workers:             cfg.Workers,
			SupportedExtensions: cfg.SupportedExtensions,
		},
		Files:  profiles,
		Errors: scanErrors,
	}

	if cfg.HideSampleValues {
		for i := range result.Files {
			for j := range result.Files[i].Columns {
				result.Files[i].Columns[j].SampleValues = nil
			}
		}
	}

	if len(trustedSet) > 0 {
		for i := range result.Files {
			if len(result.Files[i].ACLs) > 0 {
				result.Files[i].TrustedPrincipals, result.Files[i].OutsidePrincipals = utils.ClassifyPrincipals(result.Files[i].ACLs, trustedSet)
				result.Files[i].ComplianceStatus = utils.ComplianceStatus(result.Files[i].TrustedPrincipals, result.Files[i].OutsidePrincipals)
			}
		}
	}

	if err := exporter.ExportAllWithBaseName(*output, outputBaseName, result); err != nil {
		log.Fatalf("export error: %v", err)
	}
	logIf(logLvl, logBasic, "analysis completed: %d files, %d failed, output=%s (%.2fs)", len(profiles), failed, *output, duration)
}

func applyProfileTransforms(p *models.FileProfile, hideSamples bool, trustedSet map[string]struct{}) {
	if hideSamples {
		for j := range p.Columns {
			p.Columns[j].SampleValues = nil
		}
	}
	if len(trustedSet) > 0 && len(p.ACLs) > 0 {
		p.TrustedPrincipals, p.OutsidePrincipals = utils.ClassifyPrincipals(p.ACLs, trustedSet)
		p.ComplianceStatus = utils.ComplianceStatus(p.TrustedPrincipals, p.OutsidePrincipals)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

const (
	logError  = 0
	logBasic  = 1
	logDetail = 2
)

func parseLogLevel(s string) int {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "error":
		return logError
	case "basic":
		return logBasic
	case "detail":
		return logDetail
	default:
		return logBasic
	}
}

func logIf(level, minLevel int, format string, args ...interface{}) {
	if level >= minLevel {
		log.Printf(format, args...)
	}
}

func loadPreviousAnalysis(path string) (*models.AnalysisResult, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var result models.AnalysisResult
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// loadPartialFromJSONL reads a JSONL file (e.g. from an interrupted --stream run) and returns
// all file profiles found. Used by --incremental to resume when no complete analysis.json exists.
func loadPartialFromJSONL(path string) ([]models.FileProfile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var profiles []models.FileProfile
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var line struct {
			Type string            `json:"type"`
			File models.FileProfile `json:"file"`
		}
		if err := json.Unmarshal(sc.Bytes(), &line); err != nil {
			continue
		}
		if line.Type == "file" {
			profiles = append(profiles, line.File)
		}
	}
	return profiles, sc.Err()
}

func mergeIncrementalProfiles(previous, newProfiles []models.FileProfile, reAnalyzed []models.FileDiscovered) []models.FileProfile {
	reAnalyzedSet := make(map[string]bool)
	for _, f := range reAnalyzed {
		reAnalyzedSet[filepath.ToSlash(f.Path)] = true
	}
	var merged []models.FileProfile
	for _, p := range previous {
		pPath := filepath.ToSlash(p.Path)
		keep := true
		for base := range reAnalyzedSet {
			if pPath == base || strings.HasPrefix(pPath, base+"/") {
				keep = false
				break
			}
		}
		if keep {
			merged = append(merged, p)
		}
	}
	return append(merged, newProfiles...)
}
