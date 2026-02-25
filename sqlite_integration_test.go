package cel2sql_test

import (
	"database/sql"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	"github.com/spandigital/cel2sql/v3"
	sqliteDialect "github.com/spandigital/cel2sql/v3/dialect/sqlite"
	"github.com/spandigital/cel2sql/v3/pg"
)

// TestSQLiteOperatorsIntegration validates operator conversions against a real SQLite database.
// This uses an in-memory SQLite database (no Docker required).
func TestSQLiteOperatorsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Open in-memory SQLite database
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Logf("failed to close db: %v", closeErr)
		}
	}()

	// Create test table (SQLite uses INTEGER for booleans, REAL for floats)
	_, err = db.Exec(`
		CREATE TABLE test_data (
			id INTEGER PRIMARY KEY,
			text_val TEXT,
			int_val INTEGER,
			float_val REAL,
			bool_val INTEGER,
			nullable_text TEXT,
			nullable_int INTEGER
		)
	`)
	require.NoError(t, err)

	// Insert test data (using 1/0 for boolean values)
	_, err = db.Exec(`
		INSERT INTO test_data VALUES
		(1, 'hello', 10, 10.5, 1, 'present', 100),
		(2, 'world', 20, 20.5, 0, NULL, NULL),
		(3, 'test', 30, 30.5, 1, 'here', 200),
		(4, 'hello world', 5, 5.5, 0, 'value', 50),
		(5, 'testing', 15, 15.5, 1, 'test', 150)
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

	dialectOpt := cel2sql.WithDialect(sqliteDialect.New())

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
			expectedRows: 2, // rows 2 and 4 (bool_val == false)
			description:  "Logical NOT operator",
		},
		{
			name:         "Complex logical expression",
			celExpr:      `(int_val > 10 && bool_val) || int_val < 10`,
			expectedRows: 3, // rows 3, 5 (>10 && true), row 4 (<10)
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
			description:  "String concatenation operator (||)",
		},
		{
			name:         "String contains",
			celExpr:      `text_val.contains("world")`,
			expectedRows: 2, // "world", "hello world"
			description:  "String contains function (INSTR)",
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

		// Complex combined operators
		{
			name:         "Complex multi-operator expression",
			celExpr:      `int_val > 10 && bool_val && text_val.contains("test")`,
			expectedRows: 2, // rows 3 (30, true, "test") and 5 (15, true, "testing")
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
			expectedRows: 2, // rows 2 and 4 (bool_val == false)
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

			// Convert to SQL with SQLite dialect
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

// TestSQLiteEdgeCasesIntegration validates edge cases against a real SQLite database.
func TestSQLiteEdgeCasesIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("sqlite", ":memory:")
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
			zero_float REAL,
			negative_int INTEGER,
			negative_float REAL,
			large_int INTEGER
		)
	`)
	require.NoError(t, err)

	_, err = db.Exec(`
		INSERT INTO edge_cases VALUES
		(1, '', 0, 0.0, -10, -5.5, 9223372036854775807),
		(2, 'value', 1, 1.0, -1, -0.1, 123456789),
		(3, 'another', 0, 0.0, 0, 0.0, 0)
	`)
	require.NoError(t, err)

	env, err := cel.NewEnv(
		cel.Variable("id", cel.IntType),
		cel.Variable("empty_string", cel.StringType),
		cel.Variable("zero_int", cel.IntType),
		cel.Variable("zero_float", cel.DoubleType),
		cel.Variable("negative_int", cel.IntType),
		cel.Variable("negative_float", cel.DoubleType),
		cel.Variable("large_int", cel.IntType),
	)
	require.NoError(t, err)

	dialectOpt := cel2sql.WithDialect(sqliteDialect.New())

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
			expectedRows: 2,
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
			expectedRows: 2,
			description:  "Zero float should be handled correctly",
		},
		{
			name:         "Negative float comparison",
			celExpr:      `negative_float < 0.0`,
			expectedRows: 2, // -5.5 and -0.1
			description:  "Negative floats should work correctly",
		},
		{
			name:         "Arithmetic with zero",
			celExpr:      `zero_int + 10 == 10`,
			expectedRows: 2, // 0 + 10 = 10
			description:  "Arithmetic with zero should work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.celExpr)
			if issues != nil && issues.Err() != nil {
				t.Fatalf("CEL compilation failed: %v", issues.Err())
			}

			sqlCondition, err := cel2sql.Convert(ast, dialectOpt)
			require.NoError(t, err, "Conversion should succeed for: %s", tt.description)

			t.Logf("CEL Expression: %s", tt.celExpr)
			t.Logf("Generated SQL WHERE clause: %s", sqlCondition)

			// #nosec G202 - This is a test validating SQL generation, not a security risk
			query := "SELECT COUNT(*) FROM edge_cases WHERE " + sqlCondition

			var actualRows int
			err = db.QueryRow(query).Scan(&actualRows)
			require.NoError(t, err, "Generated SQL should execute successfully. %s",
				tt.description)

			require.Equal(t, tt.expectedRows, actualRows,
				"Query should return expected number of rows. %s", tt.description)

			t.Logf("OK: %s", tt.description)
		})
	}
}

// TestSQLiteJSONIntegration validates JSON operations against a real SQLite database.
func TestSQLiteJSONIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Logf("failed to close db: %v", closeErr)
		}
	}()

	// Create test table with JSON column (stored as TEXT in SQLite)
	_, err = db.Exec(`
		CREATE TABLE products (
			id INTEGER PRIMARY KEY,
			name TEXT,
			price REAL,
			metadata TEXT
		)
	`)
	require.NoError(t, err)

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

	dialectOpt := cel2sql.WithDialect(sqliteDialect.New())
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
			description:  "JSON field access with json_extract",
		},
		{
			name:         "JSON field access different value",
			celExpr:      `product.metadata.color == "blue"`,
			expectedRows: 1,
			description:  "JSON field access with different value",
		},
		{
			name:         "JSON nested field access",
			celExpr:      `product.metadata.brand == "Acme" && product.price > 30.0`,
			expectedRows: 1, // Doohickey (Acme, 39.99)
			description:  "JSON field combined with regular field comparison",
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
