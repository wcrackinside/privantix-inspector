# Website Update Prompt
Product: Privantix Catalog

This prompt instructs the AI to update the Privantix website to include the new product.

The AI must read the existing product page and add the new product without breaking layout or styling.

---

# Instructions

Open the file:

product.html

The page must be updated to include a new product card for **Privantix Catalog**.

The existing visual structure must not be modified.

Navigation, CSS classes and layout must be preserved.

---

# New Product Section

Add a product card similar to the existing ones.

Title:

Privantix Catalog

Description:

Transforms Privantix Inspector analysis results into a structured dataset inventory, providing visibility into datasets, columns, metadata and governance indicators.

Features list:

- Build catalog from inspector analysis
- Dataset and column inventory
- Sensitive data indicators
- Security metadata indicators
- JSON, CSV and HTML catalog outputs

---

# CLI Example Block

The product page must include a CLI example:

```bash
privantix-catalog build \
  --analysis ./output/analysis.json \
  --output ./catalog