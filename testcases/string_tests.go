package testcases

import "github.com/spandigital/cel2sql/v3/dialect"

// StringTests returns test cases for string functions.
func StringTests() []ConvertTestCase {
	return []ConvertTestCase{
		{
			Name:     "starts_with",
			CELExpr:  `name.startsWith("a")`,
			Category: CategoryString,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "name LIKE 'a%' ESCAPE E'\\\\'",
				dialect.MySQL:      "name LIKE 'a%' ESCAPE '\\\\'",
				dialect.SQLite:     "name LIKE 'a%' ESCAPE '\\'",
				dialect.DuckDB:     "name LIKE 'a%' ESCAPE '\\\\'",
				dialect.BigQuery:   "name LIKE 'a%'",
			},
		},
		{
			Name:     "ends_with",
			CELExpr:  `name.endsWith("z")`,
			Category: CategoryString,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "name LIKE '%z' ESCAPE E'\\\\'",
				dialect.MySQL:      "name LIKE '%z' ESCAPE '\\\\'",
				dialect.SQLite:     "name LIKE '%z' ESCAPE '\\'",
				dialect.DuckDB:     "name LIKE '%z' ESCAPE '\\\\'",
				dialect.BigQuery:   "name LIKE '%z'",
			},
		},
		{
			Name:     "contains",
			CELExpr:  `name.contains("abc")`,
			Category: CategoryString,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "POSITION('abc' IN name) > 0",
				dialect.MySQL:      "LOCATE('abc', name) > 0",
				dialect.SQLite:     "INSTR(name, 'abc') > 0",
				dialect.DuckDB:     "CONTAINS(name, 'abc')",
				dialect.BigQuery:   "STRPOS(name, 'abc') > 0",
			},
		},
		{
			Name:     "size_string",
			CELExpr:  `size("test")`,
			Category: CategoryString,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "LENGTH('test')",
				dialect.MySQL:      "LENGTH('test')",
				dialect.SQLite:     "LENGTH('test')",
				dialect.DuckDB:     "LENGTH('test')",
				dialect.BigQuery:   "LENGTH('test')",
			},
		},
		{
			Name:     "starts_with_and_ends_with",
			CELExpr:  `name.startsWith("a") && name.endsWith("z")`,
			Category: CategoryString,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "name LIKE 'a%' ESCAPE E'\\\\' AND name LIKE '%z' ESCAPE E'\\\\'",
				dialect.MySQL:      "name LIKE 'a%' ESCAPE '\\\\' AND name LIKE '%z' ESCAPE '\\\\'",
				dialect.SQLite:     "name LIKE 'a%' ESCAPE '\\' AND name LIKE '%z' ESCAPE '\\'",
				dialect.DuckDB:     "name LIKE 'a%' ESCAPE '\\\\' AND name LIKE '%z' ESCAPE '\\\\'",
				dialect.BigQuery:   "name LIKE 'a%' AND name LIKE '%z'",
			},
		},
	}
}
