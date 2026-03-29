package csvanalyzer

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"strings"

	"privantix-source-inspector/detectors"
	"privantix-source-inspector/models"
	"privantix-source-inspector/profiler"
	"privantix-source-inspector/utils"
)

func Analyze(discovered models.FileDiscovered, maxSampleRows int) models.FileProfile {
	f, err := os.Open(discovered.Path)
	if err != nil {
		profile := models.FileProfile{
			Path: discovered.Path, Name: discovered.Name, Extension: discovered.Ext,
			SizeBytes: discovered.Size, ModifiedAt: discovered.Modified, CreatedAt: discovered.Created,
			Depth: discovered.Depth, Owner: discovered.Owner, Permissions: discovered.Permissions,
			Errors: []string{err.Error()},
		}
		return profile
	}
	defer f.Close()
	return analyzeFromReader(f, discovered, maxSampleRows)
}

// AnalyzeFromReader analyzes CSV/text content from an io.Reader (e.g. file inside archive).
func AnalyzeFromReader(r io.Reader, discovered models.FileDiscovered, maxSampleRows int) models.FileProfile {
	return analyzeFromReader(r, discovered, maxSampleRows)
}

func analyzeFromReader(r io.Reader, discovered models.FileDiscovered, maxSampleRows int) models.FileProfile {
	profile := models.FileProfile{
		Path:        discovered.Path,
		Name:        discovered.Name,
		Extension:   discovered.Ext,
		SizeBytes:   discovered.Size,
		ModifiedAt:  discovered.Modified,
		CreatedAt:   discovered.Created,
		Depth:       discovered.Depth,
		Owner:       discovered.Owner,
		Permissions: discovered.Permissions,
	}

	sample := make([]byte, 4096)
	var scanSrc io.Reader = r
	if readAt, ok := r.(io.ReaderAt); ok {
		n, _ := readAt.ReadAt(sample, 0)
		profile.Encoding = utils.DetectEncoding(sample[:n])
	} else {
		buf := &bytes.Buffer{}
		tr := io.TeeReader(r, buf)
		n, _ := tr.Read(sample)
		profile.Encoding = utils.DetectEncoding(sample[:n])
		scanSrc = io.MultiReader(buf, r)
	}

	scn := bufio.NewScanner(scanSrc)
	scn.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	var lines []string
	lineCount := 0
	widthSum := 0
	for scn.Scan() {
		text := scn.Text()
		lineCount++
		widthSum += len(text)
		if len(lines) < maxSampleRows+2 {
			lines = append(lines, text)
		}
	}
	if err := scn.Err(); err != nil {
		profile.Errors = append(profile.Errors, err.Error())
	}
	profile.RowCountEstimate = lineCount
	if lineCount > 0 {
		profile.AverageRowWidth = widthSum / lineCount
	}
	if len(lines) == 0 {
		return profile
	}

	profile.Delimiter = detectors.DetectDelimiter(lines[:min(10, len(lines))])
	linesContent := strings.Join(lines, "\n")
	if linesContent != "" && !strings.HasSuffix(linesContent, "\n") {
		linesContent += "\n"
	}
	csvReader := csv.NewReader(strings.NewReader(linesContent))
	csvReader.Comma = rune(profile.Delimiter[0])
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true

	var rows [][]string
	for i := 0; i < maxSampleRows+1; i++ {
		record, err := csvReader.Read()
		if err != nil {
			break
		}
		for j := range record {
			record[j] = strings.TrimSpace(record[j])
		}
		rows = append(rows, record)
	}
	if len(rows) == 0 {
		return profile
	}
	profile.ColumnCount = len(rows[0])
	second := []string{}
	if len(rows) > 1 {
		second = rows[1]
	}
	profile.HasHeader = detectors.LooksLikeHeader(rows[0], second)
	profile.Columns = profiler.ProfileColumns(discovered.Path, "", rows, profile.HasHeader)
	return profile
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
