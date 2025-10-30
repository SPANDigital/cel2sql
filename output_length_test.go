package cel2sql

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/spandigital/cel2sql/v3/pg"
	"github.com/stretchr/testify/require"
)

const testExprXGreaterThanZero = "x > 0"

func TestMaxSQLOutputLength(t *testing.T) {
	tests := []struct {
		name         string
		buildExpr    func() string
		maxOutputLen int
		shouldErr    bool
	}{
		{
			name: "within default limit - simple expression",
			buildExpr: func() string {
				return "x > 0 && y < 100"
			},
			maxOutputLen: 0, // use default (50000)
			shouldErr:    false,
		},
		{
			name: "within custom limit - small",
			buildExpr: func() string {
				return "x == 1"
			},
			maxOutputLen: 100,
			shouldErr:    false,
		},
		{
			name: "exceeds custom limit - moderate expression",
			buildExpr: func() string {
				// Build expression that generates ~150 chars of SQL
				parts := make([]string, 20)
				for i := range parts {
					parts[i] = testExprXGreaterThanZero
				}
				return strings.Join(parts, " && ")
			},
			maxOutputLen: 50, // Very small limit
			shouldErr:    true,
		},
		{
			name: "large array literal within limit",
			buildExpr: func() string {
				// Build array with 100 elements
				parts := make([]string, 100)
				for i := range parts {
					parts[i] = "1"
				}
				return "[" + strings.Join(parts, ", ") + "].all(x, x > 0)"
			},
			maxOutputLen: 0, // use default
			shouldErr:    false,
		},
		{
			name: "very large expression exceeds default limit",
			buildExpr: func() string {
				// Build expression that will generate >50000 chars
				// Each repetition adds roughly 10 chars
				parts := make([]string, 6000)
				for i := range parts {
					parts[i] = testExprXGreaterThanZero
				}
				return strings.Join(parts, " && ")
			},
			maxOutputLen: 0, // use default (50000)
			shouldErr:    true,
		},
		{
			name: "high custom limit allows large expression",
			buildExpr: func() string {
				// Build moderately large expression
				parts := make([]string, 500)
				for i := range parts {
					parts[i] = testExprXGreaterThanZero
				}
				return strings.Join(parts, " || ")
			},
			maxOutputLen: 100000,
			shouldErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := pg.NewSchema([]pg.FieldSchema{
				{Name: "x", Type: "integer"},
				{Name: "y", Type: "integer"},
			})
			provider := pg.NewTypeProvider(map[string]pg.Schema{"vars": schema})

			env, err := cel.NewEnv(
				cel.CustomTypeProvider(provider),
				cel.Variable("x", cel.IntType),
				cel.Variable("y", cel.IntType),
			)
			require.NoError(t, err)

			expr := tt.buildExpr()
			ast, issues := env.Compile(expr)
			require.NoError(t, issues.Err())

			var opts []ConvertOption
			if tt.maxOutputLen > 0 {
				opts = append(opts, WithMaxOutputLength(tt.maxOutputLen))
			}

			sql, err := Convert(ast, opts...)

			if tt.shouldErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "output length")
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, sql)
			}
		})
	}
}

func TestMaxOutputLengthWithOtherOptions(t *testing.T) {
	// Test combining WithMaxOutputLength with other options
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "x", Type: "integer"},
		{Name: "name", Type: "text"},
	})
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
		cel.Variable("name", cel.StringType),
	)
	require.NoError(t, err)

	tests := []struct {
		name string
		expr string
		opts []ConvertOption
	}{
		{
			name: "all options together",
			expr: "x > 5 && name == 'test'",
			opts: []ConvertOption{
				WithMaxOutputLength(10000),
				WithMaxDepth(50),
				WithContext(ctx),
				WithSchemas(schemas),
				WithLogger(logger),
			},
		},
		{
			name: "maxOutputLength with context",
			expr: "x > 10",
			opts: []ConvertOption{
				WithMaxOutputLength(5000),
				WithContext(ctx),
			},
		},
		{
			name: "maxOutputLength with schemas",
			expr: "name.contains('foo')",
			opts: []ConvertOption{
				WithMaxOutputLength(1000),
				WithSchemas(schemas),
			},
		},
		{
			name: "maxOutputLength with maxDepth",
			expr: "((x + 1) + 2) + 3",
			opts: []ConvertOption{
				WithMaxOutputLength(500),
				WithMaxDepth(100),
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

func TestOutputLengthErrorMessage(t *testing.T) {
	// Verify error message format
	schema := pg.NewSchema([]pg.FieldSchema{{Name: "x", Type: "integer"}})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"vars": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("x", cel.IntType),
	)
	require.NoError(t, err)

	// Build expression that will exceed limits
	parts := make([]string, 6000)
	for i := range parts {
		parts[i] = testExprXGreaterThanZero
	}
	expr := strings.Join(parts, " && ")

	ast, issues := env.Compile(expr)
	require.NoError(t, issues.Err())

	// Test default error message
	_, err = Convert(ast)
	require.Error(t, err)
	require.Contains(t, err.Error(), "generated SQL exceeds maximum output length of 50000")

	// Test custom limit error message
	_, err = Convert(ast, WithMaxOutputLength(1000))
	require.Error(t, err)
	require.Contains(t, err.Error(), "generated SQL exceeds maximum output length of 1000")
}

func TestOutputLengthResetBetweenCalls(t *testing.T) {
	// Ensure output length check resets between Convert() calls
	schema := pg.NewSchema([]pg.FieldSchema{{Name: "x", Type: "integer"}})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"vars": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("x", cel.IntType),
	)
	require.NoError(t, err)

	// Build expression that's close to but under default limit
	parts := make([]string, 500)
	for i := range parts {
		parts[i] = testExprXGreaterThanZero
	}
	expr := strings.Join(parts, " && ")

	ast, issues := env.Compile(expr)
	require.NoError(t, issues.Err())

	// First conversion should succeed
	sql1, err := Convert(ast)
	require.NoError(t, err)
	require.NotEmpty(t, sql1)

	// Second conversion should also succeed (output was reset)
	sql2, err := Convert(ast)
	require.NoError(t, err)
	require.NotEmpty(t, sql2)

	// Third conversion should also succeed
	sql3, err := Convert(ast)
	require.NoError(t, err)
	require.NotEmpty(t, sql3)

	// All should produce identical output
	require.Equal(t, sql1, sql2)
	require.Equal(t, sql2, sql3)
}

func TestLargeArrayLiterals(t *testing.T) {
	tests := []struct {
		name         string
		arraySize    int
		maxOutputLen int
		shouldErr    bool
	}{
		{
			name:         "small array within limit",
			arraySize:    10,
			maxOutputLen: 0, // default
			shouldErr:    false,
		},
		{
			name:         "medium array within limit",
			arraySize:    100,
			maxOutputLen: 0, // default
			shouldErr:    false,
		},
		{
			name:         "large array within limit",
			arraySize:    1000,
			maxOutputLen: 0, // default
			shouldErr:    false,
		},
		{
			name:         "huge array exceeds limit",
			arraySize:    20000,
			maxOutputLen: 0, // default (50000)
			shouldErr:    true,
		},
		{
			name:         "array exceeds custom low limit",
			arraySize:    50,
			maxOutputLen: 100,
			shouldErr:    true,
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

			// Build array literal
			parts := make([]string, tt.arraySize)
			for i := range parts {
				parts[i] = "1"
			}
			expr := "[" + strings.Join(parts, ", ") + "].exists(y, y == 1)"

			ast, issues := env.Compile(expr)
			require.NoError(t, issues.Err())

			var opts []ConvertOption
			if tt.maxOutputLen > 0 {
				opts = append(opts, WithMaxOutputLength(tt.maxOutputLen))
			}

			_, err = Convert(ast, opts...)

			if tt.shouldErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "output length")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLongStringConcatenations(t *testing.T) {
	tests := []struct {
		name         string
		stringCount  int
		maxOutputLen int
		shouldErr    bool
	}{
		{
			name:         "small concatenation",
			stringCount:  5,
			maxOutputLen: 0, // default
			shouldErr:    false,
		},
		{
			name:         "medium concatenation",
			stringCount:  20,
			maxOutputLen: 0, // default
			shouldErr:    false,
		},
		{
			name:         "large concatenation within limit",
			stringCount:  50,
			maxOutputLen: 0, // default
			shouldErr:    false,
		},
		{
			name:         "very large concatenation exceeds default",
			stringCount:  70,
			maxOutputLen: 500, // Small limit to trigger error
			shouldErr:    true,
		},
		{
			name:         "concatenation exceeds custom limit",
			stringCount:  10,
			maxOutputLen: 50,
			shouldErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := pg.NewSchema([]pg.FieldSchema{{Name: "s", Type: "text"}})
			provider := pg.NewTypeProvider(map[string]pg.Schema{"vars": schema})

			env, err := cel.NewEnv(
				cel.CustomTypeProvider(provider),
				cel.Variable("s", cel.StringType),
			)
			require.NoError(t, err)

			// Build string concatenation expression
			expr := "s"
			for i := 0; i < tt.stringCount; i++ {
				expr += ` + "test"`
			}

			ast, issues := env.Compile(expr)
			require.NoError(t, issues.Err())

			var opts []ConvertOption
			if tt.maxOutputLen > 0 {
				opts = append(opts, WithMaxOutputLength(tt.maxOutputLen))
			}

			_, err = Convert(ast, opts...)

			if tt.shouldErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "output length")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMaxOutputLengthWithComprehensions(t *testing.T) {
	// Test output length limit with comprehensions
	tests := []struct {
		name         string
		expr         string
		maxOutputLen int
		shouldErr    bool
	}{
		{
			name:         "simple all comprehension",
			expr:         "[1, 2, 3].all(x, x > 0)",
			maxOutputLen: 0, // default
			shouldErr:    false,
		},
		{
			name:         "simple exists comprehension",
			expr:         "[1, 2, 3].exists(x, x == 2)",
			maxOutputLen: 0, // default
			shouldErr:    false,
		},
		{
			name:         "comprehension with low custom limit",
			expr:         "[1, 2, 3, 4, 5].filter(x, x > 2)",
			maxOutputLen: 50,
			shouldErr:    true,
		},
		{
			name:         "nested comprehension",
			expr:         "[[1, 2], [3, 4]].all(arr, arr.all(x, x > 0))",
			maxOutputLen: 0, // default
			shouldErr:    false,
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

			ast, issues := env.Compile(tt.expr)
			require.NoError(t, issues.Err())

			var opts []ConvertOption
			if tt.maxOutputLen > 0 {
				opts = append(opts, WithMaxOutputLength(tt.maxOutputLen))
			}

			_, err = Convert(ast, opts...)

			if tt.shouldErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "output length")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestParameterizedWithOutputLength(t *testing.T) {
	// Test that parameterized conversion also respects output length limits
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "x", Type: "integer"},
		{Name: "name", Type: "text"},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"vars": schema})
	schemas := provider.GetSchemas()

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("x", cel.IntType),
		cel.Variable("name", cel.StringType),
	)
	require.NoError(t, err)

	tests := []struct {
		name         string
		buildExpr    func() string
		maxOutputLen int
		shouldErr    bool
	}{
		{
			name: "parameterized within limit",
			buildExpr: func() string {
				return `x > 5 && name == "test"`
			},
			maxOutputLen: 0,
			shouldErr:    false,
		},
		{
			name: "parameterized exceeds custom limit",
			buildExpr: func() string {
				parts := make([]string, 100)
				for i := range parts {
					parts[i] = testExprXGreaterThanZero
				}
				return strings.Join(parts, " && ")
			},
			maxOutputLen: 100,
			shouldErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := tt.buildExpr()
			ast, issues := env.Compile(expr)
			require.NoError(t, issues.Err())

			var opts []ConvertOption
			opts = append(opts, WithSchemas(schemas))
			if tt.maxOutputLen > 0 {
				opts = append(opts, WithMaxOutputLength(tt.maxOutputLen))
			}

			result, err := ConvertParameterized(ast, opts...)

			if tt.shouldErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), "output length")
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.NotEmpty(t, result.SQL)
			}
		})
	}
}
