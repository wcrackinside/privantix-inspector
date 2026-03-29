# Outputs

## Data Inspector (privantix-inspector)

All outputs are written to the directory specified by `--output` (default: `./output`).

### analysis.json

Complete analysis result in JSON format.

**Structure:**

```json
{
  "run": {
    "started_at": "2026-03-06T10:36:47.8890108-03:00",
    "completed_at": "2026-03-06T10:36:47.950079-03:00",
    "path": "./examples/data",
    "output_dir": "./output",
    "supported_files": 1,
    "analyzed_files": 1,
    "failed_files": 0,
    "max_sample_rows": 200,
    "workers": 4,
    "supported_extensions": [".csv", ".txt", ".tsv", ".xlsx", "..."]
  },
  "files": [
    {
      "path": "examples/data/sample.csv",
      "name": "sample.csv",
      "extension": ".csv",
      "size_bytes": 73,
      "owner": "BUILTIN\\Administradores",
      "permissions": "-rw-rw-rw-",
      "encoding": "utf-8",
      "delimiter": ";",
      "has_header": true,
      "row_count_estimate": 4,
      "column_count": 3,
      "columns": [
        {
          "name": "name",
          "position": 1,
          "inferred_type": "string",
          "null_percentage": 0,
          "max_length": 5,
          "sample_values": ["Ana", "Luis", "Carla"]
        }
      ],
      "rules_triggered": [
        {
          "rule_name": "potential_sensitive_data",
          "severity": "info",
          "message": "Column may contain personal or sensitive identifiers",
          "target": "email"
        }
      ]
    }
  ],
  "errors": null
}
```

With `--hide-samples`, `sample_values` is `null` for all columns.

### files.csv

One row per analyzed file. Columns include: path, name, extension, owner, permissions, acls, checksum, size_bytes, modified_at, depth, encoding, delimiter, has_header, row_count_estimate, column_count, average_row_width, sheet_name, rules_triggered, errors.

**Example:**

```csv
path,name,extension,owner,permissions,acls,checksum,size_bytes,modified_at,depth,encoding,delimiter,has_header,row_count_estimate,column_count,average_row_width,sheet_name,rules_triggered,errors
examples\data\sample.csv,sample.csv,.csv,BUILTIN\Administradores,-rw-rw-rw-,,,73,2026-03-06 00:09:26,0,utf-8,;,true,4,3,17,,potential_sensitive_data,
```

### columns.csv

One row per column across all files. Columns: file_path, sheet_name, column_name, position, inferred_type, null_percentage, max_length, sample_values.

**Example:**

```csv
file_path,sheet_name,column_name,position,inferred_type,null_percentage,max_length,sample_values
examples\data\sample.csv,,name,1,string,0.00,5,Ana | Luis | Carla
examples\data\sample.csv,,email,2,email,33.33,16,ana@example.com | luis@example.com
examples\data\sample.csv,,age,3,integer,0.00,2,30 | 29 | 31
```

With `--hide-samples`, the `sample_values` column is empty.

### report.html

Interactive HTML report. Contains:
- Run metadata (path, file counts, extensions)
- File list with summaries (path, extension, row count, column count, rules triggered)
- Column-level details (name, type, null %, sample values)
- Rules triggered per file/column

Open in a web browser. No server required.

---

## ACL Audit (privantix-acl-audit)

Output files use the base name from `--output-name` or `{path_basename}_{YYYYMMDD}`.

### {basename}.json

**Structure:**

```json
{
  "started_at": "2026-03-06T12:40:48.0855273-03:00",
  "completed_at": "2026-03-06T12:40:48.3668929-03:00",
  "path": "W:\\11_CONADI\\04_DATOS",
  "output_dir": "./output",
  "total_dirs": 2,
  "trusted_groups": [
    {
      "name": "dais",
      "members": ["MDP\\ASermeno", "MDP\\carce", "MDP\\lmenesesv"]
    }
  ],
  "dirs": [
    {
      "path": "W:\\11_CONADI\\04_DATOS",
      "name": "04_DATOS",
      "depth": 0,
      "owner": "MDP\\carce",
      "permissions": "-rwxrwxrwx",
      "acls": ["NT AUTHORITY\\SYSTEM:(I)(OI)(CI)(F)", "MDP\\carce:(I)(OI)(CI)(F)", "..."],
      "trusted_principals": ["MDP\\carce", "MDP\\ASermeno", "..."],
      "outside_principals": ["NT AUTHORITY\\SYSTEM", "MDP\\admin04", "..."],
      "compliance_status": "non_compliant"
    }
  ],
  "errors": []
}
```

`trusted_principals`, `outside_principals`, and `compliance_status` appear only when `--trusted-groups` is used.

### {basename}.csv

One row per directory. Columns: path, name, depth, owner, permissions, acls, and (when trusted groups are used) trusted_principals, outside_principals, compliance_status.
