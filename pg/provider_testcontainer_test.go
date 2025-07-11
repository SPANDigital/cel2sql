package pg_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/spandigital/cel2sql"
	"github.com/spandigital/cel2sql/pg"
)

func TestLoadTableSchema_WithPostgresContainer(t *testing.T) {
	ctx := context.Background()

	// Create a PostgreSQL container
	container, err := postgres.Run(ctx,
		"postgres:15",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.WithInitScripts("create_test_table.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(time.Second*60),
		),
	)
	require.NoError(t, err)

	// Cleanup container after test
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Errorf("failed to terminate container: %v", err)
		}
	}()

	// Get connection string
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// Create type provider with database connection
	provider, err := pg.NewTypeProviderWithConnection(ctx, connStr)
	require.NoError(t, err)
	defer provider.Close()

	// Test LoadTableSchema
	err = provider.LoadTableSchema(ctx, "users")
	require.NoError(t, err)

	// Verify the schema was loaded correctly
	// Test FindStructType
	foundType, found := provider.FindStructType("users")
	assert.True(t, found, "users type should be found")
	assert.NotNil(t, foundType, "users type should not be nil")

	// Test FindStructFieldNames
	fieldNames, found := provider.FindStructFieldNames("users")
	assert.True(t, found, "users field names should be found")
	assert.Contains(t, fieldNames, "id")
	assert.Contains(t, fieldNames, "name")
	assert.Contains(t, fieldNames, "email")
	assert.Contains(t, fieldNames, "age")
	assert.Contains(t, fieldNames, "created_at")
	assert.Contains(t, fieldNames, "is_active")

	// Test FindStructFieldType for each expected field
	testCases := []struct {
		fieldName string
	}{
		{"id"},
		{"name"},
		{"email"},
		{"age"},
		{"created_at"},
		{"is_active"},
	}

	for _, tc := range testCases {
		t.Run("field_"+tc.fieldName, func(t *testing.T) {
			fieldType, found := provider.FindStructFieldType("users", tc.fieldName)
			assert.True(t, found, "field %s should be found", tc.fieldName)
			assert.NotNil(t, fieldType, "field %s type should not be nil", tc.fieldName)
		})
	}
}

func TestLoadTableSchema_WithArrayTypes(t *testing.T) {
	ctx := context.Background()

	// Create a PostgreSQL container with array types
	container, err := postgres.Run(ctx,
		"postgres:15",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.WithInitScripts("create_array_table.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(time.Second*60),
		),
	)
	require.NoError(t, err)

	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Errorf("failed to terminate container: %v", err)
		}
	}()

	// Get connection string
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// Create type provider with database connection
	provider, err := pg.NewTypeProviderWithConnection(ctx, connStr)
	require.NoError(t, err)
	defer provider.Close()

	// Test LoadTableSchema for array types
	err = provider.LoadTableSchema(ctx, "products")
	require.NoError(t, err)

	// Test array field type
	fieldType, found := provider.FindStructFieldType("products", "tags")
	assert.True(t, found, "tags field should be found")
	assert.NotNil(t, fieldType, "tags field type should not be nil")

	// Test scores array field
	scoresFieldType, found := provider.FindStructFieldType("products", "scores")
	assert.True(t, found, "scores field should be found")
	assert.NotNil(t, scoresFieldType, "scores field type should not be nil")
}

func TestLoadTableSchema_NonExistentTable(t *testing.T) {
	ctx := context.Background()

	// Create a PostgreSQL container
	container, err := postgres.Run(ctx,
		"postgres:15",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(time.Second*60),
		),
	)
	require.NoError(t, err)

	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Errorf("failed to terminate container: %v", err)
		}
	}()

	// Get connection string
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// Create type provider with database connection
	provider, err := pg.NewTypeProviderWithConnection(ctx, connStr)
	require.NoError(t, err)
	defer provider.Close()

	// Test LoadTableSchema for non-existent table
	err = provider.LoadTableSchema(ctx, "non_existent_table")
	require.NoError(t, err) // Should not error, just return empty schema

	// Verify the table type is found but has no fields
	foundType, found := provider.FindStructType("non_existent_table")
	assert.True(t, found, "non_existent_table type should be found")
	assert.NotNil(t, foundType, "non_existent_table type should not be nil")
}

func TestLoadTableSchema_WithoutConnection(t *testing.T) {
	// Create type provider without database connection
	provider := pg.NewTypeProvider(make(map[string]pg.Schema))

	// Test LoadTableSchema without connection should return error
	err := provider.LoadTableSchema(context.Background(), "any_table")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no database connection available")
}

// TestCELToSQL_ComprehensiveIntegration tests the complete workflow:
// 1. Load table schemas from PostgreSQL
// 2. Convert CEL expressions to SQL conditions  
// 3. Execute queries with date arithmetic and array manipulation
func TestCELToSQL_ComprehensiveIntegration(t *testing.T) {
	ctx := context.Background()

	// Create a PostgreSQL container with comprehensive test data
	container, err := postgres.Run(ctx,
		"postgres:15",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.WithInitScripts("create_comprehensive_test_data.sql"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(time.Second*60),
		),
	)
	require.NoError(t, err)

	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Errorf("failed to terminate container: %v", err)
		}
	}()

	// Get connection string
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// Create type provider with database connection
	provider, err := pg.NewTypeProviderWithConnection(ctx, connStr)
	require.NoError(t, err)
	defer provider.Close()

	// Load schemas for both tables
	err = provider.LoadTableSchema(ctx, "users")
	require.NoError(t, err)
	err = provider.LoadTableSchema(ctx, "products")
	require.NoError(t, err)

	// Create CEL environment with loaded schemas
	celEnv, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("users", cel.ObjectType("users")),
		cel.Variable("products", cel.ObjectType("products")),
	)
	require.NoError(t, err)

	// Test cases combining schema loading, CEL-to-SQL conversion, and query execution
	testCases := []struct {
		name          string
		table         string
		celExpression string
		expectedRows  int
		description   string
	}{
		{
			name:          "basic_equality",
			table:         "users",
			celExpression: `users.name == "John Doe"`,
			expectedRows:  1,
			description:   "Basic string equality test",
		},
		{
			name:          "numeric_comparison",
			table:         "users",
			celExpression: `users.age > 28`,
			expectedRows:  5,
			description:   "Numeric comparison test",
		},
		{
			name:          "boolean_filter",
			table:         "users",
			celExpression: `users.is_active == true`,
			expectedRows:  6,
			description:   "Boolean field filtering",
		},
		{
			name:          "date_arithmetic_recent",
			table:         "users",
			celExpression: `users.created_at > timestamp("2024-01-01T00:00:00Z")`,
			expectedRows:  8,
			description:   "Date arithmetic - recent users",
		},
		{
			name:          "array_contains",
			table:         "products",
			celExpression: `"electronics" in products.tags`,
			expectedRows:  3,
			description:   "Array contains check",
		},
		{
			name:          "array_size",
			table:         "products",
			celExpression: `size(products.tags) > 2`,
			expectedRows:  2,
			description:   "Array size comparison",
		},
		{
			name:          "complex_condition",
			table:         "users",
			celExpression: `users.age >= 25 && users.is_active == true`,
			expectedRows:  4,
			description:   "Complex condition with AND",
		},
		{
			name:          "numeric_array_filter",
			table:         "products",
			celExpression: `95 in products.scores`,
			expectedRows:  2,
			description:   "Numeric array membership test",
		},
		{
			name:          "string_operations",
			table:         "users",
			celExpression: `users.email.contains("@example.com")`,
			expectedRows:  9,
			description:   "String contains operation",
		},
		{
			name:          "or_condition",
			table:         "users",
			celExpression: `users.age < 25 || users.age > 40`,
			expectedRows:  4,
			description:   "OR condition test",
		},
		{
			name:          "json_field_access",
			table:         "users",
			celExpression: `users.preferences.theme == "dark"`,
			expectedRows:  4,
			description:   "JSONB field access test",
		},
		{
			name:          "json_boolean_field",
			table:         "users",
			celExpression: `users.preferences.notifications == "true"`,
			expectedRows:  5,
			description:   "JSONB boolean field test (as string)",
		},
		{
			name:          "json_string_field",
			table:         "users",
			celExpression: `users.profile.location == "New York"`,
			expectedRows:  1,
			description:   "JSON string field test",
		},
		{
			name:          "product_json_price_string",
			table:         "products",
			celExpression: `products.metadata.price == "999.99"`,
			expectedRows:  1,
			description:   "JSONB numeric field comparison (as string)",
		},
		{
			name:          "product_json_category",
			table:         "products",
			celExpression: `products.metadata.category == "electronics"`,
			expectedRows:  2,
			description:   "JSONB string field equality",
		},
		{
			name:          "json_complex_condition",
			table:         "users",
			celExpression: `users.preferences.theme == "dark" && users.age > 25`,
			expectedRows:  2,
			description:   "Complex condition with JSON field",
		},
	}

	// Create database connection for executing queries
	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	defer pool.Close()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Parse CEL expression
			ast, issues := celEnv.Parse(tc.celExpression)
			require.NoError(t, issues.Err(), "Failed to parse CEL expression: %s", tc.celExpression)

			// Check CEL expression
			ast, issues = celEnv.Check(ast)
			require.NoError(t, issues.Err(), "Failed to check CEL expression: %s", tc.celExpression)

			// Convert CEL to SQL
			sqlCondition, err := cel2sql.Convert(ast)
			require.NoError(t, err, "Failed to convert CEL to SQL: %s", tc.celExpression)

			t.Logf("CEL: %s", tc.celExpression)
			t.Logf("SQL: %s", sqlCondition)

			// Execute query to validate the generated SQL
			query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", tc.table, sqlCondition)
			var count int
			err = pool.QueryRow(ctx, query).Scan(&count)
			require.NoError(t, err, "Failed to execute query: %s", query)

			// Verify expected results
			assert.Equal(t, tc.expectedRows, count,
				"Expected %d rows but got %d for test case: %s\nCEL: %s\nSQL: %s\nQuery: %s",
				tc.expectedRows, count, tc.description, tc.celExpression, sqlCondition, query)
		})
	}

	// Additional test for complex date arithmetic
	t.Run("date_arithmetic_complex", func(t *testing.T) {
		celExpression := `users.created_at > timestamp("2024-06-01T00:00:00Z") && users.created_at < timestamp("2024-12-31T23:59:59Z")`

		ast, issues := celEnv.Parse(celExpression)
		require.NoError(t, issues.Err())

		ast, issues = celEnv.Check(ast)
		require.NoError(t, issues.Err())

		sqlCondition, err := cel2sql.Convert(ast)
		require.NoError(t, err)

		t.Logf("Complex date CEL: %s", celExpression)
		t.Logf("Complex date SQL: %s", sqlCondition)

		query := "SELECT COUNT(*) FROM users WHERE " + sqlCondition
		var count int
		err = pool.QueryRow(ctx, query).Scan(&count)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, count, 0, "Date arithmetic query should execute without error")
	})

	// Test array manipulation with complex conditions
	t.Run("array_manipulation_complex", func(t *testing.T) {
		celExpression := `size(products.tags) >= 2 && "electronics" in products.tags`

		ast, issues := celEnv.Parse(celExpression)
		require.NoError(t, issues.Err())

		ast, issues = celEnv.Check(ast)
		require.NoError(t, issues.Err())

		sqlCondition, err := cel2sql.Convert(ast)
		require.NoError(t, err)

		t.Logf("Complex array CEL: %s", celExpression)
		t.Logf("Complex array SQL: %s", sqlCondition)

		query := "SELECT COUNT(*) FROM products WHERE " + sqlCondition
		var count int
		err = pool.QueryRow(ctx, query).Scan(&count)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, count, 0, "Array manipulation query should execute without error")
	})
}
