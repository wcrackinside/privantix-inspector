# Usage

## Data Inspector (privantix-inspector)

### Basic Command

```bash
privantix-inspector scan --path <path> [options]
```

The `scan` subcommand is required. All other parameters are optional.

### Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `--path` | string | *(required)* | Root path to analyze (local or network) |
| `--output` | string | `./output` | Output directory for reports |
| `--config` | string | — | Path to YAML config file |
| `--max-sample-rows` | int | from config | Max rows to sample per file |
| `--extensions` | string | from config | Comma-separated extensions (e.g. `csv,xlsx,parquet`) |
| `--security` | flag | false | Extract file ACLs via OS commands |
| `--checksum` | flag | false | Calculate SHA256 checksum per file |
| `--hide-samples` | flag | false | Omit sample_values from output (data governance) |
| `--recursive` | flag | true | Scan subdirectories |
| `--log` | string | basic | Log level: `error`, `basic`, `detail`. Use `detail` to see each file with size/ext. |
| `--workers` | int | 4 (config: workers) | Número de workers en paralelo. `0` = usar valor del config. |
| `--incremental` | flag | false | Solo analizar archivos nuevos/modificados; reanuda desde run anterior o desde JSONL parcial si la run se cortó. |
| `--stream` | flag | false | Escribir resultados en JSONL/CSV sobre la marcha (menos memoria; permite reanudar con --incremental si se interrumpe). |

### ¿Cuántos workers usar?

Por defecto el inspector usa **4 workers** (o el valor de `workers` en el config si pasas `--config`). Cada worker analiza un archivo a la vez; varios workers permiten aprovechar varios núcleos y solapar lectura de disco con análisis.

| Situación | Recomendación |
|----------|----------------|
| **Equipo con varios núcleos, disco local rápido** | Dejar 4 o subir a **número de núcleos** (ej. 8). Puedes probar `--workers 8` y comparar tiempo. |
| **Disco lento o ruta de red (UNC / mapeo)** | Bajar a **2–4**. Muchos workers en red suelen saturar I/O y no mejoran. |
| **Poca RAM** | Bajar a **2** o **1**. Cada worker puede usar memoria (p. ej. archivos dentro de ZIP o Parquet grandes). |
| **Debug o proceso que se cuelga** | Usar **1** (`--workers 1`) para ver qué archivo está analizando y localizar el que bloquea. |
| **Muchos archivos pequeños (CSV/Excel)** | 4–8 suele ir bien en disco local. |
| **Archivos grandes o muchos ZIP/Parquet** | No subir de 4–6; más workers = más memoria a la vez. |

**Regla práctica:** empieza con el valor por defecto (4). Si el equipo tiene muchos núcleos y el disco no es el cuello de botella, prueba con `--workers 6` u `8`. Si va lento en red o con poca RAM, baja a `--workers 2` o `1`.

Para fijar el valor por comando: `--workers 8`. Para dejarlo fijo en un config: en tu `config.yaml` pon `workers: 6` (o el número que quieras).

### Diagnóstico si el proceso se queda colgado

Si se queda en "Analizando N/M", vuelve a ejecutar con **`--workers 1 --log detail`**; la última línea "Analyzing: ..." (o el nombre de archivo en la barra de progreso) será el archivo que está tardando o bloqueando. Ver más en [FAQ — El proceso se queda en Analizando N/M](faq.md#el-proceso-se-queda-en-analizando-nm-parece-colgado-qué-hacer).

### Examples

**Basic scan:**

```bash
privantix-inspector scan --path ./examples/data --output ./output
```

**With data governance mode (no sample values):**

```bash
privantix-inspector scan --path C:\datos --output ./report --hide-samples
```

**With security and checksum:**

```bash
privantix-inspector scan --path W:\repositorio --output ./audit --security --checksum
```

**Custom extensions and config:**

```bash
privantix-inspector scan --path ./data --config config.yaml --extensions csv,xlsx,parquet --max-sample-rows 500
```

**Non-recursive (root folder only):**

```bash
privantix-inspector scan --path ./data --recursive=false
```

### Configuration File (config.yaml)

Example `config.yaml`:

```yaml
supported_extensions:
  - .csv
  - .txt
  - .tsv
  - .xlsx
  - .rdat
  - .dat
  - .parquet
  - .7z
  - .zip
  - .rar
max_sample_rows: 200
workers: 4
# hide_sample_values: true  # omit sample_values (for data governance)
```

Command-line flags override config values.

### Modo incremental (`--incremental`)

Con `--incremental` el inspector **solo analiza archivos nuevos o modificados** desde la última ejecución y combina el resultado con el de esa ejecución. Así puedes actualizar el informe sin volver a analizar todo.

**Comportamiento:**

1. **Si existe un `analysis.json` completo** en el directorio de salida (misma `--output` y mismo `--output-name`): se toma la fecha de fin de esa ejecución y solo se analizan archivos con **fecha de modificación posterior**. El resultado final es la unión de los perfiles antiguos (sin tocar) y los nuevos/actualizados.

2. **Si no hay `analysis.json` pero sí un `analysis.jsonl`** (por ejemplo de una ejecución con `--stream` que se interrumpió): se considera una ejecución **parcial**. Se leen los perfiles ya escritos en el JSONL y **solo se analizan los archivos que faltan**. Al terminar se genera un resultado completo (JSON + CSV) con lo que ya estaba hecho más lo nuevo. Así puedes **continuar donde se quedó** tras un cierre o fallo.

**Requisitos:**

- Usar el **mismo** `--path`, `--output` y `--output-name` que en la run anterior (o al menos el mismo directorio de salida y nombre base).
- Para poder **reanudar** tras una ejecución cortada, la run anterior debe haber usado **`--stream`** (para que exista el `.jsonl` con los perfiles ya escritos).

**Ejemplos:**

```bash
# Primera vez (o run completo)
privantix-inspector scan --path W:\datos --output ./out --stream

# Siguiente día: solo archivos modificados desde ayer
privantix-inspector scan --path W:\datos --output ./out --incremental

# Se cortó en 100/106; vuelves a lanzar con incremental (lee analysis.jsonl y analiza solo los 6 que faltan)
privantix-inspector scan --path W:\datos --output ./out --incremental
```

**Nota:** El parámetro es un flag: `--incremental` (no hace falta `--incremental=true`). En el config puedes poner `incremental: true`.

---

## ACL Audit (privantix-acl-audit)

### Basic Command

```bash
privantix-acl-audit --path <path> [options]
```

### Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `--path` | string | *(required)* | Root path to audit |
| `--output` | string | `./output` | Output directory |
| `--output-name` | string | path basename + date | Base name for output files (e.g. `MiProyecto_20250306`) |
| `--trusted-groups` | string | — | JSON file defining trusted users/groups |
| `--revoked-users` | string | — | JSON file with revoked/deactivated users (principals that should not have access) |

### Examples

**Basic audit:**

```bash
privantix-acl-audit --path W:\11_CONADI\04_DATOS --output ./output
```

**With trusted groups:**

```bash
privantix-acl-audit --path W:\datos --output ./audit --trusted-groups ./examples/trusted-groups.json
```

**With revoked users (deactivated accounts):**

```bash
privantix-acl-audit --path W:\datos --output ./audit --revoked-users ./examples/revoked-users.json
```

**Combined trusted groups and revoked users:**

```bash
privantix-acl-audit --path W:\datos --output ./audit --trusted-groups ./trusted-groups.json --revoked-users ./revoked-users.json
```

**Custom output filename:**

```bash
privantix-acl-audit --path W:\proyecto --output ./reports --output-name CONADI_ACL_20250306
```

### Trusted Groups File (JSON)

Example `trusted-groups.json`:

```json
{
  "groups": [
    {
      "name": "dais",
      "members": [
        "MDP\\ASermeno",
        "MDP\\hmacuna",
        "MDP\\carce",
        "MDP\\lmenesesv"
      ]
    }
  ]
}
```

- **name**: Group identifier (informational)
- **members**: List of principals (e.g. `DOMAIN\username`). Matching is case-insensitive and supports both full principal and username-only.

When provided, the audit classifies each folder's ACL principals as:
- **trusted_principals**: Principals in the trusted groups
- **outside_principals**: Principals not in any trusted group
- **compliance_status**: `compliant` (no outsiders) or `non_compliant` (has outsiders)

### Revoked Users File (JSON)

Example `revoked-users.json`:

```json
{
  "users": [
    "CORP\\jperez",
    "CORP\\mgarcia",
    "lrodriguez"
  ]
}
```

- **users**: List of principals (e.g. `DOMAIN\username`) that are deactivated/revoked and should not have access. Matching is case-insensitive and supports full principal or username-only.

When provided, the audit adds:
- **revoked_principals**: Principals in the ACL that appear in the revoked users list
- **compliance_status**: `critical` when any revoked principal has access (overrides compliant/non_compliant)
