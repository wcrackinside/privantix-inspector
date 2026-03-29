# Roadmap

## MVP 1 (v0.1.x — Current)

### Data Inspector (`privantix-inspector`)

Repository scanning and dataset profiling.

Core capabilities:

- [x] CSV, TXT, TSV, XLSX scan
- [x] RDAT, DAT, Parquet support
- [x] Archive support (ZIP, 7z, RAR) — in-memory analysis
- [x] Technical metadata extraction
  - encoding detection
  - delimiter detection
  - row and column count
- [x] Column profiling with inferred types
- [x] JSON, CSV and HTML export
- [x] Hardcoded rules:
  - encoding validation
  - header detection
  - sensitive data detection
  - null ratio alerts

Governance and audit features:

- [x] `--hide-samples` for data governance safe mode
- [x] `--checksum` (SHA256) for file integrity
- [x] `--security` (ACL extraction)

---

### ACL Audit (`privantix-acl-audit`)

Filesystem access control analysis.

- [x] Folder permission audit
- [x] ACL inspection
- [x] Trusted groups configuration (JSON)
- [x] Compliance classification
- [x] JSON and CSV export

---

# v0.2 (Planned)

Quality improvements and reporting enhancements.

- [ ] Improved encoding detection
- [ ] Enhanced HTML report
- [ ] Multi-sheet XLSX summary
- [ ] Unit test coverage improvements
- [ ] Extended ACL support for Linux/macOS

---

# v0.3 (Planned)

Governance and risk analysis capabilities.

- [ ] YAML-driven rules engine
- [ ] Risk scoring per file and column
- [ ] Access exposure mapping (ACL + classification)
- [ ] Improved sensitive data detection

---

# Future

Platform expansion.

- [ ] Database connectivity (profiling tables and views)
- [ ] API / programmatic usage
- [ ] Integration with data catalogs
- [ ] Data lineage discovery
- [ ] Continuous repository monitoring

---

# Product Direction

Privantix aims to provide lightweight tools that help organizations understand, audit and govern their data repositories without requiring complex infrastructure.

The long-term vision is to evolve from repository inspection to **data exposure mapping and governance automation**.
