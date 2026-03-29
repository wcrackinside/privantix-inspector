# Privantix ACL Audit — Documentación

**privantix-acl-audit** audita permisos y listas de control de acceso (ACL) de carpetas y subcarpetas. Genera informes en JSON y CSV con propietario, permisos y lista de entidades con acceso a cada directorio. Con un archivo de grupos de confianza (trusted groups) clasifica cada directorio como compliant o non_compliant según si hay entidades no autorizadas.

---

## Descripción

- **Qué hace:** Recorre un árbol de directorios y, para cada carpeta, obtiene propietario, permisos del sistema y ACL (en Windows usa la API del sistema para obtener ACL detalladas).
- **Grupos de confianza:** Si se pasa `--trusted-groups` con un JSON, cada principal de las ACL se clasifica como “dentro de grupos de confianza” o “fuera”. Un directorio con al menos un principal fuera se marca como `non_compliant`.
- **Usuarios de bajas:** Si se pasa `--offboarded-users` con un JSON de usuarios dados de baja (por ejemplo, ex–empleados), el informe marca qué directorios siguen otorgando acceso a estos usuarios.
- **Salida:** Un archivo JSON y otro CSV con el mismo nombre base (por defecto `{nombre_ruta}_{YYYYMMDD}` o el que indiques con `--output-name`).

**Público objetivo:** Equipos de seguridad, cumplimiento e IT que necesitan auditar acceso a carpetas, verificar políticas de acceso o detectar exposición a entidades no autorizadas.

---

## Uso

### Comando básico

```bash
privantix-acl-audit --path <ruta> [opciones]
```

`--path` es obligatorio. El resto de opciones son opcionales.

### Parámetros

| Parámetro | Tipo | Por defecto | Descripción |
|-----------|------|-------------|-------------|
| `--path` | string | *(requerido)* | Ruta raíz a auditar (local o red) |
| `--output` | string | `./output` | Directorio donde se escriben el JSON y el CSV |
| `--output-name` | string | nombre de la ruta + fecha | Nombre base de los archivos (ej. `CONADI_ACL_20250306`) |
| `--trusted-groups` | string | — | Ruta a un archivo JSON con la definición de grupos de confianza |
| `--offboarded-users` | string | — | Ruta a un archivo JSON con la lista de usuarios de bajas (principales que ya no deberían tener acceso) |

### Nombres de archivos de salida

| Modo | JSON | CSV |
|------|------|-----|
| Por defecto | `{basename_path}_{YYYYMMDD}.json` | `{basename_path}_{YYYYMMDD}.csv` |
| Con `--output-name X` | `X.json` | `X.csv` |

Ejemplo: si `--path` es `W:\11_CONADI\04_DATOS` y no usas `--output-name`, se generan `04_DATOS_20250311.json` y `04_DATOS_20250311.csv` en el directorio indicado por `--output`.

---

## Archivo de grupos de confianza (trusted groups)

Formato del JSON (ej. `trusted-groups.json`):

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
    },
    {
      "name": "admin",
      "members": ["Administrators", "DOMAIN\\admin_group"]
    }
  ]
}
```

- **groups:** Array de grupos.
- **name:** Identificador del grupo (informativo).
- **members:** Lista de principales (usuario o grupo). En Windows suele ser `DOMINIO\usuario` o `DOMINIO\grupo`. El matching es case-insensitive y acepta tanto el principal completo como solo el nombre de usuario (parte tras el último `\`).

Cuando se usa `--trusted-groups`, en cada directorio del resultado aparecen:

- **trusted_principals:** Principales que coinciden con algún miembro de los grupos de confianza.
- **outside_principals:** Principales que no están en ningún grupo de confianza.
- **compliance_status:** `compliant` si no hay outside_principals; `non_compliant` si hay al menos uno.

---

## Archivo de usuarios de bajas (offboarded users)

El archivo de usuarios de bajas define los **principales que ya no deberían tener acceso** a ningún directorio (por ejemplo, ex–empleados o cuentas desactivadas).

Formato sugerido del JSON (ej. `offboarded-users.json`):

```json
{
  "users": [
    "MDP\\usuario1",
    "MDP\\usuario2",
    "MDP\\carce_old"
  ]
}
```

- **users:** Lista de principales (usuario o grupo) considerados de baja. El matching sigue las mismas reglas que para grupos de confianza: es case-insensitive y acepta tanto el principal completo como solo el nombre de usuario (parte tras el último `\`).

Cuando se usa `--offboarded-users`:

- A nivel raíz se incluye `offboarded_users` con la lista de principales de bajas cargados.
- En cada directorio se añaden:
  - **offboarded_principals:** Principales en las ACL que coinciden con la lista de usuarios de bajas.
  - **has_offboarded_principals:** `true` si al menos un usuario de bajas tiene acceso al directorio; `false` en caso contrario.

---

## Formato de salida

### JSON

Estructura raíz:

- **started_at** (string, RFC3339): Inicio de la auditoría.
- **completed_at** (string, RFC3339): Fin de la auditoría.
- **duration_seconds** (number): Duración en segundos.
- **path** (string): Ruta absoluta auditada.
- **output_dir** (string): Directorio de salida usado.
- **total_dirs** (int): Número de directorios incluidos.
- **trusted_groups** (array, opcional): Copia de los grupos cargados desde el JSON (cada uno con `name` y `members`).
- **dirs** (array): Un objeto por directorio (ver abajo).
- **errors** (array de strings): Errores durante el recorrido (ej. permisos denegados).

Cada elemento de **dirs** (DirACL):

| Campo | Tipo | Descripción |
|-------|------|-------------|
| path | string | Ruta absoluta del directorio |
| name | string | Nombre del directorio |
| depth | int | Profundidad desde la raíz (0 = raíz) |
| owner | string | Propietario (ej. `DOMAIN\user`) |
| permissions | string | Permisos estilo Unix (ej. `drwxr-xr-x`) |
| acls | string[] | Entradas ACL (en Windows ej. `DOMAIN\user:(I)(F)`) |
| trusted_principals | string[] | Presente si se usó --trusted-groups |
| outside_principals | string[] | Presente si se usó --trusted-groups |
| compliance_status | string | `compliant` o `non_compliant` (si se usó --trusted-groups) |
| offboarded_principals | string[] | Principales de baja (usuarios de bajas) con acceso a este directorio (si se usó `--offboarded-users`) |
| has_offboarded_principals | bool | `true` si al menos un usuario de bajas tiene acceso a este directorio |

### CSV

- Codificación: UTF-8.
- Separador: coma.
- Cabeceras: `path`, `name`, `depth`, `owner`, `permissions`, `acls`. Si se usó trusted groups se añaden `trusted_principals`, `outside_principals`, `compliance_status`. Si se usó `--offboarded-users` se añaden también `offboarded_principals`, `has_offboarded_principals`.
- Valores múltiples en una celda (ej. ACL o listas de principales) se unen con ` | ` (espacio-pipe-espacio).

Formato típico de una entrada ACL en Windows: `DOMINIO\usuario:(I)(F)` (principal + `:` + flags de permisos).

---

## Ejemplos

**Auditoría básica:**

```bash
privantix-acl-audit --path W:\11_CONADI\04_DATOS --output ./output
```

**Con grupos de confianza:**

```bash
privantix-acl-audit --path W:\datos --output ./audit --trusted-groups ./examples/trusted-groups.json
```

**Nombre de salida fijo (ej. para integración):**

```bash
privantix-acl-audit --path W:\proyecto --output ./reports --output-name CONADI_ACL_20250306
```

**Ruta UNC:**

```bash
privantix-acl-audit --path "\\servidor\recurso\carpeta" --output ./out
```

---

## Convenciones

- **Rutas:** Windows usa backslash; rutas UNC como `\\servidor\recurso\path`.
- **Fechas en JSON:** Siempre RFC3339 (ej. `2026-03-11T14:30:00-03:00`).
- **Campos opcionales:** Si no se usa `--trusted-groups`, `trusted_principals`, `outside_principals` y `compliance_status` no aparecen en los objetos de directorio (o pueden ser arrays/string vacíos según versión).

---

## Esquema y documentación para IA

- **Especificación del JSON para otra IA:** `docs/ACL_AUDIT_OUTPUT_AI.md`
- **JSON Schema del resultado:** `docs/acl-audit-result.schema.json`
