# Privantix ACL Audit — Formato de salida para IA

Este documento describe los archivos JSON y CSV generados por **privantix-acl-audit** para que una IA o sistema automatizado pueda parsearlos e interpretarlos correctamente.

**JSON Schema:** `docs/acl-audit-result.schema.json` (JSON Schema draft-07 para validación)

---

## Instrucciones para otra IA: cómo leer el JSON de ACL Audit

**Archivo a leer:** `{baseName}.json` en la carpeta de salida (por defecto `{nombre_ruta}_{YYYYMMDD}.json`, o el indicado con `--output-name`).

1. **Estructura raíz:** El JSON es un único objeto con: `started_at`, `completed_at`, `duration_seconds`, `path`, `output_dir`, `total_dirs`, `trusted_groups` (opcional), `offboarded_users` (opcional), `dirs`, `errors`.
2. **Metadatos:** `path` es la ruta absoluta auditada; `total_dirs` es el número de directorios en `dirs`; `errors` es un array de strings con fallos del recorrido.
3. **Cada elemento de `dirs`** es un **DirACL:** `path`, `name`, `depth`, `owner`, `permissions`, `acls` (array de strings). Si se usó `--trusted-groups` también tendrá `trusted_principals`, `outside_principals`, `compliance_status` (`"compliant"` o `"non_compliant"`).
4. **Fechas:** `started_at` y `completed_at` en RFC3339.
5. **ACL (Windows):** Cada string en `acls` suele tener la forma `DOMINIO\usuario:(flags)`; el principal es la parte antes de los dos puntos.
6. **Compliance:** Directorios con `compliance_status == "non_compliant"` tienen al menos un principal en `outside_principals`.

Para validación estricta del formato, usar el esquema en `docs/acl-audit-result.schema.json`.

---

## 1. Archivos de salida

| Base Name | JSON | CSV |
|-----------|------|-----|
| Por defecto | `{path_basename}_{YYYYMMDD}.json` | `{path_basename}_{YYYYMMDD}.csv` |
| Con `--output-name X` | `X.json` | `X.csv` |

Los archivos se escriben en el directorio indicado por `--output`.

---

## 2. Esquema JSON

Raíz (**AuditResult**):

```json
{
  "started_at": "string (RFC3339)",
  "completed_at": "string (RFC3339)",
  "duration_seconds": "number",
  "path": "string - ruta absoluta auditada",
  "output_dir": "string",
  "total_dirs": "int",
  "trusted_groups": [ { "name": "string", "members": ["string"] } ],
  "offboarded_users": ["string"],
  "dirs": [ "DirACL[]" ],
  "errors": ["string"]
}
```

### DirACL

| Campo | Tipo | Descripción |
|-------|------|-------------|
| path | string | Ruta absoluta del directorio |
| name | string | Nombre del directorio |
| depth | int | Profundidad desde la raíz (0 = raíz) |
| owner | string | Propietario (ej. DOMAIN\user) |
| permissions | string | Permisos estilo Unix |
| acls | string[] | Entradas ACL (ej. DOMAIN\user:(I)(F) en Windows) |
| trusted_principals | string[] | Principales en grupos de confianza (opcional) |
| outside_principals | string[] | Principales fuera de grupos de confianza (opcional) |
| compliance_status | string | "compliant" \| "non_compliant" (opcional) |
| offboarded_principals | string[] | Principales marcados como usuarios de bajas que aún tienen acceso (opcional) |
| has_offboarded_principals | bool | `true` si al menos un usuario de bajas tiene acceso a ese directorio (opcional) |

---

## 3. Esquema CSV

- Delimiter: coma.
- Cabeceras: `path`, `name`, `depth`, `owner`, `permissions`, `acls`; si hay trusted groups: `trusted_principals`, `outside_principals`, `compliance_status`. Si se usan usuarios de bajas: `offboarded_principals`, `has_offboarded_principals`.
- Listas (acls, trusted_principals, outside_principals, offboarded_principals) unidas en una celda con ` | ` (espacio-pipe-espacio).

---

## 4. Formato de ACL (Windows)

Cada entrada en `acls` suele ser: `PRINCIPAL:(FLAGS)`, por ejemplo `MDP\carce:(I)(F)`. El principal es la parte antes del primer `:`; el resto son flags de permisos. Para comparar con trusted groups se usa el principal (usuario o grupo).

---

## 5. Ejemplo mínimo (JSON)

```json
{
  "started_at": "2026-03-11T14:00:00-03:00",
  "completed_at": "2026-03-11T14:01:05-03:00",
  "duration_seconds": 65.2,
  "path": "W:\\11_CONADI\\04_DATOS",
  "output_dir": "./output",
  "total_dirs": 42,
  "trusted_groups": [
    { "name": "dais", "members": ["MDP\\carce", "MDP\\lmenesesv"] }
  ],
  "dirs": [
    {
      "path": "W:\\11_CONADI\\04_DATOS",
      "name": "04_DATOS",
      "depth": 0,
      "owner": "MDP\\admin",
      "permissions": "drwxr-xr-x",
      "acls": ["MDP\\admin:(F)", "MDP\\carce:(RX)"],
      "trusted_principals": ["MDP\\carce"],
      "outside_principals": ["MDP\\admin"],
      "compliance_status": "non_compliant"
    }
  ],
  "errors": []
}
```

---

## 6. Consejos de interpretación para IA

1. **Errores:** Revisar `errors` para fallos de acceso o rutas no encontradas.
2. **Cumplimiento:** Filtrar `dirs` por `compliance_status == "non_compliant"` para directorios con acceso no autorizado.
3. **Outside principals:** En cada directorio, `outside_principals` indica quiénes están fuera de los grupos de confianza.
4. **Usuarios de bajas:** Si existe `offboarded_users`, revisar `offboarded_principals` y `has_offboarded_principals` para detectar directorios donde aún tienen acceso usuarios dados de baja.
5. **Profundidad:** `depth` permite ordenar o agrupar por nivel en el árbol.
6. **Trusted groups:** Si `trusted_groups` está vacío o ausente, no habrá `trusted_principals` / `outside_principals` / `compliance_status` en los DirACL.
