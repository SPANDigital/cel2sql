// Copyright 2026 Observe, Inc. Licensed under Apache 2.0.
// Integration tests for list- and map-valued JSON variables, added by Observe, Inc.
package cel2sql_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/observeinc/cel2sql/v3"
)

// jsonContainerRows are object-tag documents chosen to tell the forms apart: an
// absent key, an empty list, a list beside other keys, keys that are not
// identifiers, and values holding a quote or a backslash.
var jsonContainerRows = []struct {
	id  int
	doc string
}{
	{1, `{"env":["prod","dev"],"team":["sre"]}`},
	{2, `{"env":["staging"]}`},
	{3, `{}`},
	{4, `{"team":["platform","sre"],"k8s.cluster":["us-west-2"]}`},
	{5, `{"env":[]}`},
	{6, `{"it's":["o'brien"],"back\\slash":["a\\b"]}`},
	{7, `{"env":["prod"],"owner":["dev team"]}`},
	{8, `{"env":["prod","prod"],"team":[]}`},
}

// jsonContainerMistypedRow breaks the declared type, holding a string where a
// list belongs. cel-go does not hold an activation to its declared type, so it
// would read the value as a string -- size("prod") is 4 -- and there is no CEL
// answer to compare with. It is seeded beside the others so every query runs
// over it, which shows the guard keeps it away from jsonb_array_elements_text
// and jsonb_array_length, both of which raise an error on a string.
const jsonContainerMistypedRow = `{"env":"prod"}`

const jsonContainerMistypedID = 99

// jsonContainerExprs are the filters a caller writes over object tags. The test
// does not state their answers: cel-go evaluates each over every row, and the
// SQL has to select exactly the rows it says are true.
var jsonContainerExprs = []string{
	exprEnvKeyPresent,
	`!("env" in tags)`,
	`has(tags.team)`,
	`!has(tags.team)`,
	`"k8s.cluster" in tags`,
	exprProdInEnv,
	`"prod" in tags["env"]`,
	`!("prod" in tags.env)`,
	`"prod" in tags.env || "sre" in tags.team`,
	`"prod" in tags.env && "sre" in tags.team`,
	`"prod" in tags.env || !("env" in tags)`,
	`has(tags.env) && !("prod" in tags.env)`,
	`"o'brien" in tags["it's"]`,
	`"a\\b" in tags["back\\slash"]`,
	`"dev team" in tags.owner`,
	`tags.env.exists(v, v.startsWith("pr"))`,
	`tags.env.exists(v, v.matches("^(prod|staging)$"))`,
	`!tags.env.exists(v, v == "dev")`,
	`tags.env.all(v, v == "prod")`,
	`!tags.env.all(v, v == "prod")`,
	`tags.team.exists_one(v, v == "sre")`,
	`tags.exists(k, k.startsWith("k8s."))`,
	`tags.all(k, k == "env")`,
	`size(tags) == 0`,
	`size(tags) >= 2`,
	`size(tags.env) == 0`,
	`size(tags.env) > 1`,
	`size(tags.env) != 1`,
}

// celMatches is the in-memory answer: the rows for which the expression is true.
// An evaluation error -- reading an absent key, or a string where a list belongs
// -- is not a match, which is what an error's SQL counterpart, NULL, gives too.
func celMatches(t *testing.T, env *cel.Env, ast *cel.Ast) []int {
	t.Helper()
	prg, err := env.Program(ast)
	require.NoError(t, err)
	var ids []int
	for _, row := range jsonContainerRows {
		var doc map[string]any
		require.NoError(t, json.Unmarshal([]byte(row.doc), &doc))
		out, _, err := prg.Eval(map[string]any{"tags": doc})
		if err == nil && out == types.True {
			ids = append(ids, row.id)
		}
	}
	return ids
}

func sqlMatches(t *testing.T, db *sql.DB, condition string, args []any) []int {
	t.Helper()
	// #nosec G202 - test asserting on generated SQL, not user input
	rows, err := db.Query("SELECT id FROM obj WHERE "+condition+" ORDER BY id", args...)
	require.NoError(t, err, "SQL must execute: %s", condition)
	defer func() { _ = rows.Close() }()
	var ids []int
	for rows.Next() {
		var id int
		require.NoError(t, rows.Scan(&id))
		ids = append(ids, id)
	}
	require.NoError(t, rows.Err())
	return ids
}

// TestJSONVariableContainers_PostgreSQL_MatchesCEL converts each filter both ways
// a caller can and checks that PostgreSQL selects the rows cel-go says are true.
func TestJSONVariableContainers_PostgreSQL_MatchesCEL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:17-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_PASSWORD": "password",
				"POSTGRES_DB":       "testdb",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer func() {
		if termErr := container.Terminate(ctx); termErr != nil {
			t.Logf("failed to terminate container: %v", termErr)
		}
	}()

	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	db, err := sql.Open("postgres",
		"host="+host+" port="+port.Port()+" user=postgres password=password dbname=testdb sslmode=disable")
	require.NoError(t, err)
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Logf("failed to close db: %v", closeErr)
		}
	}()

	_, err = db.Exec(`CREATE TABLE obj (id INT, tags JSONB)`)
	require.NoError(t, err)
	for _, row := range jsonContainerRows {
		_, err = db.Exec(`INSERT INTO obj VALUES ($1, $2::jsonb)`, row.id, row.doc)
		require.NoError(t, err)
	}
	_, err = db.Exec(`INSERT INTO obj VALUES ($1, $2::jsonb)`, jsonContainerMistypedID, jsonContainerMistypedRow)
	require.NoError(t, err)

	env, err := cel.NewEnv(cel.Variable("tags", cel.MapType(cel.StringType, cel.ListType(cel.StringType))))
	require.NoError(t, err)

	for _, expr := range jsonContainerExprs {
		t.Run(expr, func(t *testing.T) {
			ast, issues := env.Compile(expr)
			require.NoError(t, issues.Err())
			want := celMatches(t, env, ast)

			inlined, err := cel2sql.Convert(ast, cel2sql.WithJSONVariables("tags"))
			require.NoError(t, err)
			bound, err := cel2sql.ConvertParameterized(ast, cel2sql.WithJSONVariables("tags"))
			require.NoError(t, err)

			assert.Equal(t, want, withoutMistyped(sqlMatches(t, db, inlined, nil)), "Convert\nSQL: %s", inlined)
			assert.Equal(t, want, withoutMistyped(sqlMatches(t, db, bound.SQL, bound.Parameters)), "ConvertParameterized\nSQL: %s", bound.SQL)
		})
	}

	// A value that is not a list reads as NULL wherever a list is expected, the
	// counterpart of the type error CEL would raise for it, so it matches none of
	// these -- while the key itself is still present.
	for _, expr := range []string{
		exprProdInEnv,
		`!("prod" in tags.env)`,
		`tags.env.exists(v, v == "prod")`,
		`tags.env.all(v, v == "prod")`,
		`size(tags.env) > 0`,
		`size(tags.env) == 0`,
	} {
		t.Run("mistyped value/"+expr, func(t *testing.T) {
			ast, issues := env.Compile(expr)
			require.NoError(t, issues.Err())
			bound, err := cel2sql.ConvertParameterized(ast, cel2sql.WithJSONVariables("tags"))
			require.NoError(t, err)
			assert.NotContains(t, sqlMatches(t, db, bound.SQL, bound.Parameters), jsonContainerMistypedID)
		})
	}
	t.Run("mistyped value/key is present", func(t *testing.T) {
		ast, issues := env.Compile(exprEnvKeyPresent)
		require.NoError(t, issues.Err())
		bound, err := cel2sql.ConvertParameterized(ast, cel2sql.WithJSONVariables("tags"))
		require.NoError(t, err)
		assert.Contains(t, sqlMatches(t, db, bound.SQL, bound.Parameters), jsonContainerMistypedID)
	})
}

// withoutMistyped drops jsonContainerMistypedID, which has no CEL answer to compare with.
func withoutMistyped(ids []int) []int {
	var out []int
	for _, id := range ids {
		if id != jsonContainerMistypedID {
			out = append(out, id)
		}
	}
	return out
}
