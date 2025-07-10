package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
	"github.com/spandigital/cel2sql"
	"github.com/spandigital/cel2sql/pg"
	"github.com/spandigital/cel2sql/sqltypes"
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
		sqltypes.SQLTypeDeclarations,
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

func exampleWithDynamicSchema(ctx context.Context) {
	fmt.Println("=== Example 2: Dynamic Schema Loading ===")

	// Connect to PostgreSQL database
	// Update connection string with your database credentials
	connStr := "postgres://user:password@localhost:5432/mydb?sslmode=disable"

	provider, err := pg.NewTypeProviderWithConnection(ctx, connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer provider.Close()

	// Load table schema dynamically
	tableName := "users"
	err = provider.LoadTableSchema(ctx, tableName)
	if err != nil {
		log.Fatal("Failed to load table schema:", err)
	}

	// Create CEL environment with dynamically loaded schema
	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		sqltypes.SQLTypeDeclarations,
		cel.Variable("user", cel.ObjectType(tableName)),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Example CEL expressions using dynamically loaded schema
	expressions := []string{
		`user.name == "Jane Smith"`,
		`user.age >= 25 && user.age <= 65`,
		`user.is_active == true`,
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

	fmt.Println("Note: The exact SQL output depends on the actual table schema in your database.")
}
