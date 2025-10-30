package cel2sql_test

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
	"github.com/stretchr/testify/require"
)

func TestNestedComprehensionDepthLimit(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		shouldErr   bool
		description string
	}{
		{
			name:        "single_level_all",
			expression:  `[1, 2, 3].all(x, x > 0)`,
			shouldErr:   false,
			description: "Single level comprehension (depth 1) should succeed",
		},
		{
			name:        "single_level_exists",
			expression:  `[1, 2, 3].exists(x, x == 2)`,
			shouldErr:   false,
			description: "Single level exists (depth 1) should succeed",
		},
		{
			name:        "single_level_filter",
			expression:  `[1, 2, 3].filter(x, x > 1)`,
			shouldErr:   false,
			description: "Single level filter (depth 1) should succeed",
		},
		{
			name:        "single_level_map",
			expression:  `[1, 2, 3].map(x, x * 2)`,
			shouldErr:   false,
			description: "Single level map (depth 1) should succeed",
		},
		{
			name:        "double_nested_exists_in_all",
			expression:  `[[1, 2], [3, 4]].all(list, list.exists(x, x > 0))`,
			shouldErr:   false,
			description: "Double nested comprehension (depth 2) should succeed",
		},
		{
			name:        "double_nested_filter_in_map",
			expression:  `[[1, 2, 3], [4, 5, 6]].map(list, list.filter(x, x > 2))`,
			shouldErr:   false,
			description: "Double nested map/filter (depth 2) should succeed",
		},
		{
			name:        "triple_nested_at_limit",
			expression:  `[[[1, 2]], [[3, 4]]].all(layer1, layer1.all(layer2, layer2.all(x, x > 0)))`,
			shouldErr:   false,
			description: "Triple nested comprehension (depth 3) should succeed at limit",
		},
		{
			name:        "triple_nested_mixed_types",
			expression:  `[[1, 2, 3], [4, 5]].map(list1, list1.filter(x, x > 1).exists(z, z > 2))`,
			shouldErr:   false,
			description: "Triple nested with different types (depth 3) should succeed at limit",
		},
		{
			name:        "quadruple_nested_exceeds_limit",
			expression:  `[[[[1, 2]]]].all(a, a.all(b, b.all(c, c.all(d, d > 0))))`,
			shouldErr:   true,
			description: "Quadruple nested comprehension (depth 4) should fail",
		},
		{
			name:        "complex_nested_exceeds_limit",
			expression:  `[[[1, 2]]].map(x, x.map(y, y.filter(z, z > 0).exists(w, w > 0)))`,
			shouldErr:   true,
			description: "Complex nested exceeding depth 3 should fail",
		},
		{
			name:        "deeply_nested_exists_chain",
			expression:  `[[[[[1, 2]]]]].exists(a, a.exists(b, b.exists(c, c.exists(d, d.all(e, e > 0)))))`,
			shouldErr:   true,
			description: "Deeply nested exists chain (depth 5) should fail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv()
			require.NoError(t, err)

			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			_, err = cel2sql.Convert(ast)

			if tt.shouldErr {
				require.Error(t, err, tt.description)
				require.Contains(t, err.Error(), "comprehension depth",
					"Error should mention comprehension depth")
				require.Contains(t, err.Error(), "exceeds limit",
					"Error should mention exceeding maximum")
			} else {
				require.NoError(t, err, tt.description)
			}
		})
	}
}

func TestComprehensionDepthWithSchemas(t *testing.T) {
	// Test with realistic schemas similar to the issue example
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "list1", Type: "integer", Repeated: true},
		{Name: "list2", Type: "integer", Repeated: true},
		{Name: "list3", Type: "integer", Repeated: true},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"data": schema})
	schemas := provider.GetSchemas()

	tests := []struct {
		name        string
		expression  string
		shouldErr   bool
		description string
	}{
		{
			name:        "issue_example_similar",
			expression:  `data.list1.map(x, data.list2.filter(y, y > x))`,
			shouldErr:   false,
			description: "Similar to issue #35 example (depth 2) should succeed",
		},
		{
			name:        "issue_example_triple_nested",
			expression:  `data.list1.map(x, data.list2.filter(y, y > x).map(z, z * 2))`,
			shouldErr:   false,
			description: "Issue #35 example pattern (depth 3) should succeed at limit",
		},
		{
			name:        "too_deep_version",
			expression:  `data.list1.map(x, data.list2.map(y, data.list3.filter(z, z > y).exists(w, w == x)))`,
			shouldErr:   true,
			description: "Adding fourth level should fail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv(
				cel.CustomTypeProvider(provider),
				cel.Variable("data", cel.ObjectType("data")),
			)
			require.NoError(t, err)

			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed")

			_, err = cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))

			if tt.shouldErr {
				require.Error(t, err, tt.description)
				require.Contains(t, err.Error(), "comprehension depth")
				require.Contains(t, err.Error(), "exceeds limit")
			} else {
				require.NoError(t, err, tt.description)
			}
		})
	}
}

func TestComprehensionDepthResetBetweenCalls(t *testing.T) {
	// Ensure depth counter resets between Convert() calls
	env, err := cel.NewEnv()
	require.NoError(t, err)

	// Triple nested comprehension at the limit
	expr := `[[[1, 2]]].all(a, a.all(b, b.all(c, c > 0)))`

	ast, issues := env.Compile(expr)
	require.NoError(t, issues.Err())

	// First conversion should succeed
	_, err = cel2sql.Convert(ast)
	require.NoError(t, err, "First call should succeed")

	// Second conversion should also succeed (depth was reset)
	_, err = cel2sql.Convert(ast)
	require.NoError(t, err, "Second call should succeed - depth should have reset")

	// Third conversion should also succeed
	_, err = cel2sql.Convert(ast)
	require.NoError(t, err, "Third call should succeed - depth should have reset")
}

func TestComprehensionDepthErrorMessage(t *testing.T) {
	// Verify error message format and content
	env, err := cel.NewEnv()
	require.NoError(t, err)

	expr := `[[[[1, 2]]]].all(a, a.all(b, b.all(c, c.all(d, d > 0))))`

	ast, issues := env.Compile(expr)
	require.NoError(t, issues.Err())

	_, err = cel2sql.Convert(ast)
	require.Error(t, err)
	require.Contains(t, err.Error(), "comprehension depth",
		"Error message should mention comprehension depth")
	require.Contains(t, err.Error(), "exceeds limit",
		"Error message should mention exceeding maximum")
}

func TestComprehensionDepthWithMixedExpressions(t *testing.T) {
	// Test that comprehension depth only tracks comprehensions, not other expressions
	tests := []struct {
		name        string
		expression  string
		shouldErr   bool
		description string
	}{
		{
			name:        "comprehension_with_deep_arithmetic",
			expression:  `[1, 2, 3].all(x, ((((x + 1) + 1) + 1) + 1) > 0)`,
			shouldErr:   false,
			description: "Deep arithmetic inside comprehension doesn't count toward comprehension depth",
		},
		{
			name:        "comprehension_with_deep_conditionals",
			expression:  `[1, 2, 3].exists(x, (x > 0 ? (x > 1 ? (x > 2 ? 1 : 0) : 0) : 0) == 1)`,
			shouldErr:   false,
			description: "Deep conditionals inside comprehension don't count toward comprehension depth",
		},
		{
			name:        "nested_comprehension_with_deep_expressions",
			expression:  `[[1, 2]].all(list, list.exists(x, ((x + 1) * 2) > 0))`,
			shouldErr:   false,
			description: "Nested comprehension with deep expressions (depth 2) should succeed",
		},
		{
			name:        "triple_nested_with_expressions",
			expression:  `[[[1, 2]]].all(a, a.all(b, b.exists(c, (c + 1) > 0)))`,
			shouldErr:   false,
			description: "Triple nested with expressions (depth 3) should succeed",
		},
		{
			name:        "quadruple_nested_even_with_simple_predicates",
			expression:  `[[[[1, 2]]]].all(a, a.all(b, b.all(c, c.exists(d, true))))`,
			shouldErr:   true,
			description: "Quadruple nested (depth 4) should fail even with simple predicates",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv()
			require.NoError(t, err)

			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err())

			_, err = cel2sql.Convert(ast)

			if tt.shouldErr {
				require.Error(t, err, tt.description)
				require.Contains(t, err.Error(), "comprehension depth")
				require.Contains(t, err.Error(), "exceeds limit")
			} else {
				require.NoError(t, err, tt.description)
			}
		})
	}
}

func TestComprehensionDepthBoundaryConditions(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		depth      int
		shouldErr  bool
	}{
		{
			name:       "depth_1",
			expression: `[1, 2].all(a, a > 0)`,
			depth:      1,
			shouldErr:  false,
		},
		{
			name:       "depth_2",
			expression: `[[1, 2]].all(a, a.all(b, b > 0))`,
			depth:      2,
			shouldErr:  false,
		},
		{
			name:       "depth_3_at_limit",
			expression: `[[[1, 2]]].all(a, a.all(b, b.all(c, c > 0)))`,
			depth:      3,
			shouldErr:  false,
		},
		{
			name:       "depth_4_exceeds",
			expression: `[[[[1, 2]]]].all(a, a.all(b, b.all(c, c.all(d, d > 0))))`,
			depth:      4,
			shouldErr:  true,
		},
		{
			name:       "depth_5_exceeds",
			expression: `[[[[[1, 2]]]]].all(a, a.all(b, b.all(c, c.all(d, d.all(e, e > 0)))))`,
			depth:      5,
			shouldErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := cel.NewEnv()
			require.NoError(t, err)

			ast, issues := env.Compile(tt.expression)
			require.NoError(t, issues.Err())

			_, err = cel2sql.Convert(ast)

			if tt.shouldErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "comprehension depth")
				require.Contains(t, err.Error(), "exceeds limit")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
