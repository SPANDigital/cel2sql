// Package main demonstrates basic usage of cel2sql with a predefined schema.
package main

import (
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
	"github.com/spandigital/cel2sql"
	"github.com/spandigital/cel2sql/pg"
)

func main() {
	// Define a PostgreSQL table schema
	employeeSchema := pg.Schema{
		{Name: "name", Type: "text", Repeated: false},
		{Name: "age", Type: "integer", Repeated: false},
		{Name: "department", Type: "text", Repeated: false},
		{Name: "hired_at", Type: "timestamp with time zone", Repeated: false},
		{Name: "active", Type: "boolean", Repeated: false},
	}

	// Create CEL environment with PostgreSQL type provider
	env, err := cel.NewEnv(
		cel.CustomTypeProvider(pg.NewTypeProvider(map[string]pg.Schema{
			"Employee": employeeSchema,
		})),
		cel.Variable("employee", cel.ObjectType("Employee")),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Compile a CEL expression
	ast, issues := env.Compile(`employee.name == "John Doe" && employee.age >= 25 && employee.active`)
	if issues != nil && issues.Err() != nil {
		log.Fatal(issues.Err())
	}

	// Convert to SQL
	sqlCondition, err := cel2sql.Convert(ast)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("CEL Expression: employee.name == \"John Doe\" && employee.age >= 25 && employee.active")
	fmt.Println("PostgreSQL SQL:", sqlCondition)
}
