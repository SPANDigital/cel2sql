package testcases

import "github.com/spandigital/cel2sql/v3/dialect"

// JSONTests returns test cases for JSON/JSONB field access and operations.
// These tests require the "json_schema" environment setup.
func JSONTests() []ConvertTestCase {
	return []ConvertTestCase{
		{
			Name:     "json_field_access",
			CELExpr:  `product.metadata.brand == "Acme"`,
			Category: CategoryJSON,
			EnvSetup: EnvWithJSON,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "product.metadata->>'brand' = 'Acme'",
				dialect.MySQL:      "product.metadata->>'$.brand' = 'Acme'",
				dialect.SQLite:     "json_extract(product.metadata, '$.brand') = 'Acme'",
				dialect.DuckDB:     "product.metadata->>'brand' = 'Acme'",
				dialect.BigQuery:   "JSON_VALUE(product.metadata, '$.brand') = 'Acme'",
			},
		},
		{
			Name:     "json_nested_access",
			CELExpr:  `product.metadata.specs.color == "red"`,
			Category: CategoryJSON,
			EnvSetup: EnvWithJSON,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "product.metadata->'specs'->>'color' = 'red'",
				dialect.MySQL:      "product.metadata->'$.specs'->>'$.color' = 'red'",
				dialect.SQLite:     "json_extract(json_extract(product.metadata, '$.specs'), '$.color') = 'red'",
				dialect.DuckDB:     "product.metadata->'specs'->>'color' = 'red'",
				dialect.BigQuery:   "JSON_VALUE(JSON_QUERY(product.metadata, '$.specs'), '$.color') = 'red'",
			},
		},
		{
			Name:     "json_has_field",
			CELExpr:  `has(product.metadata.brand)`,
			Category: CategoryJSON,
			EnvSetup: EnvWithJSON,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "product.metadata ? 'brand'",
				dialect.MySQL:      "JSON_CONTAINS_PATH(product.metadata, 'one', '$.brand')",
				dialect.SQLite:     "json_type(product.metadata, '$.brand') IS NOT NULL",
				dialect.DuckDB:     "json_exists(product.metadata, '$.brand')",
				dialect.BigQuery:   "JSON_VALUE(product.metadata, '$.brand') IS NOT NULL",
			},
		},
	}
}
