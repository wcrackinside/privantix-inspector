# Architecture

## Project Structure

```
privantix-source-inspector/
├── cmd/
│   ├── privantix-inspector/   # Data inspector entry point
│   └── privantix-acl-audit/   # ACL audit entry point
├── config/                    # Configuration loading (YAML-like)
├── scanner/                   # File discovery
├── inspector/                 # Analysis dispatch
├── csvanalyzer/               # CSV, TXT, TSV, RDAT, DAT
├── xlsxanalyzer/              # Excel (.xlsx)
├── parquetanalyzer/           # Parquet
├── archiveanalyzer/           # ZIP, 7z, RAR (in-memory)
├── profiler/                  # Column profiling
├── detectors/                 # Type inference, delimiter, header detection
├── analyzer/                  # Rule engine
├── exporter/                  # JSON, CSV, HTML export
├── report/                    # HTML template
├── models/                    # Data structures
├── utils/                     # Encoding, ACLs, owner (OS-specific)
├── examples/                  # Sample data and configs
└── docs/                      # Documentation
```

## Module Responsibilities

| Module | Responsibility |
|--------|----------------|
| **cmd/privantix-inspector** | CLI, flag parsing, orchestrates scan → inspect → export |
| **cmd/privantix-acl-audit** | CLI, directory walk, ACL collection, trusted group classification |
| **config** | Load YAML-like config (extensions, workers, max_sample_rows, hide_sample_values) |
| **scanner** | Walk directory tree, filter by extension, collect file metadata (owner, permissions) |
| **inspector** | Dispatch files to analyzers by extension; apply checksum/ACL when requested |
| **csvanalyzer** | Parse CSV/text; detect delimiter, encoding, header; profile columns |
| **xlsxanalyzer** | Parse XLSX (ZIP/XML); first worksheet; shared strings resolution |
| **parquetanalyzer** | Read Parquet schema and sample rows; profile columns |
| **archiveanalyzer** | Open ZIP/7z/RAR in-memory; analyze inner CSV/XLSX/Parquet |
| **profiler** | Build column profiles (inferred type, null %, max length, sample values) |
| **detectors** | Infer data type, detect delimiter, detect header row |
| **analyzer** | Apply rules (encoding, header, column count, null ratio, sensitive data) |
| **exporter** | Write analysis.json, files.csv, columns.csv, report.html |
| **report** | HTML template for report |
| **models** | FileDiscovered, FileProfile, ColumnProfile, AnalysisResult, etc. |
| **utils** | Encoding detection, ACL retrieval, owner resolution (OS-specific) |

## Execution Flow

### Data Inspector (privantix-inspector)

```
┌─────────────┐     ┌──────────┐     ┌───────────┐     ┌─────────┐
│   Scanner   │────▶│ Inspector│────▶│  Analyzer  │────▶│ Exporter│
│ (file walk) │     │(dispatch)│     │  (rules)   │     │(JSON,   │
└─────────────┘     └──────────┘     └───────────┘     │ CSV,    │
       │                    │                          │ HTML)   │
       │                    │                          └─────────┘
       ▼                    ▼
  FileDiscovered      FileProfile[]
  (path, ext,         (columns, types,
   owner, perms)       rules, etc.)
```

1. **Scanner**: Walks the directory tree, filters by extension, collects file metadata (owner, permissions). Returns `[]FileDiscovered`.
2. **Inspector**: For each file, dispatches to the appropriate analyzer (csv, xlsx, parquet, archive). Optionally adds checksum and ACLs. Returns `[]FileProfile`.
3. **Analyzer**: Applies hardcoded rules (encoding, headers, sensitive data, null ratio) to each profile.
4. **Exporter**: Writes `analysis.json`, `files.csv`, `columns.csv`, `report.html`.

Files are processed in parallel via a worker pool (default 4 workers).

### ACL Audit (privantix-acl-audit)

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────┐
│ filepath.Walk   │────▶│ For each dir:    │────▶│ Exporter │
│ (directories    │     │ - Get owner/perms│     │ (JSON,   │
│  only)          │     │ - Get ACLs       │     │  CSV)    │
└─────────────────┘     │ - Classify vs    │     └─────────┘
                        │   trusted groups │
                        └──────────────────┘
```

1. Walk directory tree (directories only).
2. For each directory: get owner, permissions, ACLs (via `icacls` on Windows).
3. If `--trusted-groups` provided: classify principals as trusted or outside; set compliance_status.
4. Export JSON and CSV.

## Supported File Types

| Extension | Analyzer | Notes |
|-----------|----------|-------|
| .csv, .txt, .tsv, .rdat, .dat | csvanalyzer | Delimiter auto-detection |
| .xlsx | xlsxanalyzer | First worksheet only |
| .parquet | parquetanalyzer | Schema + sample rows |
| .zip, .7z, .rar | archiveanalyzer | Reads in-memory; analyzes inner CSV/XLSX/Parquet |

## Concurrency

- **Data Inspector**: Uses a worker pool (default 4 workers) to analyze files in parallel.
- **ACL Audit**: Sequential (each directory requires an `icacls` call on Windows).

## Platform-Specific Code

- **utils/owner_windows.go**, **utils/owner_unix.go**: File owner resolution
- **utils/acl_windows.go**, **utils/acl_unix.go**: ACL retrieval (Windows: `icacls`; Unix: limited)
