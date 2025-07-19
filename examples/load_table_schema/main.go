// Package main demonstrates loading table schema dynamically from PostgreSQL.
package main

import (
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
	"github.com/spandigital/cel2sql/v2"
	"github.com/spandigital/cel2sql/v2/pg"
)

func main() {
	// Example 1: Using pre-defined schema
	exampleWithPredefinedSchema()

	// Example 2: Loading schema dynamically from database
	// Note: This requires a running PostgreSQL database
	// Uncomment the line below and update the connection string
	// ctx := context.Background()
	// exampleWithDynamicSchema(ctx)
}

func exampleWithPredefinedSchema() {
	fmt.Println("=== Example 1: Pre-defined Schema ===")

	// Define schema manually
	userSchema := pg.Schema{
		{Name: "id", Type: "integer", Repeated: false},
		{Name: "name", Type: "text", Repeated: false},
		{Name: "email", Type: "text", Repeated: false},
		{Name: "age", Type: "integer", Repeated: false},
		{Name: "created_at", Type: "timestamp with time zone", Repeated: false},
		{Name: "is_active", Type: "boolean", Repeated: false},
		{Name: "tags", Type: "text", Repeated: true}, // Array field
	}

	// Create type provider with predefined schema
	provider := pg.NewTypeProvider(map[string]pg.Schema{
		"users": userSchema,
	})

	// Create CEL environment
	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("user", cel.ObjectType("users")),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Example CEL expressions
	expressions := []string{
		`user.name == "John Doe"`,
		`user.age > 30 && user.is_active`,
		`user.email.contains("@example.com")`,
		`"admin" in user.tags`,
		`user.created_at > timestamp("2023-01-01T00:00:00Z")`,
	}

	for _, expr := range expressions {
		ast, issues := env.Compile(expr)
		if issues != nil && issues.Err() != nil {
			log.Printf("Error compiling %s: %v", expr, issues.Err())
			continue
		}

		sqlCondition, err := cel2sql.Convert(ast)
		if err != nil {
			log.Printf("Error converting %s: %v", expr, err)
			continue
		}

		fmt.Printf("CEL: %s\n", expr)
		fmt.Printf("SQL: %s\n\n", sqlCondition)
	}
}
