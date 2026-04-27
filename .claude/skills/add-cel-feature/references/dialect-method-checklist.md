# Dialect Method Checklist

Steps when adding a new method to the `dialect.Dialect` interface. Skip the new-method route entirely if the SQL output is identical across all six dialects (just inline in the converter).

## Contents

- Step-by-step
- Compile-time safety net
- Naming conventions
- Error returns
- Worked example: WriteUnixSeconds

## Step-by-step

1. **Declare the method on the interface** in `dialect/dialect.go`. Group it with related methods (Literals / Operators / Type Casting / Arrays / JSON / Timestamps / String Functions / Comprehensions / Struct / Validation / Regex / Capabilities). Add a one-line doc comment naming the canonical PostgreSQL form ("For PostgreSQL: …") to anchor reviewers.
2. **Implement in all six dialect packages.** The canonical six are:
   - `dialect/postgres/dialect.go`
   - `dialect/mysql/dialect.go`
   - `dialect/sqlite/dialect.go`
   - `dialect/duckdb/dialect.go`
   - `dialect/bigquery/dialect.go`
   - `dialect/spark/dialect.go`
3. **Verify the compile-time assertion fires.** Each dialect file ends with `var _ dialect.Dialect = (*Dialect)(nil)` — running `go build ./...` will fail if any dialect is missing the new method.
4. **Wire the method into the converter** in `cel2sql.go` (or whichever visitor file owns the surface area — see the converter-file-map reference). The converter calls `con.dialect.<NewMethod>(...)`.
5. **Add a representative test case** to `testcases/<category>_tests.go`. Set `WantSQL[dialect.<Name>]` for every dialect or document a `SkipDialect[dialect.<Name>]` reason.
6. **Run the testcase coverage check** for every dialect:
   ```
   for d in postgresql mysql sqlite duckdb bigquery spark; do
     python .claude/skills/add-sql-dialect/scripts/check_testcase_coverage.py "$d"
   done
   ```

## Compile-time safety net

The `var _ dialect.Dialect = (*Dialect)(nil)` line at the top of every dialect file is critical. It's the only mechanism that guarantees all dialects implement every method. If you add a new dialect, add this line first.

## Naming conventions

Methods on the `Dialect` interface that emit SQL fragments are prefixed `Write…`. They take a `*strings.Builder` (the converter's output buffer) plus callbacks `func() error` for any sub-expressions the dialect needs to interleave. They return `error` so dialects that don't support a feature can return `dialect.ErrUnsupportedFeature`.

Capability check methods are `Supports…` and return `bool`.

## Error returns

`dialect.ErrUnsupportedFeature` is the sentinel for "this dialect cannot express this CEL construct." Wrap it with `fmt.Errorf("%w: <reason>", dialect.ErrUnsupportedFeature, ...)` so callers can `errors.Is()` it.

Examples already in the codebase:
- `dialect/sqlite/dialect.go` `WriteRegexMatch` returns `ErrUnsupportedFeature` (SQLite has no native regex).
- `dialect/spark/dialect.go` `WriteArrayLength` returns `ErrUnsupportedFeature` for `dimension > 1` (Spark doesn't support multi-dim arrays).

## Worked example: WriteUnixSeconds

Suppose CEL's `int(timestamp)` should emit a per-dialect Unix-seconds expression.

1. **Interface** (`dialect/dialect.go`):
   ```go
   // WriteEpochExtract writes extraction of epoch from a timestamp.
   // For PostgreSQL: EXTRACT(EPOCH FROM expr)::bigint.
   WriteEpochExtract(w *strings.Builder, writeExpr func() error) error
   ```
2. **PostgreSQL** (`dialect/postgres/dialect.go`):
   ```go
   func (d *Dialect) WriteEpochExtract(w *strings.Builder, writeExpr func() error) error {
       w.WriteString("EXTRACT(EPOCH FROM ")
       if err := writeExpr(); err != nil { return err }
       w.WriteString(")::bigint")
       return nil
   }
   ```
3. **MySQL**: `UNIX_TIMESTAMP(<expr>)`.
4. **SQLite**: `CAST(strftime('%s', <expr>) AS INTEGER)`.
5. **DuckDB**: `EXTRACT(EPOCH FROM <expr>)::BIGINT`.
6. **BigQuery**: `UNIX_SECONDS(<expr>)`.
7. **Spark**: `UNIX_TIMESTAMP(<expr>)`.

(All six already implement this method — it's a real example from the codebase.)

8. **Test case** in `testcases/cast_tests.go`:
   ```go
   {
       Name:     "cast_int_epoch",
       CELExpr:  `int(created_at)`,
       Category: CategoryCast,
       WantSQL: map[dialect.Name]string{
           dialect.PostgreSQL: "EXTRACT(EPOCH FROM created_at)::bigint",
           dialect.MySQL:      "UNIX_TIMESTAMP(created_at)",
           dialect.SQLite:     "CAST(strftime('%s', created_at) AS INTEGER)",
           dialect.DuckDB:     "EXTRACT(EPOCH FROM created_at)::BIGINT",
           dialect.BigQuery:   "UNIX_SECONDS(created_at)",
           dialect.Spark:      "UNIX_TIMESTAMP(created_at)",
       },
   },
   ```
