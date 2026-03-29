package exporter

import (
	"encoding/csv"
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"privantix-source-inspector/models"
	"privantix-source-inspector/report"
	"privantix-source-inspector/utils"
)

// StreamingWriter writes results incrementally to reduce memory usage.
type StreamingWriter struct {
	outputDir    string
	baseName     string
	filesCSV     *os.File
	filesW       *csv.Writer
	colsCSV      *os.File
	colsW        *csv.Writer
	jsonlFile    *os.File
	hasCompliance bool
	writeCount    int
	mu            sync.Mutex
}

// NewStreamingWriter creates a streaming writer. Call Close when done.
func NewStreamingWriter(outputDir, baseName string, hasCompliance bool) (*StreamingWriter, error) {
	if err := utils.EnsureDir(outputDir); err != nil {
		return nil, err
	}
	_, filesName, colsName, _ := outputFilenames(baseName)
	jsonlName := "analysis.jsonl"
	if baseName != "" {
		jsonlName = baseName + ".jsonl"
	}

	filesPath := filepath.Join(outputDir, filesName)
	f, err := os.Create(filesPath)
	if err != nil {
		return nil, err
	}
	filesW := csv.NewWriter(f)

	headers := []string{"path", "name", "extension", "owner", "permissions", "acls", "checksum", "size_bytes", "modified_at", "depth", "encoding", "delimiter", "has_header", "row_count_estimate", "column_count", "average_row_width", "sheet_name", "rules_triggered", "errors"}
	if hasCompliance {
		headers = append(headers, "trusted_principals", "outside_principals", "compliance_status")
	}
	if err := filesW.Write(headers); err != nil {
		f.Close()
		return nil, err
	}

	colsPath := filepath.Join(outputDir, colsName)
	cf, err := os.Create(colsPath)
	if err != nil {
		f.Close()
		return nil, err
	}
	colsW := csv.NewWriter(cf)
	if err := colsW.Write([]string{"file_path", "sheet_name", "column_name", "position", "inferred_type", "null_percentage", "max_length", "sample_values"}); err != nil {
		f.Close()
		cf.Close()
		return nil, err
	}

	jsonlPath := filepath.Join(outputDir, jsonlName)
	jf, err := os.Create(jsonlPath)
	if err != nil {
		f.Close()
		cf.Close()
		return nil, err
	}

	return &StreamingWriter{
		outputDir:     outputDir,
		baseName:      baseName,
		filesCSV:      f,
		filesW:        filesW,
		colsCSV:       cf,
		colsW:         colsW,
		jsonlFile:     jf,
		hasCompliance: hasCompliance,
	}, nil
}

// WriteRunStart writes the initial run metadata to JSONL.
func (w *StreamingWriter) WriteRunStart(run models.RunMetadata) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	obj := map[string]interface{}{
		"type":                 "run_start",
		"started_at":           run.StartedAt,
		"path":                 run.Path,
		"output_dir":           run.OutputDir,
		"supported_files":      run.SupportedFiles,
		"supported_extensions": run.SupportedExtensions,
		"max_sample_rows":      run.MaxSampleRows,
		"workers":              run.Workers,
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = w.jsonlFile.Write(append(b, '\n'))
	return err
}

// WriteProfile appends one file profile to CSV and JSONL.
func (w *StreamingWriter) WriteProfile(p models.FileProfile) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	rules := make([]string, 0, len(p.RulesTriggered))
	for _, r := range p.RulesTriggered {
		rules = append(rules, r.RuleName)
	}
	row := []string{
		p.Path, p.Name, p.Extension, p.Owner, p.Permissions,
		strings.Join(p.ACLs, " | "), p.Checksum,
		strconv.FormatInt(p.SizeBytes, 10),
		p.ModifiedAt.Format("2006-01-02 15:04:05"),
		strconv.Itoa(p.Depth), p.Encoding, p.Delimiter,
		strconv.FormatBool(p.HasHeader),
		strconv.Itoa(p.RowCountEstimate), strconv.Itoa(p.ColumnCount),
		strconv.Itoa(p.AverageRowWidth), p.SheetName,
		strings.Join(rules, ";"), strings.Join(p.Errors, ";"),
	}
	if w.hasCompliance {
		row = append(row,
			strings.Join(p.TrustedPrincipals, " | "),
			strings.Join(p.OutsidePrincipals, " | "),
			p.ComplianceStatus,
		)
	}
	if err := w.filesW.Write(row); err != nil {
		return err
	}

	for _, c := range p.Columns {
		if err := w.colsW.Write([]string{
			c.FilePath, c.SheetName, c.Name,
			strconv.Itoa(c.Position), c.InferredType,
			strconv.FormatFloat(c.NullPercentage, 'f', 2, 64),
			strconv.Itoa(c.MaxLength),
			strings.Join(c.SampleValues, " | "),
		}); err != nil {
			return err
		}
	}

	b, err := json.Marshal(map[string]interface{}{"type": "file", "file": p})
	if err != nil {
		return err
	}
	_, err = w.jsonlFile.Write(append(b, '\n'))
	if err != nil {
		return err
	}
	w.writeCount++
	if w.writeCount%50 == 0 {
		w.filesW.Flush()
		w.colsW.Flush()
	}
	return nil
}

// WriteRunEnd writes the final run metadata to JSONL.
func (w *StreamingWriter) WriteRunEnd(completedAt time.Time, analyzedFiles, failedFiles int, errors []string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	obj := map[string]interface{}{
		"type":            "run_end",
		"completed_at":    completedAt,
		"analyzed_files":  analyzedFiles,
		"failed_files":    failedFiles,
		"errors":          errors,
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = w.jsonlFile.Write(append(b, '\n'))
	return err
}

// Close flushes and closes all files.
func (w *StreamingWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	var errs []error
	w.filesW.Flush()
	if e := w.filesCSV.Close(); e != nil {
		errs = append(errs, e)
	}
	w.colsW.Flush()
	if e := w.colsCSV.Close(); e != nil {
		errs = append(errs, e)
	}
	if e := w.jsonlFile.Close(); e != nil {
		errs = append(errs, e)
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func ExportAll(outputDir string, result models.AnalysisResult) error {
	return ExportAllWithBaseName(outputDir, "", result)
}

// ExportAllWithBaseName exports to outputDir. If baseName is empty, uses default names (analysis.json, files.csv, etc.).
// If baseName is set (e.g. "CONADI_20250306"), uses baseName.json, baseName_files.csv, baseName_columns.csv, baseName_report.html.
func ExportAllWithBaseName(outputDir string, baseName string, result models.AnalysisResult) error {
	if err := utils.EnsureDir(outputDir); err != nil {
		return err
	}
	jsonName, filesName, colsName, htmlName := outputFilenames(baseName)
	if err := exportJSON(filepath.Join(outputDir, jsonName), result); err != nil {
		return err
	}
	if err := exportFilesCSV(filepath.Join(outputDir, filesName), result.Files); err != nil {
		return err
	}
	if err := exportColumnsCSV(filepath.Join(outputDir, colsName), result.Files); err != nil {
		return err
	}
	if err := exportHTML(filepath.Join(outputDir, htmlName), result); err != nil {
		return err
	}
	return nil
}

func outputFilenames(baseName string) (json, files, csv, html string) {
	if baseName == "" {
		return "analysis.json", "files.csv", "columns.csv", "report.html"
	}
	return baseName + ".json", baseName + "_files.csv", baseName + "_columns.csv", baseName + "_report.html"
}

func exportJSON(path string, result models.AnalysisResult) error {
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func exportFilesCSV(path string, files []models.FileProfile) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	headers := []string{"path", "name", "extension", "owner", "permissions", "acls", "checksum", "size_bytes", "modified_at", "depth", "encoding", "delimiter", "has_header", "row_count_estimate", "column_count", "average_row_width", "sheet_name", "rules_triggered", "errors"}
	hasCompliance := false
	for _, f := range files {
		if len(f.TrustedPrincipals) > 0 || len(f.OutsidePrincipals) > 0 {
			hasCompliance = true
			break
		}
	}
	if hasCompliance {
		headers = append(headers, "trusted_principals", "outside_principals", "compliance_status")
	}
	_ = w.Write(headers)
	for _, item := range files {
		rules := make([]string, 0, len(item.RulesTriggered))
		for _, r := range item.RulesTriggered {
			rules = append(rules, r.RuleName)
		}
		row := []string{
			item.Path,
			item.Name,
			item.Extension,
			item.Owner,
			item.Permissions,
			strings.Join(item.ACLs, " | "),
			item.Checksum,
			strconv.FormatInt(item.SizeBytes, 10),
			item.ModifiedAt.Format("2006-01-02 15:04:05"),
			strconv.Itoa(item.Depth),
			item.Encoding,
			item.Delimiter,
			strconv.FormatBool(item.HasHeader),
			strconv.Itoa(item.RowCountEstimate),
			strconv.Itoa(item.ColumnCount),
			strconv.Itoa(item.AverageRowWidth),
			item.SheetName,
			strings.Join(rules, ";"),
			strings.Join(item.Errors, ";"),
		}
		if hasCompliance {
			row = append(row,
				strings.Join(item.TrustedPrincipals, " | "),
				strings.Join(item.OutsidePrincipals, " | "),
				item.ComplianceStatus,
			)
		}
		_ = w.Write(row)
	}
	return nil
}

func exportColumnsCSV(path string, files []models.FileProfile) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	_ = w.Write([]string{"file_path", "sheet_name", "column_name", "position", "inferred_type", "null_percentage", "max_length", "sample_values"})
	for _, item := range files {
		for _, c := range item.Columns {
			_ = w.Write([]string{
				c.FilePath,
				c.SheetName,
				c.Name,
				strconv.Itoa(c.Position),
				c.InferredType,
				strconv.FormatFloat(c.NullPercentage, 'f', 2, 64),
				strconv.Itoa(c.MaxLength),
				strings.Join(c.SampleValues, " | "),
			})
		}
	}
	return nil
}

func exportHTML(path string, result models.AnalysisResult) error {
	tpl, err := template.New("report").Parse(report.HTML)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return tpl.Execute(f, result)
}
