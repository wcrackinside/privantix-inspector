# Privantix — ecosistema en GitHub

Este documento resume la arquitectura de repositorios descrita en [`PRIVANTIX_GITHUB_STRUCTURE_PROMPT.md`](PRIVANTIX_GITHUB_STRUCTURE_PROMPT.md) y su estado actual.

## Repositorios

| Repositorio | Rol |
|-------------|-----|
| [privantix-inspector](https://github.com/wcrackinside/privantix-inspector) | Monorepo principal: Inspector, ACL Audit, Catalog, Compare, documentación de producto y sitio estático (`privantix_site/`). Releases con GoReleaser. |
| [privantix-web](https://github.com/wcrackinside/privantix-web) | Sitio público en producción **https://www.privantix.io** (FTP desde `privantix_site/`). El repo GitHub puede seguir como copia o despliegue alternativo. |
| [privantix-docs](https://github.com/wcrackinside/privantix-docs) | Documentación transversal: arquitectura, estándares, governanza, whitepapers. |
| [privantix-acl-audit](https://github.com/wcrackinside/privantix-acl-audit) | Punto de entrada / roadmap para la herramienta ACL; el código fuente vive hoy en el monorepo (`cmd/privantix-acl-audit`). |
| [privantix-datalens](https://github.com/wcrackinside/privantix-datalens) | Producto planificado (exploración / DuckDB / visor). |

## Sitio web y `products.json`

- Archivo de datos: `privantix_site/products.json` (producto, descripción, `repo` `owner/name`, enlaces a documentación, changelog y roadmap).
- La home carga las tarjetas con `assets/js/products-loader.js`, que consulta **`GET https://api.github.com/repos/{owner}/{repo}/releases/latest`** para obtener `tag_name`, `name`, `published_at`, `assets[].browser_download_url`, `assets[].download_count`, `html_url` (notas de versión en GitHub).
- Si la API falla (p. ej. límite de peticiones), se puede usar la caché generada por `scripts/sync_releases_cache.py` → `privantix_site/data/releases-cache.json`.
- Carpetas de compatibilidad con el prompt: `css/`, `js/`, `img/`, más `docs/`, `roadmap/`, `downloads/` bajo `privantix_site/`.

## Versionado y releases

- [Semantic Versioning](https://semver.org/): etiquetas `vMAJOR.MINOR.PATCH`.
- Push de tag `v*` → [`.github/workflows/release.yml`](../.github/workflows/release.yml) → GoReleaser → [Releases](https://github.com/wcrackinside/privantix-inspector/releases).

## Automatización (GitHub Actions)

| Ámbito | Workflow |
|--------|----------|
| Binarios + artefactos | `release.yml` (tags `v*`) |
| Validación `products.json` | `validate-products-json.yml` |

Los metadatos de release **no** se duplican en `products.json`: se leen en tiempo de carga desde la API (o desde `releases-cache.json` como respaldo).
