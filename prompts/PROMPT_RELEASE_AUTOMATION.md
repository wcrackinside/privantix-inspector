# Release Automation Prompt
Project: Privantix Source Inspector

This document defines the instructions for generating release automation assets for the project.

The AI must treat this document as the authoritative specification for setting up automated builds and releases.

---

# 1 Purpose

Generate the files required to automate cross-platform releases for Privantix Source Inspector.

The release pipeline must:

- build binaries for Windows, Linux and macOS
- package release artifacts
- generate checksums
- publish releases through GitHub Releases

---

# 2 Required Deliverables

The AI must generate:

- `.goreleaser.yml`
- `.github/workflows/release.yml`

---

# 3 Constraints

The automation must:

- build from `./cmd/inspector`
- use the binary name `privantix-inspector`
- use semantic version tags such as `v0.1.0`
- avoid CGO when possible
- generate portable binaries

---

# 4 Platforms

Supported release targets:

- windows/amd64
- linux/amd64
- linux/arm64
- darwin/amd64
- darwin/arm64

---

# 5 Validation Rules

Before finalizing, verify that:

- workflow triggers on version tags
- GoReleaser references the correct main package
- archive names are consistent
- checksums are generated

---

# End of Specification
