package testcases

import "github.com/spandigital/cel2sql/v3/dialect"

// BasicTests returns test cases for basic comparisons and expressions.
func BasicTests() []ConvertTestCase {
	return []ConvertTestCase{
		{
			Name:     "equality_string",
			CELExpr:  `name == "a"`,
			Category: CategoryBasic,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "name = 'a'",
				dialect.MySQL:      "name = 'a'",
				dialect.SQLite:     "name = 'a'",
				dialect.DuckDB:     "name = 'a'",
				dialect.BigQuery:   "name = 'a'",
			},
		},
		{
			Name:     "inequality_int",
			CELExpr:  `age != 20`,
			Category: CategoryBasic,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "age != 20",
				dialect.MySQL:      "age != 20",
				dialect.SQLite:     "age != 20",
				dialect.DuckDB:     "age != 20",
				dialect.BigQuery:   "age != 20",
			},
		},
		{
			Name:     "less_than",
			CELExpr:  `age < 20`,
			Category: CategoryBasic,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "age < 20",
				dialect.MySQL:      "age < 20",
				dialect.SQLite:     "age < 20",
				dialect.DuckDB:     "age < 20",
				dialect.BigQuery:   "age < 20",
			},
		},
		{
			Name:     "greater_equal_float",
			CELExpr:  `height >= 1.6180339887`,
			Category: CategoryBasic,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "height >= 1.6180339887",
				dialect.MySQL:      "height >= 1.6180339887",
				dialect.SQLite:     "height >= 1.6180339887",
				dialect.DuckDB:     "height >= 1.6180339887",
				dialect.BigQuery:   "height >= 1.6180339887",
			},
		},
		{
			Name:     "is_null",
			CELExpr:  `null_var == null`,
			Category: CategoryBasic,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "null_var IS NULL",
				dialect.MySQL:      "null_var IS NULL",
				dialect.SQLite:     "null_var IS NULL",
				dialect.DuckDB:     "null_var IS NULL",
				dialect.BigQuery:   "null_var IS NULL",
			},
		},
		{
			Name:     "is_not_true",
			CELExpr:  `adult != true`,
			Category: CategoryBasic,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "adult IS NOT TRUE",
				dialect.MySQL:      "adult IS NOT TRUE",
				dialect.SQLite:     "adult IS NOT TRUE",
				dialect.DuckDB:     "adult IS NOT TRUE",
				dialect.BigQuery:   "adult IS NOT TRUE",
			},
		},
		{
			Name:     "not",
			CELExpr:  `!adult`,
			Category: CategoryBasic,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "NOT adult",
				dialect.MySQL:      "NOT adult",
				dialect.SQLite:     "NOT adult",
				dialect.DuckDB:     "NOT adult",
				dialect.BigQuery:   "NOT adult",
			},
		},
		{
			Name:     "negative_int",
			CELExpr:  `-1`,
			Category: CategoryBasic,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "-1",
				dialect.MySQL:      "-1",
				dialect.SQLite:     "-1",
				dialect.DuckDB:     "-1",
				dialect.BigQuery:   "-1",
			},
		},
		{
			Name:     "ternary",
			CELExpr:  `name == "a" ? "a" : "b"`,
			Category: CategoryBasic,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "CASE WHEN name = 'a' THEN 'a' ELSE 'b' END",
				dialect.MySQL:      "CASE WHEN name = 'a' THEN 'a' ELSE 'b' END",
				dialect.SQLite:     "CASE WHEN name = 'a' THEN 'a' ELSE 'b' END",
				dialect.DuckDB:     "CASE WHEN name = 'a' THEN 'a' ELSE 'b' END",
				dialect.BigQuery:   "CASE WHEN name = 'a' THEN 'a' ELSE 'b' END",
			},
		},
		{
			Name:     "field_select",
			CELExpr:  `page.title == "test"`,
			Category: CategoryFieldAccess,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "page.title = 'test'",
				dialect.MySQL:      "page.title = 'test'",
				dialect.SQLite:     "page.title = 'test'",
				dialect.DuckDB:     "page.title = 'test'",
				dialect.BigQuery:   "page.title = 'test'",
			},
		},
		{
			Name:     "in_list",
			CELExpr:  `name in ["a", "b", "c"]`,
			Category: CategoryBasic,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "name = ANY(ARRAY['a', 'b', 'c'])",
				dialect.MySQL:      "JSON_CONTAINS(JSON_ARRAY('a', 'b', 'c'), CAST(name AS JSON))",
				dialect.SQLite:     "name IN (SELECT value FROM json_each(json_array('a', 'b', 'c')))",
				dialect.DuckDB:     "name = ANY(['a', 'b', 'c'])",
				dialect.BigQuery:   "name IN UNNEST(['a', 'b', 'c'])",
			},
		},
	}
}
