package bigquery_test

import (
	"bytes"
	"context"
	_ "embed"
	"runtime"
	"strings"
	"testing"

	bq "cloud.google.com/go/bigquery"
	"github.com/google/cel-go/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcbigquery "github.com/testcontainers/testcontainers-go/modules/gcloud/bigquery"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/spandigital/cel2sql/v3/bigquery"
)

//go:embed testdata/provider_seed.yaml
var providerSeedYAML []byte

const (
	testProjectID = "test-project"
	testDataset   = "testdataset"
)

func TestNewTypeProvider(t *testing.T) {
	schemas := map[string]bigquery.Schema{
		"users": bigquery.NewSchema([]bigquery.FieldSchema{
			{Name: "id", Type: "INTEGER"},
			{Name: "name", Type: "STRING"},
		}),
	}

	provider := bigquery.NewTypeProvider(schemas)
	require.NotNil(t, provider)

	got := provider.GetSchemas()
	assert.Len(t, got, 1)
	assert.Contains(t, got, "users")
}

func TestNewTypeProviderWithClient_NilClient(t *testing.T) {
	_, err := bigquery.NewTypeProviderWithClient(context.Background(), nil, "dataset")
	require.Error(t, err)
	assert.ErrorIs(t, err, bigquery.ErrInvalidSchema)
}

func TestNewTypeProviderWithClient_EmptyDataset(t *testing.T) {
	// We can't create a real client without credentials, so test the nil case
	_, err := bigquery.NewTypeProviderWithClient(context.Background(), nil, "")
	require.Error(t, err)
	assert.ErrorIs(t, err, bigquery.ErrInvalidSchema)
}

func TestLoadTableSchema_NoClient(t *testing.T) {
	provider := bigquery.NewTypeProvider(nil)
	err := provider.LoadTableSchema(context.Background(), "test")
	require.Error(t, err)
	assert.ErrorIs(t, err, bigquery.ErrInvalidSchema)
}

func TestTypeProvider_FindStructType(t *testing.T) {
	schemas := map[string]bigquery.Schema{
		"users": bigquery.NewSchema([]bigquery.FieldSchema{
			{Name: "id", Type: "INTEGER"},
			{Name: "name", Type: "STRING"},
		}),
	}
	provider := bigquery.NewTypeProvider(schemas)

	typ, found := provider.FindStructType("users")
	assert.True(t, found)
	assert.NotNil(t, typ)

	_, found = provider.FindStructType("nonexistent")
	assert.False(t, found)
}

func TestTypeProvider_FindStructFieldNames(t *testing.T) {
	schemas := map[string]bigquery.Schema{
		"users": bigquery.NewSchema([]bigquery.FieldSchema{
			{Name: "id", Type: "INTEGER"},
			{Name: "name", Type: "STRING"},
			{Name: "email", Type: "STRING"},
		}),
	}
	provider := bigquery.NewTypeProvider(schemas)

	names, found := provider.FindStructFieldNames("users")
	assert.True(t, found)
	assert.ElementsMatch(t, []string{"id", "name", "email"}, names)

	_, found = provider.FindStructFieldNames("nonexistent")
	assert.False(t, found)
}

func TestTypeProvider_FindStructFieldType(t *testing.T) {
	schemas := map[string]bigquery.Schema{
		"test_table": bigquery.NewSchema([]bigquery.FieldSchema{
			{Name: "str_field", Type: "STRING"},
			{Name: "int_field", Type: "INTEGER"},
			{Name: "int64_field", Type: "INT64"},
			{Name: "float_field", Type: "FLOAT64"},
			{Name: "bool_field", Type: "BOOL"},
			{Name: "bytes_field", Type: "BYTES"},
			{Name: "json_field", Type: "JSON"},
			{Name: "ts_field", Type: "TIMESTAMP"},
			{Name: "str_lower", Type: "string"},
			{Name: "int_lower", Type: "integer"},
		}),
	}
	provider := bigquery.NewTypeProvider(schemas)

	tests := []struct {
		fieldName string
		wantType  *types.Type
		wantFound bool
	}{
		{"str_field", types.StringType, true},
		{"int_field", types.IntType, true},
		{"int64_field", types.IntType, true},
		{"float_field", types.DoubleType, true},
		{"bool_field", types.BoolType, true},
		{"bytes_field", types.BytesType, true},
		{"json_field", types.DynType, true},
		{"ts_field", types.TimestampType, true},
		{"str_lower", types.StringType, true},
		{"int_lower", types.IntType, true},
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
	provider := bigquery.NewTypeProvider(nil)
	// Close should not panic
	provider.Close()
}

// setupBigQueryContainer starts a BigQuery emulator container and returns a client.
func setupBigQueryContainer(ctx context.Context, t *testing.T) (*tcbigquery.Container, *bq.Client) {
	t.Helper()

	container, err := tcbigquery.Run(ctx,
		"ghcr.io/goccy/bigquery-emulator:0.6.6",
		tcbigquery.WithProjectID(testProjectID),
		tcbigquery.WithDataYAML(bytes.NewReader(providerSeedYAML)),
		testcontainers.WithImagePlatform("linux/amd64"),
	)
	if err != nil {
		if runtime.GOARCH == "arm64" || strings.Contains(err.Error(), "no image found") || strings.Contains(err.Error(), "container exited") {
			t.Skipf("Skipping BigQuery integration test: emulator not available on %s/%s: %v", runtime.GOOS, runtime.GOARCH, err)
		}
		t.Fatalf("Failed to start BigQuery emulator container: %v", err)
	}

	opts := []option.ClientOption{
		option.WithEndpoint(container.URI()),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		option.WithoutAuthentication(),
		internaloption.SkipDialSettingsValidation(),
	}

	client, err := bq.NewClient(ctx, container.ProjectID(), opts...)
	if err != nil {
		if termErr := container.Terminate(ctx); termErr != nil {
			t.Logf("failed to terminate container: %v", termErr)
		}
		t.Fatalf("Failed to create BigQuery client: %v", err)
	}

	return container, client
}

func TestLoadTableSchema_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	container, client := setupBigQueryContainer(ctx, t)
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			t.Logf("failed to close client: %v", closeErr)
		}
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	provider, err := bigquery.NewTypeProviderWithClient(ctx, client, testDataset)
	require.NoError(t, err)

	err = provider.LoadTableSchema(ctx, "test_data")
	require.NoError(t, err)

	// Verify schema was loaded
	schemas := provider.GetSchemas()
	assert.Contains(t, schemas, "test_data")

	// Verify FindStructType
	typ, found := provider.FindStructType("test_data")
	assert.True(t, found)
	assert.NotNil(t, typ)

	// Verify FindStructFieldNames
	names, found := provider.FindStructFieldNames("test_data")
	assert.True(t, found)
	assert.Contains(t, names, "id")
	assert.Contains(t, names, "text_val")
	assert.Contains(t, names, "int_val")

	// Verify type mappings from loaded schema
	tests := []struct {
		fieldName string
		wantType  *types.Type
	}{
		{"id", types.IntType},
		{"text_val", types.StringType},
		{"int_val", types.IntType},
		{"float_val", types.DoubleType},
		{"bool_val", types.BoolType},
	}

	for _, tt := range tests {
		t.Run("type_"+tt.fieldName, func(t *testing.T) {
			got, found := provider.FindStructFieldType("test_data", tt.fieldName)
			assert.True(t, found, "field %q should be found", tt.fieldName)
			if found {
				assert.Equal(t, tt.wantType, got.Type, "field %q type mismatch", tt.fieldName)
			}
		})
	}
}

func TestLoadTableSchema_NonexistentTable(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	container, client := setupBigQueryContainer(ctx, t)
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			t.Logf("failed to close client: %v", closeErr)
		}
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	provider, err := bigquery.NewTypeProviderWithClient(ctx, client, testDataset)
	require.NoError(t, err)

	err = provider.LoadTableSchema(ctx, "nonexistent_table")
	require.Error(t, err)
	assert.ErrorIs(t, err, bigquery.ErrInvalidSchema)
}
