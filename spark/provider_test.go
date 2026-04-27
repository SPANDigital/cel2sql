package spark_test

import (
	"context"
	"testing"

	"github.com/google/cel-go/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3/spark"
)

func TestNewTypeProvider(t *testing.T) {
	schemas := map[string]spark.Schema{
		"users": spark.NewSchema([]spark.FieldSchema{
			{Name: "id", Type: "bigint"},
			{Name: "name", Type: "string"},
		}),
	}
	provider := spark.NewTypeProvider(schemas)
	require.NotNil(t, provider)

	got := provider.GetSchemas()
	assert.Len(t, got, 1)
	assert.Contains(t, got, "users")
}

func TestNewTypeProviderWithConnection_NilDB(t *testing.T) {
	_, err := spark.NewTypeProviderWithConnection(context.Background(), nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, spark.ErrInvalidSchema)
}

func TestLoadTableSchema_NoDB(t *testing.T) {
	provider := spark.NewTypeProvider(nil)
	err := provider.LoadTableSchema(context.Background(), "test")
	require.Error(t, err)
	assert.ErrorIs(t, err, spark.ErrInvalidSchema)
}

func TestTypeProvider_FindStructType(t *testing.T) {
	schemas := map[string]spark.Schema{
		"users": spark.NewSchema([]spark.FieldSchema{
			{Name: "id", Type: "bigint"},
			{Name: "name", Type: "string"},
		}),
	}
	provider := spark.NewTypeProvider(schemas)

	typ, found := provider.FindStructType("users")
	assert.True(t, found)
	assert.NotNil(t, typ)

	_, found = provider.FindStructType("nonexistent")
	assert.False(t, found)
}

func TestTypeProvider_FindStructFieldNames(t *testing.T) {
	schemas := map[string]spark.Schema{
		"users": spark.NewSchema([]spark.FieldSchema{
			{Name: "id", Type: "bigint"},
			{Name: "name", Type: "string"},
			{Name: "email", Type: "string"},
		}),
	}
	provider := spark.NewTypeProvider(schemas)

	names, found := provider.FindStructFieldNames("users")
	assert.True(t, found)
	assert.ElementsMatch(t, []string{"id", "name", "email"}, names)

	_, found = provider.FindStructFieldNames("nonexistent")
	assert.False(t, found)
}

func TestTypeProvider_FindStructFieldType(t *testing.T) {
	schemas := map[string]spark.Schema{
		"test_table": spark.NewSchema([]spark.FieldSchema{
			{Name: "str_field", Type: "string"},
			{Name: "varchar_field", Type: "varchar"},
			{Name: "int_field", Type: "int"},
			{Name: "bigint_field", Type: "bigint"},
			{Name: "tinyint_field", Type: "tinyint"},
			{Name: "double_field", Type: "double"},
			{Name: "float_field", Type: "float"},
			{Name: "decimal_field", Type: "decimal"},
			{Name: "bool_field", Type: "boolean"},
			{Name: "binary_field", Type: "binary"},
			{Name: "ts_field", Type: "timestamp"},
			{Name: "ts_ntz_field", Type: "timestamp_ntz"},
			{Name: "tags", Type: "string", Repeated: true, ElementType: "string"},
			{Name: "scores", Type: "bigint", Repeated: true, ElementType: "bigint"},
		}),
	}
	provider := spark.NewTypeProvider(schemas)

	tests := []struct {
		fieldName string
		wantType  *types.Type
		wantFound bool
	}{
		{"str_field", types.StringType, true},
		{"varchar_field", types.StringType, true},
		{"int_field", types.IntType, true},
		{"bigint_field", types.IntType, true},
		{"tinyint_field", types.IntType, true},
		{"double_field", types.DoubleType, true},
		{"float_field", types.DoubleType, true},
		{"decimal_field", types.DoubleType, true},
		{"bool_field", types.BoolType, true},
		{"binary_field", types.BytesType, true},
		{"ts_field", types.TimestampType, true},
		{"ts_ntz_field", types.TimestampType, true},
		{"tags", types.NewListType(types.StringType), true},
		{"scores", types.NewListType(types.IntType), true},
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
	provider := spark.NewTypeProvider(nil)
	// Close should not panic.
	provider.Close()
}

func TestLoadTableSchema_RejectsInvalidNames(t *testing.T) {
	// We don't need a real DB: invalid names short-circuit before any query.
	provider := spark.NewTypeProvider(nil)
	for _, bad := range []string{
		"users; DROP TABLE users",
		"users--",
		"' OR '1'='1",
		"",
	} {
		err := provider.LoadTableSchema(context.Background(), bad)
		require.Error(t, err, "expected error for input %q", bad)
		assert.ErrorIs(t, err, spark.ErrInvalidSchema)
	}
}
