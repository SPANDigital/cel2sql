package opensearch_test

import (
	"context"
	_ "embed"
	"testing"

	"cel.dev/cel-go/cel"
	"cel.dev/cel-go/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3"
	opensearchDialect "github.com/spandigital/cel2sql/v3/dialect/opensearch"
	"github.com/spandigital/cel2sql/v3/opensearch"
)

//go:embed testdata/mapping.json
var productsMapping []byte

func productsProvider(t *testing.T) opensearch.TypeProvider {
	t.Helper()
	fields, err := opensearch.ParseMappingForTest(productsMapping)
	require.NoError(t, err)
	return opensearch.NewTypeProvider(map[string]opensearch.Schema{
		"products": opensearch.NewSchema(fields),
	})
}

func TestNewTypeProvider(t *testing.T) {
	provider := opensearch.NewTypeProvider(map[string]opensearch.Schema{
		"users": opensearch.NewSchema([]opensearch.FieldSchema{
			{Name: "id", Type: "long"},
			{Name: "name", Type: "keyword"},
		}),
	})
	require.NotNil(t, provider)
	assert.Contains(t, provider.GetSchemas(), "users")
	provider.Close()
}

func TestNewTypeProviderWithClient_NilClient(t *testing.T) {
	_, err := opensearch.NewTypeProviderWithClient(context.Background(), nil)
	require.ErrorIs(t, err, opensearch.ErrInvalidSchema)
}

func TestLoadTableSchema_NoClient(t *testing.T) {
	provider := opensearch.NewTypeProvider(nil)
	err := provider.LoadTableSchema(context.Background(), "products")
	require.ErrorIs(t, err, opensearch.ErrInvalidSchema)
}

func TestLoadTableSchema_RejectsInvalidNames(t *testing.T) {
	provider := opensearch.NewTypeProvider(nil)
	for _, name := range []string{"", "_hidden", "-leading", "Upper", "a/b", "idx; DROP", "logs-*", "a b", "a,b"} {
		t.Run(name, func(t *testing.T) {
			err := provider.LoadTableSchema(context.Background(), name)
			require.ErrorIs(t, err, opensearch.ErrInvalidSchema)
			assert.Contains(t, err.Error(), "invalid index name")
		})
	}
}

func TestParseMapping(t *testing.T) {
	fields, err := opensearch.ParseMappingForTest(productsMapping)
	require.NoError(t, err)

	byName := make(map[string]opensearch.FieldSchema, len(fields))
	names := make([]string, 0, len(fields))
	for _, f := range fields {
		byName[f.Name] = f
		names = append(names, f.Name)
	}
	assert.Equal(t, []string{"active", "created_at", "description", "dims", "extra", "name", "price", "stock", "variants"}, names)

	assert.Equal(t, "keyword", byName["name"].Type)
	assert.Equal(t, "text", byName["description"].Type)
	assert.Empty(t, byName["description"].Schema, "multi-fields are not schema fields")
	assert.Equal(t, "double", byName["price"].Type)
	assert.Equal(t, "long", byName["stock"].Type)
	assert.Equal(t, "boolean", byName["active"].Type)
	assert.Equal(t, "date", byName["created_at"].Type)

	dims := byName["dims"]
	assert.Equal(t, "object", dims.Type)
	assert.False(t, dims.Repeated)
	require.Len(t, dims.Schema, 2)
	assert.Equal(t, "height", dims.Schema[0].Name)
	assert.Equal(t, "width", dims.Schema[1].Name)

	variants := byName["variants"]
	assert.Equal(t, "nested", variants.Type)
	assert.True(t, variants.Repeated)
	assert.Equal(t, 1, variants.Dimensions)
	require.Len(t, variants.Schema, 3)
	assert.Equal(t, "attrs", variants.Schema[0].Name)
	assert.Equal(t, "object", variants.Schema[0].Type)
	require.Len(t, variants.Schema[0].Schema, 1)
	assert.Equal(t, "color", variants.Schema[0].Schema[0].Name)

	assert.True(t, byName["extra"].IsJSON)
}

func TestParseMapping_Errors(t *testing.T) {
	_, err := opensearch.ParseMappingForTest(nil)
	require.Error(t, err)

	_, err = opensearch.ParseMappingForTest([]byte(`{"properties": `))
	require.Error(t, err)

	fields, err := opensearch.ParseMappingForTest([]byte(`{}`))
	require.NoError(t, err)
	assert.Empty(t, fields)
}

func TestTypeProvider_FindStructType(t *testing.T) {
	provider := productsProvider(t)

	for _, name := range []string{"products", "products.dims", "products.variants", "products.variants.attrs"} {
		typ, found := provider.FindStructType(name)
		assert.True(t, found, name)
		assert.Equal(t, name, typ.TypeName())
	}
	for _, name := range []string{"missing", "products.name", "products.dims.width", "products.nope"} {
		_, found := provider.FindStructType(name)
		assert.False(t, found, name)
	}
}

func TestTypeProvider_FindStructFieldNames(t *testing.T) {
	provider := productsProvider(t)

	names, found := provider.FindStructFieldNames("products.dims")
	require.True(t, found)
	assert.Equal(t, []string{"height", "width"}, names)

	_, found = provider.FindStructFieldNames("missing")
	assert.False(t, found)
}

func TestTypeProvider_FindStructFieldType(t *testing.T) {
	provider := productsProvider(t)

	tests := []struct {
		structType string
		field      string
		want       *types.Type
	}{
		{"products", "name", types.StringType},
		{"products", "description", types.StringType},
		{"products", "price", types.DoubleType},
		{"products", "stock", types.IntType},
		{"products", "active", types.BoolType},
		{"products", "created_at", types.TimestampType},
		{"products", "extra", types.DynType},
		{"products", "dims", types.NewObjectType("products.dims")},
		{"products", "variants", types.NewListType(types.NewObjectType("products.variants"))},
		{"products.dims", "width", types.DoubleType},
		{"products.variants", "sku", types.StringType},
		{"products.variants", "attrs", types.NewObjectType("products.variants.attrs")},
		{"products.variants.attrs", "color", types.StringType},
	}
	for _, tc := range tests {
		t.Run(tc.structType+"."+tc.field, func(t *testing.T) {
			ft, found := provider.FindStructFieldType(tc.structType, tc.field)
			require.True(t, found)
			assert.True(t, tc.want.IsExactType(ft.Type), "got %s want %s", ft.Type, tc.want)
		})
	}

	_, found := provider.FindStructFieldType("products", "nope")
	assert.False(t, found)
	_, found = provider.FindStructFieldType("missing", "name")
	assert.False(t, found)
}

func TestTypeProvider_OtherTypes(t *testing.T) {
	provider := opensearch.NewTypeProvider(map[string]opensearch.Schema{
		"t": opensearch.NewSchema([]opensearch.FieldSchema{
			{Name: "ip", Type: "ip"},
			{Name: "b", Type: "binary"},
			{Name: "f", Type: "half_float"},
			{Name: "ts", Type: "date_nanos"},
			{Name: "geo", Type: "geo_point"},
			{Name: "tags", Type: "keyword", Repeated: true},
		}),
	})
	want := map[string]*types.Type{
		"ip": types.StringType, "b": types.BytesType, "f": types.DoubleType,
		"ts": types.TimestampType, "geo": types.DynType, "tags": types.NewListType(types.StringType),
	}
	for field, w := range want {
		ft, found := provider.FindStructFieldType("t", field)
		require.True(t, found, field)
		assert.True(t, w.IsExactType(ft.Type), "%s: got %s want %s", field, ft.Type, w)
	}
}

func TestEndToEndConversion(t *testing.T) {
	provider := productsProvider(t)
	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("p", cel.ObjectType("products")),
	)
	require.NoError(t, err)

	tests := []struct {
		cel  string
		want string
	}{
		{`p.dims.width > 10.0 && p.variants.exists(v, v.sku == "A")`,
			`p.dims.width > 10 AND EXISTS (SELECT * FROM p.variants AS v WHERE v.sku = 'A')`},
		{`p.variants.exists(v, v.qty > 0 && v.attrs.color == "red")`,
			`EXISTS (SELECT * FROM p.variants AS v WHERE v.qty > 0 AND v.attrs.color = 'red')`},
		{`!p.active || p.stock == 0`, `NOT p.active = true OR p.stock = 0`},
		{`p.name in ["a", "b"] && p.created_at < timestamp("2024-06-01T00:00:00Z")`,
			`p.name IN ('a', 'b') AND p.created_at < CAST('2024-06-01T00:00:00Z' AS TIMESTAMP)`},
		{`has(p.dims.height)`, `p.dims.height IS NOT NULL`},
	}
	for _, tc := range tests {
		t.Run(tc.cel, func(t *testing.T) {
			ast, issues := env.Compile(tc.cel)
			require.NoError(t, issues.Err())
			got, err := cel2sql.Convert(ast, cel2sql.WithDialect(opensearchDialect.New()))
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}

	t.Run("all is rejected", func(t *testing.T) {
		ast, issues := env.Compile(`p.variants.all(v, v.qty > 0)`)
		require.NoError(t, issues.Err())
		_, err := cel2sql.Convert(ast, cel2sql.WithDialect(opensearchDialect.New()))
		require.Error(t, err)
	})
}
