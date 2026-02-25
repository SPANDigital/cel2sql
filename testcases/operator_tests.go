package testcases

import "github.com/spandigital/cel2sql/v3/dialect"

// OperatorTests returns test cases for logical and arithmetic operators.
func OperatorTests() []ConvertTestCase {
	return []ConvertTestCase{
		{
			Name:     "logical_and",
			CELExpr:  `name == "a" && age > 20`,
			Category: CategoryOperator,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "name = 'a' AND age > 20",
				dialect.MySQL:      "name = 'a' AND age > 20",
				dialect.SQLite:     "name = 'a' AND age > 20",
				dialect.DuckDB:     "name = 'a' AND age > 20",
				dialect.BigQuery:   "name = 'a' AND age > 20",
			},
		},
		{
			Name:     "logical_or",
			CELExpr:  `name == "a" || age > 20`,
			Category: CategoryOperator,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "name = 'a' OR age > 20",
				dialect.MySQL:      "name = 'a' OR age > 20",
				dialect.SQLite:     "name = 'a' OR age > 20",
				dialect.DuckDB:     "name = 'a' OR age > 20",
				dialect.BigQuery:   "name = 'a' OR age > 20",
			},
		},
		{
			Name:     "parenthesized",
			CELExpr:  `age >= 10 && (name == "a" || name == "b")`,
			Category: CategoryOperator,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "age >= 10 AND (name = 'a' OR name = 'b')",
				dialect.MySQL:      "age >= 10 AND (name = 'a' OR name = 'b')",
				dialect.SQLite:     "age >= 10 AND (name = 'a' OR name = 'b')",
				dialect.DuckDB:     "age >= 10 AND (name = 'a' OR name = 'b')",
				dialect.BigQuery:   "age >= 10 AND (name = 'a' OR name = 'b')",
			},
		},
		{
			Name:     "addition",
			CELExpr:  `1 + 2 == 3`,
			Category: CategoryOperator,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "1 + 2 = 3",
				dialect.MySQL:      "1 + 2 = 3",
				dialect.SQLite:     "1 + 2 = 3",
				dialect.DuckDB:     "1 + 2 = 3",
				dialect.BigQuery:   "1 + 2 = 3",
			},
		},
		{
			Name:     "modulo",
			CELExpr:  `5 % 3 == 2`,
			Category: CategoryOperator,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "MOD(5, 3) = 2",
				dialect.MySQL:      "MOD(5, 3) = 2",
				dialect.SQLite:     "MOD(5, 3) = 2",
				dialect.DuckDB:     "MOD(5, 3) = 2",
				dialect.BigQuery:   "MOD(5, 3) = 2",
			},
		},
		{
			Name:     "string_concat",
			CELExpr:  `"a" + "b" == "ab"`,
			Category: CategoryOperator,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "'a' || 'b' = 'ab'",
				dialect.MySQL:      "CONCAT('a', 'b') = 'ab'",
				dialect.SQLite:     "'a' || 'b' = 'ab'",
				dialect.DuckDB:     "'a' || 'b' = 'ab'",
				dialect.BigQuery:   "'a' || 'b' = 'ab'",
			},
		},
		{
			Name:     "list_concat_in",
			CELExpr:  `1 in [1] + [2, 3]`,
			Category: CategoryOperator,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "1 = ANY(ARRAY[1] || ARRAY[2, 3])",
				dialect.DuckDB:     "1 = ANY([1] || [2, 3])",
				dialect.BigQuery:   "1 IN UNNEST([1] || [2, 3])",
			},
		},
	}
}
