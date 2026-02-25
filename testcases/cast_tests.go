package testcases

import "github.com/spandigital/cel2sql/v3/dialect"

// CastTests returns test cases for type casting operations.
func CastTests() []ConvertTestCase {
	return []ConvertTestCase{
		{
			Name:     "cast_bool",
			CELExpr:  `bool(0) == false`,
			Category: CategoryCast,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "CAST(0 AS BOOLEAN) IS FALSE",
				dialect.MySQL:      "CAST(0 AS UNSIGNED) IS FALSE",
				dialect.SQLite:     "CAST(0 AS INTEGER) IS FALSE",
				dialect.DuckDB:     "CAST(0 AS BOOLEAN) IS FALSE",
				dialect.BigQuery:   "CAST(0 AS BOOL) IS FALSE",
			},
		},
		{
			Name:     "cast_bytes",
			CELExpr:  `bytes("test")`,
			Category: CategoryCast,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "CAST('test' AS BYTEA)",
				dialect.MySQL:      "CAST('test' AS BINARY)",
				dialect.SQLite:     "CAST('test' AS BLOB)",
				dialect.DuckDB:     "CAST('test' AS BLOB)",
				dialect.BigQuery:   "CAST('test' AS BYTES)",
			},
		},
		{
			Name:     "cast_int",
			CELExpr:  `int(true) == 1`,
			Category: CategoryCast,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "CAST(TRUE AS BIGINT) = 1",
				dialect.MySQL:      "CAST(TRUE AS SIGNED) = 1",
				dialect.SQLite:     "CAST(TRUE AS INTEGER) = 1",
				dialect.DuckDB:     "CAST(TRUE AS BIGINT) = 1",
				dialect.BigQuery:   "CAST(TRUE AS INT64) = 1",
			},
		},
		{
			Name:     "cast_string",
			CELExpr:  `string(true) == "true"`,
			Category: CategoryCast,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "CAST(TRUE AS TEXT) = 'true'",
				dialect.MySQL:      "CAST(TRUE AS CHAR) = 'true'",
				dialect.SQLite:     "CAST(TRUE AS TEXT) = 'true'",
				dialect.DuckDB:     "CAST(TRUE AS VARCHAR) = 'true'",
				dialect.BigQuery:   "CAST(TRUE AS STRING) = 'true'",
			},
		},
		{
			Name:     "cast_string_from_timestamp",
			CELExpr:  `string(created_at)`,
			Category: CategoryCast,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "CAST(created_at AS TEXT)",
				dialect.MySQL:      "CAST(created_at AS CHAR)",
				dialect.SQLite:     "CAST(created_at AS TEXT)",
				dialect.DuckDB:     "CAST(created_at AS VARCHAR)",
				dialect.BigQuery:   "CAST(created_at AS STRING)",
			},
		},
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
			},
		},
	}
}
