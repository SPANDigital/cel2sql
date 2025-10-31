package cel2sql_test

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/ext"

	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
	"github.com/spandigital/cel2sql/v3/sqltypes"
)

// setupSimpleBenchmarkEnv creates a simple CEL environment without schemas
func setupSimpleBenchmarkEnv(b *testing.B) *cel.Env {
	b.Helper()

	env, err := cel.NewEnv(
		cel.Variable("name", cel.StringType),
		cel.Variable("email", cel.StringType),
		cel.Variable("age", cel.IntType),
		cel.Variable("active", cel.BoolType),
		cel.Variable("score", cel.DoubleType),
		cel.Variable("tags", cel.ListType(cel.StringType)),
		cel.Variable("scores", cel.ListType(cel.IntType)),
		ext.Strings(), // Enable string extension functions
	)
	if err != nil {
		b.Fatal(err)
	}

	return env
}

// setupSchemaBenchmarkEnv creates a CEL environment with schema support
func setupSchemaBenchmarkEnv(b *testing.B) (*cel.Env, map[string]pg.Schema) {
	b.Helper()

	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
		{Name: "email", Type: "text"},
		{Name: "age", Type: "integer"},
		{Name: "active", Type: "boolean"},
		{Name: "score", Type: "double precision"},
		{Name: "tags", Type: "text", Repeated: true},
		{Name: "scores", Type: "integer", Repeated: true},
		{Name: "metadata", Type: "jsonb"},
		{Name: "created_at", Type: "timestamp"},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{
		"person": schema,
	})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("person", cel.ObjectType("person")),
	)
	if err != nil {
		b.Fatal(err)
	}

	schemas := map[string]pg.Schema{
		"person": schema,
	}

	return env, schemas
}

// setupTimestampBenchmarkEnv creates a CEL environment with timestamp support
func setupTimestampBenchmarkEnv(b *testing.B) *cel.Env {
	b.Helper()

	env, err := cel.NewEnv(
		cel.Variable("created_at", cel.TimestampType),
		cel.Types(
			sqltypes.Date, sqltypes.Time, sqltypes.DateTime, sqltypes.Interval, sqltypes.DatePart,
		),
		cel.Function("date",
			cel.Overload("date_string", []*cel.Type{cel.StringType}, cel.ObjectType("DATE")),
			cel.Overload("date_int_int_int", []*cel.Type{cel.IntType, cel.IntType, cel.IntType}, cel.ObjectType("DATE"))),
		cel.Function("time", cel.Overload("time_string", []*cel.Type{cel.StringType}, cel.ObjectType("TIME"))),
		cel.Function("datetime",
			cel.Overload("datetime_string", []*cel.Type{cel.StringType}, cel.ObjectType("DATETIME")),
			cel.Overload("datetime_date_time", []*cel.Type{cel.ObjectType("DATE"), cel.ObjectType("TIME")}, cel.ObjectType("DATETIME"))),
		cel.Function("timestamp",
			cel.Overload("timestamp_datetime_string", []*cel.Type{cel.ObjectType("DATETIME"), cel.StringType}, cel.TimestampType)),
		cel.Function("interval", cel.Overload("interval_int_datepart", []*cel.Type{cel.IntType, cel.ObjectType("date_part")}, cel.ObjectType("INTERVAL"))),
		cel.Variable("YEAR", cel.ObjectType("date_part")),
		cel.Variable("MONTH", cel.ObjectType("date_part")),
		cel.Variable("DAY", cel.ObjectType("date_part")),
	)
	if err != nil {
		b.Fatal(err)
	}

	return env
}

// BenchmarkConvertSimple benchmarks simple field comparisons
func BenchmarkConvertSimple(b *testing.B) {
	env := setupSimpleBenchmarkEnv(b)

	tests := []struct {
		name       string
		expression string
	}{
		{"equality", `age == 18`},
		{"greater_than", `age > 18`},
		{"string_equality", `name == "test"`},
		{"boolean_check", `active`},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			ast, issues := env.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				b.Fatal(issues.Err())
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := cel2sql.Convert(ast)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkConvertOperators benchmarks various operators
func BenchmarkConvertOperators(b *testing.B) {
	env := setupSimpleBenchmarkEnv(b)

	tests := []struct {
		name       string
		expression string
	}{
		{"logical_and", `age > 18 && active`},
		{"logical_or", `age < 18 || age > 65`},
		{"arithmetic_add", `age + 5 > 30`},
		{"arithmetic_sub", `age - 5 < 10`},
		{"arithmetic_mul", `score * 2.0 > 100.0`},
		{"arithmetic_div", `score / 2.0 < 50.0`},
		{"modulo", `age % 2 == 0`},
		{"string_concat", `name + " " + email`},
		{"complex_expression", `(age > 18 && active) || (score > 90.0 && name.startsWith("A"))`},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			ast, issues := env.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				b.Fatal(issues.Err())
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := cel2sql.Convert(ast)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkConvertComprehensions benchmarks CEL comprehensions
func BenchmarkConvertComprehensions(b *testing.B) {
	env := setupSimpleBenchmarkEnv(b)

	tests := []struct {
		name       string
		expression string
	}{
		{"all_simple", `scores.all(s, s > 0)`},
		{"all_complex", `scores.all(s, s > 60 && s < 100)`},
		{"exists_simple", `tags.exists(t, t == "premium")`},
		{"exists_complex", `tags.exists(t, t.startsWith("prod") && t.endsWith("ion"))`},
		{"exists_one", `tags.exists_one(t, t == "vip")`},
		{"filter", `scores.filter(s, s > 80)`},
		{"map", `scores.map(s, s * 2)`},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			ast, issues := env.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				b.Fatal(issues.Err())
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := cel2sql.Convert(ast)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkConvertJSONPath benchmarks JSON field access and operations
func BenchmarkConvertJSONPath(b *testing.B) {
	env, schemas := setupSchemaBenchmarkEnv(b)

	tests := []struct {
		name       string
		expression string
	}{
		{"simple_access", `person.metadata.verified == true`},
		{"nested_access", `person.metadata.profile.name == "test"`},
		{"json_has", `has(person.metadata.verified)`},
		{"nested_json_has", `has(person.metadata.profile.name)`},
		{"json_comparison", `person.metadata.score > 50`},
		{"complex_json", `person.metadata.profile.verified && person.metadata.profile.score > 80`},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			ast, issues := env.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				b.Fatal(issues.Err())
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkConvertRegex benchmarks regex pattern conversion
func BenchmarkConvertRegex(b *testing.B) {
	env, schemas := setupSchemaBenchmarkEnv(b)

	tests := []struct {
		name       string
		expression string
	}{
		{"simple_pattern", `person.email.matches(r"@example\.com$")`},
		{"case_insensitive", `person.name.matches(r"(?i)^test")`},
		{"complex_pattern", `person.email.matches(r"[a-z0-9]+@[a-z]+\.[a-z]+")`},
		{"with_digit_class", `person.name.matches(r"\d{4}")`},
		{"with_word_class", `person.name.matches(r"\w+")`},
		{"with_word_boundary", `person.email.matches(r"\btest\b")`},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			ast, issues := env.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				b.Fatal(issues.Err())
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkConvertDeeplyNested benchmarks deeply nested expressions
func BenchmarkConvertDeeplyNested(b *testing.B) {
	env := setupSimpleBenchmarkEnv(b)

	tests := []struct {
		name       string
		expression string
	}{
		{
			"nested_and_5",
			`age > 0 && age > 1 && age > 2 && age > 3 && age > 4`,
		},
		{
			"nested_and_10",
			`age > 0 && age > 1 && age > 2 && age > 3 && age > 4 && age > 5 && age > 6 && age > 7 && age > 8 && age > 9`,
		},
		{
			"nested_parentheses_5",
			`(((((age > 18)))))`,
		},
		{
			"nested_ternary",
			`age > 18 ? (score > 90.0 ? "A" : "B") : (score > 50.0 ? "C" : "D")`,
		},
		{
			"nested_arithmetic",
			`((age + 5) * 2 - 3) / 4 + 1`,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			ast, issues := env.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				b.Fatal(issues.Err())
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := cel2sql.Convert(ast)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkConvertLargeExpression benchmarks very large CEL expressions
func BenchmarkConvertLargeExpression(b *testing.B) {
	env := setupSimpleBenchmarkEnv(b)

	// Large AND expression
	largeAnd := `age > 0`
	for i := 1; i < 20; i++ {
		largeAnd += ` && name != "test` + string(rune('0'+i)) + `"`
	}

	tests := []struct {
		name       string
		expression string
	}{
		{"mixed_conditions_20", `age > 18 && active && score > 80.0 && name.startsWith("A") && email.contains("@") && tags.exists(t, t == "premium") && age < 65 && score < 100.0 && name.endsWith("z") && email.endsWith(".com") && age != 30 && age != 40 && age != 50 && score != 0.0 && name != "" && email != "" && active == true && age >= 18 && score >= 0.0 && age <= 100`},
		{"large_and_chain", largeAnd},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			ast, issues := env.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				b.Fatal(issues.Err())
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := cel2sql.Convert(ast)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkConvertTimestamps benchmarks timestamp and date operations
func BenchmarkConvertTimestamps(b *testing.B) {
	env := setupTimestampBenchmarkEnv(b)

	tests := []struct {
		name       string
		expression string
	}{
		{"timestamp_comparison", `created_at > timestamp(datetime("2024-01-01T12:00:00"), "UTC")`},
		{"date_function", `date("2024-01-01")`},
		{"datetime_function", `datetime("2024-01-01T12:00:00")`},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			ast, issues := env.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				b.Fatal(issues.Err())
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := cel2sql.Convert(ast)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkConvertStringOperations benchmarks string operations
func BenchmarkConvertStringOperations(b *testing.B) {
	env := setupSimpleBenchmarkEnv(b)

	tests := []struct {
		name       string
		expression string
	}{
		{"startsWith", `name.startsWith("test")`},
		{"endsWith", `email.endsWith(".com")`},
		{"contains", `name.contains("admin")`},
		{"concatenation", `name + " " + email`},
		{"multiple_string_ops", `name.startsWith("A") && email.endsWith(".com") && name.contains("test")`},
		{"split_basic", `name.split(",").size() > 0`},
		{"split_with_limit", `name.split(",", 3).size() == 3`},
		{"join_basic", `["a", "b", "c"].join(",")`},
		{"join_no_delimiter", `["a", "b", "c"].join()`},
		{"format_simple", `"%s: %s".format([name, email])`},
		{"format_multiple_args", `"%s is %d years old".format(["John", 30])`},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			ast, issues := env.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				b.Fatal(issues.Err())
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := cel2sql.Convert(ast)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkAnalyzeQuery benchmarks the query analysis feature
func BenchmarkAnalyzeQuery(b *testing.B) {
	env, schemas := setupSchemaBenchmarkEnv(b)

	tests := []struct {
		name       string
		expression string
	}{
		{"simple_comparison", `person.age > 18`},
		{"json_path", `person.metadata.verified == true`},
		{"regex_pattern", `person.email.matches(r"@example\.com$")`},
		{"array_operation", `"premium" in person.tags`},
		{"complex_query", `person.age > 18 && person.metadata.verified && person.email.matches(r"@example\.com$")`},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			ast, issues := env.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				b.Fatal(issues.Err())
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, err := cel2sql.AnalyzeQuery(ast, cel2sql.WithSchemas(schemas))
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkConvertWithOptions benchmarks conversion with various options
func BenchmarkConvertWithOptions(b *testing.B) {
	env, schemas := setupSchemaBenchmarkEnv(b)

	tests := []struct {
		name       string
		expression string
		opts       []cel2sql.ConvertOption
	}{
		{"no_options", `person.age > 18`, nil},
		{"with_schemas", `person.age > 18`, []cel2sql.ConvertOption{cel2sql.WithSchemas(schemas)}},
		{"with_max_depth", `person.age > 18`, []cel2sql.ConvertOption{cel2sql.WithMaxDepth(150)}},
		{"with_max_output", `person.age > 18`, []cel2sql.ConvertOption{cel2sql.WithMaxOutputLength(100000)}},
		{"all_options", `person.age > 18`, []cel2sql.ConvertOption{
			cel2sql.WithSchemas(schemas),
			cel2sql.WithMaxDepth(150),
			cel2sql.WithMaxOutputLength(100000),
		}},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			ast, issues := env.Compile(tt.expression)
			if issues != nil && issues.Err() != nil {
				b.Fatal(issues.Err())
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := cel2sql.Convert(ast, tt.opts...)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
