package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"privantix-source-inspector/utils"
)

type TrustedGroup struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

type TrustedGroupsFile struct {
	Groups []TrustedGroup `json:"groups"`
}

type DirACL struct {
	Path              string   `json:"path"`
	Name              string   `json:"name"`
	Depth             int      `json:"depth"`
	Owner             string   `json:"owner"`
	Permissions       string   `json:"permissions"`
	ACLs              []string `json:"acls"`
	TrustedPrincipals []string `json:"trusted_principals,omitempty"`
	OutsidePrincipals []string `json:"outside_principals,omitempty"`
	RevokedPrincipals []string `json:"revoked_principals,omitempty"`
	ComplianceStatus  string   `json:"compliance_status,omitempty"`
}

type RevokedUsersFile struct {
	Users []string `json:"users"`
}

type AuditResult struct {
	StartedAt       time.Time     `json:"started_at"`
	CompletedAt     time.Time     `json:"completed_at"`
	DurationSeconds float64       `json:"duration_seconds"`
	Path            string        `json:"path"`
	OutputDir       string        `json:"output_dir"`
	TotalDirs       int           `json:"total_dirs"`
	TrustedGroups   []TrustedGroup `json:"trusted_groups,omitempty"`
	RevokedUsers    []string      `json:"revoked_users,omitempty"`
	Dirs            []DirACL      `json:"dirs"`
	Errors          []string      `json:"errors,omitempty"`
}

func main() {
	path := flag.String("path", "", "Root path to audit (required)")
	output := flag.String("output", "./output", "Output directory")
	outputName := flag.String("output-name", "", "Base name for output files. If empty, uses path basename + _yyyymmdd (e.g. MiProyecto_20250306)")
	trustedGroupsPath := flag.String("trusted-groups", "", "JSON file with trusted groups (e.g. {\"groups\":[{\"name\":\"dais\",\"members\":[\"carce\",\"lmenesesv\"]}]})")
	revokedUsersPath := flag.String("revoked-users", "", "JSON file with revoked/deactivated users (e.g. {\"users\":[\"CORP\\\\jperez\",\"mgarcia\"]})")
	flag.Parse()

	if *path == "" {
		fmt.Println("Usage: privantix-acl-audit --path <path> [--output <dir>] [--output-name <name>] [--trusted-groups <file.json>] [--revoked-users <file.json>]")
		fmt.Println("Audits ACLs and permissions for all folders and subfolders.")
		fmt.Println("With --output-name: base name for output files (default: path basename + _yyyymmdd).")
		fmt.Println("With --trusted-groups: flags principals in/outside trusted groups.")
		fmt.Println("With --revoked-users: flags principals that are deactivated/revoked and should not have access.")
		os.Exit(1)
	}

	trustedSet, loadErr := loadTrustedGroups(*trustedGroupsPath)
	if *trustedGroupsPath != "" && len(trustedSet) == 0 {
		if loadErr != nil {
			log.Printf("warning: could not load trusted groups from %s: %v", *trustedGroupsPath, loadErr)
		} else {
			log.Printf("warning: no trusted groups loaded from %s (file empty or no valid members)", *trustedGroupsPath)
		}
	} else if len(trustedSet) > 0 {
		log.Printf("loaded %d trusted principals from %s", len(trustedSet), *trustedGroupsPath)
	}

	revokedSet, revokedList, revokedErr := loadRevokedUsers(*revokedUsersPath)
	if *revokedUsersPath != "" && len(revokedSet) == 0 {
		if revokedErr != nil {
			log.Printf("warning: could not load revoked users from %s: %v", *revokedUsersPath, revokedErr)
		} else {
			log.Printf("warning: no revoked users loaded from %s (file empty or no valid users)", *revokedUsersPath)
		}
	} else if len(revokedSet) > 0 {
		log.Printf("loaded %d revoked users from %s", len(revokedSet), *revokedUsersPath)
	}

	absPath, err := filepath.Abs(*path)
	if err != nil {
		log.Fatalf("invalid path: %v", err)
	}
	if _, err := os.Stat(absPath); err != nil {
		log.Fatalf("path not found: %v", err)
	}

	started := time.Now()
	log.Printf("auditing ACLs: %s", absPath)

	// First pass: count directories for progress bar
	log.Printf("counting directories...")
	totalDirs := countDirs(absPath)
	log.Printf("found %d directories, analyzing...", totalDirs)
	bar := progressbar.NewOptions(totalDirs,
		progressbar.OptionSetDescription("Analizando ACLs"),
		progressbar.OptionSetWidth(40),
		progressbar.OptionShowCount(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	var dirs []DirACL
	var errs []string

	err = filepath.WalkDir(absPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", path, err))
			return nil
		}
		if !d.IsDir() {
			return nil
		}

		info, statErr := d.Info()
		if statErr != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", path, statErr))
			_ = bar.Add(1)
			return nil
		}

		owner, perms := utils.GetFileOwnerAndPerms(info, path)
		acls := utils.GetFileACLs(path)

		dir := DirACL{
			Path:        path,
			Name:        info.Name(),
			Depth:       fileDepth(absPath, path),
			Owner:       owner,
			Permissions: perms,
			ACLs:        acls,
		}
		if len(trustedSet) > 0 {
			dir.TrustedPrincipals, dir.OutsidePrincipals = classifyPrincipals(acls, trustedSet)
		}
		if len(revokedSet) > 0 {
			dir.RevokedPrincipals = findRevokedPrincipals(acls, revokedSet)
		}
		dir.ComplianceStatus = complianceStatus(dir.TrustedPrincipals, dir.OutsidePrincipals, dir.RevokedPrincipals)
		dirs = append(dirs, dir)
		_ = bar.Add(1)
		return nil
	})

	if err != nil {
		log.Fatalf("walk error: %v", err)
	}
	_ = bar.Finish()
	_ = bar.Clear()

	var trustedGroups []TrustedGroup
	if *trustedGroupsPath != "" {
		if tg, err := readTrustedGroupsFile(*trustedGroupsPath); err == nil {
			trustedGroups = tg.Groups
		}
	}

	completedAt := time.Now()
	result := AuditResult{
		StartedAt:       started,
		CompletedAt:     completedAt,
		DurationSeconds: completedAt.Sub(started).Seconds(),
		Path:            absPath,
		OutputDir:       *output,
		TotalDirs:       len(dirs),
		TrustedGroups:   trustedGroups,
		RevokedUsers:    revokedList,
		Dirs:            dirs,
		Errors:          errs,
	}

	if err := utils.EnsureDir(*output); err != nil {
		log.Fatalf("cannot create output dir: %v", err)
	}

	// Default output name: path basename + _yyyymmdd
	baseName := strings.TrimSpace(*outputName)
	if baseName == "" {
		baseName = sanitizeFilename(filepath.Base(absPath)) + "_" + time.Now().Format("20060102")
	}

	jsonPath := filepath.Join(*output, baseName+".json")
	if err := exportJSON(jsonPath, result); err != nil {
		log.Fatalf("export JSON: %v", err)
	}

	csvPath := filepath.Join(*output, baseName+".csv")
	if err := exportCSV(csvPath, dirs); err != nil {
		log.Fatalf("export CSV: %v", err)
	}

	log.Printf("audit completed: %d directories, output=%s (%.2fs)", len(dirs), *output, result.DurationSeconds)
}

func loadTrustedGroups(path string) (map[string]struct{}, error) {
	tg, err := readTrustedGroupsFile(path)
	if err != nil {
		return nil, err
	}
	return buildPrincipalSet(tg.Groups), nil
}

func buildPrincipalSet(groups []TrustedGroup) map[string]struct{} {
	set := make(map[string]struct{})
	for _, g := range groups {
		for _, m := range g.Members {
			m = strings.TrimSpace(m)
			if m == "" {
				continue
			}
			keyFull := strings.ToLower(m)
			set[keyFull] = struct{}{}
			if idx := strings.LastIndex(m, "\\"); idx != -1 && idx+1 < len(m) {
				keyUser := strings.ToLower(m[idx+1:])
				if keyUser != "" {
					set[keyUser] = struct{}{}
				}
			}
		}
	}
	return set
}

func loadRevokedUsers(path string) (map[string]struct{}, []string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	var f RevokedUsersFile
	if err := json.Unmarshal(b, &f); err != nil {
		return nil, nil, err
	}
	set := make(map[string]struct{})
	var list []string
	for _, u := range f.Users {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}
		keyFull := strings.ToLower(u)
		set[keyFull] = struct{}{}
		list = append(list, u)
		if idx := strings.LastIndex(u, "\\"); idx != -1 && idx+1 < len(u) {
			keyUser := strings.ToLower(u[idx+1:])
			if keyUser != "" {
				set[keyUser] = struct{}{}
			}
		}
	}
	return set, list, nil
}

func findRevokedPrincipals(acls []string, revokedSet map[string]struct{}) []string {
	var revoked []string
	seen := make(map[string]bool)
	for _, acl := range acls {
		principal := extractPrincipal(acl)
		if principal == "" || seen[principal] {
			continue
		}
		if strings.Contains(principal, "archivos") || strings.Contains(principal, "processed") {
			continue
		}
		seen[principal] = true
		principalNorm := strings.ToLower(principal)
		memberNorm := strings.ToLower(principalToMember(principal))
		_, fullMatch := revokedSet[principalNorm]
		_, memberMatch := revokedSet[memberNorm]
		if fullMatch || memberMatch {
			revoked = append(revoked, principal)
		}
	}
	return revoked
}

func readTrustedGroupsFile(path string) (*TrustedGroupsFile, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var tg TrustedGroupsFile
	if err := json.Unmarshal(b, &tg); err != nil {
		return nil, err
	}
	return &tg, nil
}

// extractPrincipal gets the user/group from ACL string "DOMAIN\user:(I)(F)" -> "DOMAIN\user"
func extractPrincipal(acl string) string {
	idx := strings.Index(acl, ":")
	if idx == -1 {
		return strings.TrimSpace(acl)
	}
	return strings.TrimSpace(acl[:idx])
}

// principalToMember extracts the member part for matching: "DOMAIN\carce" -> "carce", "carce" -> "carce"
func principalToMember(principal string) string {
	principal = strings.TrimSpace(principal)
	if idx := strings.LastIndex(principal, "\\"); idx != -1 && idx+1 < len(principal) {
		return principal[idx+1:]
	}
	return principal
}

func classifyPrincipals(acls []string, trustedSet map[string]struct{}) (trusted, outside []string) {
	seen := make(map[string]bool)
	for _, acl := range acls {
		principal := extractPrincipal(acl)
		if principal == "" {
			continue
		}
		if strings.Contains(principal, "archivos") || strings.Contains(principal, "processed") {
			continue
		}
		if seen[principal] {
			continue
		}
		seen[principal] = true
		// Match by full principal (MDP\carce) or by user part (carce)
		principalNorm := strings.ToLower(principal)
		memberNorm := strings.ToLower(principalToMember(principal))
		_, fullMatch := trustedSet[principalNorm]
		_, memberMatch := trustedSet[memberNorm]
		if fullMatch || memberMatch {
			trusted = append(trusted, principal)
		} else {
			outside = append(outside, principal)
		}
	}
	return trusted, outside
}

func complianceStatus(trusted, outside, revoked []string) string {
	if len(revoked) > 0 {
		return "critical"
	}
	if len(outside) == 0 {
		return "compliant"
	}
	return "non_compliant"
}

func fileDepth(root, path string) int {
	rel, err := filepath.Rel(root, path)
	if err != nil || rel == "." {
		return 0
	}
	parts := strings.Split(rel, string(filepath.Separator))
	return len(parts)
}

func countDirs(root string) int {
	var n int
	_ = filepath.WalkDir(root, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			n++
		}
		return nil
	})
	return n
}

// sanitizeFilename removes characters invalid for filenames (e.g. : \ / * ? " < > |)
func sanitizeFilename(s string) string {
	invalid := `:\/*?"<>|`
	for _, r := range invalid {
		s = strings.ReplaceAll(s, string(r), "_")
	}
	if s == "" {
		s = "acl-audit"
	}
	return s
}

func exportJSON(path string, result AuditResult) error {
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func exportCSV(path string, dirs []DirACL) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	headers := []string{"path", "name", "depth", "owner", "permissions", "acls"}
	if len(dirs) > 0 && (len(dirs[0].TrustedPrincipals) > 0 || len(dirs[0].OutsidePrincipals) > 0 || len(dirs[0].RevokedPrincipals) > 0) {
		headers = append(headers, "trusted_principals", "outside_principals", "revoked_principals", "compliance_status")
	}
	_ = w.Write(headers)

	for _, d := range dirs {
		row := []string{
			d.Path,
			d.Name,
			fmt.Sprintf("%d", d.Depth),
			d.Owner,
			d.Permissions,
			strings.Join(d.ACLs, " | "),
		}
		if len(d.TrustedPrincipals) > 0 || len(d.OutsidePrincipals) > 0 || len(d.RevokedPrincipals) > 0 {
			row = append(row,
				strings.Join(d.TrustedPrincipals, " | "),
				strings.Join(d.OutsidePrincipals, " | "),
				strings.Join(d.RevokedPrincipals, " | "),
				d.ComplianceStatus,
			)
		}
		_ = w.Write(row)
	}
	return nil
}
