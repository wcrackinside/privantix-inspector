package profiler

import (
	"fmt"
	"strings"

	"privantix-source-inspector/detectors"
	"privantix-source-inspector/models"
)

func ProfileColumns(filePath, sheetName string, rows [][]string, hasHeader bool) []models.ColumnProfile {
	if len(rows) == 0 {
		return nil
	}
	start := 0
	headers := rows[0]
	if hasHeader {
		start = 1
	} else {
		headers = make([]string, len(rows[0]))
		for i := range headers {
			headers[i] = fmt.Sprintf("column_%d", i+1)
		}
	}

	profiles := make([]models.ColumnProfile, len(headers))
	for i, h := range headers {
		profiles[i] = models.ColumnProfile{
			FilePath:  filePath,
			SheetName: sheetName,
			Name:      strings.TrimSpace(h),
			Position:  i + 1,
		}
	}

	counts := make([]map[string]int, len(headers))
	nulls := make([]int, len(headers))
	totals := make([]int, len(headers))
	for i := range counts {
		counts[i] = map[string]int{}
	}

	for _, row := range rows[start:] {
		for i := range headers {
			val := ""
			if i < len(row) {
				val = strings.TrimSpace(row[i])
			}
			totals[i]++
			if detectors.IsNull(val) {
				nulls[i]++
				continue
			}
			t := detectors.InferType(val)
			counts[i][t]++
			if len(val) > profiles[i].MaxLength {
				profiles[i].MaxLength = len(val)
			}
			if len(profiles[i].SampleValues) < 3 {
				profiles[i].SampleValues = append(profiles[i].SampleValues, val)
			}
		}
	}

	for i := range profiles {
		profiles[i].InferredType = dominantType(counts[i])
		if totals[i] > 0 {
			profiles[i].NullPercentage = float64(nulls[i]) / float64(totals[i]) * 100
		}
		if profiles[i].Name == "" {
			profiles[i].Name = fmt.Sprintf("column_%d", i+1)
		}
	}
	return profiles
}

func dominantType(m map[string]int) string {
	bestType := "string"
	bestCount := -1
	for t, c := range m {
		if c > bestCount {
			bestType = t
			bestCount = c
		}
	}
	return bestType
}
