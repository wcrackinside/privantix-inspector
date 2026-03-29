# `data/`

- **`releases-cache.json`** (opcional): salida de [`scripts/sync_releases_cache.py`](../../scripts/sync_releases_cache.py). La web intenta primero la API pública de GitHub; si falla (p. ej. límite de peticiones), puede usar esta caché.

Generar o actualizar:

```bash
python scripts/sync_releases_cache.py
```

Con token (más cuota): `set GITHUB_TOKEN=...` (Windows) o `export GITHUB_TOKEN=...` (Linux/macOS).
