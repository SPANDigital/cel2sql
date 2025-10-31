package cel2sql

import (
	"log/slog"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/spandigital/cel2sql/v3/pg"
)

// TestComprehensionPatternDetectionOrder verifies that pattern matching
// detects comprehension types in the correct order and chooses the right type
func TestComprehensionPatternDetectionOrder(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "numbers", Type: "integer", Repeated: true},
		{Name: "items", Type: "text", Repeated: true},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"TestTable": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("data", cel.ObjectType("TestTable")),
	)
	if err != nil {
		t.Fatalf("failed to create CEL environment: %v", err)
	}

	tests := []struct {
		name     string
		cel      string
		expected string
		wantType ComprehensionType
	}{
		{
			name:     "map with identity transform looks different from filter",
			cel:      `data.numbers.map(x, x)`,
			expected: `ARRAY(SELECT x FROM UNNEST(data.numbers) AS x)`,
			wantType: ComprehensionMap,
		},
		{
			name:     "filter with simple predicate",
			cel:      `data.numbers.filter(x, x > 10)`,
			expected: `ARRAY(SELECT x FROM UNNEST(data.numbers) AS x WHERE x > 10)`,
			wantType: ComprehensionFilter,
		},
		{
			name:     "all with multiple AND conditions",
			cel:      `data.numbers.all(x, x > 0 && x < 100)`,
			expected: `NOT EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE NOT (x > 0 AND x < 100))`,
			wantType: ComprehensionAll,
		},
		{
			name:     "exists with multiple OR conditions",
			cel:      `data.numbers.exists(x, x < 0 || x > 100)`,
			expected: `EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE x < 0 OR x > 100)`,
			wantType: ComprehensionExists,
		},
		{
			name:     "exists_one with complex predicate",
			cel:      `data.numbers.exists_one(x, x > 50 && x < 60)`,
			expected: `(SELECT COUNT(*) FROM UNNEST(data.numbers) AS x WHERE x > 50 AND x < 60) = 1`,
			wantType: ComprehensionExistsOne,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.cel)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("failed to compile CEL: %v", issues.Err())
			}

			// Verify the comprehension type is identified correctly
			conv := &converter{
				schemas: map[string]pg.Schema{"TestTable": schema},
				logger:  slog.New(slog.DiscardHandler),
			}

			// Get the comprehension expression from the AST
			checkedExpr, _ := cel.AstToCheckedExpr(ast)
			compExpr := checkedExpr.GetExpr().GetComprehensionExpr()
			if compExpr == nil {
				t.Fatal("expected comprehension expression")
			}

			info, err := conv.analyzeComprehensionPattern(compExpr)
			if err != nil {
				t.Fatalf("failed to analyze comprehension: %v", err)
			}

			if info.Type != tt.wantType {
				t.Errorf("wrong comprehension type detected: got %v, want %v", info.Type, tt.wantType)
			}

			// Verify SQL generation
			sql, err := Convert(ast, WithSchemas(map[string]pg.Schema{"TestTable": schema}))
			if err != nil {
				t.Fatalf("failed to convert to SQL: %v", err)
			}

			if sql != tt.expected {
				t.Errorf("unexpected SQL:\ngot:  %s\nwant: %s", sql, tt.expected)
			}
		})
	}
}

// TestComprehensionWithComplexNestedExpressions tests comprehensions
// with deeply nested or complex expressions in predicates/transforms
func TestComprehensionWithComplexNestedExpressions(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "numbers", Type: "integer", Repeated: true},
		{Name: "values", Type: "integer", Repeated: true},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"TestTable": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("data", cel.ObjectType("TestTable")),
	)
	if err != nil {
		t.Fatalf("failed to create CEL environment: %v", err)
	}

	tests := []struct {
		name     string
		cel      string
		expected string
	}{
		{
			name:     "map with nested ternary transform",
			cel:      `data.numbers.map(x, x > 0 ? x * 2 : x)`,
			expected: `ARRAY(SELECT CASE WHEN x > 0 THEN x * 2 ELSE x END FROM UNNEST(data.numbers) AS x)`,
		},
		{
			name:     "filter with nested logical expression",
			cel:      `data.numbers.filter(x, (x > 10 && x < 20) || (x > 30 && x < 40))`,
			expected: `ARRAY(SELECT x FROM UNNEST(data.numbers) AS x WHERE x > 10 AND x < 20 OR x > 30 AND x < 40)`,
		},
		{
			name:     "all with nested parentheses",
			cel:      `data.numbers.all(x, ((x > 0) && (x < 100)) || (x == 0))`,
			expected: `NOT EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE NOT (x > 0 AND x < 100 OR x = 0))`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.cel)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("failed to compile CEL: %v", issues.Err())
			}

			sql, err := Convert(ast, WithSchemas(map[string]pg.Schema{"TestTable": schema}))
			if err != nil {
				t.Fatalf("failed to convert to SQL: %v", err)
			}

			if sql != tt.expected {
				t.Errorf("unexpected SQL:\ngot:  %s\nwant: %s", sql, tt.expected)
			}
		})
	}
}

// TestComprehensionEdgeCasesWithEmptyLists tests comprehension patterns
// with empty lists and edge case values
func TestComprehensionEdgeCasesWithEmptyLists(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "numbers", Type: "integer", Repeated: true},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"TestTable": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("data", cel.ObjectType("TestTable")),
	)
	if err != nil {
		t.Fatalf("failed to create CEL environment: %v", err)
	}

	tests := []struct {
		name     string
		cel      string
		expected string
	}{
		{
			name:     "all on potentially empty list",
			cel:      `data.numbers.all(x, x > 0)`,
			expected: `NOT EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE NOT (x > 0))`,
		},
		{
			name:     "exists on potentially empty list",
			cel:      `data.numbers.exists(x, x == 42)`,
			expected: `EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE x = 42)`,
		},
		{
			name:     "filter on potentially empty list",
			cel:      `data.numbers.filter(x, x != 0)`,
			expected: `ARRAY(SELECT x FROM UNNEST(data.numbers) AS x WHERE x != 0)`,
		},
		{
			name:     "map on potentially empty list",
			cel:      `data.numbers.map(x, x * 2)`,
			expected: `ARRAY(SELECT x * 2 FROM UNNEST(data.numbers) AS x)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.cel)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("failed to compile CEL: %v", issues.Err())
			}

			sql, err := Convert(ast, WithSchemas(map[string]pg.Schema{"TestTable": schema}))
			if err != nil {
				t.Fatalf("failed to convert to SQL: %v", err)
			}

			if sql != tt.expected {
				t.Errorf("unexpected SQL:\ngot:  %s\nwant: %s", sql, tt.expected)
			}
		})
	}
}

// TestComprehensionWithChainedOperations tests comprehensions that are
// chained together or combined with other operations
func TestComprehensionWithChainedOperations(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "numbers", Type: "integer", Repeated: true},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"TestTable": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("data", cel.ObjectType("TestTable")),
	)
	if err != nil {
		t.Fatalf("failed to create CEL environment: %v", err)
	}

	tests := []struct {
		name     string
		cel      string
		expected string
	}{
		{
			name:     "filter then map (chained)",
			cel:      `data.numbers.filter(x, x > 0).map(y, y * 2)`,
			expected: `ARRAY(SELECT y * 2 FROM UNNEST(ARRAY(SELECT x FROM UNNEST(data.numbers) AS x WHERE x > 0)) AS y)`,
		},
		{
			name:     "map then filter (chained)",
			cel:      `data.numbers.map(x, x * 2).filter(y, y > 10)`,
			expected: `ARRAY(SELECT y FROM UNNEST(ARRAY(SELECT x * 2 FROM UNNEST(data.numbers) AS x)) AS y WHERE y > 10)`,
		},
		{
			name:     "exists with negation",
			cel:      `!data.numbers.exists(x, x < 0)`,
			expected: `NOT EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE x < 0)`,
		},
		{
			name:     "all with negation",
			cel:      `!data.numbers.all(x, x > 0)`,
			expected: `NOT NOT EXISTS (SELECT 1 FROM UNNEST(data.numbers) AS x WHERE NOT (x > 0))`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.cel)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("failed to compile CEL: %v", issues.Err())
			}

			sql, err := Convert(ast, WithSchemas(map[string]pg.Schema{"TestTable": schema}))
			if err != nil {
				t.Fatalf("failed to convert to SQL: %v", err)
			}

			if sql != tt.expected {
				t.Errorf("unexpected SQL:\ngot:  %s\nwant: %s", sql, tt.expected)
			}
		})
	}
}

// TestComprehensionWithVariableNameEdgeCases tests comprehensions with
// unusual but valid variable names
func TestComprehensionWithVariableNameEdgeCases(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "items", Type: "integer", Repeated: true},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"TestTable": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("data", cel.ObjectType("TestTable")),
	)
	if err != nil {
		t.Fatalf("failed to create CEL environment: %v", err)
	}

	tests := []struct {
		name     string
		cel      string
		expected string
	}{
		{
			name:     "single letter variable name",
			cel:      `data.items.all(i, i > 0)`,
			expected: `NOT EXISTS (SELECT 1 FROM UNNEST(data.items) AS i WHERE NOT (i > 0))`,
		},
		{
			name:     "longer variable name",
			cel:      `data.items.filter(item, item > 5)`,
			expected: `ARRAY(SELECT item FROM UNNEST(data.items) AS item WHERE item > 5)`,
		},
		{
			name:     "underscore in variable name",
			cel:      `data.items.map(item_val, item_val * 2)`,
			expected: `ARRAY(SELECT item_val * 2 FROM UNNEST(data.items) AS item_val)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.cel)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("failed to compile CEL: %v", issues.Err())
			}

			sql, err := Convert(ast, WithSchemas(map[string]pg.Schema{"TestTable": schema}))
			if err != nil {
				t.Fatalf("failed to convert to SQL: %v", err)
			}

			if sql != tt.expected {
				t.Errorf("unexpected SQL:\ngot:  %s\nwant: %s", sql, tt.expected)
			}
		})
	}
}

// TestComprehensionWithMapFilter tests the map comprehension with filter
// to ensure it's distinguished from regular filter and map
func TestComprehensionWithMapFilter(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "numbers", Type: "integer", Repeated: true},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"TestTable": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("data", cel.ObjectType("TestTable")),
	)
	if err != nil {
		t.Fatalf("failed to create CEL environment: %v", err)
	}

	tests := []struct {
		name     string
		cel      string
		expected string
		wantType ComprehensionType
	}{
		{
			name:     "regular map without filter",
			cel:      `data.numbers.map(x, x * 2)`,
			expected: `ARRAY(SELECT x * 2 FROM UNNEST(data.numbers) AS x)`,
			wantType: ComprehensionMap,
		},
		{
			name:     "regular filter",
			cel:      `data.numbers.filter(x, x > 0)`,
			expected: `ARRAY(SELECT x FROM UNNEST(data.numbers) AS x WHERE x > 0)`,
			wantType: ComprehensionFilter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.cel)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("failed to compile CEL: %v", issues.Err())
			}

			// Verify the comprehension type
			conv := &converter{
				schemas: map[string]pg.Schema{"TestTable": schema},
				logger:  slog.New(slog.DiscardHandler),
			}

			checkedExpr, _ := cel.AstToCheckedExpr(ast)
			compExpr := checkedExpr.GetExpr().GetComprehensionExpr()
			if compExpr == nil {
				t.Fatal("expected comprehension expression")
			}

			info, err := conv.analyzeComprehensionPattern(compExpr)
			if err != nil {
				t.Fatalf("failed to analyze comprehension: %v", err)
			}

			if info.Type != tt.wantType {
				t.Errorf("wrong comprehension type: got %v, want %v", info.Type, tt.wantType)
			}

			// Verify SQL generation
			sql, err := Convert(ast, WithSchemas(map[string]pg.Schema{"TestTable": schema}))
			if err != nil {
				t.Fatalf("failed to convert to SQL: %v", err)
			}

			if sql != tt.expected {
				t.Errorf("unexpected SQL:\ngot:  %s\nwant: %s", sql, tt.expected)
			}
		})
	}
}
