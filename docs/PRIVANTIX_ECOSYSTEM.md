# Privantix — ecosistema en GitHub

Este documento resume la arquitectura de repositorios descrita en [`PRIVANTIX_GITHUB_STRUCTURE_PROMPT.md`](PRIVANTIX_GITHUB_STRUCTURE_PROMPT.md) y su estado actual.

## Repositorios

| Repositorio | Rol |
|-------------|-----|
| [privantix-inspector](https://github.com/wcrackinside/privantix-inspector) | Monorepo principal: Inspector, ACL Audit, Catalog, Compare, documentación de producto y sitio estático (`privantix_site/`). Releases con GoReleaser. |
| [privantix-web](https://github.com/wcrackinside/privantix-web) | Sitio público (GitHub Pages). Se sincroniza desde `privantix_site/` del monorepo. |
| [privantix-docs](https://github.com/wcrackinside/privantix-docs) | Documentación transversal: arquitectura, estándares, governanza, whitepapers. |
| [privantix-acl-audit](https://github.com/wcrackinside/privantix-acl-audit) | Punto de entrada / roadmap para la herramienta ACL; el código fuente vive hoy en el monorepo (`cmd/privantix-acl-audit`). |
| [privantix-datalens](https://github.com/wcrackinside/privantix-datalens) | Producto planificado (exploración / DuckDB / visor). |

## Sitio web y `products.json`

- Archivo de datos: `privantix_site/products.json` (listado de productos, repos, versiones, enlaces).
- La home carga las tarjetas vía `assets/js/products-loader.js`.
- Carpetas de compatibilidad con el prompt: `css/`, `js/`, `img/`, más `docs/`, `roadmap/`, `downloads/` bajo `privantix_site/`.

## Versionado y releases

- [Semantic Versioning](https://semver.org/): etiquetas `vMAJOR.MINOR.PATCH`.
- Push de tag `v*` → [`.github/workflows/release.yml`](../.github/workflows/release.yml) → GoReleaser → [Releases](https://github.com/wcrackinside/privantix-inspector/releases).

## Automatización (GitHub Actions)

| Ámbito | Workflow |
|--------|----------|
| Binarios + artefactos | `release.yml` (tags `v*`) |
| Validación `products.json` | `validate-products-json.yml` |

La generación automática de `products.json` desde la API de GitHub queda como mejora futura.
