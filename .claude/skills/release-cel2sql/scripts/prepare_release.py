#!/usr/bin/env python3
"""Prepare CHANGELOG.md for a cel2sql release.

Usage: python prepare_release.py vX.Y.Z   (or just X.Y.Z)

What it does:
  1. Reads CHANGELOG.md from the repo root.
  2. Finds the `## [Unreleased]` section.
  3. Errors out if [Unreleased] is empty.
  4. Renames `## [Unreleased]` → `## [X.Y.Z] - <today>` and inserts a fresh
     empty `## [Unreleased]` heading above it.
  5. Writes CHANGELOG.md back in place.
  6. Prints a checklist of commits since the last tag, plus the suggested
     follow-up commands (branch, commit, PR, tag, push).

Exits 0 on success. Use --check to dry-run (print the would-be output without
modifying CHANGELOG.md).
"""
from __future__ import annotations

import argparse
import datetime
import re
import subprocess
import sys
from pathlib import Path


def find_repo_root() -> Path:
    cur = Path.cwd().resolve()
    for d in [cur, *cur.parents]:
        if (d / "go.mod").exists() and (d / "CHANGELOG.md").exists():
            return d
    raise SystemExit("error: run from cel2sql repo root (go.mod and CHANGELOG.md must be siblings)")


SEMVER_RE = re.compile(r"^v?(\d+)\.(\d+)\.(\d+)$")
UNRELEASED_HEAD_RE = re.compile(r"^## \[Unreleased\]\s*$", re.M)
NEXT_HEAD_RE = re.compile(r"^## \[", re.M)


def parse_version(arg: str) -> str:
    m = SEMVER_RE.match(arg.strip())
    if not m:
        raise SystemExit(f"error: version {arg!r} must look like X.Y.Z or vX.Y.Z")
    return f"{m.group(1)}.{m.group(2)}.{m.group(3)}"


def split_changelog(text: str) -> tuple[str, str, str]:
    """Return (head, unreleased_block, tail).

    head is the prefix before `## [Unreleased]`, unreleased_block is the
    Unreleased section content (without the heading), tail is everything
    from the next `## [...]` heading onward.
    """
    head_match = UNRELEASED_HEAD_RE.search(text)
    if not head_match:
        raise SystemExit("error: CHANGELOG.md has no `## [Unreleased]` section")

    after_heading = head_match.end()
    next_match = NEXT_HEAD_RE.search(text, after_heading)
    if not next_match:
        raise SystemExit("error: CHANGELOG.md has no version section after `## [Unreleased]`")

    head = text[: head_match.start()]
    unreleased = text[after_heading : next_match.start()]
    tail = text[next_match.start() :]
    return head, unreleased, tail


def is_empty_unreleased(block: str) -> bool:
    """The Unreleased section is empty if it contains only blank lines."""
    return all(not line.strip() for line in block.splitlines())


def rewrite_changelog(text: str, version: str, today: str) -> str:
    head, unreleased, tail = split_changelog(text)
    if is_empty_unreleased(unreleased):
        raise SystemExit(
            "error: `## [Unreleased]` is empty — nothing to release. "
            "Add changelog entries first."
        )

    new = (
        head
        + "## [Unreleased]\n"
        + "\n"
        + f"## [{version}] - {today}"
        + unreleased
        + tail
    )
    return new


def git_commits_since_last_tag() -> list[str]:
    try:
        last_tag = subprocess.check_output(
            ["git", "describe", "--tags", "--abbrev=0"], text=True
        ).strip()
    except subprocess.CalledProcessError:
        return ["(no prior tag — git describe failed)"]

    try:
        log = subprocess.check_output(
            ["git", "log", f"{last_tag}..HEAD", "--oneline"], text=True
        ).strip()
    except subprocess.CalledProcessError as exc:
        return [f"(git log failed: {exc})"]

    if not log:
        return ["(no commits since last tag)"]
    return [f"  {line}" for line in log.splitlines()]


def main() -> None:
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("version", help="target version, X.Y.Z or vX.Y.Z")
    ap.add_argument("--check", action="store_true", help="dry-run; print result without writing")
    args = ap.parse_args()

    version = parse_version(args.version)
    today = datetime.date.today().isoformat()

    root = find_repo_root()
    changelog = root / "CHANGELOG.md"
    text = changelog.read_text(encoding="utf-8")

    new_text = rewrite_changelog(text, version, today)

    if args.check:
        print("=== DRY RUN — CHANGELOG.md would become: ===")
        # Print the diff-relevant section (the new versioned block).
        excerpt_start = new_text.find(f"## [{version}] - {today}")
        if excerpt_start >= 0:
            excerpt_end = new_text.find("## [", excerpt_start + 1)
            print(new_text[excerpt_start:excerpt_end].rstrip())
        return

    changelog.write_text(new_text, encoding="utf-8")
    print(f"✓ CHANGELOG.md updated: [Unreleased] → [{version}] - {today}\n")

    print("Commits since last tag:")
    for line in git_commits_since_last_tag():
        print(line)

    print()
    print("Suggested next steps:")
    print(f"  git checkout -b chore/release-v{version}")
    print(f"  git add CHANGELOG.md")
    print(f"  git commit -m 'chore: prepare CHANGELOG for v{version} release'")
    print(f"  git push -u origin chore/release-v{version}")
    print(f"  gh pr create --title 'chore: prepare CHANGELOG for v{version} release' --body '...'")
    print()
    print("After the PR merges:")
    print(f"  git checkout main && git pull --ff-only origin main")
    print(f"  git tag -a v{version} -m 'Release v{version}'")
    print(f"  git push origin v{version}      # triggers .github/workflows/release.yml")


if __name__ == "__main__":
    main()
