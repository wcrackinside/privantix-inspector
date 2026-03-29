# Product Description Generation Prompt
Project: Privantix Tools

This document defines the instructions for generating the product descriptions used on the Privantix website.

The AI must treat this document as the authoritative specification for generating and updating product descriptions.

The generated content must be consistent with the implemented features and the roadmap of the project.

The AI must **never invent capabilities that are not present in the repository or roadmap**.

---

# 1 Purpose

Generate structured product descriptions for all Privantix tools.

The descriptions will be used on the official website and must follow a consistent structure.

Each product description must contain the following sections:

- Short Description (Hero)
- Product Description
- Key Capabilities
- Use Cases
- Example Execution
- Release Information

---

# 2 Source of Truth

The AI must extract product capabilities from the following sources:

- ROADMAP.md
- README.md
- project documentation under /docs
- CLI usage examples

If inconsistencies are detected, the AI must prioritize ROADMAP.md.

---

# 3 Product Identification

The AI must identify each product based on the tools present in the project.

Example products include:

privantix-inspector  
privantix-acl-audit  

Additional products may appear in the future.

The AI must generate one description per product.

---

# 4 Required Product Structure

Each product must include the following sections.

---

# Product Name

Example:

Privantix Source Inspector

---

## Short Description (Hero)

A concise one or two sentence description suitable for the hero section of a product page.

The hero text must:

- communicate the main capability
- be clear and concise
- avoid marketing exaggeration

Example tone:

"Analyze data repositories and understand their structure in minutes."

---

## Product Description

Explain the product in more detail.

The description must include:

- what the tool does
- what problem it solves
- how it works conceptually
- who should use it

The language must be clear and professional.

---

## Key Capabilities

Provide a bullet list of the main capabilities.

Capabilities must be derived from the roadmap and implementation.

Examples:

Repository scanning  
Technical metadata extraction  
Column profiling  
Archive inspection  
ACL analysis  

The list must reflect actual features.

---

## Use Cases

Describe typical situations where the tool is useful.

Examples may include:

- auditing data repositories
- preparing datasets for ingestion
- documenting legacy data sources
- evaluating access exposure

Each use case must be short and practical.

---

## Example Execution

Provide a CLI execution example.

Example:

```
privantix-inspector scan --path ./repository --output ./results
```

The command must match the real CLI syntax implemented in the tool.

---

## Release Information

Provide release details.

Example structure:

Current Version: v0.1.x

Status: MVP

Highlights:
- repository scanning
- metadata extraction
- dataset profiling
- report generation

The version must align with the project roadmap.

---

# 5 Writing Guidelines

All descriptions must follow these rules:

- professional tone
- clear technical language
- avoid marketing exaggeration
- avoid speculative features
- maintain consistency between products

---

# 6 Output Format

The AI must generate the product descriptions in Markdown.

Each product description should be written as a standalone section.

Example output structure:

```
products/
inspector.md
acl-audit.md
```

These files will later be integrated into the website.

---

# 7 Validation Rules

Before finishing the generation, the AI must verify:

- every product in the roadmap has a description
- CLI examples match implemented commands
- capabilities reflect actual features
- release information matches roadmap version

---

# 8 Final Deliverables

The AI must generate:

```
products/
privantix-inspector.md
privantix-acl-audit.md
```

Each file must contain:

- Short Description (Hero)
- Product Description
- Key Capabilities
- Use Cases
- Example Execution
- Release Information

The generated content must be ready for publication on the Privantix website.

---

# End of Specification
