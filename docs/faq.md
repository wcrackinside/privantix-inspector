# Frequently Asked Questions

## Data Inspector

### El proceso se queda en "Analizando N/M" (parece colgado). ¿Qué hacer?

Puede que un archivo concreto esté tardando mucho (muy grande, comprimido, red lenta) o que haya un fallo que no termina de devolver.

**1) Descubrir en qué archivo se queda**

- Vuelve a ejecutar **con un solo worker** y **log detallado**; cuando se quede colgado, la **última línea** `Analyzing: ...` será el archivo problemático:

```bash
privantix-inspector scan --path "TU_RUTA" --output ./output --workers 1 --log detail
```

- Con `--log detail` verás cada archivo con tamaño y extensión. Con `--workers 1` solo se analiza un archivo a la vez, así que la última línea impresa es la que está bloqueando.

**2) Causas habituales**

- **Archivo muy grande** (Parquet, CSV o Excel con muchas filas): el análisis puede tardar minutos.
- **Archivo dentro de un ZIP/7z/RAR** muy grande o corrupto: descomprimir en memoria puede tardar o bloquearse.
- **Ruta de red** (UNC, mapeo de red) con `--security` o `--checksum`: llamadas al sistema (ACL, hash) pueden ser lentas o quedarse colgadas.
- **Disco o recurso bloqueado** por otro proceso.

**3) Qué hacer con el archivo problemático**

- Si es un archivo que puedes excluir, muévelo fuera de la ruta escaneada o quita su extensión de `--extensions`.
- Si es un archivo necesario, prueba a analizarlo solo (por ejemplo en una carpeta con solo ese archivo) para ver si termina o da error.
- Si usas `--checksum` o `--security`, prueba sin ellos para ver si el bloqueo desaparece (sobre todo en rutas de red).

**4) Resumen de opciones de diagnóstico**

| Opción | Uso |
|--------|-----|
| `--workers 1` | Un solo archivo a la vez; la última línea "Analyzing" es la que se queda colgada. |
| `--log detail` | Muestra ruta, tamaño y extensión de cada archivo al analizarlo. |
| `--log error` | Solo errores; útil si quieres que la barra de progreso no se mezcle con muchos mensajes. |

---

### What file encodings are detected?

The tool detects UTF-8 (with or without BOM), UTF-16 LE/BE, and falls back to `windows-1252` for non-UTF-8 content. A rule flags files that are not UTF-8.

### Does it read files inside ZIP/7z/RAR without extracting?

Yes. Archives are opened in memory and inner files (CSV, XLSX, Parquet, etc.) are analyzed without writing to disk. Paths in the report use the format `archive.zip/inner.csv`.

### What data types are inferred?

Supported types: `string`, `integer`, `float`, `boolean`, `date`, `datetime`, `email`, `phone`, `id_like`. The detector uses patterns and format heuristics.

### What rules are applied?

- **encoding_not_utf8**: File encoding is not UTF-8
- **missing_header**: Header was not confidently detected
- **too_many_columns**: File has more than 100 columns
- **high_null_ratio**: Column has null ratio above 50%
- **potential_sensitive_data**: Column may contain personal or sensitive identifiers (email, phone, id_like types)

### When should I use --hide-samples?

Use `--hide-samples` when sharing reports with data governance or compliance teams who should not see actual data values. Column metadata (types, null %, max length) is still included.

### Can I analyze only certain extensions?

Yes. Use `--extensions csv,xlsx,parquet` or set `supported_extensions` in `config.yaml`.

---

## ACL Audit

### What does "compliant" vs "non_compliant" mean?

- **compliant**: All ACL principals on the folder are in the trusted groups.
- **non_compliant**: At least one principal is not in any trusted group.

### How does trusted group matching work?

Matching is case-insensitive. A principal like `MDP\carce` matches if the trusted file contains `MDP\carce` or `carce`. Both full principal and username-only are supported.

### Does it work on network paths?

Yes. As long as the path is accessible (e.g. `W:\share\folder`), the tool can audit it. On Windows, `icacls` is used to retrieve ACLs.

### Why are SIDs (e.g. S-1-5-21-...) shown in outside_principals?

When a user or group cannot be resolved to a name, Windows returns the Security Identifier (SID). These are treated as outside principals unless you add a mechanism to map SIDs to trusted names (not implemented in the current version).

---

## General

### Is there a GUI?

No. Both tools are command-line only.

### Can I run this in CI/CD?

Yes. The tools exit with code 0 on success and non-zero on fatal errors. Outputs can be parsed (JSON/CSV) for automation.

### What Go version is required?

Go 1.24 or later. Check with `go version`.
