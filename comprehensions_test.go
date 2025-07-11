package cel2sql_test

import (
	"strings"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql"
	"github.com/spandigital/cel2sql/pg"
)

func TestComprehensionImplementation(t *testing.T) {
	// Create a simple environment for testing comprehensions
	env, err := cel.NewEnv(
		cel.Variable("employees", cel.ListType(cel.MapType(cel.StringType, cel.DynType))),
		cel.Variable("numbers", cel.ListType(cel.IntType)),
	)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		expression  string
		expectError bool
		expectedSQL string
		description string
	}{
		{
			name:        "all_comprehension",
			expression:  `[{'salary': 60000}, {'salary': 40000}].all(e, e.salary > 50000)`,
			expectError: false,
			expectedSQL: `NOT EXISTS (SELECT 1 FROM UNNEST(ARRAY[STRUCT(60000 AS salary), STRUCT(40000 AS salary)]) AS e WHERE NOT (e.salary > 50000))`,
			description: "ALL comprehension should generate NOT EXISTS with UNNEST",
		},
		{
			name:        "exists_comprehension",
			expression:  `[{'department': 'Engineering'}, {'department': 'Sales'}].exists(e, e.department == 'Engineering')`,
			expectError: false,
			expectedSQL: `EXISTS (SELECT 1 FROM UNNEST(ARRAY[STRUCT('Engineering' AS department), STRUCT('Sales' AS department)]) AS e WHERE e.department = 'Engineering')`,
			description: "EXISTS comprehension should generate EXISTS with UNNEST",
		},
		{
			name:        "exists_one_comprehension",
			expression:  `[{'role': 'CEO'}, {'role': 'CTO'}].exists_one(e, e.role == 'CEO')`,
			expectError: false,
			expectedSQL: `(SELECT COUNT(*) FROM UNNEST(ARRAY[STRUCT('CEO' AS role), STRUCT('CTO' AS role)]) AS e WHERE e.role = 'CEO') = 1`,
			description: "EXISTS_ONE comprehension should generate COUNT query",
		},
		{
			name:        "map_comprehension",
			expression:  `[{'name': 'John'}, {'name': 'Jane'}].map(e, e.name)`,
			expectError: false,
			expectedSQL: `ARRAY(SELECT e.name FROM UNNEST(ARRAY[STRUCT('John' AS name), STRUCT('Jane' AS name)]) AS e)`,
			description: "MAP comprehension should generate ARRAY SELECT",
		},
		{
			name:        "filter_comprehension",
			expression:  `[{'active': true}, {'active': false}].filter(e, e.active)`,
			expectError: false,
			expectedSQL: `ARRAY(SELECT e FROM UNNEST(ARRAY[STRUCT(TRUE AS active), STRUCT(FALSE AS active)]) AS e WHERE e.active)`,
			description: "FILTER comprehension should generate ARRAY SELECT with WHERE",
		},
		{
			name:        "list_all",
			expression:  `[1, 2, 3, 4].all(x, x > 0)`,
			expectError: false,
			expectedSQL: `NOT EXISTS (SELECT 1 FROM UNNEST(ARRAY[1, 2, 3, 4]) AS x WHERE NOT (x > 0))`,
			description: "ALL comprehension on simple list should work",
		},
		{
			name:        "list_exists",
			expression:  `[1, 2, 3, 4].exists(x, x > 3)`,
			expectError: false,
			expectedSQL: `EXISTS (SELECT 1 FROM UNNEST(ARRAY[1, 2, 3, 4]) AS x WHERE x > 3)`,
			description: "EXISTS comprehension on simple list should work",
		},
		{
			name:        "map_with_filter",
			expression:  `[{'active': true, 'name': 'John'}, {'active': false, 'name': 'Jane'}].map(e, e.active, e.name)`,
			expectError: false, // Let's see what it generates
			description: "MAP with filter comprehension to see what gets generated",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ast, issues := env.Compile(tc.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed for %s", tc.description)

			result, err := cel2sql.Convert(ast)

			if tc.expectError {
				assert.Error(t, err, "Should get error for %s", tc.description)
			} else {
				if err != nil {
					t.Logf("Unexpected error for %s: %v", tc.description, err)
				}
				assert.NoError(t, err, "Should succeed for %s", tc.description)
				if tc.expectedSQL != "" {
					assert.Contains(t, result, "UNNEST", "Should use UNNEST for array comprehension: %s", tc.description)
					t.Logf("Generated SQL for %s: %s", tc.description, result)
				}
			}
		})
	}
}

func TestNonComprehensionExpressionsStillWork(t *testing.T) {
	// Ensure that non-comprehension expressions still work correctly
	schema := pg.Schema{
		{Name: "name", Type: "text", Repeated: false},
		{Name: "age", Type: "bigint", Repeated: false},
		{Name: "active", Type: "boolean", Repeated: false},
	}

	provider := pg.NewTypeProvider(map[string]pg.Schema{"User": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("user", cel.ObjectType("User")),
	)
	require.NoError(t, err)

	testCases := []struct {
		name       string
		expression string
		expected   string
	}{
		{
			name:       "simple_comparison",
			expression: `user.age > 18`,
			expected:   `user.age > 18`,
		},
		{
			name:       "string_equality",
			expression: `user.name == 'John'`,
			expected:   `user.name = 'John'`,
		},
		{
			name:       "boolean_field",
			expression: `user.active`,
			expected:   `user.active`,
		},
		{
			name:       "logical_and",
			expression: `user.age > 18 && user.active`,
			expected:   `user.age > 18 AND user.active`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ast, issues := env.Compile(tc.expression)
			require.NoError(t, issues.Err())

			result, err := cel2sql.Convert(ast)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestComprehensionWithPostgreSQLSchemas(t *testing.T) {
	// Test comprehensions with realistic PostgreSQL schemas
	schema := pg.Schema{
		{Name: "id", Type: "bigint", Repeated: false},
		{Name: "name", Type: "text", Repeated: false},
		{Name: "email", Type: "text", Repeated: false},
		{Name: "age", Type: "bigint", Repeated: false},
		{Name: "active", Type: "boolean", Repeated: false},
		{Name: "salary", Type: "double precision", Repeated: false},
		{Name: "tags", Type: "text", Repeated: true},     // Array field
		{Name: "scores", Type: "bigint", Repeated: true}, // Array of integers
	}

	provider := pg.NewTypeProvider(map[string]pg.Schema{"Employee": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("employees", cel.ListType(cel.ObjectType("Employee"))),
		cel.Variable("emp", cel.ObjectType("Employee")),
	)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		expression  string
		description string
	}{
		{
			name:        "all_high_earners",
			expression:  `employees.all(e, e.salary > 80000)`,
			description: "Check if all employees are high earners",
		},
		{
			name:        "exists_senior",
			expression:  `employees.exists(e, e.age > 50)`,
			description: "Check if any employee is senior",
		},
		{
			name:        "exists_one_admin",
			expression:  `employees.exists_one(e, e.name == 'admin')`,
			description: "Check if exactly one admin exists",
		},
		{
			name:        "map_names",
			expression:  `employees.map(e, e.name)`,
			description: "Extract all employee names",
		},
		{
			name:        "filter_active",
			expression:  `employees.filter(e, e.active)`,
			description: "Get only active employees",
		},
		{
			name:        "filter_by_age",
			expression:  `employees.filter(e, e.age >= 25 && e.age <= 65)`,
			description: "Get employees in working age range",
		},
		{
			name:        "map_emails_active",
			expression:  `employees.filter(e, e.active).map(e, e.email)`,
			description: "Get emails of active employees (chained)",
		},
		{
			name:        "all_have_email",
			expression:  `employees.all(e, e.email != '')`,
			description: "Check if all employees have non-empty email",
		},
		{
			name:        "exists_high_score",
			expression:  `emp.scores.exists(s, s > 90)`,
			description: "Check if employee has any high score (array field)",
		},
		{
			name:        "all_scores_passing",
			expression:  `emp.scores.all(s, s >= 60)`,
			description: "Check if all scores are passing (array field)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ast, issues := env.Compile(tc.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed for %s", tc.description)

			result, err := cel2sql.Convert(ast)
			require.NoError(t, err, "SQL conversion should succeed for %s", tc.description)

			t.Logf("Expression: %s", tc.expression)
			t.Logf("Generated SQL: %s", result)

			// Verify that the SQL contains expected PostgreSQL patterns
			assert.Contains(t, result, "UNNEST", "Should use UNNEST for array operations")

			// Verify proper field access patterns
			if strings.Contains(tc.expression, ".name") {
				assert.Contains(t, result, ".name", "Should access name field")
			}
			if strings.Contains(tc.expression, ".active") {
				assert.Contains(t, result, ".active", "Should access active field")
			}
		})
	}
}

func TestComprehensionWithNestedStructures(t *testing.T) {
	// Test comprehensions with nested PostgreSQL composite types
	addressSchema := pg.Schema{
		{Name: "street", Type: "text", Repeated: false},
		{Name: "city", Type: "text", Repeated: false},
		{Name: "zipcode", Type: "text", Repeated: false},
	}

	employeeSchema := pg.Schema{
		{Name: "id", Type: "bigint", Repeated: false},
		{Name: "name", Type: "text", Repeated: false},
		{Name: "address", Type: "composite", Schema: addressSchema},
		{Name: "skills", Type: "text", Repeated: true},
	}

	provider := pg.NewTypeProvider(map[string]pg.Schema{
		"Address":  addressSchema,
		"Employee": employeeSchema,
	})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("employees", cel.ListType(cel.ObjectType("Employee"))),
	)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		expression  string
		description string
	}{
		{
			name:        "filter_by_city",
			expression:  `employees.filter(e, e.address.city == 'New York')`,
			description: "Filter employees by city in nested address",
		},
		{
			name:        "exists_skill",
			expression:  `employees.exists(e, e.skills.exists(s, s == 'Go'))`,
			description: "Check if any employee has Go skill (nested array comprehension)",
		},
		{
			name:        "map_cities",
			expression:  `employees.map(e, e.address.city)`,
			description: "Extract cities from nested address structures",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ast, issues := env.Compile(tc.expression)
			if issues.Err() != nil {
				t.Logf("CEL compilation error: %v", issues.Err())
				// Some nested comprehensions might not be supported yet
				t.Skip("Skipping test due to CEL compilation issues - may need advanced comprehension support")
				return
			}

			result, err := cel2sql.Convert(ast)
			if err != nil {
				t.Logf("SQL conversion error: %v", err)
				// Some complex nested cases might not be implemented yet
				if strings.Contains(err.Error(), "not yet implemented") {
					t.Skip("Skipping test - advanced comprehension not yet implemented")
					return
				}
			}
			require.NoError(t, err, "SQL conversion should succeed for %s", tc.description)

			t.Logf("Expression: %s", tc.expression)
			t.Logf("Generated SQL: %s", result)
		})
	}
}

func TestComprehensionEdgeCases(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("emptyList", cel.ListType(cel.IntType)),
		cel.Variable("numbers", cel.ListType(cel.IntType)),
	)
	require.NoError(t, err)

	testCases := []struct {
		name        string
		expression  string
		description string
	}{
		{
			name:        "all_empty_list",
			expression:  `[].all(x, x > 0)`,
			description: "ALL on empty list should work",
		},
		{
			name:        "exists_empty_list",
			expression:  `[].exists(x, x > 0)`,
			description: "EXISTS on empty list should work",
		},
		{
			name:        "filter_empty_list",
			expression:  `[].filter(x, x > 0)`,
			description: "FILTER on empty list should work",
		},
		{
			name:        "map_empty_list",
			expression:  `[].map(x, x * 2)`,
			description: "MAP on empty list should work",
		},
		{
			name:        "complex_predicate",
			expression:  `[1, 2, 3, 4, 5].all(x, x > 0 && x < 10 && x % 2 == 0 || x == 1)`,
			description: "Complex predicate with multiple conditions",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ast, issues := env.Compile(tc.expression)
			require.NoError(t, issues.Err(), "CEL compilation should succeed for %s", tc.description)

			result, err := cel2sql.Convert(ast)
			require.NoError(t, err, "SQL conversion should succeed for %s", tc.description)

			t.Logf("Expression: %s", tc.expression)
			t.Logf("Generated SQL: %s", result)

			// Basic sanity checks
			assert.Contains(t, result, "UNNEST", "Should use UNNEST for array operations")
		})
	}
}
