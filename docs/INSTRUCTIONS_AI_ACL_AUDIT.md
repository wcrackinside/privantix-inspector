# Cómo leer el JSON de salida de privantix-acl-audit (para IA)

Usa estas instrucciones como contexto al dar a otra IA el archivo JSON generado por **privantix-acl-audit**, para que lo interprete correctamente.

---

## Texto para pegar como instrucción a la otra IA

```
El archivo que te paso es la salida JSON de privantix-acl-audit (auditoría de ACL y permisos de carpetas).

**Estructura del JSON:**
- Raíz: un objeto con:
  - "started_at", "completed_at" (RFC3339), "duration_seconds", "path" (ruta auditada), "output_dir", "total_dirs".
  - "trusted_groups" (opcional): array de { "name", "members" } con los grupos de confianza cargados.
  - "offboarded_users" (opcional): array de strings con los usuarios de bajas (principales que ya no deberían tener acceso).
  - "dirs": array de objetos, uno por directorio auditado (DirACL).
  - "errors": array de strings con errores del recorrido.

**Cada elemento de "dirs" (DirACL) contiene:**
- path, name, depth (profundidad desde raíz), owner, permissions.
- "acls": array de strings; cada uno es una entrada ACL (en Windows suele ser "DOMINIO\\usuario:(flags)").
- Si se usó --trusted-groups: "trusted_principals", "outside_principals" (arrays de strings), "compliance_status" ("compliant" o "non_compliant").
- Si se usó --offboarded-users: "offboarded_principals" (array de strings con usuarios de bajas que aún tienen acceso) y "has_offboarded_principals" (booleano).

**Interpretación:** 
- Un directorio con compliance_status "non_compliant" tiene al menos un principal en outside_principals (acceso no autorizado según los grupos de confianza).
- Si "has_offboarded_principals" es true, el directorio sigue otorgando acceso a uno o más usuarios de bajas. 
Fechas en RFC3339.
```

---

## Crear módulo de importación

Para pedir a una IA que **genere un módulo que importe/parsee estos JSON** (inspector + ACL audit), usa: **`docs/PROMPT_CREATE_IMPORT_JSON_MODULE.md`**.

## Referencia completa

- **Documentación de uso:** `docs/acl-audit.md`
- **Especificación del formato para IA:** `docs/ACL_AUDIT_OUTPUT_AI.md`
- **Esquema JSON (validación):** `docs/acl-audit-result.schema.json`
