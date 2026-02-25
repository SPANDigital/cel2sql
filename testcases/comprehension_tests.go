package testcases

import "github.com/spandigital/cel2sql/v3/dialect"

// ComprehensionTests returns test cases for CEL comprehension operations.
func ComprehensionTests() []ConvertTestCase {
	return []ConvertTestCase{
		{
			Name:     "all",
			CELExpr:  `string_list.all(x, x != "bad")`,
			Category: CategoryComprehension,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "NOT EXISTS (SELECT 1 FROM UNNEST(string_list) AS x WHERE NOT (x != 'bad'))",
				dialect.SQLite:     "NOT EXISTS (SELECT 1 FROM json_each(string_list) AS x WHERE NOT (x != 'bad'))",
				dialect.DuckDB:     "NOT EXISTS (SELECT 1 FROM UNNEST(string_list) AS x WHERE NOT (x != 'bad'))",
				dialect.BigQuery:   "NOT EXISTS (SELECT 1 FROM UNNEST(string_list) AS x WHERE NOT (x != 'bad'))",
			},
		},
		{
			Name:     "exists",
			CELExpr:  `string_list.exists(x, x == "good")`,
			Category: CategoryComprehension,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "EXISTS (SELECT 1 FROM UNNEST(string_list) AS x WHERE x = 'good')",
				dialect.SQLite:     "EXISTS (SELECT 1 FROM json_each(string_list) AS x WHERE x = 'good')",
				dialect.DuckDB:     "EXISTS (SELECT 1 FROM UNNEST(string_list) AS x WHERE x = 'good')",
				dialect.BigQuery:   "EXISTS (SELECT 1 FROM UNNEST(string_list) AS x WHERE x = 'good')",
			},
		},
		{
			Name:     "exists_one",
			CELExpr:  `string_list.exists_one(x, x == "unique")`,
			Category: CategoryComprehension,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "(SELECT COUNT(*) FROM UNNEST(string_list) AS x WHERE x = 'unique') = 1",
				dialect.SQLite:     "(SELECT COUNT(*) FROM json_each(string_list) AS x WHERE x = 'unique') = 1",
				dialect.DuckDB:     "(SELECT COUNT(*) FROM UNNEST(string_list) AS x WHERE x = 'unique') = 1",
				dialect.BigQuery:   "(SELECT COUNT(*) FROM UNNEST(string_list) AS x WHERE x = 'unique') = 1",
			},
		},
		{
			Name:     "filter",
			CELExpr:  `string_list.filter(x, x != "bad")`,
			Category: CategoryComprehension,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "ARRAY(SELECT x FROM UNNEST(string_list) AS x WHERE x != 'bad')",
				dialect.SQLite:     "(SELECT json_group_array(x) FROM json_each(string_list) AS x WHERE x != 'bad')",
				dialect.DuckDB:     "ARRAY(SELECT x FROM UNNEST(string_list) AS x WHERE x != 'bad')",
				dialect.BigQuery:   "ARRAY(SELECT x FROM UNNEST(string_list) AS x WHERE x != 'bad')",
			},
		},
		{
			Name:     "map_transform",
			CELExpr:  `string_list.map(x, x + "_suffix")`,
			Category: CategoryComprehension,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "ARRAY(SELECT x || '_suffix' FROM UNNEST(string_list) AS x)",
				dialect.SQLite:     "(SELECT json_group_array(x || '_suffix') FROM json_each(string_list) AS x)",
				dialect.DuckDB:     "ARRAY(SELECT x || '_suffix' FROM UNNEST(string_list) AS x)",
				dialect.BigQuery:   "ARRAY(SELECT x || '_suffix' FROM UNNEST(string_list) AS x)",
			},
		},
	}
}
