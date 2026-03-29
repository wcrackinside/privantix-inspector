# Privantix Source Inspector

Portable command-line tools for analyzing tabular data repositories and auditing folder permissions.

## Tools

| Tool | Purpose |
|------|---------|
| **privantix-inspector** | Scans and profiles tabular data files (CSV, Excel, Parquet, archives) |
| **privantix-acl-audit** | Audits ACLs and permissions for folders and subfolders |
| **privantix-catalog** | Transforms inspector output into a structured, navigable dataset catalog |
| **privantix-compare** | Compares two analysis.json runs (added, removed, modified files) |

## Quick Start

```bash
# Build
go build -o privantix-inspector.exe ./cmd/privantix-inspector
go build -o privantix-acl-audit.exe ./cmd/privantix-acl-audit
go build -o privantix-catalog.exe ./cmd/privantix-catalog
go build -o privantix-compare.exe ./cmd/privantix-compare

# Data Inspector — analyze a data folder
privantix-inspector scan --path ./examples/data --output ./output

# Catalog — generate catalog from inspector output
privantix-catalog --input ./output/analysis.json --output ./catalog

# Compare — diff two inspector runs
privantix-compare --left ./output/run1.json --right ./output/run2.json --output compare.json

# ACL Audit — audit folder permissions
privantix-acl-audit --path W:\datos --output ./output --trusted-groups ./examples/trusted-groups.json
```

## Features

- **Data Inspector**: CSV, TXT, TSV, XLSX, RDAT, DAT, Parquet, ZIP, 7z, RAR (archives read in-memory)
- **Column profiling**: Inferred types, null %, sample values
- **Rule engine**: Encoding, headers, sensitive data patterns, null ratios
- **Data governance mode**: `--hide-samples` to omit sample values from reports
- **ACL Audit**: Folder permissions, trusted groups, compliance status

## Outputs

- **Data Inspector**: `analysis.json`, `files.csv`, `columns.csv`, `report.html`
- **privantix-catalog**: `catalog.json`, `catalog_datasets.csv`, `catalog_columns.csv`, `catalog.html`
- **privantix-compare**: text report or JSON diff
- **ACL Audit**: `{name}.json`, `{name}.csv`

## Documentation

Full documentation is in the [`docs/`](docs/) folder:

- [Overview](docs/overview.md)
- [Installation](docs/installation.md)
- [Usage](docs/usage.md)
- [Architecture](docs/architecture.md)
- [Outputs](docs/outputs.md)
- [Output Format (AI-readable)](docs/OUTPUT_FORMAT_AI.md) — schema and field descriptions for parsing results
- [FAQ](docs/faq.md)
- [Changelog](docs/changelog.md)

## Requirements

- Go 1.24+ for building
- Windows: full ACL support via `icacls`
- Linux/macOS: data inspector supported; ACL audit uses basic POSIX permissions
