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
				dialect.Spark:      "get_json_object(product.metadata, '$.brand') = 'Acme'",
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
				dialect.Spark:      "get_json_object(get_json_object(product.metadata, '$.specs'), '$.color') = 'red'",
			},
		},
		{
			Name:     "json_array_membership",
			CELExpr:  `"electronics" in product.tags`,
			Category: CategoryJSON,
			EnvSetup: EnvWithJSON,
			WantSQL: map[dialect.Name]string{
				dialect.PostgreSQL: "'electronics' = ANY(ARRAY(SELECT jsonb_array_elements_text(product.tags)))",
				dialect.MySQL:      "JSON_OVERLAPS(JSON_ARRAY('electronics'), product.tags)",
				dialect.SQLite:     "EXISTS (SELECT 1 FROM json_each(product.tags) WHERE value = 'electronics')",
				dialect.DuckDB:     "EXISTS (SELECT 1 FROM json_each(product.tags) WHERE value = 'electronics')",
				dialect.BigQuery:   "'electronics' IN UNNEST(JSON_VALUE_ARRAY(product.tags))",
				dialect.Spark:      "array_contains(from_json(product.tags, 'ARRAY<STRING>'), 'electronics')",
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
				dialect.Spark:      "get_json_object(product.metadata, '$.brand') IS NOT NULL",
			},
		},
	}
}
