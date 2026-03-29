# Prompt for construction AI

Act as a senior Go software architect and developer. Continue building the project in this repository called **Privantix Source Inspector**.

## Product goal
Create a portable executable that scans a repository path recursively and profiles supported tabular files. The tool must analyze the content of a repository and report metadata, encoding, separators, header presence and inferred data types.

## MVP 1 requirements
- Portable executable in Go
- Recursive scan over local or network paths
- Supported formats: CSV, TXT, TSV, XLSX
- Per-file metadata: full path, file name, extension, size, modified date, folder depth
- Tabular metadata: encoding, delimiter, header presence, estimated row count, column count, average row width
- Per-column profiling: column name, inferred type, null percentage, max length, sample values
- Inferred types: string, integer, float, boolean, date, datetime, email, phone, id_like
- Non-blocking error handling
- Outputs: analysis.json, files.csv, columns.csv, report.html
- Parallel processing with configurable workers

## Technical constraints
- Keep the binary portable
- Prefer standard library when possible
- Keep the architecture modular
- Keep comments in English
- Preserve the existing contracts and models unless there is a compelling reason to evolve them

## Current repository intent
This baseline already contains:
- CLI command
- scanner
- inspector
- csv analyzer
- xlsx analyzer
- profiler
- detectors
- rule engine
- exporters
- docs

## Next recommended tasks
1. Improve encoding detection beyond BOM/UTF-8 fallback
2. Improve header detection with better confidence scoring
3. Add unit tests for delimiter detection and type inference
4. Add HTML styling and summary charts
5. Add optional YAML-configurable rules in MVP 2
6. Add JSON and Parquet analyzers in later versions

## Expected output from the construction AI
- Keep the project buildable
- Extend tests
- Improve analyzers incrementally
- Do not replace the architecture with a monolith
