package cel2sql_test

import (
	"context"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/dialect"
	_ "github.com/spandigital/cel2sql/v3/dialect/mysql"
	"github.com/spandigital/cel2sql/v3/schema"
)

// TestMySQLComprehensionRunsAgainstMySQL converts a comprehension and then
// executes it, mirroring TestSQLiteComprehensionRunsAgainstSQLite: JSON_TABLE
// is a table-valued function, so a reference to the iteration variable that
// names the alias rather than its value column is SQL the converter reports
// as a success and MySQL rejects at query time.
//
// Requires MySQL 8.4+: on 8.0 (EOL April 2026) a correlated JSON_TABLE inside
// EXISTS silently matches nothing when run as a prepared statement, so these
// queries return no error and no rows there.
func TestMySQLComprehensionRunsAgainstMySQL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	container, db := setupMySQLContainer(ctx, t)
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Logf("failed to close db: %v", closeErr)
		}
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	_, err := db.Exec(`CREATE TABLE reviews (id INTEGER PRIMARY KEY, approvals JSON)`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO reviews (id, approvals) VALUES (1, '["alice","bob"]'), (2, '["carol"]')`)
	require.NoError(t, err)

	env, err := cel.NewEnv(cel.Variable("approvals", cel.DynType))
	require.NoError(t, err)

	fields := []schema.FieldSchema{{
		Name: "approvals", Type: "json",
		IsJSON: true, Repeated: true, Dimensions: 1, ElementType: "text",
	}}
	schemas := map[string]schema.Schema{"reviews": schema.NewSchema(fields)}

	mysql, err := dialect.Get(dialect.MySQL)
	require.NoError(t, err)

	tests := []struct {
		name    string
		celExpr string
		wantIDs []int
	}{
		{name: "exists", celExpr: `approvals.exists(a, a == "alice")`, wantIDs: []int{1}},
		{name: "exists none", celExpr: `approvals.exists(a, a == "nobody")`},
		{name: "all", celExpr: `approvals.all(a, a != "carol")`, wantIDs: []int{1}},
		{name: "exists_one", celExpr: `approvals.exists_one(a, a == "carol")`, wantIDs: []int{2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, iss := env.Compile(tt.celExpr)
			require.NoError(t, iss.Err())

			converted, err := cel2sql.ConvertParameterized(ast,
				cel2sql.WithDialect(mysql), cel2sql.WithSchemas(schemas))
			require.NoError(t, err)

			rows, err := db.Query(
				`SELECT id FROM reviews WHERE `+converted.SQL+` ORDER BY id`, converted.Parameters...) //nolint:gosec // converted.SQL is parameterized output under test, not user input
			require.NoError(t, err, "generated SQL: %s", converted.SQL)
			defer func() { _ = rows.Close() }()

			var ids []int
			for rows.Next() {
				var id int
				require.NoError(t, rows.Scan(&id))
				ids = append(ids, id)
			}
			require.NoError(t, rows.Err())
			require.Equal(t, tt.wantIDs, ids, "generated SQL: %s", converted.SQL)
		})
	}
}
