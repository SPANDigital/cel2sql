# Semver Decision Tree for cel2sql

cel2sql uses Go semantic-import-versioning (SIV): the major version is encoded in the module path (`github.com/spandigital/cel2sql/v3`). A major bump is a module-path change that all downstream consumers must update.

## Contents

- Decision tree
- The SIV /v4 caveat
- Examples from history
- Breaking changes that stayed within v3

## Decision tree

| Change shape | Bump |
|---|---|
| Dependency upgrade with no behavior change | Patch |
| Doc-only fix | Patch |
| Test-only fix (e.g. resolving spurious CI annotations like #115) | Patch |
| Bug fix in SQL output that *narrows* output to be correct | Patch |
| New option (e.g. `WithJSONVariables`, `WithColumnAliases`) | Minor |
| New CEL feature exposed in default env | Minor |
| New dialect | Minor |
| New exported function on the public API | Minor |
| Removing or renaming any exported symbol | Major *only if* SIV path change is justified — see caveat |
| Changing the signature of an exported function | Major *only if* SIV path change is justified |
| Removing a registered dialect | Major |
| Module path change | Major |

## The SIV /v4 caveat

A major version bump in Go means renaming the module path from `github.com/spandigital/cel2sql/v3` to `github.com/spandigital/cel2sql/v4`. Every consumer must update their imports. This is **expensive churn** for the ecosystem — bump major only when the breakage is large enough to justify it.

In practice this means:
- Niche behavior fixes, even if technically breaking, can stay within v3 if accompanied by a prominent `### Changed (BREAKING)` CHANGELOG entry.
- Genuine API redesigns (the v3.0 BigQuery → PostgreSQL migration; a hypothetical re-shape of `dialect.Dialect`) warrant a major bump.

## Examples from history

- **v3.0.0**: BigQuery → PostgreSQL migration. Removed `cloud.google.com/go/bigquery` dependency. Required all consumers to refactor. Major.
- **v3.6.0**: Multi-dialect support added (MySQL, SQLite, DuckDB, BigQuery). Additive. Minor.
- **v3.7.0**: WithJSONVariables, WithColumnAliases, WithParamStartIndex options added; name-based numeric-cast heuristic removed. Last item was BREAKING but narrow (only affected callers using variables literally named `score` / `value` / `num` / `amount` / `count` / `level` for non-JSON usage). Stayed minor with a `### Changed (BREAKING)` flag.
- **v3.7.1**: Spurious `##[error]` annotation fix (#115) + dependabot dep bumps. Patch.
- **v3.8.0** *(forthcoming when Spark merges)*: New Spark dialect. Additive. Minor.

## Breaking changes that stayed within v3

The Observe heuristic-removal (#113) is the canonical example. The CHANGELOG entry:

```markdown
### Changed
- **BREAKING: Removed name-based numeric-cast heuristic in `visitIdent`** —
  identifiers named `score`, `value`, `num`, `amount`, `count`, or `level`
  are no longer auto-cast to `::numeric`. The heuristic was a footgun that
  incorrectly cast plain (non-JSON) variables that happened to share these
  names. Numeric casting is now driven solely by JSON text-extraction
  context (handled in `visitCall`). Callers that relied on the implicit
  cast should add an explicit cast in CEL or use the JSON comprehension
  paths. Backported from observeinc/cel2sql fork PR #1.
```

The `### Changed (BREAKING)` heading and the migration-guidance final sentence are load-bearing — they signal to consumers what to do. Without that note, the `/v4` bump becomes the only safe option.
