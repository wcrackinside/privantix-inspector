# `src/`

This repository uses a **Go layout** with entrypoints under [`cmd/`](../cmd/) and library packages at the **repository root** (`analyzer/`, `detectors/`, `inspector/`, etc.), not under a single `src/` tree.

This folder exists to match the [Privantix standard product layout](docs/PRIVANTIX_GITHUB_STRUCTURE_PROMPT.md). Treat **`cmd/`** and root packages as the effective source tree.
