package pg_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
)

// TestPostgreSQL17Compatibility tests all the PostgreSQL 17 compatibility fixes
// with a real PostgreSQL 17 database to ensure all generated SQL is valid
func TestPostgreSQL17Compatibility(t *testing.T) {
	ctx := context.Background()

	// Create a PostgreSQL 17 container
	container, err := postgres.Run(ctx,
		"postgres:17",
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

	// Get connection string and create pool
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	defer pool.Close()

	// Create test table with various data types
	_, err = pool.Exec(ctx, `
		CREATE TABLE test_compat (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			age INTEGER NOT NULL,
			score DOUBLE PRECISION,
			active BOOLEAN NOT NULL,
			data BYTEA,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			tags TEXT[]
		)
	`)
	require.NoError(t, err)

	// Insert test data
	_, err = pool.Exec(ctx, `
		INSERT INTO test_compat (name, age, score, active, data, created_at, tags) VALUES
		('Alice', 25, 95.5, true, E'\\xDEADBEEF', '2024-01-15 10:30:00+00', ARRAY['admin', 'user']),
		('Bob', 30, 87.3, false, E'\\x1234', '2024-06-20 14:45:00+00', ARRAY['user']),
		('Charlie', 35, 92.1, true, E'\\xABCD', '2024-03-10 08:00:00+00', ARRAY['admin', 'moderator']),
		('David', 28, 78.9, false, E'\\xBEEF', '2024-09-05 16:15:00+00', ARRAY['user', 'guest'])
	`)
	require.NoError(t, err)

	// Define schema for CEL
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "name", Type: "text"},
		{Name: "age", Type: "integer"},
		{Name: "score", Type: "double precision"},
		{Name: "active", Type: "boolean"},
		{Name: "data", Type: "bytea"},
		{Name: "created_at", Type: "timestamp with time zone"},
		{Name: "tags", Type: "text", Repeated: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"test_compat": schema})

	tests := []struct {
		name          string
		celExpr       string
		expectedSQL   string
		expectedCount int
		description   string
	}{
		// CASE WHEN tests (ternary operator)
		{
			name:          "case_when_ternary_true",
			celExpr:       `(test_compat.age > 30 ? "senior" : "junior") == "senior"`,
			expectedSQL:   `(CASE WHEN test_compat.age > 30 THEN 'senior' ELSE 'junior' END) = 'senior'`,
			expectedCount: 1, // Charlie (35)
			description:   "Test CASE WHEN for ternary operator (true case)",
		},
		{
			name:          "case_when_ternary_false",
			celExpr:       `(test_compat.age > 30 ? "senior" : "junior") == "junior"`,
			expectedSQL:   `(CASE WHEN test_compat.age > 30 THEN 'senior' ELSE 'junior' END) = 'junior'`,
			expectedCount: 3, // Alice (25), Bob (30), David (28)
			description:   "Test CASE WHEN for ternary operator (false case)",
		},
		{
			name:          "case_when_nested",
			celExpr:       `(test_compat.active ? (test_compat.age > 30 ? "senior_active" : "junior_active") : "inactive") == "junior_active"`,
			expectedSQL:   `(CASE WHEN test_compat.active THEN CASE WHEN test_compat.age > 30 THEN 'senior_active' ELSE 'junior_active' END ELSE 'inactive' END) = 'junior_active'`,
			expectedCount: 1, // Alice (25, active)
			description:   "Test nested CASE WHEN",
		},

		// Type casting tests
		{
			name:          "cast_to_boolean",
			celExpr:       `bool(test_compat.age > 30) == true`,
			expectedSQL:   `CAST(test_compat.age > 30 AS BOOLEAN) IS TRUE`,
			expectedCount: 1, // Charlie (35)
			description:   "Test CAST to BOOLEAN",
		},
		{
			name:          "cast_to_bigint",
			celExpr:       `int(test_compat.score) >= 95`,
			expectedSQL:   `CAST(test_compat.score AS BIGINT) >= 95`,
			expectedCount: 1, // Alice (95.5 casts to 95, but we use >= to catch it)
			description:   "Test CAST to BIGINT",
		},
		{
			name:          "cast_to_text",
			celExpr:       `string(test_compat.age) == "30"`,
			expectedSQL:   `CAST(test_compat.age AS TEXT) = '30'`,
			expectedCount: 1, // Bob
			description:   "Test CAST to TEXT",
		},

		// EXTRACT EPOCH tests
		{
			name:          "extract_epoch",
			celExpr:       `int(test_compat.created_at) > 1704000000`,
			expectedSQL:   `EXTRACT(EPOCH FROM test_compat.created_at)::bigint > 1704000000`,
			expectedCount: 4, // All records are after 2024-01-01
			description:   "Test EXTRACT(EPOCH FROM timestamp) for Unix timestamp conversion",
		},

		// EXTRACT DOY and DOW tests
		{
			name:          "extract_doy",
			celExpr:       `test_compat.created_at.getDayOfYear() > 100`,
			expectedSQL:   `EXTRACT(DOY FROM test_compat.created_at) - 1 > 100`,
			expectedCount: 2, // Bob (June 20 = day 172), David (Sept 5 = day 249)
			description:   "Test EXTRACT(DOY ...) for day of year",
		},
		{
			name:          "extract_dow",
			celExpr:       `test_compat.created_at.getDayOfWeek() == 0`,
			expectedSQL:   `EXTRACT(DOW FROM test_compat.created_at) - 1 = 0`,
			expectedCount: 1, // Monday (2024-01-15) - CEL uses 0=Monday
			description:   "Test EXTRACT(DOW ...) for day of week",
		},

		// AT TIME ZONE test - simplified to just test the conversion happens
		{
			name:          "at_time_zone",
			celExpr:       `test_compat.created_at > timestamp("2024-01-01T00:00:00Z")`,
			expectedSQL:   `test_compat.created_at > CAST('2024-01-01T00:00:00Z' AS TIMESTAMP WITH TIME ZONE)`,
			expectedCount: 4, // All records after 2024-01-01
			description:   "Test timestamp() function converts to CAST AS TIMESTAMP WITH TIME ZONE",
		},

		// Complex combinations
		{
			name:          "complex_case_when_with_cast",
			celExpr:       `(test_compat.active ? string(test_compat.age) : "0") == "25"`,
			expectedSQL:   `(CASE WHEN test_compat.active THEN CAST(test_compat.age AS TEXT) ELSE '0' END) = '25'`,
			expectedCount: 1, // Alice
			description:   "Test complex expression combining CASE WHEN and type casting",
		},
		{
			name:          "epoch_with_conditional",
			celExpr:       `(int(test_compat.created_at) > 1720000000 ? "recent" : "old") == "recent"`,
			expectedSQL:   `(CASE WHEN EXTRACT(EPOCH FROM test_compat.created_at)::bigint > 1720000000 THEN 'recent' ELSE 'old' END) = 'recent'`,
			expectedCount: 1, // David (Sept 2024)
			description:   "Test EXTRACT EPOCH combined with CASE WHEN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create CEL environment
			env, err := cel.NewEnv(
				cel.CustomTypeProvider(provider),
				cel.Variable("test_compat", cel.ObjectType("test_compat")),
			)
			require.NoError(t, err, "Failed to create CEL environment")

			// Compile CEL expression
			ast, issues := env.Compile(tt.celExpr)
			require.Nil(t, issues, "CEL compilation should succeed")

			// Convert to SQL
			sqlCondition, err := cel2sql.Convert(ast)
			require.NoError(t, err, "CEL to SQL conversion should succeed")

			t.Logf("CEL Expression: %s", tt.celExpr)
			t.Logf("Generated SQL: %s", sqlCondition)
			t.Logf("Expected SQL pattern: %s", tt.expectedSQL)
			t.Logf("Description: %s", tt.description)

			// Verify SQL pattern matches expected
			assert.Equal(t, tt.expectedSQL, sqlCondition, "Generated SQL should match expected pattern")

			// Execute query and verify results
			query := "SELECT COUNT(*) FROM test_compat WHERE " + sqlCondition
			t.Logf("Executing query: %s", query)

			var count int
			err = pool.QueryRow(ctx, query).Scan(&count)
			require.NoError(t, err, "Query should execute successfully against PostgreSQL 17")

			assert.Equal(t, tt.expectedCount, count, "Expected count should match actual results")

			// For debugging: show sample matching records
			if count > 0 && count <= 5 {
				sampleQuery := "SELECT id, name, age, active FROM test_compat WHERE " + sqlCondition + " LIMIT 3"
				rows, err := pool.Query(ctx, sampleQuery)
				require.NoError(t, err)
				defer rows.Close()

				t.Logf("Sample matching records:")
				for rows.Next() {
					var id, age int
					var name string
					var active bool
					err := rows.Scan(&id, &name, &age, &active)
					require.NoError(t, err)
					t.Logf("  - ID: %d, Name: %s, Age: %d, Active: %v", id, name, age, active)
				}
			}
		})
	}
}
