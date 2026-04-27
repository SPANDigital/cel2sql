# Skill Review Checklist

Walk this list before merging any new skill or edit. Lint runs first; the checklist catches things the linter can't.

## Frontmatter

- [ ] `name` present, lowercase-hyphenated, matches the folder name exactly.
- [ ] `description` present, third person, ≤ 2 sentences.
- [ ] Description includes **what** the skill does AND **when** to use it.
- [ ] Description includes concrete keywords an agent can match on (tool names, formats, file types, domain terms) — not just vague categories.

## SKILL.md body

- [ ] `wc -l SKILL.md` ≤ 500. Target ≤ 200.
- [ ] No time-sensitive wording in the main flow ("after August 2025", "until we migrate"). Move legacy content to a collapsed section.
- [ ] All relative links use forward slashes.
- [ ] All relative links stay within the skill folder (no `../` into other skills). Duplicate a short snippet if needed; let the other skill trigger separately for longer content.
- [ ] Every script reference is explicit about intent — **Run** (execute) vs **See** (read as reference).

## References

- [ ] References are one level deep — SKILL.md links to them; references do NOT link to further reference files.
- [ ] Reference files over ~100 lines open with a Contents / TOC block.
- [ ] Filenames describe content (`collision-authoring.md`, `finance.md`), not position (`doc2.md`, `ref1.md`).

## Scripts

- [ ] Each script starts with a docstring or one-line comment stating its purpose.
- [ ] SKILL.md shows the exact invocation (`python scripts/foo.py <args>`) and the expected output shape.
- [ ] Scripts take inputs from args or stdin — no paths hardcoded to the author's machine.
- [ ] Scripts exit with non-zero on failure; success is silent or produces the documented output.

## Assets

- [ ] Binary assets are small enough to ship in the repo (target ≤ 100 KB each).
- [ ] Asset filenames describe content.

## Evaluation

- [ ] At least one concrete scenario (a real prompt + expected behaviour) exercises this skill's main path.
- [ ] Scenarios were written **before** SKILL.md prose, or at minimum drove the prose content.

## Lint

- [ ] `python <path-to>/skill-authoring/scripts/lint_skill.py <path-to>/<skill>` exits 0.
- [ ] Any warnings are either fixed or explicitly justified in the PR description.
