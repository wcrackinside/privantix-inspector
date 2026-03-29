# Website Update Prompt
Task: Add DAMA-DMBOK alignment section to Privantix website

Read first:
- prompts/PROMPT_EDIT_WEBSITE.md
- content/products/privantix-inspector.md
- content/products/privantix-acl-audit.md
- content/products/privantix-catalog.md (if available)
- content/vision/VISION.md
- content/vision/POSITIONING.md
- content/roadmap/ROADMAP.md

Then open:
- product.html

Update the HTML directly.

## Objective

Incorporate a new website section that explains how Privantix products align with DAMA-DMBOK practices.

The goal is to strengthen product positioning and show that the tools are grounded in recognized data governance practices.

## Important rules

- Keep the page in Spanish
- Preserve the current layout and Bootstrap structure
- Do not redesign the page
- Do not remove existing sections
- Do not change navigation, footer, CSS classes or JS behavior
- Add the new content as a new section inside `product.html`
- Use the same visual style already present in the site
- Do not invent products not already defined in the content files or roadmap

## New Section Title

Use this title:

**Privantix y las prácticas DAMA-DMBOK**

Subtitle:

**Cada herramienta cubre capacidades concretas de descubrimiento, metadata, seguridad y gobierno de datos.**

## Placement

Insert the new section after the current product/tools section and before the methodologies section.

If the exact section names differ, place it after the section that presents the available Privantix tools.

## Content to include

The section must explain the alignment between products and DAMA practices.

### Product mapping

Include at least these mappings:

#### Privantix Inspector
Covers:
- Metadata Management
- Data Discovery
- Data Quality
- Data Profiling

Short explanation:
Analiza repositorios, detecta datasets, perfila columnas, extrae metadata técnica e identifica reglas y posibles datos sensibles.

#### Privantix ACL Audit
Covers:
- Data Security Management
- Data Governance
- Access Control

Short explanation:
Audita permisos y ACL sobre carpetas y archivos para identificar exposición, accesos amplios y hallazgos de seguridad.

#### Privantix Catalog
Covers:
- Metadata Management
- Data Governance
- Data Inventory
- Data Architecture (partial)

Short explanation:
Convierte el análisis técnico en un inventario estructurado de datasets, columnas, indicadores de sensibilidad y metadata de gobierno.

## Visual format

Use one of these formats, preserving the current site style:

### Preferred option
Three product cards in a row (or stacked on mobile), one per product.

Each card should include:
- Product name
- DAMA areas covered
- Short explanation

### Optional addition
Add a summary block below the cards with something like:

**Cobertura actual de Privantix**
- Gobierno de datos
- Metadata
- Seguridad
- Calidad de datos
- Inventario de activos de datos

## Optional diagram

If possible, add a lightweight visual block or inline SVG showing:

Privantix Inspector → Discovery / Metadata / Profiling  
Privantix ACL Audit → Security / Access Governance  
Privantix Catalog → Inventory / Metadata / Governance

Do not use external images.
Prefer inline SVG or HTML cards.

## Suggested section copy

Use Spanish text similar to this:

Privantix se alinea con prácticas centrales del marco DAMA-DMBOK para traducir principios de gobierno de datos en herramientas operativas y portables.

Inspector aporta descubrimiento, perfilado y metadata técnica sobre datasets.  
ACL Audit aporta visibilidad sobre permisos y exposición de acceso.  
Catalog transforma los resultados técnicos en un inventario estructurado y reutilizable para procesos de gobierno de datos.

## Validation checklist

Before finishing, verify that:

- The new section is in Spanish
- The new section matches the existing visual style
- Bootstrap classes remain valid
- Navigation and footer are unchanged
- No existing section was broken
- The products mentioned are consistent with the roadmap and content files
- The page still renders correctly on mobile

## Final instruction

Edit `product.html` directly and return the updated HTML.
Do not provide only a summary.
Do not redesign the page.