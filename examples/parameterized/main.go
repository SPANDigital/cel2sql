// Package main demonstrates parameterized query support in cel2sql.
//
// Parameterized queries provide:
//   - Query plan caching for better performance
//   - Additional SQL injection protection
//   - Better monitoring (same query pattern in logs)
//
// Run with: go run main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/google/cel-go/cel"
	_ "github.com/lib/pq"
	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func main() {
	ctx := context.Background()

	// Start a PostgreSQL 17 container for demonstration
	fmt.Println("Starting PostgreSQL 17 container...")
	pgContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2)),
	)
	if err != nil {
		log.Fatalf("Failed to start PostgreSQL container: %v", err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(pgContainer); err != nil {
			log.Printf("Failed to terminate container: %v", err)
		}
	}()

	// Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = testcontainers.TerminateContainer(pgContainer)
		//nolint:gocritic // cleanup handled explicitly before exit
		log.Fatalf("Failed to get connection string: %v", err)
	}

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	// Setup schema
	if err := setupSchema(db); err != nil {
		log.Fatalf("Failed to setup schema: %v", err)
	}

	// Insert test data
	if err := insertTestData(db); err != nil {
		log.Fatalf("Failed to insert test data: %v", err)
	}

	// Define schema for CEL
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "id", Type: "bigint"},
		{Name: "name", Type: "text"},
		{Name: "email", Type: "text"},
		{Name: "age", Type: "int"},
		{Name: "salary", Type: "double precision"},
		{Name: "department", Type: "text"},
		{Name: "active", Type: "boolean"},
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

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("PARAMETERIZED QUERIES DEMONSTRATION")
	fmt.Println(strings.Repeat("=", 80))

	// Example 1: Simple parameterized query
	demonstrateSimpleParameterized(ctx, env, db)

	// Example 2: Complex parameterized query
	demonstrateComplexParameterized(ctx, env, db)

	// Example 3: Prepared statements with parameterized queries
	demonstratePreparedStatement(ctx, env, db)

	// Example 4: Query plan caching demonstration
	demonstrateQueryPlanCaching(ctx, env, db)

	// Example 5: Comparison with non-parameterized queries
	demonstrateComparison(ctx, env, db)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("✓ All examples completed successfully!")
	fmt.Println(strings.Repeat("=", 80))
}

func setupSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE users (
			id BIGINT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			age INT NOT NULL,
			salary DOUBLE PRECISION NOT NULL,
			department TEXT NOT NULL,
			active BOOLEAN DEFAULT TRUE
		)
	`)
	return err
}

func insertTestData(db *sql.DB) error {
	data := []struct {
		id         int64
		name       string
		email      string
		age        int
		salary     float64
		department string
		active     bool
	}{
		{1, "Alice Smith", "alice@example.com", 30, 75000.00, "Engineering", true},
		{2, "Bob Johnson", "bob@example.com", 25, 60000.00, "Marketing", true},
		{3, "Carol Williams", "carol@example.com", 35, 85000.00, "Engineering", true},
		{4, "David Brown", "david@example.com", 28, 70000.00, "Sales", true},
		{5, "Eve Davis", "eve@example.com", 32, 80000.00, "Engineering", false},
	}

	for _, u := range data {
		_, err := db.Exec(
			`INSERT INTO users (id, name, email, age, salary, department, active)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			u.id, u.name, u.email, u.age, u.salary, u.department, u.active,
		)
		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", u.name, err)
		}
	}

	fmt.Println("✓ Inserted 5 test users")
	return nil
}

func demonstrateSimpleParameterized(_ context.Context, env *cel.Env, db *sql.DB) {
	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("Example 1: Simple Parameterized Query")
	fmt.Println(strings.Repeat("-", 80))

	// CEL expression: Find users older than 28
	celExpr := `users.age > 28`
	fmt.Printf("CEL Expression: %s\n\n", celExpr)

	// Compile CEL expression
	ast, issues := env.Compile(celExpr)
	if issues != nil && issues.Err() != nil {
		log.Fatalf("Failed to compile CEL: %v", issues.Err())
	}

	// Convert to parameterized SQL
	result, err := cel2sql.ConvertParameterized(ast)
	if err != nil {
		log.Fatalf("Failed to convert to SQL: %v", err)
	}

	fmt.Printf("Generated SQL:    %s\n", result.SQL)
	fmt.Printf("Parameters:       %v\n", result.Parameters)
	fmt.Printf("Parameter Types:  ")
	for i, p := range result.Parameters {
		fmt.Printf("%T", p)
		if i < len(result.Parameters)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println()

	// Execute query
	// Note: result.SQL is generated by cel2sql.ConvertParameterized(), not from user input
	//nolint:perfsprint // fmt.Sprintf preferred over concatenation for SQL security
	// #nosec G201 - SQL string is from trusted conversion function
	query := fmt.Sprintf("SELECT id, name, age FROM users WHERE %s", result.SQL)
	rows, err := db.Query(query, result.Parameters...)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}()

	fmt.Println("\nResults:")
	fmt.Println(strings.Repeat("-", 50))
	for rows.Next() {
		var id int64
		var name string
		var age int
		if err := rows.Scan(&id, &name, &age); err != nil {
			_ = rows.Close()
			//nolint:gocritic // cleanup handled explicitly before exit
			log.Fatalf("Failed to scan row: %v", err)
		}
		fmt.Printf("ID: %d, Name: %-20s Age: %d\n", id, name, age)
	}
}

func demonstrateComplexParameterized(_ context.Context, env *cel.Env, db *sql.DB) {
	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("Example 2: Complex Parameterized Query")
	fmt.Println(strings.Repeat("-", 80))

	// CEL expression with multiple conditions
	celExpr := `users.department == "Engineering" && users.salary > 70000.0 && users.active == true`
	fmt.Printf("CEL Expression: %s\n\n", celExpr)

	ast, issues := env.Compile(celExpr)
	if issues != nil && issues.Err() != nil {
		log.Fatalf("Failed to compile CEL: %v", issues.Err())
	}

	result, err := cel2sql.ConvertParameterized(ast)
	if err != nil {
		log.Fatalf("Failed to convert to SQL: %v", err)
	}

	fmt.Printf("Generated SQL:    %s\n", result.SQL)
	fmt.Printf("Parameters:       %v\n", result.Parameters)
	fmt.Println("\nNote: TRUE/FALSE/NULL are kept inline for query plan optimization")

	// Execute query
	// Note: result.SQL is generated by cel2sql.ConvertParameterized(), not from user input
	//nolint:perfsprint // fmt.Sprintf preferred over concatenation for SQL security
	// #nosec G201 - SQL string is from trusted conversion function
	query := fmt.Sprintf("SELECT id, name, department, salary FROM users WHERE %s", result.SQL)
	rows, err := db.Query(query, result.Parameters...)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}()

	fmt.Println("\nResults:")
	fmt.Println(strings.Repeat("-", 70))
	for rows.Next() {
		var id int64
		var name string
		var department string
		var salary float64
		if err := rows.Scan(&id, &name, &department, &salary); err != nil {
			_ = rows.Close()
			//nolint:gocritic // cleanup handled explicitly before exit
			log.Fatalf("Failed to scan row: %v", err)
		}
		fmt.Printf("ID: %d, Name: %-20s Dept: %-12s Salary: $%.2f\n", id, name, department, salary)
	}
}

func demonstratePreparedStatement(_ context.Context, env *cel.Env, db *sql.DB) {
	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("Example 3: Prepared Statement with Different Parameters")
	fmt.Println(strings.Repeat("-", 80))

	celExpr := `users.age > 25`
	ast, issues := env.Compile(celExpr)
	if issues != nil && issues.Err() != nil {
		log.Fatalf("Failed to compile CEL: %v", issues.Err())
	}

	result, err := cel2sql.ConvertParameterized(ast)
	if err != nil {
		log.Fatalf("Failed to convert to SQL: %v", err)
	}

	// Note: result.SQL is generated by cel2sql.ConvertParameterized(), not from user input
	//nolint:perfsprint // fmt.Sprintf preferred over concatenation for SQL security
	// #nosec G201 - SQL string is from trusted conversion function
	query := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE %s", result.SQL)
	fmt.Printf("Prepared Query: %s\n\n", query)

	// Prepare the statement once
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("Failed to close statement: %v", err)
		}
	}()

	// Execute with different age thresholds
	ages := []int{25, 28, 30, 32}
	fmt.Println("Same query, different parameters (demonstrates plan caching):")
	fmt.Println(strings.Repeat("-", 50))

	for _, age := range ages {
		var count int
		if err := stmt.QueryRow(age).Scan(&count); err != nil {
			_ = stmt.Close()
			//nolint:gocritic // cleanup handled explicitly before exit
			log.Fatalf("Failed to execute query: %v", err)
		}
		fmt.Printf("Users older than %d: %d\n", age, count)
	}
}

func demonstrateQueryPlanCaching(_ context.Context, env *cel.Env, db *sql.DB) {
	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("Example 4: Query Plan Caching Benefits")
	fmt.Println(strings.Repeat("-", 80))

	celExpr := `users.salary > 65000.0`
	ast, issues := env.Compile(celExpr)
	if issues != nil && issues.Err() != nil {
		log.Fatalf("Failed to compile CEL: %v", issues.Err())
	}

	result, err := cel2sql.ConvertParameterized(ast)
	if err != nil {
		log.Fatalf("Failed to convert to SQL: %v", err)
	}

	// Note: result.SQL is generated by cel2sql.ConvertParameterized(), not from user input
	//nolint:perfsprint // fmt.Sprintf preferred over concatenation for SQL security
	// #nosec G201 - SQL string is from trusted conversion function
	query := fmt.Sprintf("SELECT id, name, salary FROM users WHERE %s", result.SQL)

	fmt.Println("Executing the same query structure with different salary thresholds:")
	fmt.Println("PostgreSQL will cache and reuse the query plan for better performance.")
	fmt.Println(strings.Repeat("-", 70))

	salaries := []float64{65000.0, 75000.0, 80000.0}
	for _, salary := range salaries {
		fmt.Printf("\nSalary threshold: $%.2f\n", salary)
		rows, err := db.Query(query, salary)
		if err != nil {
			log.Fatalf("Failed to execute query: %v", err)
		}

		count := 0
		for rows.Next() {
			var id int64
			var name string
			var sal float64
			if err := rows.Scan(&id, &name, &sal); err != nil {
				log.Fatalf("Failed to scan row: %v", err)
			}
			count++
			fmt.Printf("  %s: $%.2f\n", name, sal)
		}
		if err := rows.Close(); err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
		fmt.Printf("  Found %d users\n", count)
	}

	fmt.Println("\n✓ All executions used the same cached query plan!")
}

func demonstrateComparison(_ context.Context, env *cel.Env, db *sql.DB) {
	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("Example 5: Parameterized vs Non-Parameterized Comparison")
	fmt.Println(strings.Repeat("-", 80))

	celExpr := `users.age > 28 && users.name == "Alice Smith"`
	ast, issues := env.Compile(celExpr)
	if issues != nil && issues.Err() != nil {
		log.Fatalf("Failed to compile CEL: %v", issues.Err())
	}

	// Non-parameterized
	nonParamSQL, err := cel2sql.Convert(ast)
	if err != nil {
		log.Fatalf("Failed to convert (non-param): %v", err)
	}

	// Parameterized
	paramResult, err := cel2sql.ConvertParameterized(ast)
	if err != nil {
		log.Fatalf("Failed to convert (param): %v", err)
	}

	fmt.Println("Non-Parameterized (inline values):")
	fmt.Printf("  SQL: %s\n", nonParamSQL)
	fmt.Println("  Parameters: none (values embedded in SQL)")
	fmt.Println("  ✗ Different query for each value set")
	fmt.Println("  ✗ No query plan caching")
	fmt.Println("  ✗ Each query looks different in logs")

	fmt.Println("\nParameterized (placeholders + parameters):")
	fmt.Printf("  SQL: %s\n", paramResult.SQL)
	fmt.Printf("  Parameters: %v\n", paramResult.Parameters)
	fmt.Println("  ✓ Same query structure, different values")
	fmt.Println("  ✓ Query plan caching enabled")
	fmt.Println("  ✓ Same query pattern in logs/metrics")
	fmt.Println("  ✓ Additional SQL injection protection")

	// Execute both and verify same results
	// Note: nonParamSQL is generated by cel2sql.Convert(), not from user input
	//nolint:perfsprint // fmt.Sprintf preferred over concatenation for SQL security
	// #nosec G201 - SQL string is from trusted conversion function
	query1 := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE %s", nonParamSQL)
	var count1 int
	if err := db.QueryRow(query1).Scan(&count1); err != nil {
		log.Fatalf("Failed to execute non-param query: %v", err)
	}

	// Note: paramResult.SQL is generated by cel2sql.ConvertParameterized(), not from user input
	//nolint:perfsprint // fmt.Sprintf preferred over concatenation for SQL security
	// #nosec G201 - SQL string is from trusted conversion function
	query2 := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE %s", paramResult.SQL)
	var count2 int
	if err := db.QueryRow(query2, paramResult.Parameters...).Scan(&count2); err != nil {
		log.Fatalf("Failed to execute param query: %v", err)
	}

	fmt.Printf("\nBoth queries returned the same result: %d row(s)\n", count1)
	if count1 != count2 {
		log.Fatalf("Results differ! Non-param: %d, Param: %d", count1, count2)
	}
}
