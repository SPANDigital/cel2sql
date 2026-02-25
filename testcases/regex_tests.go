package testcases

import "github.com/spandigital/cel2sql/v3/dialect"

// RegexTests returns test cases for regex pattern matching.
func RegexTests() []ConvertTestCase {
	return []ConvertTestCase{
		{
			Name:     "simple_match",
			CELExpr:  `name.matches("a+")`,
			Category: CategoryRegex,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "name ~ 'a+'",
				dialect.MySQL:      "name REGEXP 'a+'",
				dialect.DuckDB:     "name ~ 'a+'",
				dialect.BigQuery:   "REGEXP_CONTAINS(name, 'a+')",
			},
			SkipDialect: map[dialect.Name]string{
				dialect.SQLite: "SQLite does not support regex",
			},
		},
		{
			Name:     "function_style",
			CELExpr:  `matches(name, "^[0-9]+$")`,
			Category: CategoryRegex,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "name ~ '^[0-9]+$'",
				dialect.MySQL:      "name REGEXP '^[0-9]+$'",
				dialect.DuckDB:     "name ~ '^[0-9]+$'",
				dialect.BigQuery:   "REGEXP_CONTAINS(name, '^[0-9]+$')",
			},
			SkipDialect: map[dialect.Name]string{
				dialect.SQLite: "SQLite does not support regex",
			},
		},
		{
			Name:     "word_boundary",
			CELExpr:  `name.matches("\\btest\\b")`,
			Category: CategoryRegex,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "name ~ '\\ytest\\y'",
				dialect.MySQL:      "name REGEXP '\\btest\\b'",
				dialect.DuckDB:     "name ~ '\\btest\\b'",
				dialect.BigQuery:   "REGEXP_CONTAINS(name, '\\btest\\b')",
			},
			SkipDialect: map[dialect.Name]string{
				dialect.SQLite: "SQLite does not support regex",
			},
		},
		{
			Name:     "digit_class",
			CELExpr:  `name.matches("\\d{3}-\\d{4}")`,
			Category: CategoryRegex,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "name ~ '[[:digit:]]{3}-[[:digit:]]{4}'",
				dialect.MySQL:      "name REGEXP '\\d{3}-\\d{4}'",
				dialect.DuckDB:     "name ~ '\\d{3}-\\d{4}'",
				dialect.BigQuery:   "REGEXP_CONTAINS(name, '\\d{3}-\\d{4}')",
			},
			SkipDialect: map[dialect.Name]string{
				dialect.SQLite: "SQLite does not support regex",
			},
		},
		{
			Name:     "word_class",
			CELExpr:  `name.matches("\\w+@\\w+\\.\\w+")`,
			Category: CategoryRegex,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "name ~ '[[:alnum:]_]+@[[:alnum:]_]+\\.[[:alnum:]_]+'",
				dialect.MySQL:      "name REGEXP '\\w+@\\w+\\.\\w+'",
				dialect.DuckDB:     "name ~ '\\w+@\\w+\\.\\w+'",
				dialect.BigQuery:   "REGEXP_CONTAINS(name, '\\w+@\\w+\\.\\w+')",
			},
			SkipDialect: map[dialect.Name]string{
				dialect.SQLite: "SQLite does not support regex",
			},
		},
		{
			Name:     "complex_pattern",
			CELExpr:  `name.matches(".*pattern.*")`,
			Category: CategoryRegex,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "name ~ '.*pattern.*'",
				dialect.MySQL:      "name REGEXP '.*pattern.*'",
				dialect.DuckDB:     "name ~ '.*pattern.*'",
				dialect.BigQuery:   "REGEXP_CONTAINS(name, '.*pattern.*')",
			},
			SkipDialect: map[dialect.Name]string{
				dialect.SQLite: "SQLite does not support regex",
			},
		},
	}
}
