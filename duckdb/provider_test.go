package duckdb_test

import (
	"context"
	"testing"

	"github.com/google/cel-go/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3/duckdb"
)

func TestNewTypeProvider(t *testing.T) {
	schemas := map[string]duckdb.Schema{
		"users": duckdb.NewSchema([]duckdb.FieldSchema{
			{Name: "id", Type: "integer"},
			{Name: "name", Type: "varchar"},
		}),
	}

	provider := duckdb.NewTypeProvider(schemas)
	require.NotNil(t, provider)

	got := provider.GetSchemas()
	assert.Len(t, got, 1)
	assert.Contains(t, got, "users")
}

func TestNewTypeProviderWithConnection_NilDB(t *testing.T) {
	_, err := duckdb.NewTypeProviderWithConnection(context.Background(), nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, duckdb.ErrInvalidSchema)
}

func TestLoadTableSchema_NoDB(t *testing.T) {
	provider := duckdb.NewTypeProvider(nil)
	err := provider.LoadTableSchema(context.Background(), "test")
	require.Error(t, err)
	assert.ErrorIs(t, err, duckdb.ErrInvalidSchema)
}

func TestTypeProvider_FindStructType(t *testing.T) {
	schemas := map[string]duckdb.Schema{
		"users": duckdb.NewSchema([]duckdb.FieldSchema{
			{Name: "id", Type: "integer"},
			{Name: "name", Type: "varchar"},
		}),
	}
	provider := duckdb.NewTypeProvider(schemas)

	typ, found := provider.FindStructType("users")
	assert.True(t, found)
	assert.NotNil(t, typ)

	_, found = provider.FindStructType("nonexistent")
	assert.False(t, found)
}

func TestTypeProvider_FindStructFieldNames(t *testing.T) {
	schemas := map[string]duckdb.Schema{
		"users": duckdb.NewSchema([]duckdb.FieldSchema{
			{Name: "id", Type: "integer"},
			{Name: "name", Type: "varchar"},
			{Name: "email", Type: "text"},
		}),
	}
	provider := duckdb.NewTypeProvider(schemas)

	names, found := provider.FindStructFieldNames("users")
	assert.True(t, found)
	assert.ElementsMatch(t, []string{"id", "name", "email"}, names)

	_, found = provider.FindStructFieldNames("nonexistent")
	assert.False(t, found)
}

func TestTypeProvider_FindStructFieldType(t *testing.T) {
	schemas := map[string]duckdb.Schema{
		"test_table": duckdb.NewSchema([]duckdb.FieldSchema{
			{Name: "str_field", Type: "varchar"},
			{Name: "text_field", Type: "text"},
			{Name: "int_field", Type: "integer"},
			{Name: "bigint_field", Type: "bigint"},
			{Name: "smallint_field", Type: "smallint"},
			{Name: "tinyint_field", Type: "tinyint"},
			{Name: "hugeint_field", Type: "hugeint"},
			{Name: "double_field", Type: "double"},
			{Name: "float_field", Type: "float"},
			{Name: "bool_field", Type: "boolean"},
			{Name: "blob_field", Type: "blob"},
			{Name: "json_field", Type: "json"},
			{Name: "ts_field", Type: "timestamp"},
			{Name: "array_field", Type: "integer", Repeated: true},
		}),
	}
	provider := duckdb.NewTypeProvider(schemas)

	tests := []struct {
		fieldName string
		wantType  *types.Type
		wantFound bool
	}{
		{"str_field", types.StringType, true},
		{"text_field", types.StringType, true},
		{"int_field", types.IntType, true},
		{"bigint_field", types.IntType, true},
		{"smallint_field", types.IntType, true},
		{"tinyint_field", types.IntType, true},
		{"hugeint_field", types.IntType, true},
		{"double_field", types.DoubleType, true},
		{"float_field", types.DoubleType, true},
		{"bool_field", types.BoolType, true},
		{"blob_field", types.BytesType, true},
		{"json_field", types.DynType, true},
		{"ts_field", types.TimestampType, true},
		{"array_field", types.NewListType(types.IntType), true},
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
	provider := duckdb.NewTypeProvider(nil)
	// Close should not panic
	provider.Close()
}

func TestTypeProvider_ArrayDetection(t *testing.T) {
	// Test that arrays defined manually with Repeated=true work correctly
	schemas := map[string]duckdb.Schema{
		"test_table": duckdb.NewSchema([]duckdb.FieldSchema{
			{Name: "tags", Type: "varchar", Repeated: true, Dimensions: 1},
			{Name: "matrix", Type: "integer", Repeated: true, Dimensions: 2},
		}),
	}
	provider := duckdb.NewTypeProvider(schemas)

	// tags should be list of strings
	got, found := provider.FindStructFieldType("test_table", "tags")
	assert.True(t, found)
	assert.Equal(t, types.NewListType(types.StringType), got.Type)

	// matrix should be list of integers (CEL sees all array dims as list)
	got, found = provider.FindStructFieldType("test_table", "matrix")
	assert.True(t, found)
	assert.Equal(t, types.NewListType(types.IntType), got.Type)
}
