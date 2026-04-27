---
name: skill-authoring
description: Authors and reviews agent skills for this repo, enforcing progressive disclosure (metadata → instructions → resources), ≤500-line SKILL.md files, domain-split references, scripts over prose for deterministic work, and specific third-person descriptions. Use when creating a new skill, splitting an oversized SKILL.md, or reviewing an existing skill for structure.
---

# Skill Authoring

This skill codifies the authoring discipline so every skill loads cheaply and triggers accurately.

## Why progressive disclosure

An agent skill is a folder with `SKILL.md` at its root. The agent loads content across three levels — each step costs
more context, so keep lower levels lean.

| Level           | When loaded                                       | Budget                     | Content                                 |
|-----------------|---------------------------------------------------|----------------------------|-----------------------------------------|
| 1. Metadata     | Always, at startup                                | ~100 tokens per skill      | `name` + `description` from frontmatter |
| 2. Instructions | When the skill triggers                           | <5,000 tokens (≤500 lines) | `SKILL.md` body                         |
| 3. Resources    | On demand, via explicit links or script execution | Effectively unlimited      | `references/`, `scripts/`, `assets/`    |

Scripts are especially efficient: the source never enters context — only the **output** does.

## Authoring a new skill

1. **Copy the template.** Copy `templates/SKILL.md.tmpl` from this skill into your new skill's folder as `SKILL.md`,
   then fill it in. Place the skill wherever skills live in the current repo (e.g. `skills/<your-skill>/`).
2. **Write the description first.** Third person, include *what* the skill does AND *when* to use it.
   See [templates/description-examples.md](templates/description-examples.md) for good vs bad examples.
3. **Draft 3 evaluation scenarios before prose.** Capture concrete prompts the skill should handle (and ones it should
   decline), alongside the expected behaviour. Real failures drive what you document — this keeps the skill lean.
4. **Write the minimum SKILL.md needed to pass the evals.** Resist documenting edge cases you haven't seen fail.
5. **Split early if content grows.** See [references/patterns.md](references/patterns.md) for the four split patterns.
6. **Lint before committing.** **Run** the linter from this skill against your new skill directory — see the
   Resources section below for the exact script path.
7. **Walk the review checklist.** See [references/review-checklist.md](references/review-checklist.md).

## Hard authoring rules

- `SKILL.md` ≤ 500 lines. Target ≤ 200.
- Description is third person, specific, includes what + when. Vague descriptions don't trigger.
- References are one level deep — `SKILL.md` links to reference files, references do NOT link to further references.
  Agents often preview links with `head -100` and miss content past that on chained links.
- Reference files over ~100 lines start with a Contents / TOC block.
- Name files by content (`references/finance.md`), not position (`docs/file2.md`).
- Use forward slashes in paths — backslashes break on Unix.
- Avoid time-sensitive wording (`"after August 2025..."`) in the main flow. Move legacy content into a collapsed "old
  patterns" section.
- Deterministic work belongs in `scripts/`, not Markdown.

## Execution intent: be explicit

When referencing a script, state whether the agent should **run** it or **read** it:

- **Run** `python scripts/fetch_osm.py --bbox ...` — produces GeoJSON on stdout. (Agent executes.)
- **See** `scripts/fetch_osm.py` for the Overpass query it issues. (Agent reads as reference.)

Execution is almost always what you want — scripts are more reliable than generated code, cheaper in context, and
consistent across runs.

## When to split SKILL.md

Three signals:

1. **Length.** File is pushing 500 lines.
2. **Mutually-exclusive contexts.** Two user tasks never need the same sections. Example: a Bo-Kaap question never needs
   Camps Bay lore.
3. **Code vs prose.** Long inline code blocks the agent will re-type anyway — move to `scripts/`.

The four split patterns, with worked examples specific to this repo's skills, live
in [references/patterns.md](references/patterns.md).

## Folder layout

```
<skill-name>/
├── SKILL.md          # Required
├── references/       # Optional — loaded on demand via explicit link
├── scripts/          # Optional — executed, source never loaded
├── templates/        # Optional — starter files users copy
└── assets/           # Optional — images, palettes, map files
```

## Review and maintenance

- Before merging a new or edited skill: run the linter and walk the checklist.
- If a reference file grows past ~100 lines, add a TOC.
- If `SKILL.md` grows past ~300 lines, plan a split before it hits 500.
- If a skill stops triggering when expected, tighten its description — add concrete keywords from real user prompts.

## Resources

- [references/patterns.md](references/patterns.md) — the four split patterns with repo-specific examples.
- [references/review-checklist.md](references/review-checklist.md) — pre-merge checklist.
- [templates/SKILL.md.tmpl](templates/SKILL.md.tmpl) — starter template to copy.
- [templates/description-examples.md](templates/description-examples.md) — good vs bad description examples.
- **Run** `python scripts/lint_skill.py <skill-dir>` — lints a skill directory for the rules above.
