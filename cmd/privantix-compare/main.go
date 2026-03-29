package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"privantix-source-inspector/catalog"
	"privantix-source-inspector/models"
)

func main() {
	left := flag.String("left", "", "Path to first analysis.json (baseline)")
	right := flag.String("right", "", "Path to second analysis.json (comparison)")
	output := flag.String("output", "", "Output file for JSON report (default: print to stdout)")
	format := flag.String("format", "text", "Output format: text | json")
	flag.Parse()

	if *left == "" || *right == "" {
		fmt.Println("Usage: privantix-compare --left <analysis1.json> --right <analysis2.json> [--output <file>] [--format text|json]")
		fmt.Println("")
		fmt.Println("Compares two privantix-inspector runs. Reports:")
		fmt.Println("  - Added:   files in right, not in left")
		fmt.Println("  - Removed: files in left, not in right")
		fmt.Println("  - Modified: files in both with different size, date, rows, or columns")
		os.Exit(1)
	}

	for _, p := range []string{*left, *right} {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			log.Fatalf("file not found: %s", p)
		}
	}

	started := time.Now()
	log.Printf("loading left:  %s", *left)
	leftResult, err := catalog.ParseAnalysis(*left)
	if err != nil {
		log.Fatalf("parse left: %v", err)
	}

	log.Printf("loading right: %s", *right)
	rightResult, err := catalog.ParseAnalysis(*right)
	if err != nil {
		log.Fatalf("parse right: %v", err)
	}

	diff := Compare(leftResult, rightResult)
	diff.LeftInput = *left
	diff.RightInput = *right
	diff.DurationSeconds = time.Since(started).Seconds()

	var out *os.File = os.Stdout
	if *output != "" {
		if err := os.MkdirAll(filepath.Dir(*output), 0o755); err != nil && filepath.Dir(*output) != "." {
			log.Fatalf("create output dir: %v", err)
		}
		f, err := os.Create(*output)
		if err != nil {
			log.Fatalf("create output: %v", err)
		}
		defer f.Close()
		out = f
	}

	switch *format {
	case "json":
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		if err := enc.Encode(diff); err != nil {
			log.Fatalf("write json: %v", err)
		}
	default:
		PrintTextReport(out, diff, *left, *right)
	}

	if *output != "" {
		log.Printf("report written to %s (%.2fs)", *output, diff.DurationSeconds)
	} else {
		log.Printf("compare completed (%.2fs)", diff.DurationSeconds)
	}
}

type CompareResult struct {
	LeftInput       string   `json:"left_input,omitempty"`
	RightInput      string   `json:"right_input,omitempty"`
	LeftPath        string   `json:"left_path"`
	RightPath       string   `json:"right_path"`
	DurationSeconds float64  `json:"duration_seconds"`
	Added           []string `json:"added"`
	Removed         []string `json:"removed"`
	Modified        []ModDiff `json:"modified"`
	Summary         Summary  `json:"summary"`
}

type ModDiff struct {
	Path    string            `json:"path"`
	Changes map[string]Change `json:"changes"`
}

type Change struct {
	Left  interface{} `json:"left"`
	Right interface{} `json:"right"`
}

type Summary struct {
	AddedCount   int `json:"added_count"`
	RemovedCount int `json:"removed_count"`
	ModifiedCount int `json:"modified_count"`
	UnchangedCount int `json:"unchanged_count"`
}

func Compare(left, right *models.AnalysisResult) CompareResult {
	leftByPath := make(map[string]models.FileProfile)
	for _, f := range left.Files {
		leftByPath[normalizePath(f.Path)] = f
	}
	rightByPath := make(map[string]models.FileProfile)
	for _, f := range right.Files {
		rightByPath[normalizePath(f.Path)] = f
	}

	var added, removed []string
	var modified []ModDiff

	for path := range rightByPath {
		if _, ok := leftByPath[path]; !ok {
			added = append(added, path)
		}
	}
	for path := range leftByPath {
		if _, ok := rightByPath[path]; !ok {
			removed = append(removed, path)
		}
	}

	for path, rp := range rightByPath {
		lp, ok := leftByPath[path]
		if !ok {
			continue
		}
		changes := diffProfile(lp, rp)
		if len(changes) > 0 {
			modified = append(modified, ModDiff{Path: path, Changes: changes})
		}
	}

	unchanged := len(leftByPath) - len(removed) - len(modified)

	return CompareResult{
		LeftPath:  left.Run.Path,
		RightPath: right.Run.Path,
		Added:     sorted(added),
		Removed:   sorted(removed),
		Modified:  modified,
		Summary: Summary{
			AddedCount:     len(added),
			RemovedCount:   len(removed),
			ModifiedCount:  len(modified),
			UnchangedCount: unchanged,
		},
	}
}

func diffProfile(left, right models.FileProfile) map[string]Change {
	changes := make(map[string]Change)
	if left.SizeBytes != right.SizeBytes {
		changes["size_bytes"] = Change{Left: left.SizeBytes, Right: right.SizeBytes}
	}
	if !left.ModifiedAt.Equal(right.ModifiedAt) {
		changes["modified_at"] = Change{
			Left:  left.ModifiedAt.Format("2006-01-02 15:04:05"),
			Right: right.ModifiedAt.Format("2006-01-02 15:04:05"),
		}
	}
	if left.RowCountEstimate != right.RowCountEstimate {
		changes["row_count_estimate"] = Change{Left: left.RowCountEstimate, Right: right.RowCountEstimate}
	}
	if left.ColumnCount != right.ColumnCount {
		changes["column_count"] = Change{Left: left.ColumnCount, Right: right.ColumnCount}
	}
	if left.Encoding != right.Encoding {
		changes["encoding"] = Change{Left: left.Encoding, Right: right.Encoding}
	}
	if left.Delimiter != right.Delimiter {
		changes["delimiter"] = Change{Left: left.Delimiter, Right: right.Delimiter}
	}
	return changes
}

func normalizePath(p string) string {
	return filepath.ToSlash(filepath.Clean(p))
}

func sorted(s []string) []string {
	sort.Strings(s)
	return s
}

func PrintTextReport(out *os.File, r CompareResult, leftPath, rightPath string) {
	fmt.Fprintf(out, "Compare: %s vs %s\n\n", leftPath, rightPath)
	fmt.Fprintf(out, "Baseline path:  %s\n", r.LeftPath)
	fmt.Fprintf(out, "Compare path:   %s\n", r.RightPath)
	fmt.Fprintf(out, "\n")
	fmt.Fprintf(out, "Summary:\n")
	fmt.Fprintf(out, "  Added:    %d\n", r.Summary.AddedCount)
	fmt.Fprintf(out, "  Removed:  %d\n", r.Summary.RemovedCount)
	fmt.Fprintf(out, "  Modified: %d\n", r.Summary.ModifiedCount)
	fmt.Fprintf(out, "  Unchanged: %d\n", r.Summary.UnchangedCount)
	fmt.Fprintf(out, "\n")

	if len(r.Added) > 0 {
		fmt.Fprintf(out, "Added (%d):\n", len(r.Added))
		for _, p := range r.Added {
			fmt.Fprintf(out, "  + %s\n", p)
		}
		fmt.Fprintf(out, "\n")
	}

	if len(r.Removed) > 0 {
		fmt.Fprintf(out, "Removed (%d):\n", len(r.Removed))
		for _, p := range r.Removed {
			fmt.Fprintf(out, "  - %s\n", p)
		}
		fmt.Fprintf(out, "\n")
	}

	if len(r.Modified) > 0 {
		fmt.Fprintf(out, "Modified (%d):\n", len(r.Modified))
		for _, m := range r.Modified {
			fmt.Fprintf(out, "  ~ %s\n", m.Path)
			for k, c := range m.Changes {
				fmt.Fprintf(out, "      %s: %v -> %v\n", k, c.Left, c.Right)
			}
		}
	}
}
