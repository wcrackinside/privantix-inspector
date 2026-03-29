# Privantix ACL Audit

## Short Description (Hero)

Audit folder permissions and access control lists across directory trees, and identify principals outside your trusted groups.

---

## Product Description

Privantix ACL Audit scans folders and subfolders to collect permissions and ACLs (Access Control Lists). It reports owner, permissions, and the full list of principals with access to each directory. On Windows, it uses `icacls` to retrieve detailed ACL information.

When you provide a trusted groups file (JSON), the tool classifies each principal as trusted or outside. Folders with any outside principal are marked as `non_compliant`; folders where all principals are trusted are marked as `compliant`. This helps security and compliance teams quickly identify directories where unauthorized users or groups have access.

Additionally, you can provide an offboarded users file (JSON) with users who should no longer have access (usuarios de bajas). The tool highlights any directory where these offboarded principals still appear in the ACLs, making it easy to detect accounts that were not correctly removed during offboarding.

Reports are exported to JSON and CSV. Output filenames can be customized with `--output-name` (e.g. for dated audit runs).

**Who should use it:** Security teams, compliance officers, and IT administrators who need to audit folder access, verify access policies, or identify exposure to unauthorized principals.

---

## Key Capabilities

- Folder permission audit (directories only)
- ACL inspection (Windows: `icacls`)
- Trusted groups configuration (JSON)
- Compliance classification (compliant / non_compliant)
- Offboarded users detection (JSON list of users de bajas)
- Export to JSON and CSV
- Custom output filename for audit runs

---

## Use Cases

- Auditing folder access before compliance reviews
- Identifying directories with unauthorized principals
- Evaluating access exposure across shared drives
- Documenting current ACL state for change management
- Verifying that only trusted groups have access to sensitive folders

---

## Example Execution

```bash
privantix-acl-audit --path W:\11_CONADI\04_DATOS --output ./output
```

With trusted groups:

```bash
privantix-acl-audit --path W:\datos --output ./audit --trusted-groups ./examples/trusted-groups.json
```

With trusted groups and offboarded users:

```bash
privantix-acl-audit --path W:\datos --output ./audit \
  --trusted-groups ./examples/trusted-groups.json \
  --offboarded-users ./examples/offboarded-users.json
```

With custom output filename:

```bash
privantix-acl-audit --path W:\proyecto --output ./reports --output-name CONADI_ACL_20250306
```

---

## Release Information

**Current Version:** v0.1.x

**Status:** MVP

**Highlights:**
- Folder permission and ACL audit
- Trusted groups for compliance classification
- JSON and CSV export
- Custom output naming for audit runs
