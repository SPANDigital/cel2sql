---
name: release-cel2sql
description: Cuts a new cel2sql release by preparing CHANGELOG.md, deciding the semver bump from commit history, opening a release-prep PR, then tagging and pushing after merge to trigger .github/workflows/release.yml. Use when shipping any new patch, minor, or major version of cel2sql.
---

# Release cel2sql

The release flow is short but easy to mess up — direct push to `main` is blocked, the CHANGELOG `[Unreleased]` block has to be split into a versioned section first, and the semver decision interacts with semantic-import-versioning (`/v3` → `/v4` is a major).

## Quick start

```
git checkout main && git pull --ff-only origin main
LAST=$(git describe --tags --abbrev=0)
git log "$LAST..main" --oneline                              # commits to release
python .claude/skills/release-cel2sql/scripts/prepare_release.py X.Y.Z  # patches CHANGELOG
git checkout -b chore/release-vX.Y.Z
git add CHANGELOG.md
git commit -m "chore: prepare CHANGELOG for vX.Y.Z release"
git push -u origin chore/release-vX.Y.Z
gh pr create --title "chore: prepare CHANGELOG for vX.Y.Z release" --body "..."
# After PR merges:
git checkout main && git pull --ff-only origin main
git tag -a vX.Y.Z -m "Release vX.Y.Z"
git push origin vX.Y.Z   # triggers .github/workflows/release.yml
```

The tag push triggers the release workflow, which runs tests + builds + creates the GitHub release with auto-generated notes from PR titles.

## Picking the version

See [references/semver-decision.md](references/semver-decision.md). Quick rules:

- **Patch** — dependency bumps that don't change behavior, doc-only fixes, test fixes.
- **Minor** — new feature, new option, new dialect.
- **Major** — requires module-path change `/v3` → `/v4`. Reserve for genuinely user-disruptive breakage. The Observe heuristic-removal in #113 was BREAKING but stayed within v3 because its blast radius was niche.

When unsure between patch and minor: if any user-visible API surface changed (new exported function/option/method, new dialect, new CEL feature), bump minor.

## Common slip-ups

- **Stale `[Unreleased]` content from a prior release.** When the previous version was tagged, its content stayed under `## [Unreleased]` instead of being moved to `## [3.x.y] - <date>`. Happened with v3.6.0: the multi-dialect entry sat under `[Unreleased]` until v3.7.0, then had to be backfilled. The `prepare_release.py` script handles this by inspecting the existing `[Unreleased]` block before rewriting.
- **Direct push to main blocked.** A direct `git push origin main` for the CHANGELOG commit will fail with the project's branch-protection rules. Always use a `chore/release-vX.Y.Z` branch and a release-prep PR.
- **Forgetting the tag annotation.** Use `git tag -a` (annotated) not `git tag` (lightweight). The release workflow reads the tag message.

## Verification

After pushing the tag, monitor the release workflow:

```
gh run list --workflow Release --limit 1 --repo SPANDigital/cel2sql --json databaseId,status,headBranch
gh run watch <run-id>
```

The workflow runs unit tests, builds, then creates the GitHub release using `softprops/action-gh-release@v2` with `generate_release_notes: true`. If it fails on stdlib vulns or a transitive dep CVE, that fix is *out of scope for this skill* — flag it as a follow-up dep-bump PR. The tagged release will not exist on GitHub until the workflow succeeds.

## Resources

- [references/semver-decision.md](references/semver-decision.md) — patch/minor/major decision tree, including the SIV `/v4` caveat.
- **Run** `python scripts/prepare_release.py vX.Y.Z` (or just `X.Y.Z`) — reads `CHANGELOG.md`, splits `[Unreleased]` into a versioned section dated today, prints a checklist of commits since the last tag and the commit/tag commands.
