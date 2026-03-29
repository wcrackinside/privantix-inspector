package xlsxanalyzer

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"privantix-source-inspector/detectors"
	"privantix-source-inspector/models"
	"privantix-source-inspector/profiler"
)

type sharedStrings struct {
	SI []struct {
		T string `xml:"t"`
		R []struct {
			T string `xml:"t"`
		} `xml:"r"`
	} `xml:"si"`
}

type worksheet struct {
	Rows []sheetRow `xml:"sheetData>row"`
}

type sheetRow struct {
	Cells []sheetCell `xml:"c"`
}

type sheetCell struct {
	T  string  `xml:"t,attr"`
	V  string  `xml:"v"`
	IS inlineS `xml:"is"`
}

type inlineS struct {
	T string `xml:"t"`
}

func Analyze(discovered models.FileDiscovered, maxSampleRows int) models.FileProfile {
	zr, err := zip.OpenReader(discovered.Path)
	if err != nil {
		return models.FileProfile{
			Path: discovered.Path, Name: discovered.Name, Extension: discovered.Ext,
			SizeBytes: discovered.Size, ModifiedAt: discovered.Modified, CreatedAt: discovered.Created,
			Depth: discovered.Depth, Owner: discovered.Owner, Permissions: discovered.Permissions,
			Errors: []string{err.Error()},
		}
	}
	defer zr.Close()
	return analyzeFromZipFiles(zr.File, discovered, maxSampleRows)
}

// AnalyzeFromZipReader analyzes XLSX from a zip.Reader (e.g. xlsx inside archive).
func AnalyzeFromZipReader(zr *zip.Reader, discovered models.FileDiscovered, maxSampleRows int) models.FileProfile {
	return analyzeFromZipFiles(zr.File, discovered, maxSampleRows)
}

func analyzeFromZipFiles(files []*zip.File, discovered models.FileDiscovered, maxSampleRows int) models.FileProfile {
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
		Encoding:    "xlsx-internal",
	}

	shared := loadSharedStrings(files)
	rows, sheetName, err := loadFirstSheetRows(files, shared, maxSampleRows+1)
	if err != nil {
		profile.Errors = append(profile.Errors, err.Error())
		return profile
	}
	profile.SheetName = sheetName
	profile.RowCountEstimate = len(rows)
	if len(rows) == 0 {
		return profile
	}
	widthSum := 0
	for _, row := range rows {
		widthSum += len(strings.Join(row, ""))
	}
	profile.AverageRowWidth = widthSum / max(1, len(rows))
	profile.ColumnCount = len(rows[0])
	second := []string{}
	if len(rows) > 1 {
		second = rows[1]
	}
	profile.HasHeader = detectors.LooksLikeHeader(rows[0], second)
	profile.Columns = profiler.ProfileColumns(discovered.Path, sheetName, rows, profile.HasHeader)
	return profile
}

func loadSharedStrings(files []*zip.File) []string {
	for _, f := range files {
		if filepath.ToSlash(f.Name) == "xl/sharedStrings.xml" {
			rc, err := f.Open()
			if err != nil {
				return nil
			}
			defer rc.Close()
			b, _ := io.ReadAll(rc)
			var ss sharedStrings
			if xml.Unmarshal(b, &ss) != nil {
				return nil
			}
			result := make([]string, 0, len(ss.SI))
			for _, si := range ss.SI {
				if si.T != "" {
					result = append(result, si.T)
					continue
				}
				var sb strings.Builder
				for _, r := range si.R {
					sb.WriteString(r.T)
				}
				result = append(result, sb.String())
			}
			return result
		}
	}
	return nil
}

func loadFirstSheetRows(files []*zip.File, shared []string, limit int) ([][]string, string, error) {
	for _, f := range files {
		name := filepath.ToSlash(f.Name)
		if strings.HasPrefix(name, "xl/worksheets/") && strings.HasSuffix(name, ".xml") {
			rc, err := f.Open()
			if err != nil {
				return nil, "", err
			}
			defer rc.Close()
			b, _ := io.ReadAll(rc)
			var ws worksheet
			if err := xml.Unmarshal(b, &ws); err != nil {
				return nil, "", err
			}
			rows := make([][]string, 0, len(ws.Rows))
			for _, row := range ws.Rows {
				out := make([]string, 0, len(row.Cells))
				for _, cell := range row.Cells {
					val := strings.TrimSpace(cell.V)
					switch cell.T {
					case "s":
						idx, _ := strconv.Atoi(val)
						if idx >= 0 && idx < len(shared) {
							val = shared[idx]
						}
					case "inlineStr":
						val = cell.IS.T
					}
					out = append(out, strings.TrimSpace(val))
				}
				rows = append(rows, out)
				if len(rows) >= limit {
					break
				}
			}
			return rows, filepath.Base(name), nil
		}
	}
	return nil, "", fmt.Errorf("no worksheet xml found")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
