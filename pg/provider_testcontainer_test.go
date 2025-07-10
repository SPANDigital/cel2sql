package pg_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

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
				WithStartupTimeout(time.Second * 60),
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
	// Test FindType
	foundType, found := provider.FindType("users")
	assert.True(t, found, "users type should be found")
	assert.NotNil(t, foundType, "users type should not be nil")

	// Test FindFieldType for each expected field
	testCases := []struct {
		fieldName    string
		expectedType string
	}{
		{"id", "int"},
		{"name", "text"},
		{"email", "text"},
		{"age", "int"},
		{"created_at", "timestamp"},
		{"is_active", "boolean"},
	}

	for _, tc := range testCases {
		t.Run("field_"+tc.fieldName, func(t *testing.T) {
			fieldType, found := provider.FindFieldType("users", tc.fieldName)
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
				WithStartupTimeout(time.Second * 60),
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
	fieldType, found := provider.FindFieldType("products", "tags")
	assert.True(t, found, "tags field should be found")
	assert.NotNil(t, fieldType, "tags field type should not be nil")

	// Test scores array field
	scoresFieldType, found := provider.FindFieldType("products", "scores")
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
				WithStartupTimeout(time.Second * 60),
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
	foundType, found := provider.FindType("non_existent_table")
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
