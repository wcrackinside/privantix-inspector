# Changelog

All notable changes to Privantix Source Inspector are documented in this file.

## [Unreleased]

### Added
- Data Inspector: CSV, TXT, TSV, XLSX, RDAT, DAT, Parquet support
- Data Inspector: Archive support (ZIP, 7z, RAR) — in-memory analysis without extraction
- Data Inspector: Column profiling with inferred types (string, integer, float, date, email, phone, etc.)
- Data Inspector: Rule engine (encoding, headers, sensitive data, null ratio)
- Data Inspector: `--hide-samples` for data governance (omit sample values)
- Data Inspector: Export to JSON, CSV, HTML
- ACL Audit: Folder permissions and ACL audit
- ACL Audit: Trusted groups (JSON) for compliance classification
- ACL Audit: Export to JSON, CSV

### Technical
- Go 1.24+ required
- Worker pool for parallel file analysis (Data Inspector)
- Platform-specific ACL retrieval (Windows: `icacls`)
