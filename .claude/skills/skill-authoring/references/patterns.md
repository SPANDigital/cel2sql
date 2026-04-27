# Skill Split Patterns

## Contents

- Pattern 1: High-Level Guide with References
- Pattern 2: Domain-Specific Organisation
- Pattern 3: Conditional Detail
- Pattern 4: Scripts Instead of Prose
- Choosing a pattern
- Combining patterns

---

## Pattern 1: High-Level Guide with References

Keep `SKILL.md` as a lightweight overview and link out to deeper material. The agent loads additional files only when the specific feature is needed.

```
pdf-processing/
├── SKILL.md         # Overview + quick start
├── FORMS.md         # Form-filling guide
├── REFERENCE.md     # Full API reference
└── EXAMPLES.md      # Common patterns
```

In `SKILL.md`:

```markdown
## Advanced features

**Form filling**: see FORMS.md
**API reference**: see REFERENCE.md
**Examples**: see EXAMPLES.md
```

**When to use:** the skill covers one cohesive domain but has distinct deep-dive features most users won't need.

**Example:** a game skill's `SKILL.md` links out to `references/tilemap-loading.md`, `references/player-movement.md`, `references/npc-dialogue.md`, `references/camera.md`. A question about camera bounds pulls in `SKILL.md` + `camera.md` only.

---

## Pattern 2: Domain-Specific Organisation

Split references by domain so the agent only pulls in the relevant slice.

```
bigquery-skill/
├── SKILL.md
└── reference/
    ├── finance.md   # Revenue, ARR, billing
    ├── sales.md     # Pipeline, opportunities
    ├── product.md   # API usage, features
    └── marketing.md # Attribution, campaigns
```

A question about sales pipeline reads `SKILL.md` + `reference/sales.md`. The other three files never enter context.

**When to use:** the skill spans domains where one user's questions are irrelevant to another's.

**Example:** a city-lore skill with one reference file per neighbourhood (`bo-kaap.md`, `sea-point.md`, `camps-bay.md`, …). Bo-Kaap content never loads for a Camps Bay question, and vice versa. The SKILL.md acts as a navigation map.

---

## Pattern 3: Conditional Detail

Show the common path in `SKILL.md`; link out to rare-but-important paths.

```markdown
## Editing documents

For simple edits, modify the XML directly.

**For tracked changes**: see REDLINING.md
**For OOXML internals**: see OOXML.md
```

**When to use:** there's a clear 80/20 split between the common path and rare edge cases. Keeping edge cases inline bloats every activation for the 80% who don't need them.

**Example:** a tiled-map authoring skill's `SKILL.md` covers the happy path for adding a tile layer; `references/collision-authoring.md` covers the less common collision-shape authoring; `references/event-triggers.md` covers POI interaction objects.

---

## Pattern 4: Scripts Instead of Prose

Deterministic logic belongs in a script, not in Markdown. The script's source never enters context — only its output does.

```
pdf-skill/
├── SKILL.md
└── scripts/
    ├── analyze_form.py
    ├── fill_form.py
    └── validate.py
```

In `SKILL.md`, be explicit about execution intent:

```markdown
Run: `python scripts/analyze_form.py input.pdf > fields.json`
```

versus

```markdown
See `scripts/analyze_form.py` for the extraction algorithm.
```

The first tells the agent to execute. The second tells it to read as reference. **Execution is almost always what you want** — pre-made scripts are more reliable than generated code, cheaper in context, and consistent across runs.

**When to use:** the work is deterministic (parsing, querying, projecting, validating). If a prompt would ask the agent to re-type the same code every time, it belongs in a script.

**Example:** a map-data skill's `scripts/fetch_osm.py` issues an Overpass API query and emits GeoJSON. The agent runs it, consumes the output, never reads the script body.

---

## Choosing a pattern

| Signal in the content | Pattern |
|-----------------------|---------|
| SKILL.md has distinct feature deep-dives | 1 |
| Content splits cleanly by audience or sub-domain | 2 |
| 80/20 happy path vs rare edge cases | 3 |
| Same deterministic code re-typed across uses | 4 |

---

## Combining patterns

A mature skill often uses all four at once:

- `SKILL.md` is a thin overview (Pattern 1).
- References are split by sub-domain (Pattern 2).
- Rare edge cases are linked-out from each section (Pattern 3).
- Deterministic steps live in `scripts/` with explicit "Run:" instructions (Pattern 4).

A map-data skill is a natural candidate for all four once it grows — thin SKILL.md, references split by feature type (e.g. `overpass-queries.md`, `projection.md`), edge cases for unusual POI types linked out, and all fetch/convert work in scripts.
