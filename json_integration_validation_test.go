package cel2sql_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/cel-go/cel"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/spandigital/cel2sql/v2"
	"github.com/spandigital/cel2sql/v2/pg"
)

// TestGeneratedSQLAgainstPostgreSQL tests if our generated SQL actually works against PostgreSQL
func TestGeneratedSQLAgainstPostgreSQL(t *testing.T) {
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

	// Create test table matching our CEL schema
	_, err = db.Exec(`CREATE TABLE obj (id INT, metadata JSONB)`)
	require.NoError(t, err)

	// Insert test data
	_, err = db.Exec(`INSERT INTO obj VALUES (1, '{"user_name": "test", "settings": {"theme": "dark"}}')`)
	require.NoError(t, err)

	// Set up CEL environment
	testSchema := pg.Schema{
		{Name: "id", Type: "integer"},
		{Name: "metadata", Type: "jsonb", IsJSON: true},
	}

	provider := pg.NewTypeProvider(map[string]pg.Schema{
		"TestTable": testSchema,
	})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("obj", cel.ObjectType("TestTable")),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		celExpr     string
		shouldWork  bool
		description string
	}{
		{
			name:        "Simple JSON field access",
			celExpr:     `obj.metadata.user_name == "test"`,
			shouldWork:  true,
			description: "Access nested JSON field",
		},
		{
			name:        "Nested JSON access",
			celExpr:     `obj.metadata.settings.theme == "dark"`,
			shouldWork:  true,
			description: "Access deeply nested JSON field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compile CEL expression
			ast, issues := env.Compile(tt.celExpr)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("CEL compilation failed: %v", issues.Err())
			}

			// Convert to SQL using schema information
			schemas := map[string]pg.Schema{
				"obj": testSchema,
			}
			sqlCondition, err := cel2sql.ConvertWithSchemas(ast, schemas)
			require.NoError(t, err)

			t.Logf("CEL Expression: %s", tt.celExpr)
			t.Logf("Generated SQL WHERE clause: %s", sqlCondition)

			// Try to execute the generated SQL
			// #nosec G202 - This is a test validating SQL generation, not a security risk
			query := "SELECT * FROM obj WHERE " + sqlCondition
			t.Logf("Full SQL Query: %s", query)

			rows, err := db.Query(query)
			if tt.shouldWork {
				if err != nil {
					t.Errorf("❌ Generated SQL failed to execute: %v", err)
					t.Errorf("   This means the SQL syntax is incorrect for PostgreSQL")
					t.Errorf("   Expected it to work but got error")
				} else {
					defer func() {
						if closeErr := rows.Close(); closeErr != nil {
							t.Logf("failed to close rows: %v", closeErr)
						}
					}()
					hasRow := rows.Next()
					if hasRow {
						t.Logf("✓ Generated SQL works correctly and returns expected results")
					} else {
						t.Errorf("❌ Generated SQL executed but returned no rows (expected 1 row)")
					}
				}
			} else {
				if err != nil {
					t.Logf("✓ Generated SQL failed as expected: %v", err)
				} else {
					defer func() {
						if closeErr := rows.Close(); closeErr != nil {
							t.Logf("failed to close rows: %v", closeErr)
						}
					}()
					t.Errorf("❌ Generated SQL should have failed but succeeded")
				}
			}
		})
	}
}
