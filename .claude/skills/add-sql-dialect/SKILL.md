---
name: add-sql-dialect
description: Adds a new SQL dialect to cel2sql by registering its name, implementing the Dialect interface in dialect/<name>/, adding a sibling type provider in <name>/, populating shared test cases in testcases/, and refreshing README.md and CLAUDE.md. Use when adding any new database backend (Trino, Snowflake, ClickHouse, MS SQL, Athena, Oracle) or porting cel2sql to another analytics engine.
---

# Add SQL Dialect

Adding a SQL dialect is the largest contribution shape in this repo (~1000-line PR, ~25 file touches in lockstep). The pattern is well-established — six dialects already follow it. This skill captures the procedure so the engineer can follow the template instead of reverse-engineering the layout from existing dialects.

## Quick start

1. **Pick the closest analogue** by syntax family — see "Picking the analogue" below.
2. **Register the name** — add `<Name> Name = "<name>"` to the constants block in `dialect/dialect.go`.
3. **Implement the four files under `dialect/<name>/`**: `dialect.go` (the SQL emitter), `regex.go` (RE2 → engine-native conversion + ReDoS validators), `validation.go` (reserved keywords + identifier regex), and optionally `index_advisor.go`.
4. **Implement the type provider under `<name>/`**: `provider.go` + `provider_test.go`. Mirror `mysql/provider.go` for caller-owned `*sql.DB`, or `pg/provider.go` for owned-pool.
5. **Wire the test runner**: add `<Name>EnvFactory()` to `testutil/env.go`, add the dialect to `DialectEnvFactory()`'s switch, create `testutil/runner_<name>_test.go`.
6. **Populate shared test cases** — see "Testcase coverage" below.
7. **Refresh docs** — see [references/touchpoints.md](references/touchpoints.md) for the exhaustive file list.
8. **Run** `python .claude/skills/skill-authoring/scripts/lint_skill.py .claude/skills/<name>` if you also added a new skill alongside.

## Picking the analogue

Decide by which existing dialect's syntax shape your target most resembles:

| Question | Operator-style → use DuckDB | Function-style → use BigQuery |
|---|---|---|
| Regex match | `target ~ 'p'` (PG, DuckDB) | `REGEXP_CONTAINS(target, 'p')` (BQ); `RLIKE` (Spark) |
| JSON access | `b->>'f'` (PG, DuckDB) | `JSON_VALUE(b, '$.f')` (BQ); `get_json_object()` (Spark) |
| Array literal | `ARRAY[…]` (PG); `[…]` (DuckDB, BQ) | `array(…)` (Spark) |
| Array index | 1-indexed (PG, DuckDB) | 0-indexed (BQ via OFFSET, Spark direct) |
| Param placeholder | `$N` (PG, DuckDB) | `?` (MySQL, SQLite, Spark) or `@pN` (BQ) |
| Cast to numeric | `::numeric` postfix (PG, DuckDB) | ` + 0` arithmetic coercion (MySQL, SQLite, Spark) |

For the half-dozen most-divergent methods at a glance, see [references/dialect-method-matrix.md](references/dialect-method-matrix.md). When in doubt, copy DuckDB and patch — its layout is the cleanest.

## Critical surface to map

These methods on `dialect.Dialect` are where dialects diverge most. Plan how to implement them before writing any code:

- `WriteRegexMatch` — operator vs function call
- `WriteJSONFieldAccess` — operator (`->>`) vs function wrapper (`JSON_VALUE`, `get_json_object`); whether intermediate vs final access uses different forms
- `WriteArrayLiteralOpen/Close` — `ARRAY[`, `[`, `array(`
- `WriteListIndex` / `WriteListIndexConst` — 0-indexed vs 1-indexed; bare `[i]` vs `OFFSET(i)` vs `+ 1`
- `WriteParamPlaceholder` — `$N`, `?`, `@pN`
- `WriteExtract` for DOW — Sunday-1 (BQ, Spark) vs Sunday-0 (PG/DuckDB convention) — adjust by `+ 6) % 7` or `- 1`
- `WriteCastToNumeric` — postfix `::TYPE` vs arithmetic coercion `+ 0`
- `WriteJSONArrayElements` — must be a **set-returning expression** (used in `FROM <here> AS iter`); use the engine's `EXPLODE`/`UNNEST`/`json_each`/`from_json` form
- `WriteJSONArrayMembership` — must produce a valid RHS for `lhs = ` (subquery form, like SQLite's `(SELECT value FROM json_each(...))`)

## Testcase coverage

The shared test runner (`testutil/RunAllConvertTests`) iterates `testcases/*.go` and reads the new dialect's `WantSQL[dialect.<Name>]` entry. Tests with no entry **and** no `SkipDialect` reason are silently skipped, which masks regressions.

To populate `WantSQL` efficiently, generate the actual output once and paste it back in:

1. **Run** `cp .claude/skills/add-sql-dialect/templates/gen_expected_sql.go.tmpl testutil/gen_<name>_test.go` and substitute `{{DIALECT}}` (lowercase) and `{{Dialect}}` (PascalCase) and `{{DialectFactory}}` in the file.
2. **Run** `go test -tags=<name>gen -v -run TestGen<Dialect>Expected ./testutil/`.
3. Paste outputs into the corresponding `testcases/*.go` `WantSQL` maps.
4. Delete `testutil/gen_<name>_test.go`.
5. **Run** `python .claude/skills/add-sql-dialect/scripts/check_testcase_coverage.py <name>` from the repo root to verify every test case has either a `WantSQL[dialect.<Name>]` entry or a `SkipDialect[dialect.<Name>]` reason. The script exits non-zero if anything is missing.

## Optional: IndexAdvisor

Implement `dialect/<name>/index_advisor.go` only if the engine has user-controllable indexes (BTREE, GIN, ART, CLUSTERING). Skip for storage-layer-driven engines like Spark (Delta Z-order, Iceberg sort) — set `SupportsIndexAnalysis()` to `false`.

## Doc refresh

When the implementation is green, refresh:

- `README.md`: bump dialect count, add badge, add column to the comparison table, add row to the index-analysis table, add to the import block.
- `CLAUDE.md`: Project Overview count, new entry in Core Components, new entry in Project Structure tree.

The full file-by-file checklist is in [references/touchpoints.md](references/touchpoints.md).

## Resources

- [references/dialect-method-matrix.md](references/dialect-method-matrix.md) — the most-divergent `Dialect` methods across the existing six dialects, side by side.
- [references/touchpoints.md](references/touchpoints.md) — exhaustive file-by-file checklist for a new dialect.
- [templates/gen_expected_sql.go.tmpl](templates/gen_expected_sql.go.tmpl) — build-tagged Go test template that prints actual SQL output for every shared test case.
- **Run** `python scripts/check_testcase_coverage.py <name>` — verifies `WantSQL`/`SkipDialect` coverage for the dialect across all testcases files.
