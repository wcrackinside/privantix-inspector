# Privantix Source Inspector — Overview

## What It Is

Privantix Source Inspector is a portable command-line tool for analyzing tabular data repositories. It helps data governance teams, analysts, and engineers understand the structure, content, and quality of data stored in files and folders—without requiring access to the actual data values when needed.

The product includes two main tools:

1. **privantix-inspector** — Scans and profiles tabular data files (CSV, Excel, Parquet, archives, etc.).
2. **privantix-acl-audit** — Audits folder permissions and ACLs (Access Control Lists) across a directory tree.

## Key Capabilities

### Data Inspector (privantix-inspector)

- **Recursive scanning** of local or network paths
- **Multi-format support**: CSV, TXT, TSV, XLSX, RDAT, DAT, Parquet, and compressed archives (ZIP, 7z, RAR) — archives are read in-memory without extracting to disk
- **Technical profiling** per file: encoding, delimiter, row count, column count
- **Column profiling** with inferred data types (string, integer, float, date, email, phone, etc.)
- **Rule engine** that flags potential issues (encoding, missing headers, sensitive data patterns, high null ratios)
- **Data governance mode** (`--hide-samples`) to omit sample values from reports
- **Export** to JSON, CSV, and HTML

### ACL Audit (privantix-acl-audit)

- **Folder-only audit** of permissions and ACLs
- **Trusted groups** support: define which users/groups are authorized; the tool flags principals outside those groups
- **Compliance status** per folder (compliant / non_compliant)
- **Export** to JSON and CSV

## Target Users

- **Data governance teams** — Catalog and classify data assets without viewing sensitive content
- **Data engineers** — Understand repository structure and data quality before migration or integration
- **Security / compliance** — Audit folder access permissions and identify unauthorized principals

## Platform Support

- **Windows**: Full support including ACL audit via `icacls`
- **Linux / macOS**: Data inspector supported; ACL audit uses basic POSIX permissions (extended ACL support may be limited)

## Documentation

- **Usage:** [usage.md](usage.md) (inspector and ACL audit parameters)
- **Inspector output (for humans & AI):** [OUTPUT_FORMAT_AI.md](OUTPUT_FORMAT_AI.md), [INSTRUCTIONS_AI_READ_JSON.md](INSTRUCTIONS_AI_READ_JSON.md)
- **ACL Audit (full doc):** [acl-audit.md](acl-audit.md)
- **ACL Audit output (for AI):** [ACL_AUDIT_OUTPUT_AI.md](ACL_AUDIT_OUTPUT_AI.md), [INSTRUCTIONS_AI_ACL_AUDIT.md](INSTRUCTIONS_AI_ACL_AUDIT.md)
- **Prompt para crear módulo de importación de los JSON:** [PROMPT_CREATE_IMPORT_JSON_MODULE.md](PROMPT_CREATE_IMPORT_JSON_MODULE.md)

## License & Requirements

- Go 1.24 or later for building
- No external runtime dependencies; produces standalone executables
