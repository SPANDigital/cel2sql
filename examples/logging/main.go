// Package main demonstrates cel2sql structured logging with log/slog.
// Shows both JSON and text handlers with debug-level logging enabled.
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/google/cel-go/cel"
	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
)

func main() {
	// Create a logger with JSON output and debug level
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Define schema with JSON fields and arrays
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "name", Type: "text"},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
		{Name: "tags", Type: "text", Repeated: true, ElementType: "text"},
	})
	provider := pg.NewTypeProvider(map[string]pg.Schema{"users": schema})

	// Create CEL environment
	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("users", cel.ObjectType("users")),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating CEL environment: %v\n", err)
		os.Exit(1)
	}

	// Compile a complex CEL expression with JSON access and comprehensions
	expression := `users.metadata.active == true && users.tags.exists(t, t == "premium")`
	ast, issues := env.Compile(expression)
	if issues != nil && issues.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error compiling CEL expression: %v\n", issues.Err())
		os.Exit(1)
	}

	fmt.Println("CEL Expression:")
	fmt.Println(expression)
	fmt.Println()
	fmt.Println("Conversion Log (JSON format):")
	fmt.Println("---")

	// Convert with logging enabled
	sql, err := cel2sql.Convert(ast,
		cel2sql.WithSchemas(provider.GetSchemas()),
		cel2sql.WithLogger(logger),
	)

	fmt.Println("---")
	fmt.Println()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting to SQL: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Generated SQL:")
	fmt.Println(sql)
	fmt.Println()

	// Example with text handler for human-readable logs
	fmt.Println("=== Example with Text Handler ===")
	fmt.Println()

	textLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	expression2 := `users.name.matches("(?i)admin.*") && users.id > 100`
	ast2, issues2 := env.Compile(expression2)
	if issues2 != nil && issues2.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error compiling CEL expression: %v\n", issues2.Err())
		os.Exit(1)
	}

	fmt.Println("CEL Expression:")
	fmt.Println(expression2)
	fmt.Println()
	fmt.Println("Conversion Log (Text format):")
	fmt.Println("---")

	sql2, err := cel2sql.Convert(ast2,
		cel2sql.WithSchemas(provider.GetSchemas()),
		cel2sql.WithLogger(textLogger),
	)

	fmt.Println("---")
	fmt.Println()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting to SQL: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Generated SQL:")
	fmt.Println(sql2)
}
