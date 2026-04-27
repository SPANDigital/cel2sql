#!/usr/bin/env python3
"""Verify that a dialect has WantSQL or SkipDialect coverage for every shared test case.

Usage: python check_testcase_coverage.py <dialect-name>

Where <dialect-name> is the lowercase name registered in dialect/dialect.go
(postgresql, mysql, sqlite, duckdb, bigquery, spark, ...).

Scans testcases/*.go for ConvertTestCase / ParameterizedTestCase blocks.
For each block, verifies that either:
  - WantSQL[dialect.<Name>] is present, OR
  - WantErr[dialect.<Name>] is true, OR
  - SkipDialect[dialect.<Name>] has a reason string.

Run from the repo root. Exits 0 if coverage is complete, 1 otherwise.
"""
from __future__ import annotations

import re
import sys
from pathlib import Path

REPO_ROOT_MARKERS = ("go.mod", "testcases")


def find_repo_root() -> Path:
    cur = Path.cwd().resolve()
    for d in [cur, *cur.parents]:
        if all((d / m).exists() for m in REPO_ROOT_MARKERS):
            return d
    raise SystemExit("error: run from cel2sql repo root (go.mod and testcases/ must be siblings)")


CASE_RE = re.compile(r"^\s*\{\s*$", re.M)


def find_blocks(text: str) -> list[tuple[int, int, str]]:
    """Yield every `{ ... }` block in the file.

    Walks the source character by character, skipping string and comment
    contents, and emits one region per balanced brace pair. The caller
    filters down to actual test case literals via `is_test_case()`.
    """
    out: list[tuple[int, int, str]] = []
    in_str = False
    str_q = ""
    in_comment = False
    line_comment = False
    starts: list[int] = []  # stack of open-brace positions

    i = 0
    while i < len(text):
        ch = text[i]
        nxt = text[i + 1] if i + 1 < len(text) else ""

        if line_comment:
            if ch == "\n":
                line_comment = False
            i += 1
            continue
        if in_comment:
            if ch == "*" and nxt == "/":
                in_comment = False
                i += 2
                continue
            i += 1
            continue
        if in_str:
            if ch == "\\" and str_q != "`":
                i += 2
                continue
            if ch == str_q:
                in_str = False
                str_q = ""
            i += 1
            continue

        if ch == "/" and nxt == "/":
            line_comment = True
            i += 2
            continue
        if ch == "/" and nxt == "*":
            in_comment = True
            i += 2
            continue
        if ch in ('"', "'", "`"):
            in_str = True
            str_q = ch
            i += 1
            continue

        if ch == "{":
            starts.append(i)
        elif ch == "}":
            if starts:
                start = starts.pop()
                out.append((start, i + 1, text[start : i + 1]))
        i += 1
    return out


def has_dialect_entry(block: str, dialect: str) -> bool:
    """Return True if block contains any of:
      WantSQL[ ... dialect.<Name> ... ]
      WantErr[ ... dialect.<Name> ... ]
      SkipDialect[ ... dialect.<Name> ... ]
    The map literals span lines, so we look for `dialect.<Name>:` anywhere
    inside one of those map blocks. We also accept the bare token `dialect.<Name>`
    as a conservative match.
    """
    pascal = dialect.capitalize()
    if dialect == "postgresql":
        pascal = "PostgreSQL"
    elif dialect == "mysql":
        pascal = "MySQL"
    elif dialect == "sqlite":
        pascal = "SQLite"
    elif dialect == "duckdb":
        pascal = "DuckDB"
    elif dialect == "bigquery":
        pascal = "BigQuery"
    elif dialect == "spark":
        pascal = "Spark"

    token = re.compile(rf"\bdialect\.{pascal}\b")
    return bool(token.search(block))


def case_name(block: str) -> str:
    m = re.search(r'Name:\s*"([^"]+)"', block)
    return m.group(1) if m else "<anonymous>"


NAME_FIELD_RE = re.compile(r'\bName:\s*"')


def is_test_case(block: str) -> bool:
    """Filter to only individual ConvertTestCase / ParameterizedTestCase blocks.

    Heuristic: exactly one `Name: "…"` field (so we exclude the enclosing
    slice literal and function body, which contain many) and at least one of
    WantSQL / WantErr / WantParams / SkipDialect (so we exclude unrelated
    struct literals like fixture FieldSchemas).
    """
    if len(NAME_FIELD_RE.findall(block)) != 1:
        return False
    return bool(re.search(r"\b(WantSQL|WantErr|WantParams|SkipDialect)\b", block))


def main() -> None:
    if len(sys.argv) != 2:
        print(__doc__, file=sys.stderr)
        sys.exit(2)
    dialect = sys.argv[1].strip().lower()
    if not re.match(r"^[a-z][a-z0-9_]*$", dialect):
        print(f"error: invalid dialect name {dialect!r}", file=sys.stderr)
        sys.exit(2)

    root = find_repo_root()
    testcases_dir = root / "testcases"

    missing: list[tuple[str, str]] = []
    total = 0

    for path in sorted(testcases_dir.glob("*_tests.go")):
        text = path.read_text(encoding="utf-8")
        for _, _, block in find_blocks(text):
            if not is_test_case(block):
                continue
            total += 1
            if not has_dialect_entry(block, dialect):
                missing.append((path.name, case_name(block)))

    if missing:
        print(f"MISSING coverage for dialect.{dialect} in {len(missing)} of {total} test case(s):")
        for fname, name in missing:
            print(f"  {fname}: {name}")
        sys.exit(1)

    print(f"OK: dialect.{dialect} has coverage in all {total} shared test case(s).")
    sys.exit(0)


if __name__ == "__main__":
    main()
