package cel2sql_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/cel-go/cel"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/spandigital/cel2sql/v2"
	"github.com/spandigital/cel2sql/v2/pg"
)

// setupPostgresContainer starts a PostgreSQL 17 container for testing
func setupPostgresContainer(ctx context.Context, t *testing.T) (testcontainers.Container, *sql.DB) {
	t.Helper()

	// Start PostgreSQL 17 container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:17-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err, "Failed to start PostgreSQL container")

	// Get connection details
	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Connect to database
	connStr := fmt.Sprintf("host=%s port=%s user=postgres password=password dbname=testdb sslmode=disable",
		host, port.Port())
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err, "Failed to connect to database")

	// Wait for database to be fully ready
	err = db.Ping()
	require.NoError(t, err, "Failed to ping database")

	return container, db
}

func TestParameterizedQueriesIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	container, db := setupPostgresContainer(ctx, t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close db: %v", err)
		}
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	// Create test table
	_, err := db.Exec(`
		CREATE TABLE users (
			id BIGINT PRIMARY KEY,
			name TEXT NOT NULL,
			age INT NOT NULL,
			salary DOUBLE PRECISION,
			active BOOLEAN DEFAULT TRUE,
			metadata JSONB
		)
	`)
	require.NoError(t, err, "Failed to create users table")

	// Insert test data
	_, err = db.Exec(`
		INSERT INTO users (id, name, age, salary, active, metadata) VALUES
		(1, 'John', 25, 50000.00, true, '{"role": "developer", "level": "senior"}'),
		(2, 'Jane', 30, 75000.00, true, '{"role": "manager", "level": "lead"}'),
		(3, 'Bob', 22, 40000.00, false, '{"role": "intern", "level": "junior"}'),
		(4, 'Alice', 35, 90000.00, true, '{"role": "architect", "level": "principal"}')
	`)
	require.NoError(t, err, "Failed to insert test data")

	// Set up CEL environment
	testSchema := pg.Schema{
		{Name: "id", Type: "bigint"},
		{Name: "name", Type: "text"},
		{Name: "age", Type: "integer"},
		{Name: "salary", Type: "double precision"},
		{Name: "active", Type: "boolean"},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
	}

	provider := pg.NewTypeProvider(map[string]pg.Schema{
		"users": testSchema,
	})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("usr", cel.ObjectType("users")),
	)
	require.NoError(t, err)

	schemas := map[string]pg.Schema{
		"usr": testSchema,
	}

	tests := []struct {
		name           string
		celExpr        string
		expectedRows   int
		expectedIDs    []int64
		description    string
	}{
		{
			name:         "simple string equality",
			celExpr:      `usr.name == "John"`,
			expectedRows: 1,
			expectedIDs:  []int64{1},
			description:  "Find user by exact name match",
		},
		{
			name:         "integer comparison",
			celExpr:      `usr.age > 25`,
			expectedRows: 2,
			expectedIDs:  []int64{2, 4},
			description:  "Find users older than 25",
		},
		{
			name:         "double comparison",
			celExpr:      `usr.salary >= 75000.0`,
			expectedRows: 2,
			expectedIDs:  []int64{2, 4},
			description:  "Find users with high salary",
		},
		{
			name:         "boolean constant (inline)",
			celExpr:      `usr.active == true`,
			expectedRows: 3,
			expectedIDs:  []int64{1, 2, 4},
			description:  "Find active users",
		},
		{
			name:         "complex AND expression",
			celExpr:      `usr.age >= 25 && usr.salary > 50000.0 && usr.active == true`,
			expectedRows: 2,
			expectedIDs:  []int64{2, 4},
			description:  "Find senior active users with good salary",
		},
		{
			name:         "complex OR expression",
			celExpr:      `usr.name == "John" || usr.name == "Alice"`,
			expectedRows: 2,
			expectedIDs:  []int64{1, 4},
			description:  "Find specific users by name",
		},
		{
			name:         "JSON field comparison",
			celExpr:      `usr.metadata.role == "developer"`,
			expectedRows: 1,
			expectedIDs:  []int64{1},
			description:  "Find users by JSON field value",
		},
		{
			name:         "nested JSON field",
			celExpr:      `usr.metadata.level == "senior"`,
			expectedRows: 1,
			expectedIDs:  []int64{1},
			description:  "Find users by nested JSON field",
		},
		{
			name:         "mixed regular and JSON fields",
			celExpr:      `usr.age > 30 && usr.metadata.role == "architect"`,
			expectedRows: 1,
			expectedIDs:  []int64{4},
			description:  "Combine regular and JSON field filters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing: %s", tt.description)

			// Compile CEL expression
			ast, issues := env.Compile(tt.celExpr)
			require.Nil(t, issues, "CEL compilation should succeed")

			// Convert to parameterized SQL
			result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err, "ConvertParameterized should succeed")

			t.Logf("CEL: %s", tt.celExpr)
			t.Logf("SQL: %s", result.SQL)
			t.Logf("Parameters: %v", result.Parameters)

			// Execute query with parameters
			// Use table alias to match CEL variable name
			// Note: result.SQL is generated by cel2sql.ConvertParameterized(), not from user input
			// #nosec G201 - SQL string is from trusted conversion function
			query := fmt.Sprintf("SELECT id FROM users usr WHERE %s ORDER BY id", result.SQL)
			rows, err := db.Query(query, result.Parameters...)
			require.NoError(t, err, "Query execution should succeed")
			defer func() {
				if closeErr := rows.Close(); closeErr != nil {
					t.Logf("failed to close rows: %v", closeErr)
				}
			}()

			// Collect results
			var ids []int64
			for rows.Next() {
				var id int64
				err := rows.Scan(&id)
				require.NoError(t, err)
				ids = append(ids, id)
			}
			require.NoError(t, rows.Err())

			// Assert results
			assert.Len(t, ids, tt.expectedRows, "Should return expected number of rows")
			if tt.expectedIDs != nil {
				assert.Equal(t, tt.expectedIDs, ids, "Should return expected IDs")
			}
		})
	}
}

func TestParameterizedQueryPlanCaching(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	container, db := setupPostgresContainer(ctx, t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close db: %v", err)
		}
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	// Create test table
	_, err := db.Exec(`
		CREATE TABLE products (
			id BIGINT PRIMARY KEY,
			name TEXT NOT NULL,
			price DOUBLE PRECISION NOT NULL
		)
	`)
	require.NoError(t, err)

	// Insert test data
	_, err = db.Exec(`
		INSERT INTO products (id, name, price) VALUES
		(1, 'Widget', 19.99),
		(2, 'Gadget', 29.99),
		(3, 'Doohickey', 39.99)
	`)
	require.NoError(t, err)

	// Set up CEL environment
	testSchema := pg.Schema{
		{Name: "id", Type: "bigint"},
		{Name: "name", Type: "text"},
		{Name: "price", Type: "double precision"},
	}

	provider := pg.NewTypeProvider(map[string]pg.Schema{
		"products": testSchema,
	})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("product", cel.ObjectType("products")),
	)
	require.NoError(t, err)

	schemas := map[string]pg.Schema{
		"product": testSchema,
	}

	// Compile CEL expression
	ast, issues := env.Compile(`product.price > 20.0`)
	require.Nil(t, issues)

	// Convert to parameterized SQL
	result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithSchemas(schemas))
	require.NoError(t, err)

	// Test 1: Execute with parameter 20.0
	// Note: result.SQL is generated by cel2sql.ConvertParameterized(), not from user input
	// #nosec G201 - SQL string is from trusted conversion function
	query := fmt.Sprintf("SELECT id FROM products product WHERE %s ORDER BY id", result.SQL)
	t.Logf("Query: %s", query)
	t.Logf("Parameters: %v", result.Parameters)

	rows1, err := db.Query(query, result.Parameters...)
	require.NoError(t, err)
	defer func() {
		if closeErr := rows1.Close(); closeErr != nil {
			t.Logf("failed to close rows: %v", closeErr)
		}
	}()

	var ids1 []int64
	for rows1.Next() {
		var id int64
		require.NoError(t, rows1.Scan(&id))
		ids1 = append(ids1, id)
	}
	require.NoError(t, rows1.Err())
	assert.Equal(t, []int64{2, 3}, ids1, "First query should return products 2 and 3")

	// Test 2: Execute same SQL structure with different parameter value
	// This demonstrates query plan caching - same SQL, different parameters
	rows2, err := db.Query(query, 30.0) // Different parameter value
	require.NoError(t, err)
	defer func() {
		if closeErr := rows2.Close(); closeErr != nil {
			t.Logf("failed to close rows: %v", closeErr)
		}
	}()

	var ids2 []int64
	for rows2.Next() {
		var id int64
		require.NoError(t, rows2.Scan(&id))
		ids2 = append(ids2, id)
	}
	require.NoError(t, rows2.Err())
	assert.Equal(t, []int64{3}, ids2, "Second query should return only product 3")

	t.Log("Successfully demonstrated query plan caching with parameterized queries")
}

func TestParameterizedPreparedStatements(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	container, db := setupPostgresContainer(ctx, t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close db: %v", err)
		}
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	// Create test table
	_, err := db.Exec(`
		CREATE TABLE employees (
			id BIGINT PRIMARY KEY,
			name TEXT NOT NULL,
			department TEXT NOT NULL,
			salary DOUBLE PRECISION NOT NULL
		)
	`)
	require.NoError(t, err)

	// Insert test data
	_, err = db.Exec(`
		INSERT INTO employees (id, name, department, salary) VALUES
		(1, 'Alice', 'Engineering', 80000.00),
		(2, 'Bob', 'Engineering', 75000.00),
		(3, 'Carol', 'Sales', 70000.00),
		(4, 'Dave', 'Engineering', 85000.00)
	`)
	require.NoError(t, err)

	// Set up CEL environment
	testSchema := pg.Schema{
		{Name: "id", Type: "bigint"},
		{Name: "name", Type: "text"},
		{Name: "department", Type: "text"},
		{Name: "salary", Type: "double precision"},
	}

	provider := pg.NewTypeProvider(map[string]pg.Schema{
		"employees": testSchema,
	})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("emp", cel.ObjectType("employees")),
	)
	require.NoError(t, err)

	schemas := map[string]pg.Schema{
		"emp": testSchema,
	}

	// Compile CEL expression
	ast, issues := env.Compile(`emp.department == "Engineering" && emp.salary > 75000.0`)
	require.Nil(t, issues)

	// Convert to parameterized SQL
	result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithSchemas(schemas))
	require.NoError(t, err)

	t.Logf("SQL: %s", result.SQL)
	t.Logf("Parameters: %v", result.Parameters)

	// Prepare statement
	// Note: result.SQL is generated by cel2sql.ConvertParameterized(), not from user input
	// #nosec G201 - SQL string is from trusted conversion function
	query := fmt.Sprintf("SELECT id, name, salary FROM employees emp WHERE %s ORDER BY id", result.SQL)
	stmt, err := db.Prepare(query)
	require.NoError(t, err, "Should prepare statement successfully")
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			t.Logf("failed to close statement: %v", closeErr)
		}
	}()

	// Execute prepared statement with parameters
	rows, err := stmt.Query(result.Parameters...)
	require.NoError(t, err, "Should execute prepared statement successfully")
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			t.Logf("failed to close rows: %v", closeErr)
		}
	}()

	// Collect results
	type employee struct {
		id     int64
		name   string
		salary float64
	}
	var employees []employee
	for rows.Next() {
		var emp employee
		err := rows.Scan(&emp.id, &emp.name, &emp.salary)
		require.NoError(t, err)
		employees = append(employees, emp)
	}
	require.NoError(t, rows.Err())

	// Assert results
	assert.Len(t, employees, 2, "Should return 2 employees")
	assert.Equal(t, int64(1), employees[0].id, "First employee should be Alice")
	assert.Equal(t, "Alice", employees[0].name)
	assert.Equal(t, int64(4), employees[1].id, "Second employee should be Dave")
	assert.Equal(t, "Dave", employees[1].name)

	t.Log("Successfully used prepared statements with parameterized queries")
}

func TestParameterizedSQLInjectionPrevention(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	container, db := setupPostgresContainer(ctx, t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close db: %v", err)
		}
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	// Create test table
	_, err := db.Exec(`
		CREATE TABLE accounts (
			id BIGINT PRIMARY KEY,
			username TEXT NOT NULL,
			balance DOUBLE PRECISION NOT NULL
		)
	`)
	require.NoError(t, err)

	// Insert test data
	_, err = db.Exec(`
		INSERT INTO accounts (id, username, balance) VALUES
		(1, 'alice', 1000.00),
		(2, 'bob', 2000.00)
	`)
	require.NoError(t, err)

	// Set up CEL environment
	testSchema := pg.Schema{
		{Name: "id", Type: "bigint"},
		{Name: "username", Type: "text"},
		{Name: "balance", Type: "double precision"},
	}

	provider := pg.NewTypeProvider(map[string]pg.Schema{
		"accounts": testSchema,
	})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("account", cel.ObjectType("accounts")),
	)
	require.NoError(t, err)

	schemas := map[string]pg.Schema{
		"account": testSchema,
	}

	// This test verifies that potential SQL injection strings are safely parameterized
	// The malicious input attempts SQL injection but should be treated as a literal string
	maliciousUsername := "alice' OR '1'='1"

	// Create a CEL expression that would be vulnerable if not parameterized
	// With parameterization, the malicious string is safely passed as a parameter
	celExpr := fmt.Sprintf(`account.username == "%s"`, maliciousUsername)

	// Compile CEL expression
	ast, issues := env.Compile(celExpr)
	require.Nil(t, issues, "CEL compilation should succeed")

	// Convert to parameterized SQL
	result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithSchemas(schemas))
	require.NoError(t, err)

	t.Logf("CEL: %s", celExpr)
	t.Logf("SQL: %s", result.SQL)
	t.Logf("Parameters: %v", result.Parameters)

	// Execute query - should return 0 rows because the username doesn't match exactly
	// Note: result.SQL is generated by cel2sql.ConvertParameterized(), not from user input
	//nolint:perfsprint // fmt.Sprintf preferred over concatenation for SQL security
	// #nosec G201 - SQL string is from trusted conversion function
	query := fmt.Sprintf("SELECT id FROM accounts account WHERE %s", result.SQL)
	rows, err := db.Query(query, result.Parameters...)
	require.NoError(t, err, "Query should execute successfully without SQL injection")
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			t.Logf("failed to close rows: %v", closeErr)
		}
	}()

	// Count rows
	var count int
	for rows.Next() {
		count++
	}
	require.NoError(t, rows.Err())

	// Assert that no rows were returned (SQL injection was prevented)
	assert.Equal(t, 0, count, "Should return 0 rows - SQL injection prevented")

	t.Log("Successfully prevented SQL injection using parameterized queries")
}

func TestParameterizedDataTypeCompatibility(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	container, db := setupPostgresContainer(ctx, t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close db: %v", err)
		}
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	// Create test table with various data types
	_, err := db.Exec(`
		CREATE TABLE datatypes (
			id BIGINT PRIMARY KEY,
			text_col TEXT,
			int_col BIGINT,
			double_col DOUBLE PRECISION,
			bool_col BOOLEAN,
			bytes_col BYTEA
		)
	`)
	require.NoError(t, err)

	// Insert test data
	_, err = db.Exec(`
		INSERT INTO datatypes (id, text_col, int_col, double_col, bool_col, bytes_col) VALUES
		(1, 'test', 42, 3.14, true, '\xDEADBEEF'::bytea),
		(2, 'hello', 100, 2.718, false, '\xCAFEBABE'::bytea)
	`)
	require.NoError(t, err)

	// Set up CEL environment
	testSchema := pg.Schema{
		{Name: "id", Type: "bigint"},
		{Name: "text_col", Type: "text"},
		{Name: "int_col", Type: "bigint"},
		{Name: "double_col", Type: "double precision"},
		{Name: "bool_col", Type: "boolean"},
		{Name: "bytes_col", Type: "bytea"},
	}

	provider := pg.NewTypeProvider(map[string]pg.Schema{
		"datatypes": testSchema,
	})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("row", cel.ObjectType("datatypes")),
	)
	require.NoError(t, err)

	schemas := map[string]pg.Schema{
		"row": testSchema,
	}

	tests := []struct {
		name        string
		celExpr     string
		expectedIDs []int64
	}{
		{
			name:        "text parameter",
			celExpr:     `row.text_col == "test"`,
			expectedIDs: []int64{1},
		},
		{
			name:        "int parameter",
			celExpr:     `row.int_col == 42`,
			expectedIDs: []int64{1},
		},
		{
			name:        "double parameter",
			celExpr:     `row.double_col > 3.0`,
			expectedIDs: []int64{1},
		},
		{
			name:        "bool parameter (inline)",
			celExpr:     `row.bool_col == true`,
			expectedIDs: []int64{1},
		},
		{
			name:        "bytes parameter",
			celExpr:     `row.bytes_col == b"\xDE\xAD\xBE\xEF"`,
			expectedIDs: []int64{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compile CEL expression
			ast, issues := env.Compile(tt.celExpr)
			require.Nil(t, issues, "CEL compilation should succeed")

			// Convert to parameterized SQL
			result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err)

			t.Logf("SQL: %s", result.SQL)
			t.Logf("Parameters: %v", result.Parameters)

			// Execute query
			// Note: result.SQL is generated by cel2sql.ConvertParameterized(), not from user input
			// #nosec G201 - SQL string is from trusted conversion function
			query := fmt.Sprintf("SELECT id FROM datatypes row WHERE %s ORDER BY id", result.SQL)
			rows, err := db.Query(query, result.Parameters...)
			require.NoError(t, err, "Query should execute successfully")
			defer func() {
				if closeErr := rows.Close(); closeErr != nil {
					t.Logf("failed to close rows: %v", closeErr)
				}
			}()

			// Collect results
			var ids []int64
			for rows.Next() {
				var id int64
				require.NoError(t, rows.Scan(&id))
				ids = append(ids, id)
			}
			require.NoError(t, rows.Err())

			// Assert results
			assert.Equal(t, tt.expectedIDs, ids, "Should return expected IDs")
		})
	}
}

func BenchmarkParameterizedVsInline(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	ctx := context.Background()
	container, db := setupPostgresContainer(ctx, &testing.T{})
	defer func() {
		if err := db.Close(); err != nil {
			b.Logf("failed to close db: %v", err)
		}
		if err := container.Terminate(ctx); err != nil {
			b.Logf("failed to terminate container: %v", err)
		}
	}()

	// Create test table
	_, err := db.Exec(`
		CREATE TABLE benchmark_test (
			id BIGINT PRIMARY KEY,
			value INT NOT NULL
		)
	`)
	if err != nil {
		b.Fatal(err)
	}

	// Insert test data
	for i := 1; i <= 1000; i++ {
		_, err := db.Exec(`INSERT INTO benchmark_test (id, value) VALUES ($1, $2)`, i, i*10)
		if err != nil {
			b.Fatal(err)
		}
	}

	// Set up CEL environment
	testSchema := pg.Schema{
		{Name: "id", Type: "bigint"},
		{Name: "value", Type: "integer"},
	}

	provider := pg.NewTypeProvider(map[string]pg.Schema{
		"benchmark_test": testSchema,
	})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("row", cel.ObjectType("benchmark_test")),
	)
	if err != nil {
		b.Fatal(err)
	}

	schemas := map[string]pg.Schema{
		"row": testSchema,
	}

	// Compile CEL expression
	ast, issues := env.Compile(`row.value > 500`)
	if issues != nil && issues.Err() != nil {
		b.Fatal(issues.Err())
	}

	// Benchmark parameterized queries
	b.Run("Parameterized", func(b *testing.B) {
		result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithSchemas(schemas))
		if err != nil {
			b.Fatal(err)
		}

		// Note: result.SQL is generated by cel2sql.ConvertParameterized(), not from user input
		//nolint:perfsprint // fmt.Sprintf preferred over concatenation for SQL security
		// #nosec G201 - SQL string is from trusted conversion function
		query := fmt.Sprintf("SELECT id FROM benchmark_test row WHERE %s", result.SQL)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rows, err := db.Query(query, result.Parameters...)
			if err != nil {
				b.Fatal(err)
			}
			for rows.Next() {
				var id int64
				if err := rows.Scan(&id); err != nil {
					b.Fatal(err)
				}
			}
			if err := rows.Close(); err != nil {
				b.Fatal(err)
			}
		}
	})

	// Benchmark inline queries
	b.Run("Inline", func(b *testing.B) {
		sqlCondition, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
		if err != nil {
			b.Fatal(err)
		}

		// Note: sqlCondition is generated by cel2sql.Convert(), not from user input
		//nolint:perfsprint // fmt.Sprintf preferred over concatenation for SQL security
		// #nosec G201 - SQL string is from trusted conversion function
		query := fmt.Sprintf("SELECT id FROM benchmark_test row WHERE %s", sqlCondition)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rows, err := db.Query(query)
			if err != nil {
				b.Fatal(err)
			}
			for rows.Next() {
				var id int64
				if err := rows.Scan(&id); err != nil {
					b.Fatal(err)
				}
			}
			if err := rows.Close(); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestParameterizedComplexExpressions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	container, db := setupPostgresContainer(ctx, t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close db: %v", err)
		}
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	// Create test table
	_, err := db.Exec(`
		CREATE TABLE complex_test (
			id BIGINT PRIMARY KEY,
			name TEXT NOT NULL,
			age INT NOT NULL,
			score DOUBLE PRECISION NOT NULL,
			tags TEXT[] NOT NULL
		)
	`)
	require.NoError(t, err)

	// Insert test data
	_, err = db.Exec(`
		INSERT INTO complex_test (id, name, age, score, tags) VALUES
		(1, 'Alice', 25, 85.5, ARRAY['senior', 'golang']),
		(2, 'Bob', 30, 92.3, ARRAY['lead', 'python']),
		(3, 'Carol', 22, 78.9, ARRAY['junior', 'golang'])
	`)
	require.NoError(t, err)

	// Set up CEL environment
	testSchema := pg.Schema{
		{Name: "id", Type: "bigint"},
		{Name: "name", Type: "text"},
		{Name: "age", Type: "integer"},
		{Name: "score", Type: "double precision"},
		{Name: "tags", Type: "text", Repeated: true, ElementType: "text"},
	}

	provider := pg.NewTypeProvider(map[string]pg.Schema{
		"complex_test": testSchema,
	})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("row", cel.ObjectType("complex_test")),
	)
	require.NoError(t, err)

	schemas := map[string]pg.Schema{
		"row": testSchema,
	}

	tests := []struct {
		name        string
		celExpr     string
		expectedIDs []int64
		description string
	}{
		{
			name:        "ternary with parameters",
			celExpr:     `row.age > 25 ? row.score > 90.0 : row.score > 80.0`,
			expectedIDs: []int64{1, 2},
			description: "Conditional logic with parameterized comparisons",
		},
		{
			name:        "IN with parameterized list",
			celExpr:     `row.name in ["Alice", "Bob"]`,
			expectedIDs: []int64{1, 2},
			description: "IN operator with parameterized array",
		},
		{
			name:        "string contains with parameter",
			celExpr:     `row.name.contains("ob")`,
			expectedIDs: []int64{2},
			description: "String contains with parameterized substring",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing: %s", tt.description)

			// Compile CEL expression
			ast, issues := env.Compile(tt.celExpr)
			require.Nil(t, issues, "CEL compilation should succeed")

			// Convert to parameterized SQL
			result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithSchemas(schemas))
			require.NoError(t, err)

			t.Logf("CEL: %s", tt.celExpr)
			t.Logf("SQL: %s", result.SQL)
			t.Logf("Parameters: %v", result.Parameters)

			// Execute query
			// Note: result.SQL is generated by cel2sql.ConvertParameterized(), not from user input
			// #nosec G201 - SQL string is from trusted conversion function
			query := fmt.Sprintf("SELECT id FROM complex_test row WHERE %s ORDER BY id", result.SQL)
			rows, err := db.Query(query, result.Parameters...)
			require.NoError(t, err, "Query should execute successfully")
			defer func() {
				if closeErr := rows.Close(); closeErr != nil {
					t.Logf("failed to close rows: %v", closeErr)
				}
			}()

			// Collect results
			var ids []int64
			for rows.Next() {
				var id int64
				require.NoError(t, rows.Scan(&id))
				ids = append(ids, id)
			}
			require.NoError(t, rows.Err())

			// Assert results
			assert.Equal(t, tt.expectedIDs, ids, "Should return expected IDs")
		})
	}
}

func TestParameterizedQueryPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	container, db := setupPostgresContainer(ctx, t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close db: %v", err)
		}
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	// Create test table
	_, err := db.Exec(`
		CREATE TABLE perf_test (
			id BIGINT PRIMARY KEY,
			value INT NOT NULL
		)
	`)
	require.NoError(t, err)

	// Insert test data
	const rowCount = 10000
	for i := 1; i <= rowCount; i++ {
		_, err := db.Exec(`INSERT INTO perf_test (id, value) VALUES ($1, $2)`, i, i%100)
		require.NoError(t, err)
	}

	// Set up CEL environment
	testSchema := pg.Schema{
		{Name: "id", Type: "bigint"},
		{Name: "value", Type: "integer"},
	}

	provider := pg.NewTypeProvider(map[string]pg.Schema{
		"perf_test": testSchema,
	})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("row", cel.ObjectType("perf_test")),
	)
	require.NoError(t, err)

	schemas := map[string]pg.Schema{
		"row": testSchema,
	}

	// Compile CEL expression
	ast, issues := env.Compile(`row.value == 42`)
	require.Nil(t, issues)

	// Convert to parameterized SQL
	result, err := cel2sql.ConvertParameterized(ast, cel2sql.WithSchemas(schemas))
	require.NoError(t, err)

	// Note: result.SQL is generated by cel2sql.ConvertParameterized(), not from user input
	//nolint:perfsprint // fmt.Sprintf preferred over concatenation for SQL security
	// #nosec G201 - SQL string is from trusted conversion function
	query := fmt.Sprintf("SELECT COUNT(*) FROM perf_test row WHERE %s", result.SQL)

	// Execute multiple times to measure consistency
	const iterations = 100
	var totalDuration time.Duration

	for i := 0; i < iterations; i++ {
		start := time.Now()
		var count int
		err := db.QueryRow(query, result.Parameters...).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 100, count, "Should return expected count")
		totalDuration += time.Since(start)
	}

	avgDuration := totalDuration / iterations
	t.Logf("Average query time: %v over %d iterations", avgDuration, iterations)
	t.Logf("Total rows: %d, matching rows: %d", rowCount, 100)

	// Performance assertion - queries should be reasonably fast
	assert.Less(t, avgDuration.Milliseconds(), int64(100),
		"Average query time should be less than 100ms")
}
