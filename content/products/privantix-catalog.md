# Privantix Catalog

## Descripción corta (Hero)
Transforma la salida de Privantix Inspector en un inventario estructurado de datasets, permitiendo entender qué datos existen en los repositorios.

## Descripción del producto
Privantix Catalog construye un catálogo portable y estructurado a partir del análisis técnico generado por Privantix Inspector. Resume datasets, columnas, metadata, detección de reglas e indicadores de gobernanza en un inventario legible por humanos y máquinas.

Permite entender rápidamente: qué datasets existen, dónde están ubicados, qué columnas contienen, si se detectó información potencialmente sensible, si hay metadata de seguridad disponible y qué reglas se dispararon durante el análisis.

Privantix Catalog no escanea repositorios directamente; consume la salida generada por Privantix Inspector.

## Características principales
- Construye catálogo desde el análisis del inspector
- Resumen de datasets y columnas
- Indicadores de datos potencialmente sensibles
- Indicadores de metadata de seguridad (checksum, ACL)
- Exportación JSON, CSV y HTML

## Casos de uso
- Inventario de datos: construir inventario de datasets en repositorios compartidos
- Visibilidad de gobernanza: identificar datasets con información potencialmente sensible
- Documentación de repositorios: generar documentación para repositorios legacy
- Preparación de auditoría: reportes de catálogo para auditorías y cumplimiento

## Ejemplo de ejecución
privantix-catalog --input ./output/analysis.json --output ./catalog

## Release
Versión actual: v0.1.x (MVP)
