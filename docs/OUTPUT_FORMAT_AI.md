# Privantix Output Format — AI-Readable Documentation

This document describes the output files produced by **privantix-inspector** and **privantix-acl-audit** so that an AI or automated system can parse and interpret them correctly.

**JSON Schema:** `docs/analysis-result.schema.json` (JSON Schema draft-07 for validation)

---

## Instrucciones para otra IA: cómo leer el JSON del inspector

**Archivo a leer:** `analysis.json` (o `{baseName}.json` en la carpeta de salida del comando `privantix-inspector scan`).

1. **Estructura raíz:** El JSON tiene tres claves: `run` (metadatos del escaneo), `files` (array de perfiles por archivo), `errors` (array de strings con errores del escaneo).
2. **Metadatos:** En `run` encontrarás `started_at`, `completed_at`, `path` (ruta escaneada), `analyzed_files`, `failed_files`, `supported_extensions`.
3. **Cada elemento de `files`** es un **FileProfile**: `path`, `name`, `extension`, `size_bytes`, `modified_at`, `created_at` (opcional), `owner`, `encoding`, `column_count`, `row_count_estimate`, `columns` (array de columnas con `name`, `inferred_type`, `null_percentage`, `sample_values`), `rules_triggered`, `errors`.
4. **Formatos de fecha:** Todas las fechas en JSON son RFC3339 (ej. `2026-03-10T00:57:16.1931827-03:00`).
5. **Campos opcionales:** Pueden faltar `acls`, `trusted_principals`, `outside_principals`, `compliance_status`, `checksum`, `created_at`, `sheet_name`; si no están, no están presentes o son array/string vacío.
6. **Reglas de calidad:** En `rules_triggered` cada elemento tiene `rule_name`, `severity`, `message`, `target` (opcional). Nombres típicos: `encoding_not_utf8`, `missing_header`, `potential_sensitive_data`, etc.
7. **Tipos inferidos de columna:** En cada `columns[].inferred_type` los valores son: `string`, `integer`, `float`, `boolean`, `date`, `datetime`, `email`, `phone`, `id_like`.

Para validación estricta del formato, usa el esquema JSON en `docs/analysis-result.schema.json`. Para modo streaming (`.jsonl`), cada línea es un objeto JSON: primero uno con `"type":"run_start"`, luego varios con `"type":"file"` y `"file": <FileProfile>`, y al final uno con `"type":"run_end"`.

---

## 1. privantix-inspector Output Files

### 1.1 File Naming

| Base Name | JSON | Files CSV | Columns CSV | HTML Report |
|-----------|------|-----------|-------------|-------------|
| Default (empty) | `analysis.json` | `files.csv` | `columns.csv` | `report.html` |
| Custom (e.g. `arc_20260215`) | `arc_20260215.json` | `arc_20260215_files.csv` | `arc_20260215_columns.csv` | `arc_20260215_report.html` |

**Streaming mode (`--stream`):** Uses `analysis.jsonl` (or `{baseName}.jsonl`) instead of JSON. No HTML report. CSV files are written incrementally. Reduces memory usage for large scans.

---

### 1.2 JSON Schema (`analysis.json` or `{baseName}.json`)

```json
{
  "run": { /* RunMetadata */ },
  "files": [ /* FileProfile[] */ ],
  "errors": [ /* string[] - scan/analysis errors */ ]
}
```

#### RunMetadata

| Field | Type | Description |
|-------|------|-------------|
| `started_at` | string (RFC3339) | Scan start timestamp |
| `completed_at` | string (RFC3339) | Scan completion timestamp |
| `path` | string | Root path scanned |
| `output_dir` | string | Output directory used |
| `supported_files` | int | Number of files discovered (matching extensions) |
| `analyzed_files` | int | Number of files successfully analyzed |
| `failed_files` | int | Number of files that failed analysis |
| `max_sample_rows` | int | Max rows sampled per file (default: 200) |
| `workers` | int | Parallel workers used |
| `supported_extensions` | string[] | Extensions scanned (e.g. [".csv", ".xlsx", ".parquet"]) |

#### FileProfile

| Field | Type | Description |
|-------|------|-------------|
| `path` | string | Full file path |
| `name` | string | File name |
| `extension` | string | Extension (e.g. ".csv", ".xlsx") |
| `size_bytes` | int64 | File size in bytes |
| `modified_at` | string (RFC3339) | Last modification time |
| `created_at` | string (RFC3339, optional) | File creation time (when available) |
| `depth` | int | Directory depth from scan root (0 = root) |
| `owner` | string | File owner (e.g. "MDP\\carce") |
| `permissions` | string | Unix-style permissions (e.g. "-rw-rw-rw-") |
| `acls` | string[] | Access control list entries (when --security used) |
| `trusted_principals` | string[] | Principals in trusted groups (when --trusted-groups used) |
| `outside_principals` | string[] | Principals outside trusted groups |
| `compliance_status` | string | "compliant" or "non_compliant" |
| `checksum` | string | SHA256 hash (when --checksum used) |
| `encoding` | string | Detected encoding (utf-8, windows-1252, etc.) |
| `delimiter` | string | CSV delimiter (e.g. ";", ",") |
| `has_header` | bool | Whether header row was detected |
| `row_count_estimate` | int | Estimated row count |
| `column_count` | int | Number of columns |
| `average_row_width` | int | Average characters per row |
| `sheet_name` | string | Excel sheet name (for .xlsx) |
| `columns` | ColumnProfile[] | Per-column profiles |
| `rules_triggered` | RuleResult[] | Quality/compliance rules that fired |
| `errors` | string[] | Analysis errors for this file |

#### ColumnProfile

| Field | Type | Description |
|-------|------|-------------|
| `file_path` | string | Parent file path |
| `sheet_name` | string | Sheet name (Excel) or empty |
| `name` | string | Column name |
| `position` | int | 1-based column index |
| `inferred_type` | string | See InferredType values below |
| `null_percentage` | float | Percentage of null/empty values |
| `max_length` | int | Max value length in sample |
| `sample_values` | string[] | Up to 3 sample values (null if --hide-samples) |

**InferredType values:** `string`, `integer`, `float`, `boolean`, `date`, `datetime`, `email`, `phone`, `id_like`

#### RuleResult

| Field | Type | Description |
|-------|------|-------------|
| `rule_name` | string | See Rule names below |
| `severity` | string | "info", "warning" |
| `message` | string | Human-readable description |
| `target` | string | Column name (when rule applies to column) |

**Rule names:** `encoding_not_utf8`, `missing_header`, `too_many_columns`, `high_null_ratio`, `potential_sensitive_data`

---

### 1.2b JSONL Schema (streaming mode, `--stream`)

When using `--stream`, output is `analysis.jsonl` (or `{baseName}.jsonl`). Each line is a JSON object:

1. **First line** — `{"type":"run_start", "started_at":"...", "path":"...", "output_dir":"...", "supported_files":N, "supported_extensions":[...], "max_sample_rows":N, "workers":N}`
2. **File lines** — `{"type":"file", "file":{ /* FileProfile */ }}`
3. **Last line** — `{"type":"run_end", "completed_at":"...", "analyzed_files":N, "failed_files":M, "errors":[...]}`

---

### 1.3 Files CSV Schema (`files.csv` or `{baseName}_files.csv`)

CSV with header. Delimiter: comma. Multi-value fields use ` | ` (space-pipe-space) or `;` for rules.

| Column | Description |
|--------|-------------|
| path | Full file path |
| name | File name |
| extension | e.g. .csv |
| owner | File owner |
| permissions | Unix-style |
| acls | ACL entries joined by ` \| ` |
| checksum | SHA256 (if --checksum) |
| size_bytes | File size |
| modified_at | YYYY-MM-DD HH:MM:SS |
| depth | Directory depth |
| encoding | Detected encoding |
| delimiter | CSV delimiter |
| has_header | true/false |
| row_count_estimate | Row count |
| column_count | Column count |
| average_row_width | Avg row length |
| sheet_name | Excel sheet (or empty) |
| rules_triggered | Rule names joined by `;` |
| errors | Error messages joined by `;` |
| trusted_principals | (optional) Joined by ` \| ` |
| outside_principals | (optional) Joined by ` \| ` |
| compliance_status | (optional) compliant/non_compliant |

---

### 1.4 Columns CSV Schema (`columns.csv` or `{baseName}_columns.csv`)

| Column | Description |
|--------|-------------|
| file_path | Parent file path |
| sheet_name | Sheet name or empty |
| column_name | Column name |
| position | 1-based index |
| inferred_type | string, integer, float, etc. |
| null_percentage | Float (e.g. 0.00, 22.50) |
| max_length | Max value length |
| sample_values | Sample values joined by ` \| ` |

---

### 1.5 HTML Report (`report.html` or `{baseName}_report.html`)

Human-readable HTML summary. Contains the same data as the JSON in tabular form. Structure may vary by version.

---

## 2. privantix-acl-audit Output Files

**Documentación detallada:** `docs/acl-audit.md` · **Formato para IA:** `docs/ACL_AUDIT_OUTPUT_AI.md` · **Schema:** `docs/acl-audit-result.schema.json` · **Instrucciones cortas para IA:** `docs/INSTRUCTIONS_AI_ACL_AUDIT.md`

### 2.1 File Naming

| Base Name | JSON | CSV |
|-----------|------|-----|
| Default | `{path_basename}_{YYYYMMDD}.json` | `{path_basename}_{YYYYMMDD}.csv` |
| Custom (--output-name X) | `X.json` | `X.csv` |

### 2.2 ACL Audit JSON Schema

```json
{
  "started_at": "string (RFC3339)",
  "completed_at": "string (RFC3339)",
  "duration_seconds": "number",
  "path": "string - root path audited",
  "output_dir": "string",
  "total_dirs": "int",
  "trusted_groups": [ { "name": "string", "members": ["string"] } ],
  "dirs": [ /* DirACL[] */ ],
  "errors": ["string"]
}
```

#### DirACL

| Field | Type | Description |
|-------|------|-------------|
| path | string | Directory path |
| name | string | Directory name |
| depth | int | Depth from root |
| owner | string | Directory owner |
| permissions | string | Unix-style permissions |
| acls | string[] | ACL entries |
| trusted_principals | string[] | Principals in trusted groups |
| outside_principals | string[] | Principals outside trusted groups |
| compliance_status | string | "compliant" or "non_compliant" |

### 2.3 ACL Audit CSV

Headers: path, name, depth, owner, permissions, acls, trusted_principals, outside_principals, compliance_status

ACL format (Windows): `DOMAIN\user:(I)(F)` — principal followed by colon and permission flags.

---

## 3. Timestamp Formats

- **JSON:** RFC3339 (e.g. `2026-03-10T00:57:16.1931827-03:00`)
- **CSV:** `2006-01-02 15:04:05` (e.g. `2023-07-24 22:26:37`)

---

## 4. Path Conventions

- Windows: Backslashes (e.g. `W:\11_CONADI\04_DATOS\file.csv`)
- UNC: `\\server\share\path`
- Archive entries: `archive.zip/inner/file.csv` (forward slash)

---

## 5. Empty and Null Values

- JSON: `null` for absent optional fields
- CSV: Empty string for missing values
- `sample_values` may be `null` when `--hide-samples` is used

---

## 6. Example: Minimal JSON (0 files)

```json
{
  "run": {
    "started_at": "2026-03-10T00:57:16.1931827-03:00",
    "completed_at": "2026-03-10T00:57:16.3387091-03:00",
    "path": "\\\\nas04\\share\\path",
    "output_dir": "./output",
    "supported_files": 0,
    "analyzed_files": 0,
    "failed_files": 0,
    "max_sample_rows": 200,
    "workers": 8,
    "supported_extensions": [".csv", ".txt", ".tsv", ".xlsx", ".parquet", ".7z", ".zip", ".rar"]
  },
  "files": [],
  "errors": ["GetFileAttributesEx ...: The system cannot find the file specified."]
}
```

---

## 7. privantix-compare Output

Compares two `analysis.json` files. Use `--format json` for machine-readable output.

### 7.1 JSON Schema (compare result)

```json
{
  "left_input": "path/to/run1.json",
  "right_input": "path/to/run2.json",
  "left_path": "W:\\scanned\\path\\1",
  "right_path": "W:\\scanned\\path\\2",
  "added": ["path/to/new/file.csv"],
  "removed": ["path/to/deleted/file.csv"],
  "modified": [
    {
      "path": "path/to/changed/file.csv",
      "changes": {
        "size_bytes": {"left": 1000, "right": 1200},
        "modified_at": {"left": "2024-01-01 00:00:00", "right": "2024-01-02 00:00:00"},
        "row_count_estimate": {"left": 100, "right": 120},
        "column_count": {"left": 5, "right": 6},
        "encoding": {"left": "utf-8", "right": "windows-1252"},
        "delimiter": {"left": ",", "right": ";"}
      }
    }
  ],
  "summary": {
    "added_count": 1,
    "removed_count": 1,
    "modified_count": 1,
    "unchanged_count": 10
  }
}
```

---

## 8. Parsing Tips for AI

1. **Compare runs:** Use `privantix-compare --left run1.json --right run2.json --format json` to get a structured diff.
2. **Identify file set:** Use `run.path` and `run.supported_extensions` to understand scan scope.
2. **Check for errors:** `run.errors` and `file.errors` indicate failures.
3. **Compliance:** `compliance_status` and `outside_principals` indicate access control issues.
4. **Data quality:** `rules_triggered` and `inferred_type` help assess data quality.
5. **Sensitive data:** `potential_sensitive_data` rule flags columns with emails, phones, IDs.
6. **Encoding:** `encoding_not_utf8` rule warns about non-UTF-8 files.
