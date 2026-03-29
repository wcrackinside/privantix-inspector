# Privantix ACL Audit

## Descripción corta (Hero)
Auditoría de permisos y control de acceso sobre repositorios de archivos.

## Descripción del producto
Permite identificar exposiciones de acceso en carpetas y verificar cumplimiento de políticas de seguridad.

## Características principales
- Análisis de ACL
- Clasificación de grupos confiables
- Detección de usuarios de bajas (offboarded users) que aún tienen acceso
- Exportación JSON y CSV

## Casos de uso
- Auditoría de seguridad
- Evaluación de exposición de datos
- Cumplimiento de políticas internas

## Ejemplo de ejecución
privantix-acl-audit --path ./repository --output ./output

## Release
Versión actual: v0.1.x (MVP)