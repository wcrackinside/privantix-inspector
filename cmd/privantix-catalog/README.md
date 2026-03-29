# privantix-catalog

Portable executable that transforms **privantix-inspector** output into a structured and navigable dataset catalog.

## Purpose

Reads `analysis.json` from privantix-inspector and generates a catalog that summarizes:

- Datasets (files) with technical metadata
- Columns with inferred types, null rates, max length
- Detected rules and governance indicators
- Searchable HTML report

## Build

From the project root:

```bash
go build -o privantix-catalog ./cmd/privantix-catalog
```

For a portable Windows executable:

```bash
GOOS=windows GOARCH=amd64 go build -o privantix-catalog.exe ./cmd/privantix-catalog
```

## Usage

```bash
privantix-catalog --input <path-to-analysis.json> [--output <dir>]
```

### Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `--input` | *(required)* | Path to `analysis.json` from privantix-inspector |
| `--output` | `./catalog_output` | Output directory for catalog files |

### Example workflow

```bash
# 1. Run inspector to analyze a repository
privantix-inspector scan --path ./data --output ./output

# 2. Generate catalog from inspector output
privantix-catalog --input ./output/analysis.json --output ./catalog

# 3. Open the HTML catalog
# Open ./catalog/catalog.html in a browser
```

## Output files

| File | Description |
|------|-------------|
| `catalog.json` | Full catalog in JSON format |
| `catalog_datasets.csv` | One row per dataset (path, metadata, counts) |
| `catalog_columns.csv` | One row per column across all datasets |
| `catalog.html` | Interactive HTML report with search |

## Catalog structure

Each dataset entry includes:

- **Metadata**: path, name, extension, owner, permissions, encoding, delimiter
- **Counts**: row estimate, column count, average row width
- **Governance**: checksum presence, ACL presence, samples hidden flag
- **Columns**: name, inferred type, null %, max length
- **Rules**: triggered rules (e.g. sensitive data, encoding validation)
