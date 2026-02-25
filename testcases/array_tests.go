package testcases

import "github.com/spandigital/cel2sql/v3/dialect"

// ArrayTests returns test cases for array operations.
func ArrayTests() []ConvertTestCase {
	return []ConvertTestCase{
		{
			Name:     "list_index_literal",
			CELExpr:  `[1, 2, 3][0] == 1`,
			Category: CategoryArray,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "ARRAY[1, 2, 3][1] = 1",
				dialect.DuckDB:     "[1, 2, 3][1] = 1",
				dialect.BigQuery:   "[1, 2, 3][OFFSET(0)] = 1",
			},
		},
		{
			Name:     "list_var_index",
			CELExpr:  `string_list[0] == "a"`,
			Category: CategoryArray,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "string_list[1] = 'a'",
				dialect.DuckDB:     "string_list[1] = 'a'",
				dialect.BigQuery:   "string_list[OFFSET(0)] = 'a'",
			},
		},
		{
			Name:     "size_list",
			CELExpr:  `size(string_list)`,
			Category: CategoryArray,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "COALESCE(ARRAY_LENGTH(string_list, 1), 0)",
				dialect.DuckDB:     "COALESCE(array_length(string_list), 0)",
				dialect.BigQuery:   "ARRAY_LENGTH(string_list)",
			},
		},
		{
			Name:     "size_list_comparison",
			CELExpr:  `size(string_list) > 0`,
			Category: CategoryArray,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "COALESCE(ARRAY_LENGTH(string_list, 1), 0) > 0",
				dialect.DuckDB:     "COALESCE(array_length(string_list), 0) > 0",
				dialect.BigQuery:   "ARRAY_LENGTH(string_list) > 0",
			},
		},
		{
			Name:     "array_index_overflow",
			CELExpr:  `string_list[9223372036854775807]`,
			Category: CategoryArray,
			WantErr: map[dialect.Name]bool{
				dialect.PostgreSQL: true,
				dialect.DuckDB:     true,
				dialect.BigQuery:   true,
			},
		},
		{
			Name:     "array_index_negative",
			CELExpr:  `string_list[-1]`,
			Category: CategoryArray,
			WantErr: map[dialect.Name]bool{
				dialect.PostgreSQL: true,
				dialect.DuckDB:     true,
				dialect.BigQuery:   true,
			},
		},
	}
}
