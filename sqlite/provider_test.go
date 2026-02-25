package sqlite_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/cel-go/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	"github.com/spandigital/cel2sql/v3/sqlite"
)

func TestNewTypeProvider(t *testing.T) {
	schemas := map[string]sqlite.Schema{
		"users": sqlite.NewSchema([]sqlite.FieldSchema{
			{Name: "id", Type: "integer"},
			{Name: "name", Type: "text"},
		}),
	}

	provider := sqlite.NewTypeProvider(schemas)
	require.NotNil(t, provider)

	got := provider.GetSchemas()
	assert.Len(t, got, 1)
	assert.Contains(t, got, "users")
}

func TestNewTypeProviderWithConnection_NilDB(t *testing.T) {
	_, err := sqlite.NewTypeProviderWithConnection(context.Background(), nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, sqlite.ErrInvalidSchema)
}

func TestLoadTableSchema_NoDB(t *testing.T) {
	provider := sqlite.NewTypeProvider(nil)
	err := provider.LoadTableSchema(context.Background(), "test")
	require.Error(t, err)
	assert.ErrorIs(t, err, sqlite.ErrInvalidSchema)
}

func TestLoadTableSchema_InvalidTableName(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	ctx := context.Background()
	provider, err := sqlite.NewTypeProviderWithConnection(ctx, db)
	require.NoError(t, err)

	// SQL injection attempts should be rejected
	invalidNames := []string{
		"table; DROP TABLE users",
		"table'name",
		"table-name",
		"table.name",
		"123table",
		"",
		"table name",
	}

	for _, name := range invalidNames {
		t.Run(name, func(t *testing.T) {
			err := provider.LoadTableSchema(ctx, name)
			require.Error(t, err)
			assert.ErrorIs(t, err, sqlite.ErrInvalidSchema)
		})
	}
}

func TestTypeProvider_FindStructType(t *testing.T) {
	schemas := map[string]sqlite.Schema{
		"users": sqlite.NewSchema([]sqlite.FieldSchema{
			{Name: "id", Type: "integer"},
			{Name: "name", Type: "text"},
		}),
	}
	provider := sqlite.NewTypeProvider(schemas)

	typ, found := provider.FindStructType("users")
	assert.True(t, found)
	assert.NotNil(t, typ)

	_, found = provider.FindStructType("nonexistent")
	assert.False(t, found)
}

func TestTypeProvider_FindStructFieldNames(t *testing.T) {
	schemas := map[string]sqlite.Schema{
		"users": sqlite.NewSchema([]sqlite.FieldSchema{
			{Name: "id", Type: "integer"},
			{Name: "name", Type: "text"},
			{Name: "email", Type: "text"},
		}),
	}
	provider := sqlite.NewTypeProvider(schemas)

	names, found := provider.FindStructFieldNames("users")
	assert.True(t, found)
	assert.ElementsMatch(t, []string{"id", "name", "email"}, names)

	_, found = provider.FindStructFieldNames("nonexistent")
	assert.False(t, found)
}

func TestTypeProvider_FindStructFieldType(t *testing.T) {
	schemas := map[string]sqlite.Schema{
		"test_table": sqlite.NewSchema([]sqlite.FieldSchema{
			{Name: "text_field", Type: "text"},
			{Name: "int_field", Type: "integer"},
			{Name: "real_field", Type: "real"},
			{Name: "bool_field", Type: "boolean"},
			{Name: "blob_field", Type: "blob"},
			{Name: "json_field", Type: "json"},
			{Name: "datetime_field", Type: "datetime"},
		}),
	}
	provider := sqlite.NewTypeProvider(schemas)

	tests := []struct {
		fieldName string
		wantType  *types.Type
		wantFound bool
	}{
		{"text_field", types.StringType, true},
		{"int_field", types.IntType, true},
		{"real_field", types.DoubleType, true},
		{"bool_field", types.BoolType, true},
		{"blob_field", types.BytesType, true},
		{"json_field", types.DynType, true},
		{"datetime_field", types.TimestampType, true},
		{"nonexistent", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.fieldName, func(t *testing.T) {
			got, found := provider.FindStructFieldType("test_table", tt.fieldName)
			assert.Equal(t, tt.wantFound, found)
			if found {
				assert.Equal(t, tt.wantType, got.Type)
			}
		})
	}
}

func TestTypeProvider_Close(_ *testing.T) {
	provider := sqlite.NewTypeProvider(nil)
	// Close should not panic
	provider.Close()
}

func TestLoadTableSchema_Integration(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	ctx := context.Background()

	// Create test table with various SQLite types
	_, err = db.ExecContext(ctx, `
		CREATE TABLE schema_test (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			description VARCHAR(255),
			age INTEGER,
			score REAL,
			is_active BOOLEAN,
			avatar BLOB,
			metadata JSON,
			created_at DATETIME
		)
	`)
	require.NoError(t, err)

	provider, err := sqlite.NewTypeProviderWithConnection(ctx, db)
	require.NoError(t, err)

	err = provider.LoadTableSchema(ctx, "schema_test")
	require.NoError(t, err)

	// Verify schema was loaded
	schemas := provider.GetSchemas()
	assert.Contains(t, schemas, "schema_test")

	// Verify FindStructType
	typ, found := provider.FindStructType("schema_test")
	assert.True(t, found)
	assert.NotNil(t, typ)

	// Verify FindStructFieldNames
	names, found := provider.FindStructFieldNames("schema_test")
	assert.True(t, found)
	assert.Contains(t, names, "id")
	assert.Contains(t, names, "name")
	assert.Contains(t, names, "metadata")

	// Verify type mappings from loaded schema
	tests := []struct {
		fieldName string
		wantType  *types.Type
	}{
		{"id", types.IntType},
		{"name", types.StringType},
		{"description", types.StringType},
		{"age", types.IntType},
		{"score", types.DoubleType},
		{"is_active", types.BoolType},
		{"avatar", types.BytesType},
		{"metadata", types.DynType},
		{"created_at", types.TimestampType},
	}

	for _, tt := range tests {
		t.Run("type_"+tt.fieldName, func(t *testing.T) {
			got, found := provider.FindStructFieldType("schema_test", tt.fieldName)
			assert.True(t, found, "field %q should be found", tt.fieldName)
			if found {
				assert.Equal(t, tt.wantType, got.Type, "field %q type mismatch", tt.fieldName)
			}
		})
	}

	// Verify JSON detection
	schemaObj := schemas["schema_test"]
	metadataField, found := schemaObj.FindField("metadata")
	assert.True(t, found)
	assert.True(t, metadataField.IsJSON, "metadata should be detected as JSON")
}

func TestLoadTableSchema_NonexistentTable(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	ctx := context.Background()
	provider, err := sqlite.NewTypeProviderWithConnection(ctx, db)
	require.NoError(t, err)

	err = provider.LoadTableSchema(ctx, "nonexistent_table")
	require.Error(t, err)
	assert.ErrorIs(t, err, sqlite.ErrInvalidSchema)
}

func TestLoadTableSchema_MultipleTablesIntegration(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	ctx := context.Background()

	// Create two tables
	_, err = db.ExecContext(ctx, `CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, email TEXT)`)
	require.NoError(t, err)
	_, err = db.ExecContext(ctx, `CREATE TABLE products (id INTEGER PRIMARY KEY, title TEXT, price REAL)`)
	require.NoError(t, err)

	provider, err := sqlite.NewTypeProviderWithConnection(ctx, db)
	require.NoError(t, err)

	err = provider.LoadTableSchema(ctx, "users")
	require.NoError(t, err)
	err = provider.LoadTableSchema(ctx, "products")
	require.NoError(t, err)

	schemas := provider.GetSchemas()
	assert.Len(t, schemas, 2)
	assert.Contains(t, schemas, "users")
	assert.Contains(t, schemas, "products")

	// Verify both schemas are independent
	userNames, found := provider.FindStructFieldNames("users")
	assert.True(t, found)
	assert.ElementsMatch(t, []string{"id", "name", "email"}, userNames)

	productNames, found := provider.FindStructFieldNames("products")
	assert.True(t, found)
	assert.ElementsMatch(t, []string{"id", "title", "price"}, productNames)
}

func TestNormalizeSQLiteType(t *testing.T) {
	// Test via LoadTableSchema that various SQLite type declarations are normalized correctly
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	ctx := context.Background()

	_, err = db.ExecContext(ctx, `
		CREATE TABLE type_test (
			col_int INTEGER,
			col_varchar VARCHAR(255),
			col_char CHAR(10),
			col_text TEXT,
			col_real REAL,
			col_float FLOAT,
			col_double DOUBLE,
			col_numeric NUMERIC,
			col_blob BLOB,
			col_bool BOOLEAN,
			col_datetime DATETIME,
			col_timestamp TIMESTAMP,
			col_bigint BIGINT,
			col_smallint SMALLINT,
			col_tinyint TINYINT,
			col_json JSON
		)
	`)
	require.NoError(t, err)

	provider, err := sqlite.NewTypeProviderWithConnection(ctx, db)
	require.NoError(t, err)

	err = provider.LoadTableSchema(ctx, "type_test")
	require.NoError(t, err)

	tests := []struct {
		fieldName string
		wantType  *types.Type
	}{
		{"col_int", types.IntType},
		{"col_varchar", types.StringType},
		{"col_char", types.StringType},
		{"col_text", types.StringType},
		{"col_real", types.DoubleType},
		{"col_float", types.DoubleType},
		{"col_double", types.DoubleType},
		{"col_numeric", types.DoubleType},
		{"col_blob", types.BytesType},
		{"col_bool", types.BoolType},
		{"col_datetime", types.TimestampType},
		{"col_timestamp", types.TimestampType},
		{"col_bigint", types.IntType},
		{"col_smallint", types.IntType},
		{"col_tinyint", types.IntType},
		{"col_json", types.DynType},
	}

	for _, tt := range tests {
		t.Run(tt.fieldName, func(t *testing.T) {
			got, found := provider.FindStructFieldType("type_test", tt.fieldName)
			require.True(t, found, "field %q should be found", tt.fieldName)
			assert.Equal(t, tt.wantType, got.Type, "field %q type mismatch", tt.fieldName)
		})
	}
}
