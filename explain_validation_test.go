package cel2sql_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/google/cel-go/cel"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/spandigital/cel2sql/v3"
)

// TestEXPLAINValidation validates that generated SQL produces valid query plans
// and doesn't contain syntax errors that would only be caught during planning
func TestEXPLAINValidation(t *testing.T) {
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

	// Create comprehensive test table
	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT,
			age INTEGER,
			score DOUBLE PRECISION,
			active BOOLEAN,
			tags TEXT[],
			metadata JSONB,
			created_at TIMESTAMP
		)
	`)
	require.NoError(t, err)

	// Create index for testing index scan detection
	_, err = db.Exec(`CREATE INDEX idx_users_email ON users(email)`)
	require.NoError(t, err)

	_, err = db.Exec(`CREATE INDEX idx_users_age ON users(age)`)
	require.NoError(t, err)

	// Insert test data
	_, err = db.Exec(`
		INSERT INTO users VALUES
		(1, 'Alice', 'alice@example.com', 30, 95.5, true, ARRAY['admin', 'user'], '{"role": "admin"}', '2024-01-01 10:00:00'),
		(2, 'Bob', 'bob@example.com', 25, 88.2, true, ARRAY['user'], '{"role": "user"}', '2024-01-02 11:00:00'),
		(3, 'Charlie', 'charlie@example.com', 35, 92.1, false, ARRAY['user', 'guest'], '{"role": "guest"}', '2024-01-03 12:00:00')
	`)
	require.NoError(t, err)

	// Set up CEL environment with simple variables (no struct wrapper)
	env, err := cel.NewEnv(
		cel.Variable("id", cel.IntType),
		cel.Variable("name", cel.StringType),
		cel.Variable("email", cel.StringType),
		cel.Variable("age", cel.IntType),
		cel.Variable("score", cel.DoubleType),
		cel.Variable("active", cel.BoolType),
		cel.Variable("tags", cel.ListType(cel.StringType)),
		cel.Variable("metadata", cel.DynType), // JSON field
		cel.Variable("created_at", cel.TimestampType),
	)
	require.NoError(t, err)

	tests := []struct {
		name              string
		celExpr           string
		expectedInPlan    []string
		notExpectedInPlan []string
		description       string
	}{
		{
			name:    "Simple equality condition",
			celExpr: `email == "alice@example.com"`,
			expectedInPlan: []string{
				"users", // Just verify table is referenced, optimizer may choose Seq Scan or Index Scan
			},
			notExpectedInPlan: []string{
				"ERROR",
				"INVALID",
			},
			description: "Basic WHERE clause should produce valid query plan",
		},
		{
			name:    "Multiple AND conditions",
			celExpr: `age > 25 && active == true`,
			expectedInPlan: []string{
				"users",
			},
			notExpectedInPlan: []string{
				"ERROR",
			},
			description: "Multiple conditions should be optimizable",
		},
		{
			name:    "OR conditions",
			celExpr: `age < 30 || score > 90.0`,
			expectedInPlan: []string{
				"users",
			},
			notExpectedInPlan: []string{
				"ERROR",
			},
			description: "OR conditions should produce valid plan",
		},
		{
			name:    "Array membership",
			celExpr: `"admin" in tags`,
			expectedInPlan: []string{
				"users",
			},
			notExpectedInPlan: []string{
				"ERROR",
			},
			description: "Array operations should be valid",
		},
		{
			name:    "Complex nested condition",
			celExpr: `(age > 25 && active) || (score > 90.0 && !active)`,
			expectedInPlan: []string{
				"users",
			},
			notExpectedInPlan: []string{
				"ERROR",
			},
			description: "Complex nested boolean logic should be optimizable",
		},
		{
			name:    "String operations",
			celExpr: `name.contains("li")`,
			expectedInPlan: []string{
				"users",
			},
			notExpectedInPlan: []string{
				"ERROR",
			},
			description: "String functions should produce valid plans",
		},
		{
			name:    "Timestamp comparison",
			celExpr: `created_at > timestamp("2024-01-01T00:00:00Z")`,
			expectedInPlan: []string{
				"users",
			},
			notExpectedInPlan: []string{
				"ERROR",
			},
			description: "Timestamp operations should be valid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compile CEL expression
			ast, issues := env.Compile(tt.celExpr)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("CEL compilation failed: %v", issues.Err())
			}

			// Convert to SQL
			sqlCondition, err := cel2sql.Convert(ast)
			require.NoError(t, err)

			t.Logf("CEL Expression: %s", tt.celExpr)
			t.Logf("Generated SQL WHERE clause: %s", sqlCondition)

			// Run EXPLAIN to validate query plan
			// #nosec G202 - This is a test validating SQL generation, not a security risk
			explainQuery := "EXPLAIN SELECT * FROM users WHERE " + sqlCondition
			t.Logf("EXPLAIN Query: %s", explainQuery)

			rows, err := db.Query(explainQuery)
			require.NoError(t, err, "EXPLAIN should not fail for valid SQL")
			defer func() {
				if closeErr := rows.Close(); closeErr != nil {
					t.Logf("failed to close rows: %v", closeErr)
				}
			}()

			// Collect all query plan lines
			var planLines []string
			for rows.Next() {
				var planLine string
				err := rows.Scan(&planLine)
				require.NoError(t, err)
				planLines = append(planLines, planLine)
				t.Logf("  Plan: %s", planLine)
			}

			require.NotEmpty(t, planLines, "EXPLAIN should return query plan")

			// Combine all plan lines for pattern matching
			fullPlan := strings.Join(planLines, "\n")

			// Check for expected patterns in plan
			for _, expected := range tt.expectedInPlan {
				require.Contains(t, fullPlan, expected,
					"Query plan should contain '%s'. Description: %s", expected, tt.description)
			}

			// Check that unwanted patterns are not in plan
			for _, notExpected := range tt.notExpectedInPlan {
				require.NotContains(t, fullPlan, notExpected,
					"Query plan should not contain '%s'. Description: %s", notExpected, tt.description)
			}

			t.Logf("✓ EXPLAIN validation passed: %s", tt.description)
		})
	}
}

// TestEXPLAINAnalyzeValidation validates generated SQL with EXPLAIN ANALYZE
// This actually executes the query and validates both plan and execution
func TestEXPLAINAnalyzeValidation(t *testing.T) {
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
	_, err = db.Exec(`
		CREATE TABLE products (
			id INTEGER PRIMARY KEY,
			name TEXT,
			price DOUBLE PRECISION,
			in_stock BOOLEAN,
			categories TEXT[]
		)
	`)
	require.NoError(t, err)

	// Insert test data
	_, err = db.Exec(`
		INSERT INTO products VALUES
		(1, 'Laptop', 999.99, true, ARRAY['electronics', 'computers']),
		(2, 'Mouse', 29.99, true, ARRAY['electronics', 'accessories']),
		(3, 'Desk', 299.99, false, ARRAY['furniture']),
		(4, 'Chair', 199.99, true, ARRAY['furniture'])
	`)
	require.NoError(t, err)

	// Set up CEL environment with simple variables (no struct wrapper)
	env, err := cel.NewEnv(
		cel.Variable("id", cel.IntType),
		cel.Variable("name", cel.StringType),
		cel.Variable("price", cel.DoubleType),
		cel.Variable("in_stock", cel.BoolType),
		cel.Variable("categories", cel.ListType(cel.StringType)),
	)
	require.NoError(t, err)

	tests := []struct {
		name           string
		celExpr        string
		expectedRows   int
		description    string
	}{
		{
			name:           "Filter by price range",
			celExpr:        `price > 50.0 && price < 500.0`,
			expectedRows:   2, // Mouse and Chair
			description:    "Price range filtering should return correct count",
		},
		{
			name:           "Filter by stock status",
			celExpr:        `in_stock == true`,
			expectedRows:   3, // Laptop, Mouse, Chair
			description:    "Boolean filtering should work correctly",
		},
		{
			name:           "Array membership",
			celExpr:        `"electronics" in categories`,
			expectedRows:   2, // Laptop, Mouse
			description:    "Array membership check should return correct results",
		},
		{
			name:           "Complex condition",
			celExpr:        `in_stock && price < 300.0`,
			expectedRows:   2, // Mouse, Chair
			description:    "Complex conditions should execute correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compile CEL expression
			ast, issues := env.Compile(tt.celExpr)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("CEL compilation failed: %v", issues.Err())
			}

			// Convert to SQL
			sqlCondition, err := cel2sql.Convert(ast)
			require.NoError(t, err)

			t.Logf("CEL Expression: %s", tt.celExpr)
			t.Logf("Generated SQL WHERE clause: %s", sqlCondition)

			// Run EXPLAIN ANALYZE to validate execution
			// #nosec G201,G202 - This is a test validating SQL generation, not a security risk
			explainAnalyzeQuery := "EXPLAIN ANALYZE SELECT * FROM products WHERE " + sqlCondition
			t.Logf("EXPLAIN ANALYZE Query: %s", explainAnalyzeQuery)

			rows, err := db.Query(explainAnalyzeQuery)
			require.NoError(t, err, "EXPLAIN ANALYZE should not fail for valid SQL")
			defer func() {
				if closeErr := rows.Close(); closeErr != nil {
					t.Logf("failed to close rows: %v", closeErr)
				}
			}()

			// Collect plan lines
			var planLines []string
			rowsReturned := -1
			for rows.Next() {
				var planLine string
				err := rows.Scan(&planLine)
				require.NoError(t, err)
				planLines = append(planLines, planLine)
				t.Logf("  Plan: %s", planLine)

				// Extract actual rows from plan
				if strings.Contains(planLine, "rows=") {
					// Parse "rows=N" from plan line
					parts := strings.Split(planLine, "rows=")
					if len(parts) > 1 {
						rowsPart := strings.Fields(parts[1])[0]
						_, _ = fmt.Sscanf(rowsPart, "%d", &rowsReturned)
					}
				}
			}

			require.NotEmpty(t, planLines, "EXPLAIN ANALYZE should return query plan")

			// Verify the query actually executed
			fullPlan := strings.Join(planLines, "\n")
			require.Contains(t, fullPlan, "actual time=",
				"EXPLAIN ANALYZE should show actual execution time")

			t.Logf("✓ EXPLAIN ANALYZE validation passed: %s", tt.description)

			// Now actually run the query to verify row count
			// #nosec G201,G202 - This is a test validating SQL generation, not a security risk
			actualQuery := "SELECT COUNT(*) FROM products WHERE " + sqlCondition
			var actualRows int
			err = db.QueryRow(actualQuery).Scan(&actualRows)
			require.NoError(t, err)

			require.Equal(t, tt.expectedRows, actualRows,
				"Query should return expected number of rows. %s", tt.description)

			t.Logf("✓ Row count validation passed: expected %d, got %d", tt.expectedRows, actualRows)
		})
	}
}
