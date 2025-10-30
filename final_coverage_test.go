package cel2sql_test

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
)

// TestJSONFieldContains tests callContains with JSON fields (POSITION operator path)
func TestJSONFieldContains(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "tags", Type: "jsonb", IsJSON: true, IsJSONB: true},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"doc": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("doc", cel.ObjectType("doc")),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		expectedSQL string
		description string
	}{
		{
			name:        "json_field_contains_string",
			expression:  `doc.tags.contains("important")`,
			expectedSQL: `POSITION('important' IN doc.tags) > 0`,
			description: "Contains on JSONB field (uses POSITION)",
		},
		{
			name:        "json_nested_field_contains",
			expression:  `doc.metadata.items.contains("value")`,
			expectedSQL: `POSITION('value' IN doc.metadata->>'items') > 0`,
			description: "Contains on nested JSONB field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			schemas := provider.GetSchemas()
			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err, "Conversion should succeed for %s", tt.description)
			assert.Equal(t, tt.expectedSQL, sql, "SQL should match for %s", tt.description)
			t.Logf("Generated SQL: %s", sql)
		})
	}
}

// TestNestedJSONHasOperations tests visitNestedJSONHas and visitJSONColumnReference
func TestNestedJSONHasOperations(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
		{Name: "properties", Type: "jsonb", IsJSON: true, IsJSONB: true},
		{Name: "content", Type: "json", IsJSON: true, IsJSONB: false},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"record": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("record", cel.ObjectType("record")),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		description string
		checkSQL    func(t *testing.T, sql string)
	}{
		{
			name:        "nested_json_has_two_levels",
			expression:  `has(record.metadata.user.name)`,
			description: "has() with two-level nested JSON path",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "jsonb_extract_path_text")
				assert.Contains(t, sql, "record.metadata")
				assert.Contains(t, sql, "'user'")
				assert.Contains(t, sql, "'name'")
				assert.Contains(t, sql, "IS NOT NULL")
			},
		},
		{
			name:        "nested_json_has_three_levels",
			expression:  `has(record.metadata.settings.permissions.read)`,
			description: "has() with three-level nested JSON path",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "jsonb_extract_path_text")
				assert.Contains(t, sql, "record.metadata")
				assert.Contains(t, sql, "'settings'")
				assert.Contains(t, sql, "'permissions'")
				assert.Contains(t, sql, "'read'")
			},
		},
		{
			name:        "nested_json_has_properties",
			expression:  `has(record.properties.config.enabled)`,
			description: "has() on different JSON column",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "jsonb_extract_path_text")
				assert.Contains(t, sql, "record.properties")
				assert.Contains(t, sql, "'config'")
				assert.Contains(t, sql, "'enabled'")
			},
		},
		{
			name:        "nested_json_has_deep_path",
			expression:  `has(record.metadata.a.b.c.d.e)`,
			description: "has() with deeply nested path (5 levels)",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "jsonb_extract_path_text")
				assert.Contains(t, sql, "record.metadata")
				assert.Contains(t, sql, "'a'")
				assert.Contains(t, sql, "'e'")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			schemas := provider.GetSchemas()
			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err, "Conversion should succeed for %s", tt.description)

			tt.checkSQL(t, sql)
			t.Logf("Generated SQL for %s: %s", tt.description, sql)
		})
	}
}

// TestAdditionalFunctionTypes tests visitCallFunc with more function types
func TestAdditionalFunctionTypes(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "created", Type: "timestamp with time zone"},
		{Name: "modified", Type: "timestamp with time zone"},
		{Name: "data", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"record": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("record", cel.ObjectType("record")),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		description string
		checkSQL    func(t *testing.T, sql string)
	}{
		{
			name:        "extract_year_from_timestamp",
			expression:  `record.created.getFullYear() == 2024`,
			description: "getFullYear() on timestamp",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "EXTRACT")
				assert.Contains(t, sql, "YEAR")
				assert.Contains(t, sql, "record.created")
			},
		},
		{
			name:        "extract_month_from_timestamp",
			expression:  `record.created.getMonth() == 11`,
			description: "getMonth() on timestamp",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "EXTRACT")
				assert.Contains(t, sql, "MONTH")
			},
		},
		{
			name:        "extract_day_from_timestamp",
			expression:  `record.created.getDate() == 15`,
			description: "getDate() on timestamp",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "EXTRACT")
				assert.Contains(t, sql, "DAY")
			},
		},
		{
			name:        "extract_hours_from_timestamp",
			expression:  `record.created.getHours() >= 9`,
			description: "getHours() on timestamp",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "EXTRACT")
				assert.Contains(t, sql, "HOUR")
			},
		},
		{
			name:        "extract_minutes_from_timestamp",
			expression:  `record.created.getMinutes() < 30`,
			description: "getMinutes() on timestamp",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "EXTRACT")
				assert.Contains(t, sql, "MINUTE")
			},
		},
		{
			name:        "extract_seconds_from_timestamp",
			expression:  `record.created.getSeconds() > 0`,
			description: "getSeconds() on timestamp",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "EXTRACT")
				assert.Contains(t, sql, "SECOND")
			},
		},
		{
			name:        "extract_day_of_week",
			expression:  `record.created.getDayOfWeek() == 1`,
			description: "getDayOfWeek() on timestamp",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "EXTRACT")
				assert.Contains(t, sql, "DOW")
			},
		},
		{
			name:        "extract_day_of_year",
			expression:  `record.created.getDayOfYear() > 100`,
			description: "getDayOfYear() on timestamp",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "EXTRACT")
				assert.Contains(t, sql, "DOY")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			schemas := provider.GetSchemas()
			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err, "Conversion should succeed for %s", tt.description)

			tt.checkSQL(t, sql)
			t.Logf("Generated SQL for %s: %s", tt.description, sql)
		})
	}
}

// TestBinaryOperatorAdditionalCases tests more visitCallBinary scenarios
func TestBinaryOperatorAdditionalCases(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "score", Type: "double precision"},
		{Name: "data", Type: "bytea"},
		{Name: "tags", Type: "text", Repeated: true, ElementType: "text"},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"record": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("record", cel.ObjectType("record")),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		expectedSQL string
		description string
	}{
		{
			name:        "bytes_concatenation",
			expression:  `b"hello" + b" world" != b""`,
			expectedSQL: `'\x68656c6c6f' || '\x20776f726c64' != '\x'`,
			description: "Bytes concatenation with || operator",
		},
		{
			name:        "json_numeric_comparison_cast",
			expression:  `record.metadata.count > 100`,
			expectedSQL: `(record.metadata->>'count')::numeric > 100`,
			description: "JSON text extraction with numeric cast",
		},
		{
			name:        "json_nested_numeric_comparison",
			expression:  `record.metadata.stats.total < 1000`,
			expectedSQL: `(record.metadata->'stats'->>'total')::numeric < 1000`,
			description: "Nested JSON numeric comparison with cast",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			schemas := provider.GetSchemas()
			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err, "Conversion should succeed for %s", tt.description)
			assert.Equal(t, tt.expectedSQL, sql, "SQL should match for %s", tt.description)
			t.Logf("Generated SQL: %s", sql)
		})
	}
}

// TestIdentifierWithNumericCasting tests visitIdent with numeric casting scenarios
func TestIdentifierWithNumericCasting(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "items", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"doc": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("doc", cel.ObjectType("doc")),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		description string
		checkSQL    func(t *testing.T, sql string)
	}{
		{
			name:        "json_comprehension_with_score",
			expression:  `doc.items.all(score, score > 50)`,
			description: "Comprehension with 'score' iterator (may need numeric cast)",
			checkSQL: func(t *testing.T, sql string) {
				// Should have comprehension with score identifier
				assert.Contains(t, sql, "score")
				assert.Contains(t, sql, "> 50")
			},
		},
		{
			name:        "json_comprehension_with_value",
			expression:  `doc.items.exists(value, value >= 100)`,
			description: "Comprehension with 'value' iterator (may need numeric cast)",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "value")
				assert.Contains(t, sql, ">= 100")
			},
		},
		{
			name:        "json_comprehension_with_amount",
			expression:  `doc.items.filter(amount, amount > 0)`,
			description: "Comprehension with 'amount' iterator (may need numeric cast)",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "amount")
				assert.Contains(t, sql, "> 0")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			schemas := provider.GetSchemas()
			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err, "Conversion should succeed for %s", tt.description)

			tt.checkSQL(t, sql)
			t.Logf("Generated SQL for %s: %s", tt.description, sql)
		})
	}
}

// TestJSONColumnReferenceEdgeCases tests edge cases for JSON column handling
func TestJSONColumnReferenceEdgeCases(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
		{Name: "structure", Type: "jsonb", IsJSON: true, IsJSONB: true},
		{Name: "taxonomy", Type: "jsonb", IsJSON: true, IsJSONB: true},
		{Name: "analytics", Type: "jsonb", IsJSON: true, IsJSONB: true},
		{Name: "classification", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"asset": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("asset", cel.ObjectType("asset")),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		expression  string
		description string
		checkSQL    func(t *testing.T, sql string)
	}{
		{
			name:        "has_on_structure_column",
			expression:  `has(asset.structure.level.parent)`,
			description: "has() on 'structure' JSON column (known column name)",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "asset.structure")
				assert.Contains(t, sql, "'level'")
				assert.Contains(t, sql, "'parent'")
			},
		},
		{
			name:        "has_on_taxonomy_column",
			expression:  `has(asset.taxonomy.category.main)`,
			description: "has() on 'taxonomy' JSON column",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "asset.taxonomy")
				assert.Contains(t, sql, "'category'")
			},
		},
		{
			name:        "has_on_analytics_column",
			expression:  `has(asset.analytics.views.total)`,
			description: "has() on 'analytics' JSON column",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "asset.analytics")
				assert.Contains(t, sql, "'views'")
			},
		},
		{
			name:        "has_on_classification_column",
			expression:  `has(asset.classification.level.value)`,
			description: "has() on 'classification' JSON column",
			checkSQL: func(t *testing.T, sql string) {
				assert.Contains(t, sql, "asset.classification")
				assert.Contains(t, sql, "'level'")
				assert.Contains(t, sql, "'value'")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			schemas := provider.GetSchemas()
			sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err, "Conversion should succeed for %s", tt.description)

			tt.checkSQL(t, sql)
			t.Logf("Generated SQL for %s: %s", tt.description, sql)
		})
	}
}
