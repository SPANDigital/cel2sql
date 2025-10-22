package cel2sql_test

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestPostgreSQLJSONOperatorValidation validates the correct usage of -> vs ->> operators
// with a real PostgreSQL instance to determine the correct behavior
func TestPostgreSQLJSONOperatorValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Start PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:17-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	}()

	// Get connection string
	host, err := postgresContainer.Host(ctx)
	require.NoError(t, err)
	port, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	connStr := "host=" + host + " port=" + port.Port() + " user=postgres password=password dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Logf("failed to close db: %v", closeErr)
		}
	}()

	// Create test table
	_, err = db.Exec(`CREATE TABLE test_table (id INT, metadata JSONB)`)
	require.NoError(t, err)

	// Insert test data with nested JSON
	_, err = db.Exec(`INSERT INTO test_table VALUES (1, '{"user_name": "test", "settings": {"theme": "dark"}}')`)
	require.NoError(t, err)

	tests := []struct {
		name        string
		query       string
		shouldWork  bool
		description string
	}{
		{
			name:        "Correct: table.column->>'field'",
			query:       "SELECT * FROM test_table WHERE metadata->>'user_name' = 'test'",
			shouldWork:  true,
			description: "Standard JSON field extraction from a JSONB column",
		},
		{
			name:        "WRONG: table->>'column'->>'field'",
			query:       "SELECT * FROM test_table WHERE test_table->>'metadata'->>'user_name' = 'test'",
			shouldWork:  false,
			description: "Cannot use ->> on table name - this should fail",
		},
		{
			name:        "Correct nested: table.column->'intermediate'->>'final'",
			query:       "SELECT * FROM test_table WHERE metadata->'settings'->>'theme' = 'dark'",
			shouldWork:  true,
			description: "Use -> for intermediate JSON fields, ->> for final",
		},
		{
			name:        "WRONG nested: table.column->>'intermediate'->>'final'",
			query:       "SELECT * FROM test_table WHERE metadata->>'settings'->>'theme' = 'dark'",
			shouldWork:  false,
			description: "Cannot chain ->> because it returns text",
		},
		{
			name:        "WRONG nested: table->>'column'->'intermediate'->>'final'",
			query:       "SELECT * FROM test_table WHERE test_table->>'metadata'->'settings'->>'theme' = 'dark'",
			shouldWork:  false,
			description: "Cannot use JSON operators on table name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, err := db.Query(tt.query)
			if tt.shouldWork {
				require.NoError(t, err, "Query should succeed: %s", tt.description)
				require.NotNil(t, rows)
				defer func() {
					if closeErr := rows.Close(); closeErr != nil {
						t.Logf("failed to close rows: %v", closeErr)
					}
				}()

				// Verify we get a result
				hasRow := rows.Next()
				require.True(t, hasRow, "Query should return at least one row")

				t.Logf("✓ Query works correctly: %s", tt.query)
			} else {
				require.Error(t, err, "Query should fail: %s", tt.description)
				t.Logf("✓ Query fails as expected: %s", tt.query)
				t.Logf("  Error: %v", err)
			}
		})
	}

	// Document the findings
	t.Log("\n=== VALIDATION RESULTS ===")
	t.Log("For CEL expression: obj.metadata.user_name")
	t.Log("Where obj is a table with a JSONB column 'metadata'")
	t.Log("Correct SQL: obj.metadata->>'user_name'")
	t.Log("NOT: obj->>'metadata'->>'user_name'")
	t.Log("")
	t.Log("For CEL expression: obj.metadata.settings.theme")
	t.Log("Correct SQL: obj.metadata->'settings'->>'theme'")
	t.Log("NOT: obj->>'metadata'->'settings'->>'theme'")
	t.Log("NOT: obj.metadata->>'settings'->>'theme'")
}
