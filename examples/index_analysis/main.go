// Package main demonstrates how to use cel2sql.AnalyzeQuery() for index recommendations
package main

import (
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
)

func main() {
	fmt.Println("cel2sql - Index Analysis Example")
	fmt.Println("===================================")

	// Define a schema for a hypothetical users table
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "bigint"},
		{Name: "name", Type: "text"},
		{Name: "email", Type: "text"},
		{Name: "age", Type: "integer"},
		{Name: "active", Type: "boolean"},
		{Name: "tags", Type: "text", Repeated: true},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"users": schema})

	// Create CEL environment
	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("users", cel.ObjectType("users")),
	)
	if err != nil {
		log.Fatalf("Failed to create CEL environment: %v", err)
	}

	// Example queries demonstrating different index recommendation types
	examples := []struct {
		name        string
		description string
		expression  string
	}{
		{
			name:        "Simple Comparison",
			description: "Basic comparison operations that benefit from B-tree indexes",
			expression:  `users.age > 18 && users.active == true`,
		},
		{
			name:        "Regex Matching",
			description: "Pattern matching that benefits from GIN index with pg_trgm",
			expression:  `users.email.matches(r"@example\.com$")`,
		},
		{
			name:        "JSON Path Access",
			description: "JSONB field access that benefits from GIN indexes",
			expression:  `users.metadata.verified == true && users.metadata.tier == "premium"`,
		},
		{
			name:        "Array Operations",
			description: "Array membership tests that benefit from GIN indexes",
			expression:  `"premium" in users.tags`,
		},
		{
			name:        "Array Comprehensions",
			description: "Array comprehensions that benefit from GIN indexes",
			expression:  `users.tags.all(t, t.startsWith("valid_"))`,
		},
		{
			name:        "Complex Query",
			description: "Multiple patterns requiring different index types",
			expression: `users.age > 21 &&
			             users.email.matches(r"@(gmail|yahoo)\.com$") &&
			             users.metadata.subscription.active == true &&
			             "verified" in users.tags`,
		},
	}

	// Analyze each query and display recommendations
	for i, ex := range examples {
		fmt.Printf("%d. %s\n", i+1, ex.name)
		fmt.Printf("   Description: %s\n", ex.description)
		fmt.Printf("   CEL Expression: %s\n\n", ex.expression)

		// Compile the CEL expression
		ast, issues := env.Compile(ex.expression)
		if issues != nil && issues.Err() != nil {
			log.Printf("   ERROR: Failed to compile: %v\n\n", issues.Err())
			continue
		}

		// Analyze the query
		sql, recommendations, err := cel2sql.AnalyzeQuery(ast,
			cel2sql.WithSchemas(provider.GetSchemas()))
		if err != nil {
			log.Printf("   ERROR: Failed to analyze: %v\n\n", err)
			continue
		}

		// Display the generated SQL
		fmt.Printf("   Generated SQL:\n   %s\n\n", sql)

		// Display index recommendations
		if len(recommendations) == 0 {
			fmt.Printf("   No index recommendations (query uses constants or simple conditions)\n\n")
		} else {
			fmt.Printf("   Index Recommendations (%d):\n", len(recommendations))
			for j, rec := range recommendations {
				fmt.Printf("   [%d] Column: %s\n", j+1, rec.Column)
				fmt.Printf("       Type: %s\n", rec.IndexType)
				fmt.Printf("       Reason: %s\n", rec.Reason)
				fmt.Printf("       SQL: %s\n", rec.Expression)
				fmt.Println()
			}
		}

		fmt.Println("   " + string(make([]byte, 60)))
		fmt.Println()
	}

	// Summary
	fmt.Println("\nSummary")
	fmt.Println("=======")
	fmt.Println("Index recommendations help optimize query performance by:")
	fmt.Println("  • B-tree indexes: Fast equality and range queries on scalar columns")
	fmt.Println("  • GIN indexes: Efficient JSON path access and array operations")
	fmt.Println("  • GIN with pg_trgm: Fast regex pattern matching on text columns")
	fmt.Println()
	fmt.Println("To apply recommendations:")
	fmt.Println("  1. Review each recommendation and its reason")
	fmt.Println("  2. Adjust table_name to your actual table name")
	fmt.Println("  3. Execute the CREATE INDEX statements on your database")
	fmt.Println("  4. Monitor query performance improvements")
}
