// Package main demonstrates CEL string extension functions in cel2sql
package main

import (
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/ext"
	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/pg"
)

func main() {
	fmt.Println("CEL String Extension Functions Examples")
	fmt.Println("==========================================")

	// Create a schema for our person table
	schema := pg.NewSchema([]pg.FieldSchema{
		{Name: "name", Type: "text"},
		{Name: "email", Type: "text"},
		{Name: "csv_data", Type: "text"},
		{Name: "tags", Type: "text", Repeated: true},
	})
	schemas := map[string]pg.Schema{"person": schema}

	// Create CEL environment with string extensions
	env, err := cel.NewEnv(
		cel.CustomTypeProvider(pg.NewTypeProvider(schemas)),
		cel.Variable("person", cel.ObjectType("person")),
		ext.Strings(), // Enable string extension functions
	)
	if err != nil {
		log.Fatalf("Failed to create CEL environment: %v", err)
	}

	// Example 1: split() - Basic splitting
	fmt.Println("1. split() - Basic String Splitting")
	fmt.Println("   CEL: person.csv_data.split(',').size() > 0")
	runExample(env, schemas, `person.csv_data.split(',').size() > 0`)

	// Example 2: split() with limit
	fmt.Println("\n2. split() with Limit")
	fmt.Println("   CEL: person.csv_data.split(',', 3).size() == 3")
	runExample(env, schemas, `person.csv_data.split(',', 3).size() == 3`)

	// Example 3: split() with limit = 0 (empty array)
	fmt.Println("\n3. split() with limit=0 (empty array)")
	fmt.Println("   CEL: person.csv_data.split(',', 0).size() == 0")
	runExample(env, schemas, `person.csv_data.split(',', 0).size() == 0`)

	// Example 4: split() with limit = 1 (no split)
	fmt.Println("\n4. split() with limit=1 (no split)")
	fmt.Println("   CEL: person.csv_data.split(',', 1) == [person.csv_data]")
	runExample(env, schemas, `person.csv_data.split(',', 1).size() == 1`)

	// Example 5: split() in comprehension
	fmt.Println("\n5. split() in Comprehension")
	fmt.Println("   CEL: person.csv_data.split(',').exists(x, x == 'target')")
	runExample(env, schemas, `person.csv_data.split(',').exists(x, x == 'target')`)

	// Example 6: join() - Basic joining
	fmt.Println("\n6. join() - Basic Array Joining")
	fmt.Println("   CEL: person.tags.join(',') == 'tag1,tag2,tag3'")
	runExample(env, schemas, `person.tags.join(',') == 'tag1,tag2,tag3'`)

	// Example 7: join() without delimiter
	fmt.Println("\n7. join() without Delimiter (empty string)")
	fmt.Println("   CEL: person.tags.join() == 'tag1tag2tag3'")
	runExample(env, schemas, `person.tags.join().contains('tag1')`)

	// Example 8: join() with comprehension
	fmt.Println("\n8. join() with Filtered Array")
	fmt.Println("   CEL: person.tags.filter(t, t.startsWith('a')).join(',')")
	runExample(env, schemas, `person.tags.filter(t, t.startsWith('a')).join(',').contains('a')`)

	// Example 9: format() - Basic formatting
	fmt.Println("\n9. format() - Basic String Formatting")
	fmt.Printf("   CEL: 'Name: %%s, Email: %%s'.format([person.name, person.email])\n")
	runExample(env, schemas, `'Name: %s, Email: %s'.format([person.name, person.email]) != ''`)

	// Example 10: format() with %d and %f
	fmt.Println("\n10. format() with Different Specifiers")
	fmt.Printf("    CEL: 'User %%s is %%d years old with score %%f'.format(['John', 30, 95.5])\n")
	runExample(env, schemas, `'User %s is %d years old with score %f'.format(['John', 30, 95.5]) != ''`)

	// Example 11: format() with escaped %%
	fmt.Println("\n11. format() with Escaped Percent Sign")
	fmt.Println("    CEL: 'Progress: 100%%'.format([])")
	runExample(env, schemas, `'Progress: 100%%'.format([]) != ''`)

	// Example 12: Combining split() and join()
	fmt.Println("\n12. Combining split() and join()")
	fmt.Println("    CEL: person.csv_data.split(',').filter(x, x.size() > 0).join('|')")
	runExample(env, schemas, `person.csv_data.split(',').filter(x, x.size() > 0).join('|').contains('|')`)

	// Example 13: Complex expression with all functions
	fmt.Println("\n13. Complex Expression")
	fmt.Println("    CEL: person.tags.join(',').split(',').size() > 0")
	runExample(env, schemas, `person.tags.join(',').split(',').size() > 0`)

	fmt.Println("\n==========================================")
	fmt.Println("All examples completed successfully!")
}

func runExample(env *cel.Env, schemas map[string]pg.Schema, expression string) {
	// Compile the CEL expression
	ast, issues := env.Compile(expression)
	if issues != nil && issues.Err() != nil {
		log.Printf("   ❌ Compilation error: %v\n", issues.Err())
		return
	}

	// Convert to SQL
	sql, err := cel2sql.Convert(ast, cel2sql.WithSchemas(schemas))
	if err != nil {
		log.Printf("   ❌ Conversion error: %v\n", err)
		return
	}

	// Print the SQL
	fmt.Printf("   SQL: %s\n", sql)
	fmt.Printf("   ✓ Success\n")
}
