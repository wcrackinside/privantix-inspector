# Prompt: Crear módulo para importar los JSON de Privantix

Usa este documento (o el bloque de texto siguiente) para pedirle a una IA que **cree un módulo que importe y parsee los archivos JSON** generados por **privantix-inspector** y **privantix-acl-audit**, de forma que puedas usarlos en código (Python, Node, etc.) o en pipelines de datos.

---

## Texto para pegar a la IA

```
Necesito un módulo (en [LENGUAJE]) que importe y parsee los JSON de salida de dos herramientas:

1) **privantix-inspector** — archivo analysis.json (o .jsonl en modo stream)
2) **privantix-acl-audit** — archivo {nombre}.json de auditoría de ACL

Especificaciones para implementar el módulo:

---

### 1. Inspector (analysis.json)

- **Raíz:** objeto con "run", "files", "errors".
- **run:** started_at, completed_at (RFC3339), path, output_dir, supported_files, analyzed_files, failed_files, max_sample_rows, workers, supported_extensions.
- **files:** array de FileProfile. Cada FileProfile tiene:
  - path, name, extension, size_bytes, depth, modified_at, created_at (opcional)
  - owner, permissions, acls (opcional), trusted_principals, outside_principals, compliance_status (opcionales)
  - encoding, delimiter, has_header, row_count_estimate, column_count, average_row_width, sheet_name
  - columns: array de { file_path, sheet_name, name, position, inferred_type, null_percentage, max_length, sample_values }
  - rules_triggered: array de { rule_name, severity, message, target }
  - errors: array de strings
- **errors:** array de strings (errores globales del escaneo).
- Parsear fechas en RFC3339 a tipo fecha/hora nativo del lenguaje.
- Si el archivo es .jsonl: cada línea es un JSON; hay líneas con "type":"run_start", "type":"file" (con "file": FileProfile), "type":"run_end". El módulo debe poder reconstruir la lista de files y los metadatos de run desde el JSONL.

### 2. ACL Audit (audit result JSON)

- **Raíz:** started_at, completed_at (RFC3339), duration_seconds, path, output_dir, total_dirs, trusted_groups (opcional), dirs, errors.
- **dirs:** array de DirACL. Cada DirACL: path, name, depth, owner, permissions, acls (array de strings), y si aplica trusted_principals, outside_principals, compliance_status ("compliant" | "non_compliant").
- **errors:** array de strings.
- Parsear started_at y completed_at a tipo fecha/hora nativo.

### 3. Contrato del módulo

- Función o clase para cargar **inspector JSON**: recibe ruta al .json (o .jsonl); devuelve un objeto/estructura con run, files, errors (y en files cada item con columns y rules_triggered como listas).
- Función o clase para cargar **ACL audit JSON**: recibe ruta al .json; devuelve un objeto/estructura con metadatos (started_at, path, total_dirs, etc.) y dirs (lista de directorios con acls, compliance_status, etc.).
- Opcional: validación contra JSON Schema si el lenguaje lo permite (schemas en docs/analysis-result.schema.json y docs/acl-audit-result.schema.json).
- Opcional: exportar "files" del inspector o "dirs" del ACL audit a tabla/DataFrame para análisis (pandas, etc.) si es Python.

Genera el código del módulo, tipos/datos bien definidos (clases, interfaces o tipos) y un ejemplo mínimo de uso para cada formato.
```

---

## Variante corta (solo recordatorio)

Si la IA ya conoce el formato y solo necesitas el encargo directo:

```
Crea un módulo para importar los JSON de:
1) privantix-inspector: analysis.json (o .jsonl) con estructura run / files[] / errors; cada file tiene columns[], rules_triggered[], fechas RFC3339.
2) privantix-acl-audit: JSON con started_at, completed_at, path, total_dirs, dirs[] (path, name, depth, owner, permissions, acls[], trusted_principals, outside_principals, compliance_status), errors[].

Incluir: carga desde ruta, parseo de fechas RFC3339, tipos/estructuras claros y ejemplo de uso para cada uno.
```

---

## Referencia de esquemas

- Inspector: `docs/analysis-result.schema.json`, `docs/OUTPUT_FORMAT_AI.md`, `docs/INSTRUCTIONS_AI_READ_JSON.md`
- ACL Audit: `docs/acl-audit-result.schema.json`, `docs/ACL_AUDIT_OUTPUT_AI.md`, `docs/INSTRUCTIONS_AI_ACL_AUDIT.md`

Puedes adjuntar estos archivos (o los schemas .json) a la IA para que genere el módulo de importación alineado con el formato real.
