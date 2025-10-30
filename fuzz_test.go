package cel2sql_test

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/google/cel-go/cel"
	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
	"github.com/spandigital/cel2sql/v3/sqltypes"
)

// FuzzConvert fuzzes the main Convert function to find crashes, panics, and SQL injection vulnerabilities
func FuzzConvert(f *testing.F) {
	// Seed corpus with valid CEL expressions from existing tests
	seeds := []string{
		`name == "test"`,
		`age > 18`,
		`active && verified`,
		`name.startsWith("A")`,
		`email.endsWith("@test.com")`,
		`text.contains("search")`,
		`email.matches(r"^[a-z]+@test\.com$")`,
		`price < 100 && in_stock`,
		`"admin" in roles`,
		`size(tags) > 0`,
		`scores.all(s, s > 60)`,
		`tags.exists(t, t == "urgent")`,
		`items.filter(i, i.active)`,
		`names.map(n, n + "!")`,
		`preferences.theme == "dark"`,
		`age >= 18 ? "adult" : "minor"`,
		`created_at > timestamp("2024-01-01T00:00:00Z")`,
		`int(created_at) > 1704000000`,
		// Edge cases
		`name == null`,
		`flag == true`,
		`!deleted`,
		`(a && b) || c`,
		// Complex expressions
		`(age > 30 ? "senior" : "junior") == "senior"`,
		`user.profile.settings.language == "en"`,
		`employees.filter(e, e.age > 30).map(e, e.name)`,
		// Potential injection vectors
		`name == "'; DROP TABLE users; --"`,
		`text.contains("100% discount")`,
		`pattern == "test_underscore"`,
		`field == "back\\slash"`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	// Create CEL environment for fuzzing
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
		{Name: "email", Type: "text"},
		{Name: "text", Type: "text"},
		{Name: "age", Type: "integer"},
		{Name: "price", Type: "double precision"},
		{Name: "active", Type: "boolean"},
		{Name: "verified", Type: "boolean"},
		{Name: "deleted", Type: "boolean"},
		{Name: "in_stock", Type: "boolean"},
		{Name: "flag", Type: "boolean"},
		{Name: "created_at", Type: "timestamp with time zone"},
		{Name: "roles", Type: "text", Repeated: true},
		{Name: "tags", Type: "text", Repeated: true},
		{Name: "scores", Type: "integer", Repeated: true},
		{Name: "names", Type: "text", Repeated: true},
		{Name: "preferences", Type: "jsonb"},
		{Name: "profile", Type: "jsonb"},
		{Name: "pattern", Type: "text"},
		{Name: "field", Type: "text"},
		{Name: "a", Type: "boolean"},
		{Name: "b", Type: "boolean"},
		{Name: "c", Type: "boolean"},
	})

	itemSchema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "active", Type: "boolean"},
	})

	employeeSchema := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
		{Name: "age", Type: "integer"},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{
		"Record":   schema,
		"Item":     itemSchema,
		"Employee": employeeSchema,
	})

	env, err := cel.NewEnv(
		cel.Types(sqltypes.Date, sqltypes.Time, sqltypes.DateTime, sqltypes.Interval, sqltypes.DatePart),
		cel.CustomTypeProvider(provider),
		cel.Variable("name", cel.StringType),
		cel.Variable("email", cel.StringType),
		cel.Variable("text", cel.StringType),
		cel.Variable("age", cel.IntType),
		cel.Variable("price", cel.DoubleType),
		cel.Variable("active", cel.BoolType),
		cel.Variable("verified", cel.BoolType),
		cel.Variable("deleted", cel.BoolType),
		cel.Variable("in_stock", cel.BoolType),
		cel.Variable("flag", cel.BoolType),
		cel.Variable("created_at", cel.TimestampType),
		cel.Variable("roles", cel.ListType(cel.StringType)),
		cel.Variable("tags", cel.ListType(cel.StringType)),
		cel.Variable("scores", cel.ListType(cel.IntType)),
		cel.Variable("names", cel.ListType(cel.StringType)),
		cel.Variable("items", cel.ListType(cel.ObjectType("Item"))),
		cel.Variable("employees", cel.ListType(cel.ObjectType("Employee"))),
		cel.Variable("user", cel.ObjectType("Record")),
		cel.Variable("preferences", cel.ObjectType("Record")),
		cel.Variable("profile", cel.ObjectType("Record")),
		cel.Variable("pattern", cel.StringType),
		cel.Variable("field", cel.StringType),
		cel.Variable("a", cel.BoolType),
		cel.Variable("b", cel.BoolType),
		cel.Variable("c", cel.BoolType),
	)
	if err != nil {
		f.Fatalf("Failed to create CEL environment: %v", err)
	}

	f.Fuzz(func(t *testing.T, celExpr string) {
		// Skip empty strings and overly long expressions
		if len(celExpr) == 0 || len(celExpr) > 10000 {
			return
		}

		// Try to compile the CEL expression
		ast, issues := env.Compile(celExpr)
		if issues != nil && issues.Err() != nil {
			// Compilation errors are expected for random input
			return
		}

		// Try to convert to SQL - this should never panic or crash
		sqlOutput, err := cel2sql.Convert(ast)

		// We don't care if conversion fails with an error,
		// but it should never panic or produce invalid output
		if err != nil {
			// Error is acceptable - just ensure it's a proper error, not a panic
			return
		}

		// Basic sanity checks on generated SQL
		if len(sqlOutput) > 0 {
			// SQL should not contain null bytes
			if strings.Contains(sqlOutput, "\x00") {
				t.Errorf("Generated SQL contains null bytes: %q", sqlOutput)
			}

			// SQL should be valid UTF-8
			if !utf8.ValidString(sqlOutput) {
				t.Errorf("Generated SQL is not valid UTF-8: %q", sqlOutput)
			}
		}
	})
}

// FuzzEscapeLikePattern fuzzes the LIKE pattern escaping to find SQL injection vulnerabilities
func FuzzEscapeLikePattern(f *testing.F) {
	// Import the internal function for testing
	// Note: This requires exporting escapeLikePattern or using a test-only wrapper

	// Seed with known patterns including edge cases
	seeds := []string{
		"simple",
		"with%percent",
		"with_underscore",
		"with\\backslash",
		"with'quote",
		"with\"doublequote",
		"100% discount",
		"test_value",
		"back\\slash\\here",
		"multi'''quotes",
		"unicode: 你好",
		"emoji: 🎉",
		"newline\nhere",
		"tab\there",
		"all%special_chars\\and'quotes",
		"",
		" ",
		"%",
		"_",
		"\\",
		"'",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, pattern string) {
		// Skip overly long patterns
		if len(pattern) > 1000 {
			return
		}

		// Create a CEL expression with startsWith using the pattern
		celExpr := `name.startsWith("` + strings.ReplaceAll(pattern, `"`, `\"`) + `")`

		schema := pg.NewSchema([]pg.FieldSchema{
			{Name: "name", Type: "text"},
		})
		provider := pg.NewTypeProvider(map[string]pg.Schema{"Record": schema})

		env, err := cel.NewEnv(
			cel.CustomTypeProvider(provider),
			cel.Variable("name", cel.StringType),
		)
		if err != nil {
			t.Fatalf("Failed to create CEL environment: %v", err)
		}

		ast, issues := env.Compile(celExpr)
		if issues != nil && issues.Err() != nil {
			// Invalid CEL is fine - we're testing the escaping logic
			return
		}

		sqlOutput, err := cel2sql.Convert(ast)
		if err != nil {
			// Errors are acceptable
			return
		}

		// Security checks on the escaped pattern
		if len(sqlOutput) > 0 {
			// The SQL should contain the LIKE operator
			if !strings.Contains(sqlOutput, "LIKE") {
				t.Errorf("Expected LIKE in output for pattern %q, got: %q", pattern, sqlOutput)
			}

			// Check for proper quoting - all LIKE patterns should be quoted
			if strings.Contains(sqlOutput, "LIKE ") && !strings.Contains(sqlOutput, "LIKE '") {
				t.Errorf("LIKE pattern not properly quoted: %q", sqlOutput)
			}

			// Ensure no null bytes
			if strings.Contains(sqlOutput, "\x00") {
				t.Errorf("Generated SQL contains null bytes: %q", sqlOutput)
			}

			// Count quotes to ensure they're balanced (rough check)
			// This helps catch escaping issues
			singleQuotes := strings.Count(sqlOutput, "'")
			if singleQuotes%2 != 0 {
				// This might be a false positive if there are escaped quotes,
				// but it's worth checking
				t.Logf("Warning: Unbalanced quotes in SQL: %q", sqlOutput)
			}
		}
	})
}

// FuzzConvertRE2ToPOSIX fuzzes the regex conversion to find edge cases and crashes
func FuzzConvertRE2ToPOSIX(f *testing.F) {
	// Seed with various regex patterns
	seeds := []string{
		`^[a-z]+$`,
		`\d{3}-\d{4}`,
		`\w+@\w+\.\w+`,
		`\bword\b`,
		`(?i)case-insensitive`,
		`[[:digit:]]+`,
		`\\escaped`,
		`^start.*end$`,
		`(group|alternative)`,
		`\s+\S+`,
		`\d+\w+`,
		`\\d`,
		`\d`,
		`\\\\`,
		`[a-zA-Z0-9]`,
		`.*`,
		`.+`,
		`.?`,
		`^$`,
		``,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, regexPattern string) {
		// Skip overly long patterns
		if len(regexPattern) > 500 {
			return
		}

		// Create a CEL matches expression
		// Escape the pattern for CEL
		escapedPattern := strings.ReplaceAll(regexPattern, `\`, `\\`)
		escapedPattern = strings.ReplaceAll(escapedPattern, `"`, `\"`)
		celExpr := `email.matches(r"` + escapedPattern + `")`

		schema := pg.NewSchema([]pg.FieldSchema{
			{Name: "email", Type: "text"},
		})
		provider := pg.NewTypeProvider(map[string]pg.Schema{"Record": schema})

		env, err := cel.NewEnv(
			cel.CustomTypeProvider(provider),
			cel.Variable("email", cel.StringType),
		)
		if err != nil {
			t.Fatalf("Failed to create CEL environment: %v", err)
		}

		ast, issues := env.Compile(celExpr)
		if issues != nil && issues.Err() != nil {
			// CEL compilation errors are fine - we're fuzzing the conversion
			return
		}

		sqlOutput, err := cel2sql.Convert(ast)
		if err != nil {
			// Conversion errors are acceptable
			return
		}

		// Basic checks on converted regex
		if len(sqlOutput) > 0 {
			// Should contain the ~ operator for regex matching
			if !strings.Contains(sqlOutput, "~") {
				// Might be ~* for case-insensitive
				if !strings.Contains(sqlOutput, "~*") {
					t.Errorf("Expected ~ or ~* operator in regex output, got: %q", sqlOutput)
				}
			}

			// Ensure no null bytes
			if strings.Contains(sqlOutput, "\x00") {
				t.Errorf("Generated SQL contains null bytes: %q", sqlOutput)
			}

			// Pattern should be quoted
			if strings.Contains(sqlOutput, "~ ") && !strings.Contains(sqlOutput, "~ '") &&
				!strings.Contains(sqlOutput, "~* '") {
				t.Errorf("Regex pattern not properly quoted: %q", sqlOutput)
			}
		}
	})
}
