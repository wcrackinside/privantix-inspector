package parquetanalyzer

import (
	"io"
	"os"
	"strings"

	"github.com/parquet-go/parquet-go"

	"privantix-source-inspector/detectors"
	"privantix-source-inspector/models"
	"privantix-source-inspector/profiler"
)

// AnalyzeFromReaderAt analyzes Parquet from io.ReaderAt (e.g. file inside archive).
func AnalyzeFromReaderAt(readerAt io.ReaderAt, size int64, discovered models.FileDiscovered, maxSampleRows int) models.FileProfile {
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
		Encoding:    "parquet",
	}

	pf, err := parquet.OpenFile(readerAt, size)
	if err != nil {
		profile.Errors = append(profile.Errors, err.Error())
		return profile
	}

	schema := pf.Schema()
	rowGroups := pf.RowGroups()
	var totalRows int64
	for _, rg := range rowGroups {
		totalRows += rg.NumRows()
	}
	profile.RowCountEstimate = int(totalRows)

	fields := schema.Fields()
	colNames := make([]string, len(fields))
	for i, f := range fields {
		colNames[i] = strings.TrimSuffix(f.Name(), ".")
	}
	if len(colNames) == 0 {
		return profile
	}
	profile.ColumnCount = len(colNames)

	reader := parquet.NewReader(pf)
	defer reader.Close()

	limit := maxSampleRows + 1
	if limit <= 0 {
		limit = 201
	}
	rows := make([][]string, 0, limit)
	rowBuf := make([]parquet.Row, 1)
	widthSum := 0

	for len(rows) < limit {
		n, err := reader.ReadRows(rowBuf)
		if n == 0 || err != nil {
			break
		}
		row := make([]string, len(colNames))
		for i, val := range rowBuf[0] {
			if i < len(colNames) {
				row[i] = valueToString(val)
				widthSum += len(row[i])
			}
		}
		rows = append(rows, row)
	}

	if len(rows) > 0 {
		profile.AverageRowWidth = widthSum / (len(rows) * max(1, len(colNames)))
		second := []string{}
		if len(rows) > 1 {
			second = rows[1]
		}
		profile.HasHeader = detectors.LooksLikeHeader(rows[0], second)
		profile.Columns = profiler.ProfileColumns(discovered.Path, "", rows, profile.HasHeader)
	}

	return profile
}

func Analyze(discovered models.FileDiscovered, maxSampleRows int) models.FileProfile {
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
		Encoding:    "parquet",
	}

	f, err := os.Open(discovered.Path)
	if err != nil {
		profile.Errors = append(profile.Errors, err.Error())
		return profile
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		profile.Errors = append(profile.Errors, err.Error())
		return profile
	}

	pf, err := parquet.OpenFile(f, stat.Size())
	if err != nil {
		profile.Errors = append(profile.Errors, err.Error())
		return profile
	}

	schema := pf.Schema()
	rowGroups := pf.RowGroups()
	var totalRows int64
	for _, rg := range rowGroups {
		totalRows += rg.NumRows()
	}
	profile.RowCountEstimate = int(totalRows)

	// Get column names from schema
	fields := schema.Fields()
	colNames := make([]string, len(fields))
	for i, f := range fields {
		colNames[i] = strings.TrimSuffix(f.Name(), ".")
	}
	if len(colNames) == 0 {
		return profile
	}
	profile.ColumnCount = len(colNames)

	// Read sample rows using parquet.Row
	reader := parquet.NewReader(pf)
	defer reader.Close()

	limit := maxSampleRows + 1
	if limit <= 0 {
		limit = 201
	}
	rows := make([][]string, 0, limit)
	rowBuf := make([]parquet.Row, 1)
	widthSum := 0

	for len(rows) < limit {
		n, err := reader.ReadRows(rowBuf)
		if n == 0 || err != nil {
			break
		}
		row := make([]string, len(colNames))
		for i, val := range rowBuf[0] {
			if i < len(colNames) {
				row[i] = valueToString(val)
				widthSum += len(row[i])
			}
		}
		rows = append(rows, row)
	}

	if len(rows) > 0 {
		profile.AverageRowWidth = widthSum / (len(rows) * max(1, len(colNames)))
		second := []string{}
		if len(rows) > 1 {
			second = rows[1]
		}
		profile.HasHeader = detectors.LooksLikeHeader(rows[0], second)
		profile.Columns = profiler.ProfileColumns(discovered.Path, "", rows, profile.HasHeader)
	}

	return profile
}

func valueToString(v parquet.Value) string {
	if !v.IsNull() {
		return v.String()
	}
	return ""
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
