#!/usr/bin/env python3
"""Lint an agent skill directory against this repo's authoring rules.

Usage: python lint_skill.py <skill-dir>

Checks:
- <skill-dir>/SKILL.md exists
- Frontmatter has name (matches folder) and description
- Description is <= 2 sentences, third person, mentions "when" to use
- SKILL.md body <= 500 lines (warn at > 300)
- Reference files > 100 lines have a Contents/TOC block near the top
- SKILL.md relative links use forward slashes and stay inside the skill dir
- Reference files do not link to further reference files (one-level rule)
- Filenames are not positional (file1.md, doc2.md, ref3.md)
- SKILL.md body has no time-sensitive wording ("after <year>", etc.)

Exits 0 on clean run, 1 if any ERROR, prints WARNINGs either way.
Templates (files ending with .tmpl) are skipped.
"""
from __future__ import annotations

import re
import sys
from pathlib import Path

ERRORS: list[str] = []
WARNINGS: list[str] = []


def err(msg: str) -> None:
    ERRORS.append(msg)


def warn(msg: str) -> None:
    WARNINGS.append(msg)


FRONTMATTER_RE = re.compile(r"\A---\r?\n(.*?)\r?\n---\r?\n", re.DOTALL)


def parse_frontmatter(text: str):
    match = FRONTMATTER_RE.match(text)
    if not match:
        return None, text
    fm_block = match.group(1)
    body = text[match.end():]
    meta: dict[str, str] = {}
    current_key = None
    for line in fm_block.splitlines():
        if re.match(r"^[A-Za-z_][A-Za-z0-9_-]*:\s", line) or line.endswith(":"):
            key, _, value = line.partition(":")
            meta[key.strip()] = value.strip()
            current_key = key.strip()
        elif current_key and line.startswith(" "):
            meta[current_key] = (meta[current_key] + " " + line.strip()).strip()
    return meta, body


SENTENCE_SPLIT = re.compile(r"(?<=[.!?])\s+")
PERSON_WORDS = re.compile(r"\b(you|your|yours|we|our|ours|us|i|me|my)\b", re.I)
WHEN_PHRASE = re.compile(r"\b(use when|used when|used to|used for|use for|use this when)\b", re.I)
TIME_SENSITIVE = re.compile(r"\b(after|until|before)\s+(19|20)\d{2}\b", re.I)
LINK_RE = re.compile(r"\]\(([^)]+)\)")
POSITIONAL = re.compile(r"^(file|doc|page|ref|part|section)\d+\.(md|py|mjs|js|ts)$", re.I)


def lint_skill_md(path: Path, skill_dir: Path) -> None:
    try:
        text = path.read_text(encoding="utf-8")
    except UnicodeDecodeError as exc:
        err(f"{path}: could not be decoded as UTF-8: {exc}")
        return
    meta, body = parse_frontmatter(text)

    if meta is None:
        err(f"{path}: missing YAML frontmatter")
        return

    expected_name = path.parent.name
    name = meta.get("name", "").strip()
    if not name:
        err(f"{path}: frontmatter missing 'name'")
    elif name != expected_name:
        err(f"{path}: frontmatter name {name!r} does not match folder {expected_name!r}")

    desc = meta.get("description", "").strip().strip('"').strip("'")
    if not desc:
        err(f"{path}: frontmatter missing 'description'")
    else:
        sentences = [s for s in SENTENCE_SPLIT.split(desc) if s.strip()]
        if len(sentences) > 2:
            warn(f"{path}: description is {len(sentences)} sentences (target <= 2)")
        if PERSON_WORDS.search(desc):
            warn(f"{path}: description appears to use first/second person; prefer third person")
        if not WHEN_PHRASE.search(desc):
            warn(f"{path}: description should include 'when' to use the skill (e.g. 'Use when ...')")

    lines = body.splitlines()
    n = len(lines)
    if n > 500:
        err(f"{path}: SKILL.md body is {n} lines (hard limit 500)")
    elif n > 300:
        warn(f"{path}: SKILL.md body is {n} lines (target <= 200, warn >= 300)")

    if TIME_SENSITIVE.search(body):
        warn(f"{path}: body contains time-sensitive wording; move to a legacy section")

    skill_dir_resolved = skill_dir.resolve()
    for m in LINK_RE.finditer(body):
        link = m.group(1).strip()
        if link.startswith(("http://", "https://", "mailto:", "#")):
            continue
        if "\\" in link:
            err(f"{path}: link {link!r} uses backslashes; use forward slashes")
        target = link.split("#", 1)[0]
        if not target:
            continue
        if target.startswith("/"):
            err(f"{path}: link {link!r} is absolute; use a relative path inside the skill directory")
            continue
        resolved = (path.parent / target).resolve()
        try:
            resolved.relative_to(skill_dir_resolved)
        except ValueError:
            err(f"{path}: link {link!r} escapes the skill directory")


TOC_HEADING_RE = re.compile(r"^#{1,3}\s+(contents|table of contents|toc)\b", re.I | re.M)


def lint_reference(path: Path) -> None:
    try:
        text = path.read_text(encoding="utf-8")
    except UnicodeDecodeError as exc:
        err(f"{path}: could not be decoded as UTF-8: {exc}")
        return
    lines = text.splitlines()
    n = len(lines)

    if n > 100:
        head = "\n".join(lines[:40])
        if not TOC_HEADING_RE.search(head):
            warn(f"{path}: {n} lines without a TOC (add a Contents block for files over 100 lines)")

    for m in LINK_RE.finditer(text):
        link = m.group(1).strip()
        if link.startswith(("http://", "https://", "mailto:", "#")):
            continue
        target = link.split("#", 1)[0]
        if not target:
            continue
        if target.endswith(".md") and not target.startswith("../"):
            warn(f"{path}: reference links to another file {link!r}; references should be one level deep")


def lint_filenames(skill_dir: Path) -> None:
    for p in skill_dir.rglob("*"):
        if not p.is_file():
            continue
        if p.suffix == ".tmpl":
            continue
        if POSITIONAL.match(p.name):
            err(f"{p}: positional filename; rename to describe content")


def lint_directory(skill_dir: Path) -> None:
    if not skill_dir.is_dir():
        err(f"{skill_dir}: not a directory")
        return

    skill_md = skill_dir / "SKILL.md"
    if not skill_md.exists():
        err(f"{skill_dir}: missing SKILL.md")
        return

    lint_skill_md(skill_md, skill_dir)
    lint_filenames(skill_dir)

    for ref_dir_name in ("references", "reference"):
        ref_dir = skill_dir / ref_dir_name
        if ref_dir.exists():
            for p in sorted(ref_dir.rglob("*.md")):
                lint_reference(p)


def main() -> None:
    if len(sys.argv) != 2:
        print(__doc__, file=sys.stderr)
        sys.exit(2)

    target = Path(sys.argv[1])
    lint_directory(target)

    for w in WARNINGS:
        print(f"WARN  {w}")
    for e in ERRORS:
        print(f"ERROR {e}")

    print(f"\n{len(ERRORS)} error(s), {len(WARNINGS)} warning(s)")
    sys.exit(1 if ERRORS else 0)


if __name__ == "__main__":
    main()
