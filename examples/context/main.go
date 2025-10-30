// Package main demonstrates context usage with cel2sql for cancellation and timeout support.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
)

func main() {
	fmt.Println("CEL to SQL Context Usage Examples")
	fmt.Println("==================================")
	fmt.Println()

	// Example 1: Without context (backward compatible)
	fmt.Println("1. Without Context (backward compatible):")
	exampleWithoutContext()

	// Example 2: With active context
	fmt.Println("\n2. With Active Context:")
	exampleWithActiveContext()

	// Example 3: With timeout
	fmt.Println("\n3. With Timeout:")
	exampleWithTimeout()

	// Example 4: With cancellation
	fmt.Println("\n4. With Cancellation:")
	exampleWithCancellation()

	// Example 5: Complex expression with context and schemas
	fmt.Println("\n5. Complex Expression with Context and Schemas:")
	exampleComplexWithOptions()
}

func exampleWithoutContext() {
	env, err := cel.NewEnv(
		cel.Variable("age", cel.IntType),
		cel.Variable("score", cel.IntType),
	)
	if err != nil {
		log.Fatal(err)
	}

	ast, issues := env.Compile(`age >= 18 && score > 80`)
	if issues.Err() != nil {
		log.Fatal(issues.Err())
	}

	// Convert without context - works as before
	sql, err := cel2sql.Convert(ast)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("  CEL: age >= 18 && score > 80\n")
	fmt.Printf("  SQL: %s\n", sql)
}

func exampleWithActiveContext() {
	env, err := cel.NewEnv(
		cel.Variable("age", cel.IntType),
		cel.Variable("score", cel.IntType),
	)
	if err != nil {
		log.Fatal(err)
	}

	ast, issues := env.Compile(`age >= 18 && score > 80`)
	if issues.Err() != nil {
		log.Fatal(issues.Err())
	}

	// Convert with background context
	ctx := context.Background()
	sql, err := cel2sql.Convert(ast, cel2sql.WithContext(ctx))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("  CEL: age >= 18 && score > 80\n")
	fmt.Printf("  SQL: %s\n", sql)
	fmt.Printf("  ✓ Conversion succeeded with active context\n")
}

func exampleWithTimeout() {
	env, err := cel.NewEnv()
	if err != nil {
		log.Fatal(err)
	}

	ast, issues := env.Compile(`(1 + 2) * (3 + 4) == (5 + 6) * (7 + 8)`)
	if issues.Err() != nil {
		log.Fatal(issues.Err())
	}

	// Set a reasonable timeout for conversion
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql, err := cel2sql.Convert(ast, cel2sql.WithContext(ctx))
	if err != nil {
		log.Printf("  ✗ Conversion failed: %v\n", err)
		return
	}

	fmt.Printf("  CEL: (1 + 2) * (3 + 4) == (5 + 6) * (7 + 8)\n")
	fmt.Printf("  SQL: %s\n", sql)
	fmt.Printf("  ✓ Conversion completed within timeout\n")
}

func exampleWithCancellation() {
	env, err := cel.NewEnv()
	if err != nil {
		log.Fatal(err)
	}

	ast, issues := env.Compile(`(1 + 2) * (3 + 4) == (5 + 6) * (7 + 8)`)
	if issues.Err() != nil {
		log.Fatal(issues.Err())
	}

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Simulate cancellation (in real code, this might be triggered by user action or shutdown)
	cancel()

	sql, err := cel2sql.Convert(ast, cel2sql.WithContext(ctx))
	if err != nil {
		fmt.Printf("  ✓ Conversion was cancelled as expected: %v\n", err)
		return
	}

	// Should not reach here
	fmt.Printf("  ✗ Unexpected: Conversion succeeded despite cancellation: %s\n", sql)
}

func exampleComplexWithOptions() {
	// Define a schema with JSON fields
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "integer"},
		{Name: "name", Type: "text"},
		{Name: "age", Type: "integer"},
		{Name: "scores", Type: "integer", Repeated: true},
		{Name: "metadata", Type: "jsonb", IsJSON: true, IsJSONB: true},
	})

	provider := pg.NewTypeProvider(map[string]pg.Schema{"person": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("person", cel.ObjectType("person")),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Complex expression with comprehensions and JSON access
	ast, issues := env.Compile(`person.scores.all(s, s > 50) && person.metadata.active == "true" && person.age >= 18`)
	if issues.Err() != nil {
		log.Fatal(issues.Err())
	}

	// Use multiple options: context with timeout + schemas
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	schemas := provider.GetSchemas()
	sql, err := cel2sql.Convert(ast,
		cel2sql.WithContext(ctx),
		cel2sql.WithSchemas(schemas))

	if err != nil {
		log.Printf("  ✗ Conversion failed: %v\n", err)
		return
	}

	fmt.Printf("  CEL: person.scores.all(s, s > 50) && person.metadata.active == \"true\" && person.age >= 18\n")
	fmt.Printf("  SQL: %s\n", sql)
	fmt.Printf("  ✓ Complex conversion with context and schemas succeeded\n")
}
