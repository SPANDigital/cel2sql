package cel2sql_test

import (
	"database/sql"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	"github.com/spandigital/cel2sql/v3"
	"github.com/spandigital/cel2sql/v3/dialect"
	_ "github.com/spandigital/cel2sql/v3/dialect/sqlite"
	"github.com/spandigital/cel2sql/v3/schema"
)

// TestSQLiteComprehensionRunsAgainstSQLite converts a comprehension and then
// executes it, which is the check this needs: json_each is a table-valued
// function, so a reference to the iteration variable that names the alias
// rather than its value is SQL the converter reports as a success and SQLite
// rejects with "no such column".
func TestSQLiteComprehensionRunsAgainstSQLite(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	_, err = db.Exec(`CREATE TABLE reviews (id INTEGER PRIMARY KEY, approvals TEXT)`)
	require.NoError(t, err)
	_, err = db.Exec(`INSERT INTO reviews (id, approvals) VALUES (1, '["alice","bob"]'), (2, '["carol"]')`)
	require.NoError(t, err)

	env, err := cel.NewEnv(cel.Variable("approvals", cel.DynType))
	require.NoError(t, err)

	fields := []schema.FieldSchema{{
		Name: "approvals", Type: "text",
		IsJSON: true, Repeated: true, Dimensions: 1, ElementType: "text",
	}}
	schemas := map[string]schema.Schema{"reviews": schema.NewSchema(fields)}

	sqlite, err := dialect.Get(dialect.SQLite)
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
				cel2sql.WithDialect(sqlite), cel2sql.WithSchemas(schemas))
			require.NoError(t, err)

			rows, err := db.Query(
				`SELECT id FROM reviews WHERE `+converted.SQL, converted.Parameters...)
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
