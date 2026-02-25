// Package main demonstrates how to use cel2sql.AnalyzeQuery() for index recommendations
package main

import (
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/dialect"
	dialectbq "github.com/spandigital/cel2sql/v3/dialect/bigquery"
	dialectduckdb "github.com/spandigital/cel2sql/v3/dialect/duckdb"
	dialectmysql "github.com/spandigital/cel2sql/v3/dialect/mysql"
	dialectpg "github.com/spandigital/cel2sql/v3/dialect/postgres"
	dialectsqlite "github.com/spandigital/cel2sql/v3/dialect/sqlite"
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

	// Analyze each query and display recommendations (PostgreSQL default)
	fmt.Println("\n--- PostgreSQL (default) ---")
	for i, ex := range examples {
		analyzeExample(env, provider, i, ex.name, ex.description, ex.expression)
	}

	// Multi-dialect examples
	fmt.Println("\n===================================")
	fmt.Println("Multi-Dialect Index Recommendations")
	fmt.Println("===================================")

	// Use a simple comparison query to show dialect differences
	comparisonExpr := `users.age > 21 && users.metadata.verified == true`

	dialectExamples := []struct {
		name    string
		dialect dialect.Dialect
	}{
		{"PostgreSQL", dialectpg.New()},
		{"MySQL", dialectmysql.New()},
		{"SQLite", dialectsqlite.New()},
		{"DuckDB", dialectduckdb.New()},
		{"BigQuery", dialectbq.New()},
	}

	ast, issues := env.Compile(comparisonExpr)
	if issues != nil && issues.Err() != nil {
		log.Fatalf("Failed to compile expression: %v", issues.Err())
	}

	for _, de := range dialectExamples {
		fmt.Printf("\n--- %s ---\n", de.name)
		fmt.Printf("   CEL Expression: %s\n\n", comparisonExpr)

		_, recommendations, err := cel2sql.AnalyzeQuery(ast,
			cel2sql.WithSchemas(provider.GetSchemas()),
			cel2sql.WithDialect(de.dialect))
		if err != nil {
			log.Printf("   ERROR: Failed to analyze: %v\n\n", err)
			continue
		}

		if len(recommendations) == 0 {
			fmt.Printf("   No index recommendations\n")
		} else {
			for j, rec := range recommendations {
				fmt.Printf("   [%d] Column: %s\n", j+1, rec.Column)
				fmt.Printf("       Type: %s\n", rec.IndexType)
				fmt.Printf("       Reason: %s\n", rec.Reason)
				fmt.Printf("       DDL: %s\n", rec.Expression)
				fmt.Println()
			}
		}
	}

	// Summary
	fmt.Println("\nSummary")
	fmt.Println("=======")
	fmt.Println("Index recommendations are dialect-aware:")
	fmt.Println("  PostgreSQL: B-tree, GIN, GIN with pg_trgm")
	fmt.Println("  MySQL:      B-tree, FULLTEXT, functional JSON indexes")
	fmt.Println("  SQLite:     B-tree (limited index types)")
	fmt.Println("  DuckDB:     ART (Adaptive Radix Tree)")
	fmt.Println("  BigQuery:   Clustering keys, Search indexes")
	fmt.Println()
	fmt.Println("Use WithDialect() to get dialect-specific recommendations:")
	fmt.Println("  sql, recs, err := cel2sql.AnalyzeQuery(ast,")
	fmt.Println("      cel2sql.WithDialect(mysql.New()),")
	fmt.Println("      cel2sql.WithSchemas(schemas))")
}

func analyzeExample(env *cel.Env, provider pg.TypeProvider, idx int, name, description, expression string) {
	fmt.Printf("%d. %s\n", idx+1, name)
	fmt.Printf("   Description: %s\n", description)
	fmt.Printf("   CEL Expression: %s\n\n", expression)

	// Compile the CEL expression
	ast, issues := env.Compile(expression)
	if issues != nil && issues.Err() != nil {
		log.Printf("   ERROR: Failed to compile: %v\n\n", issues.Err())
		return
	}

	// Analyze the query
	sql, recommendations, err := cel2sql.AnalyzeQuery(ast,
		cel2sql.WithSchemas(provider.GetSchemas()))
	if err != nil {
		log.Printf("   ERROR: Failed to analyze: %v\n\n", err)
		return
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
