package cel2sql

import (
	"context"
	"strings"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/spandigital/cel2sql/v3/pg"
)

// TestIssue40_LikeEscapeClause tests that LIKE patterns include explicit ESCAPE clause
func TestIssue40_LikeEscapeClause(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"Person": schema})

	tests := []struct {
		name         string
		expression   string
		wantContains string
		description  string
	}{
		{
			name:         "startsWith with literal includes ESCAPE",
			expression:   `person.name.startsWith("test")`,
			wantContains: "ESCAPE E'\\\\'",
			description:  "startsWith should include explicit ESCAPE clause",
		},
		{
			name:         "endsWith with literal includes ESCAPE",
			expression:   `person.name.endsWith("test")`,
			wantContains: "ESCAPE E'\\\\'",
			description:  "endsWith should include explicit ESCAPE clause",
		},
		{
			name:         "startsWith with special chars includes ESCAPE",
			expression:   `person.name.startsWith("test%")`,
			wantContains: "ESCAPE E'\\\\'",
			description:  "Escaped patterns should still include ESCAPE clause",
		},
		{
			name:         "endsWith with special chars includes ESCAPE",
			expression:   `person.name.endsWith("test_")`,
			wantContains: "ESCAPE E'\\\\'",
			description:  "Escaped patterns should still include ESCAPE clause",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv(
				cel.CustomTypeProvider(provider),
				cel.Variable("person", cel.ObjectType("Person")),
			)
			if err != nil {
				t.Fatalf("Failed to create environment: %v", err)
			}

			ast, issues := env.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("Failed to compile expression: %v", issues.Err())
			}

			sql, err := Convert(ast, WithSchemas(map[string]pg.Schema{"Person": schema}))
			if err != nil {
				t.Fatalf("Failed to convert: %v", err)
			}

			if !strings.Contains(sql, tt.wantContains) {
				t.Errorf("%s\nExpression: %s\nGenerated SQL: %s\nShould contain: %s",
					tt.description, tt.expression, sql, tt.wantContains)
			}
		})
	}
}

// TestIssue43_VariableLikeEscaping tests that variables in LIKE patterns are escaped at runtime
func TestIssue43_VariableLikeEscaping(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"Person": schema})

	tests := []struct {
		name         string
		expression   string
		wantContains []string
		description  string
	}{
		{
			name:       "startsWith with variable uses REPLACE for backslash",
			expression: `person.name.startsWith(prefix)`,
			wantContains: []string{
				"REPLACE(REPLACE(REPLACE(prefix",
				"ESCAPE E'\\\\'",
			},
			description: "startsWith with variable should escape backslash at runtime",
		},
		{
			name:       "startsWith with variable uses REPLACE for percent",
			expression: `person.name.startsWith(prefix)`,
			wantContains: []string{
				"REPLACE(",
				", '%', '\\%')", // Escapes percent
				"ESCAPE E'\\\\'",
			},
			description: "startsWith with variable should escape percent at runtime",
		},
		{
			name:       "startsWith with variable uses REPLACE for underscore",
			expression: `person.name.startsWith(prefix)`,
			wantContains: []string{
				"REPLACE(",
				", '_', '\\_')", // Escapes underscore
				"ESCAPE E'\\\\'",
			},
			description: "startsWith with variable should escape underscore at runtime",
		},
		{
			name:       "endsWith with variable uses REPLACE",
			expression: `person.name.endsWith(suffix)`,
			wantContains: []string{
				"'%' || REPLACE(REPLACE(REPLACE(",
				"ESCAPE E'\\\\'",
			},
			description: "endsWith with variable should use REPLACE chain",
		},
		{
			name:       "startsWith with variable includes all escapes in order",
			expression: `person.name.startsWith(prefix)`,
			wantContains: []string{
				"REPLACE(REPLACE(REPLACE(prefix",
				") || '%' ESCAPE",
			},
			description: "REPLACE chain should escape backslash, then percent, then underscore",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv(
				cel.CustomTypeProvider(provider),
				cel.Variable("person", cel.ObjectType("Person")),
				cel.Variable("prefix", cel.StringType),
				cel.Variable("suffix", cel.StringType),
			)
			if err != nil {
				t.Fatalf("Failed to create environment: %v", err)
			}

			ast, issues := env.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("Failed to compile expression: %v", issues.Err())
			}

			sql, err := Convert(ast, WithSchemas(map[string]pg.Schema{"Person": schema}))
			if err != nil {
				t.Fatalf("Failed to convert: %v", err)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(sql, want) {
					t.Errorf("%s\nExpression: %s\nGenerated SQL: %s\nShould contain: %s",
						tt.description, tt.expression, sql, want)
				}
			}
		})
	}
}

// TestIssue48_JSONComprehensionValidation tests that JSON comprehensions no longer have type checks
func TestIssue48_JSONComprehensionValidation(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "tags", Type: "integer", Repeated: true}, // Regular array, not JSONB
		{Name: "metadata", Type: "jsonb"},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"Person": schema})

	tests := []struct {
		name            string
		expression      string
		wantNotContains []string
		wantContains    []string
		description     string
	}{
		{
			name:       "all() comprehension has no type check",
			expression: `person.tags.all(t, t > 0)`,
			wantNotContains: []string{
				"IS NOT NULL",
				"jsonb_typeof",
				"= 'array'",
			},
			wantContains: []string{
				"NOT EXISTS",
				"UNNEST",
			},
			description: "all() should not include NULL or type checks",
		},
		{
			name:       "exists() comprehension has no type check",
			expression: `person.tags.exists(t, t > 0)`,
			wantNotContains: []string{
				"IS NOT NULL",
				"jsonb_typeof",
				"= 'array'",
			},
			wantContains: []string{
				"EXISTS",
				"UNNEST",
			},
			description: "exists() should not include NULL or type checks",
		},
		{
			name:       "exists_one() comprehension has no type check",
			expression: `person.tags.exists_one(t, t > 0)`,
			wantNotContains: []string{
				"IS NOT NULL",
				"jsonb_typeof",
				"= 'array'",
			},
			wantContains: []string{
				"SELECT COUNT(*)",
				"UNNEST",
				") = 1",
			},
			description: "exists_one() should not include NULL or type checks",
		},
		{
			name:       "filter() comprehension has no type check",
			expression: `person.tags.filter(t, t > 0)`,
			wantNotContains: []string{
				"IS NOT NULL",
				"jsonb_typeof",
				"= 'array'",
			},
			wantContains: []string{
				"ARRAY(SELECT",
				"UNNEST",
			},
			description: "filter() should not include NULL or type checks",
		},
		{
			name:       "map() comprehension has no type check",
			expression: `person.tags.map(t, t * 2)`,
			wantNotContains: []string{
				"IS NOT NULL",
				"jsonb_typeof",
				"= 'array'",
			},
			wantContains: []string{
				"ARRAY(SELECT",
				"UNNEST",
			},
			description: "map() should not include NULL or type checks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv(
				cel.CustomTypeProvider(provider),
				cel.Variable("person", cel.ObjectType("Person")),
			)
			if err != nil {
				t.Fatalf("Failed to create environment: %v", err)
			}

			ast, issues := env.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("Failed to compile expression: %v", issues.Err())
			}

			sql, err := Convert(ast, WithSchemas(map[string]pg.Schema{"Person": schema}))
			if err != nil {
				t.Fatalf("Failed to convert: %v", err)
			}

			// Check that unwanted strings are NOT present
			for _, unwanted := range tt.wantNotContains {
				if strings.Contains(sql, unwanted) {
					t.Errorf("%s\nExpression: %s\nGenerated SQL: %s\nShould NOT contain: %s",
						tt.description, tt.expression, sql, unwanted)
				}
			}

			// Check that expected strings ARE present
			for _, want := range tt.wantContains {
				if !strings.Contains(sql, want) {
					t.Errorf("%s\nExpression: %s\nGenerated SQL: %s\nShould contain: %s",
						tt.description, tt.expression, sql, want)
				}
			}
		})
	}
}

// TestIssue48_ComprehensionWithoutPredicate tests comprehensions work correctly
func TestIssue48_ComprehensionWithoutPredicate(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "numbers", Type: "integer", Repeated: true},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"Data": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("data", cel.ObjectType("Data")),
	)
	if err != nil {
		t.Fatalf("Failed to create environment: %v", err)
	}

	// Test filter with true predicate
	ast, issues := env.Compile(`data.numbers.filter(n, true)`)
	if issues != nil && issues.Err() != nil {
		t.Fatalf("Failed to compile expression: %v", issues.Err())
	}

	sql, err := Convert(ast, WithSchemas(map[string]pg.Schema{"Data": schema}))
	if err != nil {
		t.Fatalf("Failed to convert: %v", err)
	}

	// Should generate filter with true predicate (PostgreSQL uses uppercase TRUE)
	if !strings.Contains(sql, "ARRAY(SELECT n FROM UNNEST(data.numbers) AS n WHERE TRUE)") {
		t.Errorf("Unexpected SQL: %s", sql)
	}
}

// TestPhase1Integration tests all three fixes work together in complex expressions
func TestPhase1Integration(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
		{Name: "tags", Type: "integer", Repeated: true},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"Person": schema})

	tests := []struct {
		name         string
		expression   string
		wantContains []string
		description  string
	}{
		{
			name:       "Combined LIKE escaping and comprehension",
			expression: `person.name.startsWith("test") && person.tags.all(t, t > 0)`,
			wantContains: []string{
				"ESCAPE E'\\\\'", // Issue #40
				"NOT EXISTS",     // Issue #48
			},
			description: "Both LIKE ESCAPE and comprehension without type checks",
		},
		{
			name:       "Variable LIKE with comprehension",
			expression: `person.name.startsWith(prefix) && person.tags.exists(t, t > 0)`,
			wantContains: []string{
				"REPLACE(REPLACE(REPLACE(", // Issue #43
				"ESCAPE E'\\\\'",           // Issue #40/#43
				"EXISTS",                   // Issue #48
			},
			description: "Variable escaping with comprehension",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv(
				cel.CustomTypeProvider(provider),
				cel.Variable("person", cel.ObjectType("Person")),
				cel.Variable("prefix", cel.StringType),
			)
			if err != nil {
				t.Fatalf("Failed to create environment: %v", err)
			}

			ast, issues := env.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("Failed to compile expression: %v", issues.Err())
			}

			sql, err := Convert(ast, WithSchemas(map[string]pg.Schema{"Person": schema}))
			if err != nil {
				t.Fatalf("Failed to convert: %v", err)
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(sql, want) {
					t.Errorf("%s\nExpression: %s\nGenerated SQL: %s\nShould contain: %s",
						tt.description, tt.expression, sql, want)
				}
			}
		})
	}
}

// TestPhase1WithContext tests all fixes work with context
func TestPhase1WithContext(t *testing.T) {
	ctx := context.Background()
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"Person": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("person", cel.ObjectType("Person")),
	)
	if err != nil {
		t.Fatalf("Failed to create environment: %v", err)
	}

	ast, issues := env.Compile(`person.name.startsWith("test")`)
	if issues != nil && issues.Err() != nil {
		t.Fatalf("Failed to compile expression: %v", issues.Err())
	}

	sql, err := Convert(ast,
		WithContext(ctx),
		WithSchemas(map[string]pg.Schema{"Person": schema}))
	if err != nil {
		t.Fatalf("Failed to convert: %v", err)
	}

	if !strings.Contains(sql, "ESCAPE E'\\\\'") {
		t.Errorf("ESCAPE clause missing in context-aware conversion: %s", sql)
	}
}
