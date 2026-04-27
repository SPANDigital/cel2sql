---
name: add-cel-feature
description: Adds new CEL operator, function, macro, or comprehension support to cel2sql across all SQL dialects, touching cel2sql.go visit*, the dialect.Dialect interface, every dialect/<name>/dialect.go, and testcases/<category>_tests.go. Use when wiring a new CEL function (e.g., string.toLowerCase, timestamp.toISOString), a new comprehension form, a new operator overload, or extending the converter to recognise a CEL extension package.
---

# Add CEL Feature

Adding a new CEL surface area (function, operator, macro, comprehension) is the second-most-common contribution shape after adding a dialect (5+ feature PRs in the recent history: string-extension functions, comprehensions, multi-dim arrays, regex, get/setDayOfWeek). The work is the same shape every time — this skill captures it.

## Quick start

1. **Classify the feature** — function call, operator, macro (like `has()`), or comprehension form (like `all`/`exists`/`filter`/`map`).
2. **Identify the converter file** that owns the surface area — see [references/converter-file-map.md](references/converter-file-map.md).
3. **Decide whether a new `dialect.Dialect` method is needed** — see "New method or inline?" below.
4. **If a new method is needed**, walk [references/dialect-method-checklist.md](references/dialect-method-checklist.md) — the new method must land in `dialect/dialect.go` plus all six dialect packages.
5. **Add a test case** to the appropriate `testcases/<category>_tests.go` with a `WantSQL` entry for every dialect (use the gen-script from `add-sql-dialect` if outputs diverge significantly).
6. **Run** `python .claude/skills/add-sql-dialect/scripts/check_testcase_coverage.py <each-dialect>` to confirm coverage is still complete.

## New method or inline?

The `dialect.Dialect` interface is already large (~47 methods). Add a new method only when needed; otherwise inline the SQL in the visitor.

| Situation | Verdict |
|---|---|
| The SQL is identical across all six dialects | Inline in the visitor (e.g., `con.str.WriteString("MOD(")` for `%`). |
| Any dialect needs a different syntax | New `dialect.Dialect` method. |
| A subset of dialects don't support the feature at all | New method + each unsupported dialect returns `dialect.ErrUnsupportedFeature`. |
| The feature is JSON-related, regex-related, or array-related | Look for an existing `Write…` method first — JSON/array/regex method coverage is broad. |

If you add a method, the compile-time `var _ dialect.Dialect = (*Dialect)(nil)` assertion at the top of every dialect file will catch missing implementations. The build won't pass until all six dialects are updated.

## Where features live

For example: adding `string.toLowerCase()`. CEL parses this as a method call on a string value. The visitor entry is `visitCall` in `cel2sql.go`. CEL's standard string functions are dispatched by overload ID — see how `startsWith`/`endsWith`/`contains` are wired and follow the same pattern. The output SQL (`LOWER(expr)`) is identical across all dialects, so this is a candidate for inline conversion (no new `Dialect` method).

For another example: adding `timestamp.toUnixSeconds()`. The output diverges (`EXTRACT(EPOCH FROM e)::bigint` PG, `UNIX_TIMESTAMP(e)` MySQL, etc.) — so add a new `Dialect.WriteUnixSeconds(w, writeExpr)` method. `cel2sql.go` is generic; per-dialect SQL goes in `dialect/<name>/dialect.go`.

For the file-by-file map, see [references/converter-file-map.md](references/converter-file-map.md).

## CEL environment registration

For a new function, the CEL environment must declare it before the converter can compile expressions that use it. Check if it's already registered in:

- `testutil/env.go` `NewDefaultEnv` / `NewTimestampEnv` — for the test environments.
- The user's CEL environment (out of scope for this skill).

If you add a function only to the test environment, the user-facing API (the documentation) needs to mention the new function so consumers know to register it themselves.

## Testcase coverage

Add a representative test to `testcases/<category>_tests.go` (or the appropriate category file). The test must have a `WantSQL[dialect.<Name>]` entry for every dialect, or a `SkipDialect[dialect.<Name>]` reason for any that legitimately can't express the feature.

If the SQL output diverges significantly across dialects, use the `gen_expected_sql.go.tmpl` template from the `add-sql-dialect` skill to generate the actual outputs in bulk:

**Run** `cp .claude/skills/add-sql-dialect/templates/gen_expected_sql.go.tmpl testutil/gen_<dialect>_test.go` for each dialect, transcribe outputs, delete the generators.

After the test case is in place, **run** the coverage check for every dialect:

```
for d in postgresql mysql sqlite duckdb bigquery spark; do
  python .claude/skills/add-sql-dialect/scripts/check_testcase_coverage.py "$d"
done
```

A new test case missing a single dialect's `WantSQL` entry silently skips for that dialect — easy regression to ship without this check.

## Resources

- [references/converter-file-map.md](references/converter-file-map.md) — which converter file owns which CEL surface area.
- [references/dialect-method-checklist.md](references/dialect-method-checklist.md) — checklist when adding a method to the `Dialect` interface.
- **Run** `python .claude/skills/add-sql-dialect/scripts/check_testcase_coverage.py <dialect>` — verifies test coverage; reuse the script across this skill and `add-sql-dialect`.
