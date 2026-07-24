# Card B6 — Replace `main_test.go`

**Package:** B (form, validation & output) · **Owner:** Spiros
**Priority:** 🟡 near the end

---

**Context:** old CLI subprocess tests no longer describe the product once
`main.go` is a server. Replace with HTTP integration tests.

## Steps
- [ ] Remove the CLI subprocess integration tests
- [ ] Add HTTP integration tests via `httptest`
- [ ] Do NOT rewrite old tests to pretend CLI is still supported

## Acceptance criteria
- [ ] Obsolete root CLI tests are replaced when `main.go` becomes a server
- [ ] HTTP behavior is testable without opening a network port

## Meta
- **Files:** `main_test.go`
- **PRD:** §4.2
- **Depends on:** `main.go` = server (Phase 0) + handlers exist

**Definition of Done:** `main_test.go` contains HTTP integration tests; no CLI
pretence remains.

---

## Design tokens (locked in Phase 0)

Style every control with these shared CSS variables in `static/style.css`.
**Do not introduce new colours.**

| Token | Hex | Use |
| --- | --- | --- |
| `--page-bg` | `#FAF8F3` | page background — soft warm off-white |
| `--card` | `#FFFFFF` | card surface |
| `--accent` | `#8F3B3B` | wine/μπορντό — logo accent, primary action, active states |
| `--text` | `#3B322A` | espresso — primary text |
| `--text-muted` | `#7A7062` | labels, secondary text, placeholders |
| `--border` | `#E8E4DA` | field and card borders |
| `--field-bg` | `#EEEDE7` | input / control fill |

**Output box (terminal):** dark surface `#2E2820`, light text `#F6EFE2`, monospace.
Colored runs use `rgb(R, G, B)` from validated integers only.
