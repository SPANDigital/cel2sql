// Package main demonstrates CEL comprehensions support in cel2sql with PostgreSQL integration.
// This example shows how to use all(), exists(), exists_one(), filter(), and map() comprehensions
// with both simple arrays and complex PostgreSQL schemas including nested structures.
package main

import (
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
	"github.com/spandigital/cel2sql/v2"
	"github.com/spandigital/cel2sql/v2/pg"
)

func main() {
	fmt.Println("CEL Comprehension to PostgreSQL SQL Examples")
	fmt.Println("===========================================")

	// Example 1: Simple comprehensions on arrays
	fmt.Println("\n1. Simple Array Comprehensions:")
	simpleExamples()

	// Example 2: PostgreSQL schema-based comprehensions
	fmt.Println("\n2. PostgreSQL Schema-based Comprehensions:")
	schemaExamples()

	// Example 3: Complex nested comprehensions
	fmt.Println("\n3. Nested and Complex Comprehensions:")
	complexExamples()
}

func simpleExamples() {
	env, err := cel.NewEnv()
	if err != nil {
		log.Fatal(err)
	}

	examples := []string{
		`[1, 2, 3, 4, 5].all(x, x > 0)`,
		`[1, 2, 3, 4, 5].exists(x, x > 3)`,
		`[1, 2, 3, 4, 5].exists_one(x, x == 3)`,
		`[1, 2, 3, 4, 5].filter(x, x % 2 == 0)`,
		`[1, 2, 3, 4, 5].map(x, x * 2)`,
	}

	for _, expr := range examples {
		ast, issues := env.Compile(expr)
		if issues.Err() != nil {
			fmt.Printf("  ❌ %s: %v\n", expr, issues.Err())
			continue
		}

		sql, err := cel2sql.Convert(ast)
		if err != nil {
			fmt.Printf("  ❌ %s: %v\n", expr, err)
			continue
		}

		fmt.Printf("  ✅ %s\n", expr)
		fmt.Printf("     → %s\n", sql)
	}
}

func schemaExamples() {
	// Define a PostgreSQL schema for employees
	schema := pg.Schema{
		{Name: "id", Type: "bigint", Repeated: false},
		{Name: "name", Type: "text", Repeated: false},
		{Name: "email", Type: "text", Repeated: false},
		{Name: "age", Type: "bigint", Repeated: false},
		{Name: "active", Type: "boolean", Repeated: false},
		{Name: "salary", Type: "double precision", Repeated: false},
		{Name: "skills", Type: "text", Repeated: true}, // Array field
	}

	provider := pg.NewTypeProvider(map[string]pg.Schema{"Employee": schema})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("employees", cel.ListType(cel.ObjectType("Employee"))),
		cel.Variable("emp", cel.ObjectType("Employee")),
	)
	if err != nil {
		log.Fatal(err)
	}

	examples := []string{
		`employees.all(e, e.salary > 50000)`,
		`employees.exists(e, e.age > 65)`,
		`employees.filter(e, e.active)`,
		`employees.map(e, e.name)`,
		`employees.filter(e, e.active).map(e, e.email)`,
		`emp.skills.exists(s, s == 'Go')`,
	}

	for _, expr := range examples {
		ast, issues := env.Compile(expr)
		if issues.Err() != nil {
			fmt.Printf("  ❌ %s: %v\n", expr, issues.Err())
			continue
		}

		sql, err := cel2sql.Convert(ast)
		if err != nil {
			fmt.Printf("  ❌ %s: %v\n", expr, err)
			continue
		}

		fmt.Printf("  ✅ %s\n", expr)
		fmt.Printf("     → %s\n", sql)
	}
}

func complexExamples() {
	// Define nested schema with addresses
	addressSchema := pg.Schema{
		{Name: "street", Type: "text", Repeated: false},
		{Name: "city", Type: "text", Repeated: false},
		{Name: "country", Type: "text", Repeated: false},
	}

	employeeSchema := pg.Schema{
		{Name: "id", Type: "bigint", Repeated: false},
		{Name: "name", Type: "text", Repeated: false},
		{Name: "address", Type: "composite", Schema: addressSchema},
		{Name: "skills", Type: "text", Repeated: true},
		{Name: "scores", Type: "bigint", Repeated: true},
	}

	provider := pg.NewTypeProvider(map[string]pg.Schema{
		"Address":  addressSchema,
		"Employee": employeeSchema,
	})

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("employees", cel.ListType(cel.ObjectType("Employee"))),
	)
	if err != nil {
		log.Fatal(err)
	}

	examples := []string{
		`employees.filter(e, e.address.city == 'New York')`,
		`employees.exists(e, e.skills.exists(s, s == 'Go'))`,
		`employees.map(e, e.address.city)`,
		`employees.filter(e, e.scores.all(s, s >= 80))`,
	}

	for _, expr := range examples {
		ast, issues := env.Compile(expr)
		if issues.Err() != nil {
			fmt.Printf("  ❌ %s: %v\n", expr, issues.Err())
			continue
		}

		sql, err := cel2sql.Convert(ast)
		if err != nil {
			fmt.Printf("  ❌ %s: %v\n", expr, err)
			continue
		}

		fmt.Printf("  ✅ %s\n", expr)
		fmt.Printf("     → %s\n", sql)
	}
}
