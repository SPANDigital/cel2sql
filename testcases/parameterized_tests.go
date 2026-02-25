package testcases

import "github.com/spandigital/cel2sql/v3/dialect"

// ParameterizedTests returns test cases for parameterized SQL conversion.
func ParameterizedTests() []ParameterizedTestCase {
	return []ParameterizedTestCase{
		{
			Name:     "simple_string_equality",
			CELExpr:  `name == "John"`,
			Category: CategoryParameterized,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "name = $1",
				dialect.SQLite:     "name = ?",
				dialect.DuckDB:     "name = $1",
				dialect.BigQuery:   "name = @p1",
			},
			WantParams: map[dialect.Name][]any{
				dialect.PostgreSQL: {"John"},
				dialect.SQLite:     {"John"},
				dialect.DuckDB:     {"John"},
				dialect.BigQuery:   {"John"},
			},
		},
		{
			Name:     "multiple_string_params",
			CELExpr:  `name == "John" && name != "Jane"`,
			Category: CategoryParameterized,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "name = $1 AND name != $2",
				dialect.SQLite:     "name = ? AND name != ?",
				dialect.DuckDB:     "name = $1 AND name != $2",
				dialect.BigQuery:   "name = @p1 AND name != @p2",
			},
			WantParams: map[dialect.Name][]any{
				dialect.PostgreSQL: {"John", "Jane"},
				dialect.SQLite:     {"John", "Jane"},
				dialect.DuckDB:     {"John", "Jane"},
				dialect.BigQuery:   {"John", "Jane"},
			},
		},
		{
			Name:     "integer_equality",
			CELExpr:  `age == 18`,
			Category: CategoryParameterized,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "age = $1",
				dialect.SQLite:     "age = ?",
				dialect.DuckDB:     "age = $1",
				dialect.BigQuery:   "age = @p1",
			},
			WantParams: map[dialect.Name][]any{
				dialect.PostgreSQL: {int64(18)},
				dialect.SQLite:     {int64(18)},
				dialect.DuckDB:     {int64(18)},
				dialect.BigQuery:   {int64(18)},
			},
		},
		{
			Name:     "integer_range",
			CELExpr:  `age > 21 && age < 65`,
			Category: CategoryParameterized,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "age > $1 AND age < $2",
				dialect.SQLite:     "age > ? AND age < ?",
				dialect.DuckDB:     "age > $1 AND age < $2",
				dialect.BigQuery:   "age > @p1 AND age < @p2",
			},
			WantParams: map[dialect.Name][]any{
				dialect.PostgreSQL: {int64(21), int64(65)},
				dialect.SQLite:     {int64(21), int64(65)},
				dialect.DuckDB:     {int64(21), int64(65)},
				dialect.BigQuery:   {int64(21), int64(65)},
			},
		},
		{
			Name:     "double_equality",
			CELExpr:  `salary == 50000.50`,
			Category: CategoryParameterized,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "salary = $1",
				dialect.SQLite:     "salary = ?",
				dialect.DuckDB:     "salary = $1",
				dialect.BigQuery:   "salary = @p1",
			},
			WantParams: map[dialect.Name][]any{
				dialect.PostgreSQL: {50000.50},
				dialect.SQLite:     {50000.50},
				dialect.DuckDB:     {50000.50},
				dialect.BigQuery:   {50000.50},
			},
		},
		{
			Name:     "boolean_true_inline",
			CELExpr:  `active == true`,
			Category: CategoryParameterized,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "active IS TRUE",
				dialect.SQLite:     "active IS TRUE",
				dialect.DuckDB:     "active IS TRUE",
				dialect.BigQuery:   "active IS TRUE",
			},
			WantParams: map[dialect.Name][]any{
				dialect.PostgreSQL: {},
				dialect.SQLite:     {},
				dialect.DuckDB:     {},
				dialect.BigQuery:   {},
			},
		},
	}
}
