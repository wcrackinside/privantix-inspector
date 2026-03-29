# Privantix – GitHub Structure Prompt

## Objetivo
Construir la estructura completa de repositorios, carpetas, archivos base y workflows para un conjunto de herramientas de software llamado **Privantix**, orientado a Data Governance, Data Engineering, auditoría de datos, inspección de datasets y herramientas de análisis.

La arquitectura debe estar pensada como una **suite de productos de software**, no como un solo proyecto, y debe permitir:
- múltiples herramientas
- releases descargables
- sitio web en GitHub Pages
- documentación técnica
- whitepapers
- roadmap de productos
- automatización con GitHub Actions
- arquitectura escalable
- futura organización GitHub

---

## Repositorios a crear

Se deben crear los siguientes repositorios:

- privantix-web
- privantix-inspector
- privantix-acl-audit
- privantix-datalens
- privantix-docs

---

## Descripción de cada repositorio

### privantix-web
Sitio web público del proyecto Privantix.

Debe publicarse con GitHub Pages y contener:
- index.html
- css/
- js/
- img/
- products.json
- docs/
- roadmap/
- downloads/

El sitio debe mostrar:
- listado de productos
- descripción de cada producto
- versiones disponibles
- links de descarga
- documentación
- roadmap
- whitepapers
- arquitectura del proyecto

Debe leer la información de productos desde un archivo:
products.json

---

### privantix-inspector
Herramienta para inspección técnica de archivos y datasets.

Funcionalidades esperadas:
- inspección CSV, Excel, Parquet, TXT
- metadata técnica
- encoding detection
- delimiter detection
- column profiling
- data type inference
- reglas de gobernanza
- detección de datos sensibles
- exportación JSON, CSV, HTML
- checksum SHA256
- análisis de estructura de datasets

Debe tener la siguiente estructura:

src/
docs/
examples/
rules/
prompts/
product/
releases/
tests/
README.md
CHANGELOG.md
ROADMAP.md
LICENSE

---

### privantix-acl-audit
Herramienta para auditoría de permisos y seguridad de carpetas.

Funcionalidades:
- lectura de ACL
- análisis de permisos
- grupos de confianza
- detección de riesgos
- exportación de reportes
- auditoría de repositorios de datos
- análisis de permisos en NAS o file servers

Debe tener la misma estructura estándar de producto.

---

### privantix-datalens
Herramienta de exploración y visualización de datos.

Funcionalidades:
- visor de archivos
- explorador de datasets
- consultas SQL con DuckDB
- exportaciones
- preview de archivos
- metadata viewer
- navegación tipo explorador de Windows
- integración futura con lakehouse

Debe usar la estructura estándar de producto.

---

### privantix-docs
Repositorio de documentación general del ecosistema Privantix.

Debe contener:

architecture/
standards/
governance/
manuals/
whitepapers/
diagrams/
presentations/

Aquí se almacenan:
- arquitectura
- estándares de datos
- gobierno de datos
- manuales
- documentos técnicos
- whitepapers
- diagramas
- presentaciones

---

## Estructura estándar para repositorios de productos

Todos los repositorios de productos deben tener la siguiente estructura:

/src
/docs
/examples
/rules
/prompts
/product
/releases
/tests
README.md
CHANGELOG.md
ROADMAP.md
LICENSE

Descripción de carpetas:

- src → código fuente
- docs → documentación del producto
- examples → ejemplos de uso
- rules → reglas de gobernanza
- prompts → prompts para IA
- product → descripción del producto para la web
- releases → scripts de build o empaquetado
- tests → pruebas
- README.md → descripción del producto
- CHANGELOG.md → historial de versiones
- ROADMAP.md → roadmap del producto
- LICENSE → licencia

---

## Releases y versionado

Los productos deben publicarse usando GitHub Releases.

El versionado debe ser Semantic Versioning:

v0.1.0  
v0.2.0  
v0.2.1  
v1.0.0  

Flujo de releases:

Desarrollo  
→ Commit  
→ Push  
→ Tag versión  
→ GitHub Release  
→ Subir ZIP / EXE  
→ Web muestra descarga  

---

## Archivo products.json para la web

La web debe usar un archivo:

products.json

Con estructura:

[
  {
    "id": "privantix-inspector",
    "name": "Privantix Inspector",
    "description": "Data inspection and governance analysis",
    "repo": "user/privantix-inspector",
    "version": "v0.1.0",
    "downloads": [],
    "documentation": "",
    "roadmap": ""
  }
]

En el futuro este archivo puede generarse automáticamente consultando la API de GitHub Releases.

---

## GitHub Actions

Se deben definir workflows para:

1. Build del producto
2. Generación de release
3. Publicación GitHub Pages
4. Actualización automática de products.json
5. Empaquetado de ejecutables
6. Publicación de documentación

---

## Objetivo de la arquitectura

La arquitectura debe permitir:

- múltiples herramientas
- releases descargables
- sitio web de productos
- documentación centralizada
- roadmap de productos
- whitepapers técnicos
- open source parcial
- automatización con GitHub Actions
- futura organización GitHub
- ecosistema de herramientas de Data Governance
- distribución de software
- documentación técnica
- plataforma de productos

---

## Flujo de trabajo esperado

Desarrollar producto  
→ git add  
→ git commit  
→ git push  
→ git tag vX.X.X  
→ git push origin tag  
→ GitHub Release  
→ Subir ZIP/EXE  
→ Web Privantix muestra descarga  

---

## Resultado esperado

El resultado debe ser una estructura completa de ecosistema de software, similar a una organización de productos de software, incluyendo:

- repositorios
- estructura de carpetas
- documentación
- releases
- sitio web
- workflows
- versionado
- arquitectura
- roadmap
- distribución de software

El objetivo final es construir una plataforma de herramientas de Data Governance y Data Engineering llamada Privantix.