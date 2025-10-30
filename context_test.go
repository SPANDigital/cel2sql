package cel2sql_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
)

// TestConvertWithContext_NilContext tests that conversion works without a context (backward compatibility)
func TestConvertWithContext_NilContext(t *testing.T) {
	env, err := cel.NewEnv()
	require.NoError(t, err)

	ast, issues := env.Compile(`1 + 2 == 3`)
	require.NoError(t, issues.Err())

	// Convert without providing context
	sql, err := cel2sql.Convert(ast)
	require.NoError(t, err)
	assert.Equal(t, "1 + 2 = 3", sql)
}

// TestConvertWithContext_ActiveContext tests that conversion works with an active context
func TestConvertWithContext_ActiveContext(t *testing.T) {
	env, err := cel.NewEnv()
	require.NoError(t, err)

	ast, issues := env.Compile(`1 + 2 == 3`)
	require.NoError(t, issues.Err())

	// Convert with an active context
	ctx := context.Background()
	sql, err := cel2sql.Convert(ast, cel2sql.WithContext(ctx))
	require.NoError(t, err)
	assert.Equal(t, "1 + 2 = 3", sql)
}

// TestConvertWithContext_CancelledContext tests that a cancelled context stops conversion
func TestConvertWithContext_CancelledContext(t *testing.T) {
	env, err := cel.NewEnv()
	require.NoError(t, err)

	// Create a complex expression with nested operations
	ast, issues := env.Compile(`(1 + 2) * (3 + 4) == (5 + 6) * (7 + 8)`)
	require.NoError(t, issues.Err())

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Conversion should fail with context cancelled error
	sql, err := cel2sql.Convert(ast, cel2sql.WithContext(ctx))
	require.Error(t, err)
	assert.Empty(t, sql)
	assert.True(t, errors.Is(err, context.Canceled), "error should wrap context.Canceled")
	assert.Contains(t, err.Error(), "conversion cancelled")
}

// TestConvertWithContext_Timeout tests that a timeout context stops conversion
func TestConvertWithContext_Timeout(t *testing.T) {
	env, err := cel.NewEnv()
	require.NoError(t, err)

	// Create a complex expression
	ast, issues := env.Compile(`(1 + 2) * (3 + 4) == (5 + 6) * (7 + 8)`)
	require.NoError(t, issues.Err())

	// Create a context with immediate timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Give the context a chance to expire
	time.Sleep(10 * time.Millisecond)

	// Conversion should fail with deadline exceeded error
	sql, err := cel2sql.Convert(ast, cel2sql.WithContext(ctx))
	require.Error(t, err)
	assert.Empty(t, sql)
	assert.True(t, errors.Is(err, context.DeadlineExceeded), "error should wrap context.DeadlineExceeded")
	assert.Contains(t, err.Error(), "conversion cancelled")
}

// TestConvertWithContext_ComplexExpression tests context cancellation with a complex expression
func TestConvertWithContext_ComplexExpression(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
		{Name: "age", Type: "integer"},
		{Name: "scores", Type: "integer", Repeated: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"Person": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("person", cel.ObjectType("Person")),
	)
	require.NoError(t, err)

	// Create a complex expression with comprehensions
	ast, issues := env.Compile(`person.scores.all(s, s > 50) && person.age >= 18 && person.name.startsWith("A")`)
	require.NoError(t, issues.Err())

	// Test with active context
	ctx := context.Background()
	schemas := provider.GetSchemas()
	sql, err := cel2sql.Convert(ast, cel2sql.WithContext(ctx), cel2sql.WithSchemas(schemas))
	require.NoError(t, err)
	assert.NotEmpty(t, sql)

	// Test with cancelled context
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	sql, err = cel2sql.Convert(ast, cel2sql.WithContext(cancelledCtx), cel2sql.WithSchemas(schemas))
	require.Error(t, err)
	assert.Empty(t, sql)
	assert.Contains(t, err.Error(), "conversion cancelled")
}

// TestConvertWithContext_MultipleOptions tests combining context with other options
func TestConvertWithContext_MultipleOptions(t *testing.T) {
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"resource": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("resource", cel.ObjectType("resource")),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`resource.metadata.name == "test"`)
	require.NoError(t, issues.Err())

	// Test with both context and schemas
	ctx := context.Background()
	schemas := provider.GetSchemas()

	sql, err := cel2sql.Convert(ast,
		cel2sql.WithContext(ctx),
		cel2sql.WithSchemas(schemas))

	require.NoError(t, err)
	assert.Contains(t, sql, "metadata")
	assert.Contains(t, sql, "->>")
}

// TestConvertWithContext_LongRunningConversion tests that context can cancel long conversions
func TestConvertWithContext_LongRunningConversion(t *testing.T) {
	// Skip in short mode as this test involves timing
	if testing.Short() {
		t.Skip("Skipping long-running test in short mode")
	}

	env, err := cel.NewEnv()
	require.NoError(t, err)

	// Create a deeply nested expression
	deepExpr := "((((1 + 2) * (3 + 4)) == ((5 + 6) * (7 + 8))) && (((9 + 10) * (11 + 12)) == ((13 + 14) * (15 + 16))))"
	ast, issues := env.Compile(deepExpr)
	require.NoError(t, issues.Err())

	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer cancel()

	// Wait for timeout
	time.Sleep(10 * time.Millisecond)

	// Should fail due to timeout
	_, err = cel2sql.Convert(ast, cel2sql.WithContext(ctx))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "conversion cancelled")
}
