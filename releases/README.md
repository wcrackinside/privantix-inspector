# Releases

- **Automation:** [`.goreleaser.yml`](../.goreleaser.yml) bundles CLI binaries with `docs/`, `examples/`, and `content/products/`.
- **CI:** [`.github/workflows/release.yml`](../.github/workflows/release.yml) runs on tags `v*`.

Tag and push to publish:

```bash
git tag -a v0.1.0 -m "Release"
git push origin v0.1.0
```
