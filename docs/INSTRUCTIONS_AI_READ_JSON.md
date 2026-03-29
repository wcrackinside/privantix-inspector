# Cómo leer el JSON de salida de privantix-inspector (para IA)

Usa estas instrucciones como contexto al dar a otra IA el archivo `analysis.json` (o `*.json`/`*.jsonl`) generado por **privantix-inspector**, para que lo interprete correctamente.

---

## Texto para pegar como instrucción a la otra IA

```
El archivo que te paso es la salida JSON de privantix-inspector (análisis de fuentes de datos: CSV, XLSX, Parquet, etc.).

**Estructura del JSON (analysis.json):**
- Raíz: objeto con 3 claves:
  - "run": metadatos del escaneo (started_at, completed_at, path, output_dir, supported_files, analyzed_files, failed_files, max_sample_rows, workers, supported_extensions).
  - "files": array de objetos, uno por archivo analizado (FileProfile).
  - "errors": array de strings con errores globales del escaneo.

**Cada elemento de "files" (FileProfile) contiene:**
- Identificación: path, name, extension, size_bytes, depth.
- Fechas: modified_at, created_at (opcional); formato RFC3339.
- Seguridad/ACL: owner, permissions, acls (opcional), trusted_principals, outside_principals, compliance_status (si se usó --trusted-groups).
- Contenido: encoding, delimiter, has_header, row_count_estimate, column_count, average_row_width, sheet_name (para Excel).
- "columns": array de columnas; cada una: file_path, sheet_name, name, position, inferred_type, null_percentage, max_length, sample_values.
- "rules_triggered": array de { rule_name, severity, message, target }; reglas típicas: encoding_not_utf8, missing_header, potential_sensitive_data, high_null_ratio.
- "errors": array de strings con errores de ese archivo.

**inferred_type** en columnas: string, integer, float, boolean, date, datetime, email, phone, id_like.

**Si el archivo es .jsonl (modo --stream):** cada línea es un JSON object. Busca "type": "run_start" (metadatos), "type": "file" (cada "file" es un FileProfile), "type": "run_end" (resumen final).
```

---

## Crear módulo de importación

Para pedir a una IA que **genere un módulo que importe/parsee estos JSON**, usa: **`docs/PROMPT_CREATE_IMPORT_JSON_MODULE.md`** (incluye especificación para inspector y ACL audit).

## Referencia completa

- **Especificación detallada:** `docs/OUTPUT_FORMAT_AI.md`
- **Esquema JSON (validación):** `docs/analysis-result.schema.json`
