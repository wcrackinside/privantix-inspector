# Privantix Source Inspector

## Short Description (Hero)

Analyze tabular data repositories and understand their structure, columns, and quality in minutes—without viewing sensitive content when needed.

---

## Product Description

Privantix Source Inspector scans local or network paths for tabular data files and produces technical profiles for each file and column. It detects encoding, delimiters, row and column counts, and infers data types (string, integer, float, date, email, phone, etc.). A built-in rule engine flags potential issues such as non-UTF-8 encoding, missing headers, sensitive data patterns, and high null ratios.

The tool supports CSV, TXT, TSV, XLSX, RDAT, DAT, Parquet, and compressed archives (ZIP, 7z, RAR). Archives are analyzed in-memory without extracting files to disk. Reports are exported to JSON, CSV, and HTML.

Data governance teams can use `--hide-samples` to omit sample values from reports, enabling cataloging and classification without exposing actual data. Optional `--checksum` and `--security` flags add file integrity and ACL information to the output.

**Who should use it:** Data governance teams, data engineers, and analysts who need to catalog, document, or assess data repositories before migration, integration, or compliance reviews.

---

## Key Capabilities

- Repository scanning (recursive, configurable extensions)
- Multi-format support: CSV, TXT, TSV, XLSX, RDAT, DAT, Parquet
- Archive inspection: ZIP, 7z, RAR (in-memory, no extraction)
- Technical metadata extraction (encoding, delimiter, row/column count)
- Column profiling with inferred types
- Rule engine (encoding, header, sensitive data, null ratio)
- Data governance mode (`--hide-samples`)
- SHA256 checksum per file (`--checksum`)
- ACL extraction (`--security`)
- Export to JSON, CSV, HTML

---

## Use Cases

- Auditing data repositories before migration or integration
- Documenting legacy data sources and column structures
- Preparing datasets for ingestion pipelines
- Cataloging data assets for governance without viewing sensitive content
- Assessing data quality (encoding, headers, null ratios) across large repositories

---

## Example Execution

```bash
privantix-inspector scan --path ./repository --output ./results
```

With data governance mode and security:

```bash
privantix-inspector scan --path C:\datos --output ./report --hide-samples --security --checksum
```

---

## Release Information

**Current Version:** v0.1.x

**Status:** MVP

**Highlights:**
- Repository scanning with recursive support
- Metadata extraction and column profiling
- Archive inspection (ZIP, 7z, RAR)
- Rule engine and governance-safe reporting
- JSON, CSV, and HTML export
