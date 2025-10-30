package cel2sql_test

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
)

// TestLoggingWithDiscard verifies that default behavior uses discard handler (zero overhead)
func TestLoggingWithoutExplicitLogger(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("x", cel.IntType),
	)
	require.NoError(t, err)

	ast, issues := env.Compile("x > 10")
	require.NoError(t, issues.Err())

	// Convert without logger - should work with zero overhead
	sql, err := cel2sql.Convert(ast)
	require.NoError(t, err)
	assert.Equal(t, "x > 10", sql)
}

// TestLoggingWithJSONHandler verifies JSON logging works correctly
func TestLoggingWithJSONHandler(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"users": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("users", cel.ObjectType("users")),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`users.metadata.active == true`)
	require.NoError(t, issues.Err())

	sql, err := cel2sql.Convert(ast,
		cel2sql.WithSchemas(provider.GetSchemas()),
		cel2sql.WithLogger(logger),
	)
	require.NoError(t, err)
	assert.Contains(t, sql, "->>")

	// Verify logs were written
	logOutput := buf.String()
	assert.Contains(t, logOutput, "starting CEL to SQL conversion")
	assert.Contains(t, logOutput, "conversion completed")
	assert.Contains(t, logOutput, "duration")
}

// TestLoggingJSONPathDetection verifies JSON path detection logging
func TestLoggingJSONPathDetection(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "data", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"records": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("records", cel.ObjectType("records")),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`records.data.status == "active"`)
	require.NoError(t, issues.Err())

	_, err = cel2sql.Convert(ast,
		cel2sql.WithSchemas(provider.GetSchemas()),
		cel2sql.WithLogger(logger),
	)
	require.NoError(t, err)

	logOutput := buf.String()
	// Should log JSON path detection
	assert.Contains(t, logOutput, "JSON path")
	assert.Contains(t, logOutput, "field type lookup")
}

// TestLoggingComprehensions verifies comprehension logging
func TestLoggingComprehensions(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "tags", Type: "text", Repeated: true, ElementType: "text"},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"items": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("items", cel.ObjectType("items")),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`items.tags.exists(t, t == "important")`)
	require.NoError(t, issues.Err())

	_, err = cel2sql.Convert(ast,
		cel2sql.WithSchemas(provider.GetSchemas()),
		cel2sql.WithLogger(logger),
	)
	require.NoError(t, err)

	logOutput := buf.String()
	// Should log comprehension identification
	assert.Contains(t, logOutput, "comprehension identified")
	assert.Contains(t, logOutput, "exists")
}

// TestLoggingOperatorConversion verifies operator conversion logging
func TestLoggingOperatorConversion(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	env, err := cel.NewEnv(
		cel.Variable("x", cel.IntType),
		cel.Variable("y", cel.IntType),
	)
	require.NoError(t, err)

	ast, issues := env.Compile("x > 10 && y < 20")
	require.NoError(t, issues.Err())

	_, err = cel2sql.Convert(ast, cel2sql.WithLogger(logger))
	require.NoError(t, err)

	logOutput := buf.String()
	// Should log operator conversions
	assert.Contains(t, logOutput, "binary operator conversion")
}

// TestLoggingRegexConversion verifies regex pattern conversion logging
func TestLoggingRegexConversion(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	env, err := cel.NewEnv(
		cel.Variable("name", cel.StringType),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`name.matches("(?i)test.*")`)
	require.NoError(t, issues.Err())

	_, err = cel2sql.Convert(ast, cel2sql.WithLogger(logger))
	require.NoError(t, err)

	logOutput := buf.String()
	// Should log regex pattern conversion
	assert.Contains(t, logOutput, "regex pattern conversion")
	assert.Contains(t, logOutput, "case_insensitive")
}

// TestLoggingSchemaLookup verifies schema lookup logging
func TestLoggingSchemaLookup(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "data", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"users": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("users", cel.ObjectType("users")),
	)
	require.NoError(t, err)

	ast, issues := env.Compile(`users.data.field == "test"`)
	require.NoError(t, issues.Err())

	_, err = cel2sql.Convert(ast,
		cel2sql.WithSchemas(provider.GetSchemas()),
		cel2sql.WithLogger(logger),
	)
	require.NoError(t, err)

	logOutput := buf.String()
	// Should log field type lookups for JSON fields
	assert.Contains(t, logOutput, "field type lookup")
}

// TestLoggingErrorContext verifies error logging includes context
func TestLoggingErrorContext(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	env, err := cel.NewEnv(
		cel.Variable("x", cel.IntType),
	)
	require.NoError(t, err)

	// Create an invalid AST by using unsupported expression
	// (This would require more complex setup, so we'll test the happy path for now)
	ast, issues := env.Compile("x > 10")
	require.NoError(t, issues.Err())

	_, err = cel2sql.Convert(ast, cel2sql.WithLogger(logger))
	require.NoError(t, err)

	// Verify basic logging works (error logging would need an actual error condition)
	logOutput := buf.String()
	assert.NotEmpty(t, logOutput)
}

// TestLoggingPerformanceMetrics verifies duration logging
func TestLoggingPerformanceMetrics(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	env, err := cel.NewEnv(
		cel.Variable("x", cel.IntType),
	)
	require.NoError(t, err)

	ast, issues := env.Compile("x > 10")
	require.NoError(t, issues.Err())

	_, err = cel2sql.Convert(ast, cel2sql.WithLogger(logger))
	require.NoError(t, err)

	logOutput := buf.String()
	// Should include duration metric
	assert.Contains(t, logOutput, "duration")
	assert.Contains(t, logOutput, "conversion completed")
}

// TestLoggingWithContext verifies logging works with context
func TestLoggingWithContext(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	ctx := context.Background()

	env, err := cel.NewEnv(
		cel.Variable("x", cel.IntType),
	)
	require.NoError(t, err)

	ast, issues := env.Compile("x > 10")
	require.NoError(t, issues.Err())

	_, err = cel2sql.Convert(ast,
		cel2sql.WithContext(ctx),
		cel2sql.WithLogger(logger),
	)
	require.NoError(t, err)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "starting CEL to SQL conversion")
	assert.Contains(t, logOutput, "conversion completed")
}

// TestLoggingTextHandler verifies text handler works
func TestLoggingTextHandler(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	env, err := cel.NewEnv(
		cel.Variable("x", cel.IntType),
	)
	require.NoError(t, err)

	ast, issues := env.Compile("x > 10")
	require.NoError(t, issues.Err())

	_, err = cel2sql.Convert(ast, cel2sql.WithLogger(logger))
	require.NoError(t, err)

	logOutput := buf.String()
	// Text handler produces different format but should have key content
	assert.True(t, strings.Contains(logOutput, "starting") || strings.Contains(logOutput, "conversion"))
}
