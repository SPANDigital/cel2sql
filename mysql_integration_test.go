package cel2sql_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"

	"github.com/spandigital/cel2sql/v3"
	mysqlDialect "github.com/spandigital/cel2sql/v3/dialect/mysql"
	"github.com/spandigital/cel2sql/v3/pg"
)

// setupMySQLContainer starts a MySQL 8 container and returns a database connection.
func setupMySQLContainer(ctx context.Context, t *testing.T) (testcontainers.Container, *sql.DB) {
	t.Helper()

	container, err := tcmysql.Run(ctx, "mysql:8.0",
		tcmysql.WithDatabase("testdb"),
		tcmysql.WithUsername("testuser"),
		tcmysql.WithPassword("testpass"),
	)
	require.NoError(t, err, "Failed to start MySQL container")

	// Get connection string
	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "3306")
	require.NoError(t, err)

	connStr := fmt.Sprintf("testuser:testpass@tcp(%s:%s)/testdb?parseTime=true",
		host, port.Port())
	db, err := sql.Open("mysql", connStr)
	require.NoError(t, err, "Failed to connect to MySQL database")

	err = db.Ping()
	require.NoError(t, err, "Failed to ping MySQL database")

	return container, db
}

// TestMySQLOperatorsIntegration validates operator conversions against a real MySQL database.
func TestMySQLOperatorsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	container, db := setupMySQLContainer(ctx, t)
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Logf("failed to close db: %v", closeErr)
		}
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	// Create test table
	_, err := db.Exec(`
		CREATE TABLE test_data (
			id INTEGER PRIMARY KEY,
			text_val TEXT,
			int_val INTEGER,
			float_val DOUBLE,
			bool_val BOOLEAN,
			nullable_text TEXT,
			nullable_int INTEGER
		)
	`)
	require.NoError(t, err)

	// Insert test data
	_, err = db.Exec(`
		INSERT INTO test_data VALUES
		(1, 'hello', 10, 10.5, true, 'present', 100),
		(2, 'world', 20, 20.5, false, NULL, NULL),
		(3, 'test', 30, 30.5, true, 'here', 200),
		(4, 'hello world', 5, 5.5, false, 'value', 50),
		(5, 'testing', 15, 15.5, true, 'test', 150)
	`)
	require.NoError(t, err)

	// Set up CEL environment with simple variables
	env, err := cel.NewEnv(
		cel.Variable("id", cel.IntType),
		cel.Variable("text_val", cel.StringType),
		cel.Variable("int_val", cel.IntType),
		cel.Variable("float_val", cel.DoubleType),
		cel.Variable("bool_val", cel.BoolType),
		cel.Variable("nullable_text", cel.StringType),
		cel.Variable("nullable_int", cel.IntType),
	)
	require.NoError(t, err)

	dialectOpt := cel2sql.WithDialect(mysqlDialect.New())

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
			expectedRows: 2, // rows 3 (30,true) and 5 (15,true)
			description:  "Logical AND operator",
		},
		{
			name:         "Logical OR",
			celExpr:      `int_val < 10 || bool_val == false`,
			expectedRows: 2, // rows 2 (20,false) and 4 (5,false)
			description:  "Logical OR operator",
		},
		{
			name:         "Logical NOT",
			celExpr:      `!bool_val`,
			expectedRows: 2, // rows 2 and 4
			description:  "Logical NOT operator",
		},
		{
			name:         "Complex logical expression",
			celExpr:      `(int_val > 10 && bool_val) || int_val < 10`,
			expectedRows: 3, // rows 3, 5, 4
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
			expectedRows: 3, // (20*2)+5=45, (30*2)+5=65, (15*2)+5=35
			description:  "Complex arithmetic expression",
		},

		// String operators
		{
			name:         "String concatenation",
			celExpr:      `text_val + "!" == "hello!"`,
			expectedRows: 1,
			description:  "String concatenation (CONCAT)",
		},
		{
			name:         "String contains",
			celExpr:      `text_val.contains("world")`,
			expectedRows: 2, // "world", "hello world"
			description:  "String contains function (LOCATE)",
		},
		{
			name:         "String startsWith",
			celExpr:      `text_val.startsWith("hello")`,
			expectedRows: 2, // "hello", "hello world"
			description:  "String startsWith function (LIKE)",
		},
		{
			name:         "String endsWith",
			celExpr:      `text_val.endsWith("world")`,
			expectedRows: 2, // "world", "hello world"
			description:  "String endsWith function (LIKE)",
		},

		// Regex (MySQL 8.0+ supports REGEXP)
		{
			name:         "Regex match",
			celExpr:      `text_val.matches(r"^hello")`,
			expectedRows: 2, // "hello", "hello world"
			description:  "Regex match (REGEXP)",
		},
		{
			name:         "Regex word boundary",
			celExpr:      `text_val.matches(r"test")`,
			expectedRows: 2, // "test", "testing"
			description:  "Regex simple pattern",
		},

		// Complex combined operators
		{
			name:         "Complex multi-operator expression",
			celExpr:      `int_val > 10 && bool_val && text_val.contains("test")`,
			expectedRows: 2, // rows 3 and 5
			description:  "Complex expression with multiple operator types",
		},
		{
			name:         "Nested parenthesized operators",
			celExpr:      `((int_val + 5) * 2 > 30) && (text_val.contains("test") || bool_val)`,
			expectedRows: 2, // rows 3 and 5
			description:  "Deeply nested operators with parentheses",
		},
		{
			name:         "Triple negation",
			celExpr:      `!!!bool_val`,
			expectedRows: 2,
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

			// Convert to SQL with MySQL dialect
			sqlCondition, err := cel2sql.Convert(ast, dialectOpt)
			require.NoError(t, err, "Conversion should succeed for: %s", tt.description)

			t.Logf("CEL Expression: %s", tt.celExpr)
			t.Logf("Generated SQL WHERE clause: %s", sqlCondition)

			// Execute query and count results
			// #nosec G202 - This is a test validating SQL generation, not a security risk
			query := "SELECT COUNT(*) FROM test_data WHERE " + sqlCondition
			t.Logf("Full SQL Query: %s", query)

			var actualRows int
			err = db.QueryRow(query).Scan(&actualRows)
			require.NoError(t, err, "Generated SQL should execute successfully. %s\nSQL: %s",
				tt.description, sqlCondition)

			require.Equal(t, tt.expectedRows, actualRows,
				"Query should return expected number of rows. %s\nCEL: %s\nSQL: %s",
				tt.description, tt.celExpr, sqlCondition)

			t.Logf("OK: %s (expected %d rows, got %d rows)",
				tt.description, tt.expectedRows, actualRows)
		})
	}
}

// TestMySQLJSONIntegration validates JSON operations against a real MySQL database.
func TestMySQLJSONIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	container, db := setupMySQLContainer(ctx, t)
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Logf("failed to close db: %v", closeErr)
		}
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	// Create test table with JSON column
	_, err := db.Exec(`
		CREATE TABLE products (
			id INTEGER PRIMARY KEY,
			name TEXT,
			price DOUBLE,
			metadata JSON
		)
	`)
	require.NoError(t, err)

	// Insert test data with JSON metadata
	_, err = db.Exec(`
		INSERT INTO products VALUES
		(1, 'Widget', 19.99, '{"brand": "Acme", "color": "red", "specs": {"weight": 100}}'),
		(2, 'Gadget', 29.99, '{"brand": "Beta", "color": "blue", "specs": {"weight": 200}}'),
		(3, 'Doohickey', 39.99, '{"brand": "Acme", "color": "green", "specs": {"weight": 150}}')
	`)
	require.NoError(t, err)

	// Set up CEL environment with schema for JSON detection
	productSchema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "name", Type: "text"},
		{Name: "price", Type: "double precision"},
		{Name: "metadata", Type: "json", IsJSON: true},
	})

	schemas := map[string]pg.Schema{
		"product": productSchema,
	}

	provider := pg.NewTypeProvider(schemas)

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("product", cel.ObjectType("product")),
	)
	require.NoError(t, err)

	dialectOpt := cel2sql.WithDialect(mysqlDialect.New())
	schemaOpt := cel2sql.WithSchemas(schemas)

	tests := []struct {
		name         string
		celExpr      string
		expectedRows int
		description  string
	}{
		{
			name:         "JSON field access",
			celExpr:      `product.metadata.brand == "Acme"`,
			expectedRows: 2,
			description:  "JSON field access with ->>",
		},
		{
			name:         "JSON field access different value",
			celExpr:      `product.metadata.color == "blue"`,
			expectedRows: 1,
			description:  "JSON field access with different value",
		},
		{
			name:         "JSON with regular field",
			celExpr:      `product.metadata.brand == "Acme" && product.price > 30.0`,
			expectedRows: 1, // Doohickey (Acme, 39.99)
			description:  "JSON field combined with regular field comparison",
		},
		{
			name:         "JSON field existence",
			celExpr:      `has(product.metadata.brand)`,
			expectedRows: 3, // All rows have 'brand'
			description:  "JSON field existence check",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.celExpr)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("CEL compilation failed: %v", issues.Err())
			}

			sqlCondition, err := cel2sql.Convert(ast, dialectOpt, schemaOpt)
			require.NoError(t, err, "Conversion should succeed for: %s", tt.description)

			t.Logf("CEL Expression: %s", tt.celExpr)
			t.Logf("Generated SQL WHERE clause: %s", sqlCondition)

			// #nosec G202 - This is a test validating SQL generation, not a security risk
			query := "SELECT COUNT(*) FROM products product WHERE " + sqlCondition
			t.Logf("Full SQL Query: %s", query)

			var actualRows int
			err = db.QueryRow(query).Scan(&actualRows)
			require.NoError(t, err, "Generated SQL should execute successfully. %s\nSQL: %s",
				tt.description, sqlCondition)

			require.Equal(t, tt.expectedRows, actualRows,
				"Query should return expected number of rows. %s\nCEL: %s\nSQL: %s",
				tt.description, tt.celExpr, sqlCondition)

			t.Logf("OK: %s", tt.description)
		})
	}
}
