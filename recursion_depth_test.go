package cel2sql

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/spandigital/cel2sql/v3/pg"
	"github.com/stretchr/testify/require"
)

func TestMaxRecursionDepth(t *testing.T) {
	tests := []struct {
		name      string
		depth     int
		maxDepth  int
		shouldErr bool
	}{
		{
			name:      "within default limit",
			depth:     50,
			maxDepth:  0, // use default (100)
			shouldErr: false,
		},
		{
			name:      "exceeds default limit",
			depth:     150,
			maxDepth:  0, // use default (100)
			shouldErr: true,
		},
		{
			name:      "within custom limit",
			depth:     20,
			maxDepth:  100,
			shouldErr: false,
		},
		{
			name:      "exceeds custom limit",
			depth:     55,
			maxDepth:  50,
			shouldErr: true,
		},
		{
			name:      "very deep with high custom limit",
			depth:     200,
			maxDepth:  250,
			shouldErr: false,
		},
		{
			name:      "shallow expression with low limit",
			depth:     5,
			maxDepth:  10,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build deeply nested expression: (((x + 1) + 1) + 1)...
			expr := "x"
			for i := 0; i < tt.depth; i++ {
				expr = "(" + expr + " + 1)"
			}

			schema := pg.NewSchema([]pg.FieldSchema{{Name: "x", Type: "integer"}})
			provider := pg.NewTypeProvider(map[string]pg.Schema{"vars": schema})

			env, err := cel.NewEnv(
				cel.CustomTypeProvider(provider),
				cel.Variable("x", cel.IntType),
			)
			require.NoError(t, err)

			ast, issues := env.Compile(expr)
			require.NoError(t, issues.Err())

			var opts []ConvertOption
			if tt.maxDepth > 0 {
				opts = append(opts, WithMaxDepth(tt.maxDepth))
			}

			_, err = Convert(ast, opts...)

			if tt.shouldErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "recursion depth")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMaxDepthWithOtherOptions(t *testing.T) {
	// Test combining WithMaxDepth with other options
	schema := pg.NewSchema([]pg.FieldSchema{{Name: "x", Type: "integer"}})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"vars": schema})
	schemas := provider.GetSchemas()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("x", cel.IntType),
	)
	require.NoError(t, err)

	tests := []struct {
		name string
		expr string
		opts []ConvertOption
	}{
		{
			name: "all options together",
			expr: "((x + 1) + 2) + 3",
			opts: []ConvertOption{
				WithMaxDepth(50),
				WithContext(ctx),
				WithSchemas(schemas),
				WithLogger(logger),
			},
		},
		{
			name: "maxDepth with context",
			expr: "(x + 1) + 2",
			opts: []ConvertOption{
				WithMaxDepth(100),
				WithContext(ctx),
			},
		},
		{
			name: "maxDepth with schemas",
			expr: "x > 5",
			opts: []ConvertOption{
				WithMaxDepth(75),
				WithSchemas(schemas),
			},
		},
		{
			name: "maxDepth with logger",
			expr: "x == 10",
			opts: []ConvertOption{
				WithMaxDepth(150),
				WithLogger(logger),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expr)
			require.NoError(t, issues.Err())

			sql, err := Convert(ast, tt.opts...)
			require.NoError(t, err)
			require.NotEmpty(t, sql)
		})
	}
}

func TestRecursionDepthErrorMessage(t *testing.T) {
	// Verify error message format
	expr := "x"
	for i := 0; i < 150; i++ {
		expr = "(" + expr + " + 1)"
	}

	schema := pg.NewSchema([]pg.FieldSchema{{Name: "x", Type: "integer"}})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"vars": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("x", cel.IntType),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(expr)
	require.NoError(t, issues.Err())

	// Test default error message
	_, err = Convert(ast)
	require.Error(t, err)
	require.Contains(t, err.Error(), "expression exceeds maximum recursion depth of 100")

	// Test custom depth error message
	_, err = Convert(ast, WithMaxDepth(50))
	require.Error(t, err)
	require.Contains(t, err.Error(), "expression exceeds maximum recursion depth of 50")
}

func TestDeeplyNestedExpressions(t *testing.T) {
	tests := []struct {
		name      string
		buildExpr func() string
		maxDepth  int
		shouldErr bool
	}{
		{
			name: "nested AND conditions",
			buildExpr: func() string {
				expr := "x > 0"
				for i := 0; i < 60; i++ {
					expr = "(" + expr + " && x < 100)"
				}
				return expr
			},
			maxDepth:  0, // default (100)
			shouldErr: false,
		},
		{
			name: "nested OR conditions",
			buildExpr: func() string {
				expr := "x == 1"
				for i := 0; i < 60; i++ {
					expr = "(" + expr + " || x == 2)"
				}
				return expr
			},
			maxDepth:  0, // default (100)
			shouldErr: false,
		},
		{
			name: "nested function calls",
			buildExpr: func() string {
				expr := "x"
				for i := 0; i < 60; i++ {
					expr = "int(" + expr + ")"
				}
				return expr + " > 0"
			},
			maxDepth:  0, // default (100)
			shouldErr: false,
		},
		{
			name: "deeply nested ternary",
			buildExpr: func() string {
				expr := "x"
				for i := 0; i < 25; i++ {
					expr = "(" + expr + " > 0 ? 1 : 0)"
				}
				return expr + " == 1"
			},
			maxDepth:  0, // default (100)
			shouldErr: false,
		},
		{
			name: "exceeds limit nested arithmetic",
			buildExpr: func() string {
				expr := "x"
				for i := 0; i < 150; i++ {
					expr = "(" + expr + " + 1)"
				}
				return expr + " > 0"
			},
			maxDepth:  0, // default (100)
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := pg.NewSchema([]pg.FieldSchema{{Name: "x", Type: "integer"}})
			provider := pg.NewTypeProvider(map[string]pg.Schema{"vars": schema})

			env, err := cel.NewEnv(
				cel.CustomTypeProvider(provider),
				cel.Variable("x", cel.IntType),
			)
			require.NoError(t, err)

			expr := tt.buildExpr()
			ast, issues := env.Compile(expr)
			require.NoError(t, issues.Err())

			var opts []ConvertOption
			if tt.maxDepth > 0 {
				opts = append(opts, WithMaxDepth(tt.maxDepth))
			}

			_, err = Convert(ast, opts...)

			if tt.shouldErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "recursion depth")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDepthResetBetweenCalls(t *testing.T) {
	// Ensure depth counter resets between Convert() calls
	schema := pg.NewSchema([]pg.FieldSchema{{Name: "x", Type: "integer"}})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"vars": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("x", cel.IntType),
	)
	require.NoError(t, err)

	// Build an expression at 75% of default limit
	expr := "x"
	for i := 0; i < 75; i++ {
		expr = "(" + expr + " + 1)"
	}

	ast, issues := env.Compile(expr)
	require.NoError(t, issues.Err())

	// First conversion should succeed
	_, err = Convert(ast)
	require.NoError(t, err)

	// Second conversion should also succeed (depth was reset)
	_, err = Convert(ast)
	require.NoError(t, err)

	// Third conversion should also succeed
	_, err = Convert(ast)
	require.NoError(t, err)
}

func TestMaxDepthWithComprehensions(t *testing.T) {
	// Test that comprehensions are counted in recursion depth
	tests := []struct {
		name      string
		expr      string
		maxDepth  int
		shouldErr bool
	}{
		{
			name:      "simple all comprehension",
			expr:      "[1, 2, 3].all(x, x > 0)",
			maxDepth:  0, // default
			shouldErr: false,
		},
		{
			name:      "simple exists comprehension",
			expr:      "[1, 2, 3].exists(x, x == 2)",
			maxDepth:  0, // default
			shouldErr: false,
		},
		{
			name:      "nested comprehension in condition",
			expr:      "x > 0 && [1, 2, 3].all(y, y < 10)",
			maxDepth:  0, // default
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := pg.NewSchema([]pg.FieldSchema{
				{Name: "x", Type: "integer"},
			})
			provider := pg.NewTypeProvider(map[string]pg.Schema{"vars": schema})

			env, err := cel.NewEnv(
				cel.CustomTypeProvider(provider),
				cel.Variable("x", cel.IntType),
			)
			require.NoError(t, err)

			ast, issues := env.Compile(tt.expr)
			require.NoError(t, issues.Err())

			var opts []ConvertOption
			if tt.maxDepth > 0 {
				opts = append(opts, WithMaxDepth(tt.maxDepth))
			}

			_, err = Convert(ast, opts...)

			if tt.shouldErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "recursion depth")
			} else {
				require.NoError(t, err)
			}
		})
	}
}
