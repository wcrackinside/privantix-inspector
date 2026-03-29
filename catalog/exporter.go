package catalog

import (
	"encoding/csv"
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
)

// ExportAll writes catalog to JSON, CSV and HTML.
func ExportAll(outputDir string, catalog *CatalogResult) error {
	if err := ensureDir(outputDir); err != nil {
		return err
	}
	if err := exportJSON(filepath.Join(outputDir, "catalog.json"), catalog); err != nil {
		return err
	}
	if err := exportDatasetsCSV(filepath.Join(outputDir, "catalog_datasets.csv"), catalog); err != nil {
		return err
	}
	if err := exportColumnsCSV(filepath.Join(outputDir, "catalog_columns.csv"), catalog); err != nil {
		return err
	}
	if err := exportHTML(filepath.Join(outputDir, "catalog.html"), catalog); err != nil {
		return err
	}
	return nil
}

func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0o755)
}

func exportJSON(path string, catalog *CatalogResult) error {
	b, err := json.MarshalIndent(catalog, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func exportDatasetsCSV(path string, catalog *CatalogResult) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	_ = w.Write([]string{"path", "name", "extension", "owner", "permissions", "encoding", "delimiter", "has_header", "row_count", "column_count", "checksum", "has_acls", "rules_count"})
	for _, d := range catalog.Datasets {
		rules := strconv.Itoa(len(d.RulesTriggered))
		_ = w.Write([]string{
			d.Path,
			d.Name,
			d.Extension,
			d.Owner,
			d.Permissions,
			d.Encoding,
			d.Delimiter,
			strconv.FormatBool(d.HasHeader),
			strconv.Itoa(d.RowCountEstimate),
			strconv.Itoa(d.ColumnCount),
			d.Checksum,
			strconv.FormatBool(d.HasACLs),
			rules,
		})
	}
	return nil
}

func exportColumnsCSV(path string, catalog *CatalogResult) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	_ = w.Write([]string{"dataset_path", "column_name", "position", "inferred_type", "null_percentage", "max_length", "has_samples"})
	for _, d := range catalog.Datasets {
		for _, c := range d.Columns {
			_ = w.Write([]string{
				d.Path,
				c.Name,
				strconv.Itoa(c.Position),
				c.InferredType,
				strconv.FormatFloat(c.NullPercentage, 'f', 2, 64),
				strconv.Itoa(c.MaxLength),
				strconv.FormatBool(c.HasSamples),
			})
		}
	}
	return nil
}

func exportHTML(path string, catalog *CatalogResult) error {
	tpl, err := template.New("catalog").Parse(HTMLTemplate)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return tpl.Execute(f, catalog)
}
