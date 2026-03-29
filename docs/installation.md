# Installation

## Prerequisites

- **Go 1.24** or later (for building from source)
- No additional runtime dependencies; the tools produce standalone executables

## Build from Source

### 1. Clone or download the repository

```bash
git clone <repository-url>
cd privantix-source-inspector
```

### 2. Download dependencies

```bash
go mod download
```

### 3. Build the executables

**Data Inspector:**

```bash
go build -o privantix-inspector.exe ./cmd/privantix-inspector
```

On Linux/macOS:

```bash
go build -o privantix-inspector ./cmd/privantix-inspector
```

**ACL Audit:**

```bash
go build -o privantix-acl-audit.exe ./cmd/privantix-acl-audit
```

On Linux/macOS:

```bash
go build -o privantix-acl-audit ./cmd/privantix-acl-audit
```

### 4. Verify installation

```bash
# Data Inspector
./privantix-inspector scan --path ./examples/data --output ./output

# ACL Audit
./privantix-acl-audit --path . --output ./output
```

## Optional: Build both at once

```bash
go build -o privantix-inspector.exe ./cmd/privantix-inspector
go build -o privantix-acl-audit.exe ./cmd/privantix-acl-audit
```

## Directory Layout After Build

```
privantix-source-inspector/
├── privantix-inspector.exe   # Data analysis tool
├── privantix-acl-audit.exe   # ACL audit tool
├── config.yaml               # Optional configuration
├── examples/
│   ├── data/
│   │   └── sample.csv
│   └── trusted-groups.json
└── docs/
```
