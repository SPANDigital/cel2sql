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

	"github.com/spandigital/cel2sql/v3"
)

// TestOperatorsIntegration validates all operator conversions against real PostgreSQL
func TestOperatorsIntegration(t *testing.T) {
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

	// Create comprehensive test table with diverse data types
	_, err = db.Exec(`
		CREATE TABLE test_data (
			id INTEGER PRIMARY KEY,
			text_val TEXT,
			int_val INTEGER,
			float_val DOUBLE PRECISION,
			bool_val BOOLEAN,
			array_val TEXT[],
			nullable_text TEXT,
			nullable_int INTEGER,
			timestamp_val TIMESTAMP
		)
	`)
	require.NoError(t, err)

	// Insert test data with variety of values
	_, err = db.Exec(`
		INSERT INTO test_data VALUES
		(1, 'hello', 10, 10.5, true, ARRAY['a', 'b', 'c'], 'present', 100, '2024-01-01 10:00:00'),
		(2, 'world', 20, 20.5, false, ARRAY['x', 'y'], NULL, NULL, '2024-01-02 11:00:00'),
		(3, 'test', 30, 30.5, true, ARRAY['p', 'q', 'r'], 'here', 200, '2024-01-03 12:00:00'),
		(4, 'hello world', 5, 5.5, false, ARRAY['a'], 'value', 50, '2024-01-04 13:00:00'),
		(5, 'testing', 15, 15.5, true, ARRAY['b', 'c'], 'test', 150, '2024-01-05 14:00:00')
	`)
	require.NoError(t, err)

	// Set up CEL environment with simple variables (no struct wrapper)
	env, err := cel.NewEnv(
		cel.Variable("id", cel.IntType),
		cel.Variable("text_val", cel.StringType),
		cel.Variable("int_val", cel.IntType),
		cel.Variable("float_val", cel.DoubleType),
		cel.Variable("bool_val", cel.BoolType),
		cel.Variable("array_val", cel.ListType(cel.StringType)),
		cel.Variable("nullable_text", cel.StringType),
		cel.Variable("nullable_int", cel.IntType),
		cel.Variable("timestamp_val", cel.TimestampType),
	)
	require.NoError(t, err)

	// Test cases organized by operator type
	tests := []struct {
		name         string
		celExpr      string
		expectedRows int
		description  string
	}{
		// Comparison operators
		{
			name:         "Equality string",
			celExpr:      `text_val == "hello"`,
			expectedRows: 1,
			description:  "String equality comparison",
		},
		{
			name:         "Equality integer",
			celExpr:      `int_val == 20`,
			expectedRows: 1,
			description:  "Integer equality comparison",
		},
		{
			name:         "Equality float",
			celExpr:      `float_val == 10.5`,
			expectedRows: 1,
			description:  "Float equality comparison",
		},
		{
			name:         "Equality boolean",
			celExpr:      `bool_val == true`,
			expectedRows: 3,
			description:  "Boolean equality comparison",
		},
		{
			name:         "Not equal",
			celExpr:      `text_val != "hello"`,
			expectedRows: 4,
			description:  "Not equal comparison",
		},
		{
			name:         "Less than",
			celExpr:      `int_val < 15`,
			expectedRows: 2, // 10, 5
			description:  "Less than comparison",
		},
		{
			name:         "Less than or equal",
			celExpr:      `int_val <= 15`,
			expectedRows: 3, // 10, 5, 15
			description:  "Less than or equal comparison",
		},
		{
			name:         "Greater than",
			celExpr:      `int_val > 15`,
			expectedRows: 2, // 20, 30
			description:  "Greater than comparison",
		},
		{
			name:         "Greater than or equal",
			celExpr:      `int_val >= 15`,
			expectedRows: 3, // 20, 30, 15
			description:  "Greater than or equal comparison",
		},

		// Logical operators
		{
			name:         "Logical AND",
			celExpr:      `int_val > 10 && bool_val == true`,
			expectedRows: 2, // int_val > 10: 20, 30, 15; bool_val == true: 10, 30, 15; intersection: 30, 15
			description:  "Logical AND operator",
		},
		{
			name:         "Logical OR",
			celExpr:      `int_val < 10 || bool_val == false`,
			expectedRows: 2, // int_val < 10: 5; bool_val == false: 20, 5; union: 20, 5
			description:  "Logical OR operator",
		},
		{
			name:         "Logical NOT",
			celExpr:      `!bool_val`,
			expectedRows: 2, // bool_val == false: 20, 5
			description:  "Logical NOT operator",
		},
		{
			name:         "Complex logical expression",
			celExpr:      `(int_val > 10 && bool_val) || int_val < 10`,
			expectedRows: 3, // (int_val > 10 && bool_val): 30, 15; int_val < 10: 5; union: 30, 15, 5
			description:  "Complex nested logical operators",
		},

		// Arithmetic operators
		{
			name:         "Addition",
			celExpr:      `int_val + 10 == 20`,
			expectedRows: 1, // 10 + 10 = 20
			description:  "Addition operator",
		},
		{
			name:         "Subtraction",
			celExpr:      `int_val - 5 == 15`,
			expectedRows: 1, // 20 - 5 = 15
			description:  "Subtraction operator",
		},
		{
			name:         "Multiplication",
			celExpr:      `int_val * 2 == 20`,
			expectedRows: 1, // 10 * 2 = 20
			description:  "Multiplication operator",
		},
		{
			name:         "Division",
			celExpr:      `int_val / 2 == 10`,
			expectedRows: 1, // 20 / 2 = 10
			description:  "Division operator",
		},
		{
			name:         "Modulo",
			celExpr:      `int_val % 10 == 0`,
			expectedRows: 3, // 10, 20, 30
			description:  "Modulo operator",
		},
		{
			name:         "Complex arithmetic",
			celExpr:      `(int_val * 2) + 5 > 30`,
			expectedRows: 3, // (10*2)+5=25, (20*2)+5=45, (30*2)+5=65, (5*2)+5=15, (15*2)+5=35; > 30: 45, 65, 35
			description:  "Complex arithmetic expression",
		},

		// String operators
		{
			name:         "String concatenation",
			celExpr:      `text_val + "!" == "hello!"`,
			expectedRows: 1,
			description:  "String concatenation operator",
		},
		{
			name:         "String contains",
			celExpr:      `text_val.contains("world")`,
			expectedRows: 2, // "world", "hello world"
			description:  "String contains function",
		},
		{
			name:         "String startsWith",
			celExpr:      `text_val.startsWith("hello")`,
			expectedRows: 2, // "hello", "hello world"
			description:  "String startsWith function",
		},
		{
			name:         "String endsWith",
			celExpr:      `text_val.endsWith("world")`,
			expectedRows: 2, // "world", "hello world"
			description:  "String endsWith function",
		},
		{
			name:         "String matches regex",
			celExpr:      `text_val.matches(r"^hello")`,
			expectedRows: 2, // "hello", "hello world"
			description:  "String regex matching",
		},

		// Array/membership operators
		{
			name:         "Array membership - in",
			celExpr:      `"a" in array_val`,
			expectedRows: 2, // arrays containing 'a'
			description:  "Array membership with 'in' operator",
		},
		{
			name:         "Array membership - multiple values",
			celExpr:      `"b" in array_val || "x" in array_val`,
			expectedRows: 3, // arrays containing 'b' or 'x'
			description:  "Array membership with OR",
		},

		// Timestamp operators
		{
			name:         "Timestamp equality",
			celExpr:      `timestamp_val == timestamp("2024-01-01T10:00:00Z")`,
			expectedRows: 1,
			description:  "Timestamp equality comparison",
		},
		{
			name:         "Timestamp greater than",
			celExpr:      `timestamp_val > timestamp("2024-01-02T00:00:00Z")`,
			expectedRows: 4, // all except first row
			description:  "Timestamp greater than comparison",
		},
		{
			name:         "Timestamp less than",
			celExpr:      `timestamp_val < timestamp("2024-01-03T00:00:00Z")`,
			expectedRows: 2, // first two rows
			description:  "Timestamp less than comparison",
		},

		// Complex combined operators
		{
			name:         "Complex multi-operator expression",
			celExpr:      `int_val > 10 && bool_val && "b" in array_val`,
			expectedRows: 1, // Only row 5: int_val=15 > 10, bool_val=true, 'b' in ['b','c']
			description:  "Complex expression with multiple operator types",
		},
		{
			name:         "Nested parenthesized operators",
			celExpr:      `((int_val + 5) * 2 > 30) && (text_val.contains("test") || bool_val)`,
			expectedRows: 2, // Rows 3 and 5: arithmetic > 30 AND (contains "test" OR bool_val)
			description:  "Deeply nested operators with parentheses",
		},
		{
			name:         "Triple negation",
			celExpr:      `!!!bool_val`,
			expectedRows: 2, // !!!false = false, !!!true = true; so bool_val == false
			description:  "Multiple NOT operators",
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
			require.NoError(t, err, "Conversion should succeed for: %s", tt.description)

			t.Logf("CEL Expression: %s", tt.celExpr)
			t.Logf("Generated SQL WHERE clause: %s", sqlCondition)

			// Execute query and count results
			// #nosec G202 - This is a test validating SQL generation, not a security risk
			query := "SELECT COUNT(*) FROM test_data WHERE " + sqlCondition
			t.Logf("Full SQL Query: %s", query)

			var actualRows int
			err = db.QueryRow(query).Scan(&actualRows)
			require.NoError(t, err, "Generated SQL should execute successfully. %s", tt.description)

			require.Equal(t, tt.expectedRows, actualRows,
				"Query should return expected number of rows. %s\nCEL: %s\nSQL: %s",
				tt.description, tt.celExpr, sqlCondition)

			t.Logf("✓ Operator validation passed: %s (expected %d rows, got %d rows)",
				tt.description, tt.expectedRows, actualRows)
		})
	}
}

// TestOperatorEdgeCases validates edge cases and boundary conditions
func TestOperatorEdgeCases(t *testing.T) {
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

	// Create test table with edge case values
	_, err = db.Exec(`
		CREATE TABLE edge_cases (
			id INTEGER PRIMARY KEY,
			empty_string TEXT,
			zero_int INTEGER,
			zero_float DOUBLE PRECISION,
			negative_int INTEGER,
			negative_float DOUBLE PRECISION,
			large_int BIGINT,
			empty_array TEXT[]
		)
	`)
	require.NoError(t, err)

	// Insert edge case data
	_, err = db.Exec(`
		INSERT INTO edge_cases VALUES
		(1, '', 0, 0.0, -10, -5.5, 9223372036854775807, ARRAY[]::TEXT[]),
		(2, 'value', 1, 1.0, -1, -0.1, 123456789, ARRAY['item']),
		(3, 'another', -0, 0.0, 0, 0.0, 0, ARRAY['a', 'b'])
	`)
	require.NoError(t, err)

	// Set up CEL environment with simple variables (no struct wrapper)
	env, err := cel.NewEnv(
		cel.Variable("id", cel.IntType),
		cel.Variable("empty_string", cel.StringType),
		cel.Variable("zero_int", cel.IntType),
		cel.Variable("zero_float", cel.DoubleType),
		cel.Variable("negative_int", cel.IntType),
		cel.Variable("negative_float", cel.DoubleType),
		cel.Variable("large_int", cel.IntType),
		cel.Variable("empty_array", cel.ListType(cel.StringType)),
	)
	require.NoError(t, err)

	tests := []struct {
		name         string
		celExpr      string
		expectedRows int
		description  string
	}{
		{
			name:         "Empty string equality",
			celExpr:      `empty_string == ""`,
			expectedRows: 1,
			description:  "Empty string should be handled correctly",
		},
		{
			name:         "Zero integer equality",
			celExpr:      `zero_int == 0`,
			expectedRows: 2, // rows with 0 and -0
			description:  "Zero should be handled correctly",
		},
		{
			name:         "Negative integer comparison",
			celExpr:      `negative_int < 0`,
			expectedRows: 2, // -10 and -1
			description:  "Negative numbers should work correctly",
		},
		{
			name:         "Large integer comparison",
			celExpr:      `large_int > 1000000`,
			expectedRows: 2, // 9223372036854775807 and 123456789
			description:  "Large integers should be handled correctly",
		},
		{
			name:         "Zero float equality",
			celExpr:      `zero_float == 0.0`,
			expectedRows: 2, // rows with 0.0
			description:  "Zero float should be handled correctly",
		},
		{
			name:         "Negative float comparison",
			celExpr:      `negative_float < 0.0`,
			expectedRows: 2, // -5.5 and -0.1
			description:  "Negative floats should work correctly",
		},
		{
			name:         "Division by non-zero",
			celExpr:      `zero_int + 10 == 10`,
			expectedRows: 2, // 0 + 10 = 10
			description:  "Arithmetic with zero should work",
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
			require.NoError(t, err, "Conversion should succeed for: %s", tt.description)

			t.Logf("CEL Expression: %s", tt.celExpr)
			t.Logf("Generated SQL WHERE clause: %s", sqlCondition)

			// Execute query and count results
			// #nosec G202 - This is a test validating SQL generation, not a security risk
			query := "SELECT COUNT(*) FROM edge_cases WHERE " + sqlCondition

			var actualRows int
			err = db.QueryRow(query).Scan(&actualRows)
			require.NoError(t, err, "Generated SQL should execute successfully. %s", tt.description)

			require.Equal(t, tt.expectedRows, actualRows,
				"Query should return expected number of rows. %s", tt.description)

			t.Logf("✓ Edge case validation passed: %s", tt.description)
		})
	}
}
