# Touchpoints for Adding a Dialect

Exhaustive file-by-file checklist. Paths are relative to the repo root. Replace `<name>` with the lowercase dialect name (e.g. `spark`) and `<Name>` with the PascalCase form (e.g. `Spark`).

## Contents

- Required files
- Optional files
- Files to modify
- Verification checklist

## Required files

### 1. Register the name

`dialect/dialect.go` — add a constant to the `Name` block:

```go
const (
    PostgreSQL Name = "postgresql"
    MySQL      Name = "mysql"
    // …
    <Name>     Name = "<name>"
)
```

### 2. SQL emitter

`dialect/<name>/dialect.go` — implements the full `dialect.Dialect` interface (~47 methods). Pattern:

```go
package <name>

type Dialect struct{}

func New() *Dialect { return &Dialect{} }

func init() {
    dialect.Register(dialect.<Name>, func() dialect.Dialect { return New() })
}

var _ dialect.Dialect = (*Dialect)(nil)

func (d *Dialect) Name() dialect.Name { return dialect.<Name> }
// … all other methods
```

The compile-time `var _ dialect.Dialect = (*Dialect)(nil)` line catches any missing methods. The `init()` registration makes the dialect findable via `dialect.Get(dialect.<Name>)`.

### 3. Regex conversion

`dialect/<name>/regex.go` — exports `convertRE2To<Name>(pattern string) (converted string, caseInsensitive bool, err error)`. Reuse the ReDoS validator pattern from `dialect/duckdb/regex.go` (length cap 500, no nested quantifiers, ≤20 capture groups, max nesting depth 10, no quantified alternation).

If the engine's regex flavor is RE2-compatible (Spark's Java regex, BigQuery's RE2): pass through unchanged after validation.
If POSIX (PostgreSQL): translate `\d` → `[[:digit:]]`, `\w` → `[[:alnum:]_]`, `\s` → `[[:space:]]`, `\b` → `\y`.
If unsupported: return `dialect.ErrUnsupportedFeature` and set `SupportsRegex()` to `false`.

### 4. Validation

`dialect/<name>/validation.go` — exports `validateFieldName(name string) error` and the `reservedSQLKeywords map[string]bool`. Field-name regex is typically `^[a-zA-Z_][a-zA-Z0-9_]*$`.

### 5. Type provider

`<name>/provider.go` — implements `types.Provider` for CEL. Mirror one of:

- `mysql/provider.go` for caller-owned `*sql.DB` (most common — works with MySQL, SQLite, DuckDB, Spark).
- `pg/provider.go` for owned `pgxpool.Pool` (PostgreSQL-specific connection-string form).
- `bigquery/provider.go` for non-`database/sql` clients.

Required: `NewSchema`, `NewTypeProvider(map[string]Schema)`, `NewTypeProviderWithConnection(ctx, db)` (or `WithClient`), `LoadTableSchema(ctx, name)`, `Close()`, plus the `types.Provider` methods (`FindStructType`, `FindStructFieldNames`, `FindStructFieldType`, etc.).

`LoadTableSchema` strategy:
- `information_schema.columns` (MySQL, DuckDB) — preferred, supports parameterized table names.
- `PRAGMA table_info(<name>)` (SQLite) — table name must be regex-validated, `// #nosec G202` suppression required.
- `DESCRIBE TABLE <name>` (Spark) — same as PRAGMA: regex-validate, `#nosec G202`.
- Vendor client API (BigQuery via `Table.Metadata`).

### 6. Provider unit tests

`<name>/provider_test.go` — cover type mapping, schema lookup (`FindStructType`, `FindStructFieldType`), name validation rejection cases (SQL injection inputs).

### 7. Test environment factory

`testutil/env.go` — add `<Name>EnvFactory()` and a case in `DialectEnvFactory()`. Pattern:

```go
func <Name>EnvFactory() func(envSetup string) (*EnvResult, error) {
    return func(envSetup string) (*EnvResult, error) {
        switch envSetup {
        case testcases.EnvDefault:    /* … */
        case testcases.EnvWithTimestamp: /* … */
        case testcases.EnvWithJSON:   /* … */
        }
        result.Opts = append(result.Opts, cel2sql.WithDialect(<name>Dialect.New()))
        return result, nil
    }
}
```

### 8. Shared test runner

`testutil/runner_<name>_test.go` — 17-line file:

```go
func Test<Name>SharedCases(t *testing.T) {
    testutil.RunAllConvertTests(t, dialect.<Name>, testutil.<Name>EnvFactory())
}
func Test<Name>ParameterizedSharedCases(t *testing.T) {
    testutil.RunAllParameterizedTests(t, dialect.<Name>, testutil.<Name>EnvFactory())
}
```

### 9. Shared test cases

`testcases/*.go` — for **every** test case across `basic_tests.go`, `operator_tests.go`, `string_tests.go`, `regex_tests.go`, `cast_tests.go`, `array_tests.go`, `timestamp_tests.go`, `comprehension_tests.go`, `json_tests.go`, `parameterized_tests.go`: add a `dialect.<Name>` entry to `WantSQL` (and `WantParams` for parameterized cases). Where the dialect cannot express a case, add `SkipDialect[dialect.<Name>] = "<reason>"`.

Use the `gen_expected_sql.go.tmpl` template to generate actual outputs in bulk.

## Optional files

### Index advisor

`dialect/<name>/index_advisor.go` — implements the `IndexAdvisor` interface. Skip and return `false` from `SupportsIndexAnalysis()` if the engine has no user-controllable indexes (Spark on Delta, Iceberg, plain Parquet).

## Files to modify

### `README.md`

- Title and intro: "N SQL dialects" → "N+1 SQL dialects".
- Badges row: add a badge for the new engine.
- Feature list bullet near the top: bump count.
- "Multi-Dialect Support" intro: bump count + add to import list + add `WithDialect()` example.
- Dialect comparison table (around the "Dialect Comparison" heading): add a new column.
- Per-dialect type providers import block: add the import line.
- Per-dialect index types table (under "Query Analysis"): add a new column.

### `CLAUDE.md`

- Project Overview line: "supports six SQL dialects" → "seven", add the new dialect name.
- Core Components numbered list: add an entry for `<name>/provider.go`.
- Per-dialect type providers section: add the new dialect to the connection-handling list.
- Project Structure tree: add `<name>/` and `dialect/<name>/`.
- Testing Guidelines: add `<name>.NewTypeProvider()` to the per-dialect provider list.

## Verification checklist

1. `make fmt && make build && make lint` clean.
2. `go test -count=1 ./dialect/<name>/... ./<name>/...` passes.
3. `go test -count=1 -run Test<Name>SharedCases ./testutil/` passes.
4. `python .claude/skills/add-sql-dialect/scripts/check_testcase_coverage.py <name>` reports full coverage.
5. `go test -count=1 -short -race ./...` — no regressions in existing dialects.
6. `gosec ./<name>/... ./dialect/<name>/...` — 0 issues (note any `#nosec` annotations).
7. PR description lists the new dialect and the index-analysis stance (supported / out-of-scope).
