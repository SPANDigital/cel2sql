package mysql_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/cel-go/common/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcmysql "github.com/testcontainers/testcontainers-go/modules/mysql"

	"github.com/spandigital/cel2sql/v3/mysql"
)

func TestNewTypeProvider(t *testing.T) {
	schemas := map[string]mysql.Schema{
		"users": mysql.NewSchema([]mysql.FieldSchema{
			{Name: "id", Type: "int"},
			{Name: "name", Type: "varchar"},
			{Name: "email", Type: "text"},
		}),
	}

	provider := mysql.NewTypeProvider(schemas)
	require.NotNil(t, provider)

	got := provider.GetSchemas()
	assert.Len(t, got, 1)
	assert.Contains(t, got, "users")
}

func TestNewTypeProviderWithConnection_NilDB(t *testing.T) {
	_, err := mysql.NewTypeProviderWithConnection(context.Background(), nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, mysql.ErrInvalidSchema)
}

func TestLoadTableSchema_NoDB(t *testing.T) {
	provider := mysql.NewTypeProvider(nil)
	err := provider.LoadTableSchema(context.Background(), "test")
	require.Error(t, err)
	assert.ErrorIs(t, err, mysql.ErrInvalidSchema)
}

func TestTypeProvider_FindStructType(t *testing.T) {
	schemas := map[string]mysql.Schema{
		"users": mysql.NewSchema([]mysql.FieldSchema{
			{Name: "id", Type: "int"},
			{Name: "name", Type: "varchar"},
		}),
	}
	provider := mysql.NewTypeProvider(schemas)

	typ, found := provider.FindStructType("users")
	assert.True(t, found)
	assert.NotNil(t, typ)

	_, found = provider.FindStructType("nonexistent")
	assert.False(t, found)
}

func TestTypeProvider_FindStructFieldNames(t *testing.T) {
	schemas := map[string]mysql.Schema{
		"users": mysql.NewSchema([]mysql.FieldSchema{
			{Name: "id", Type: "int"},
			{Name: "name", Type: "varchar"},
			{Name: "email", Type: "text"},
		}),
	}
	provider := mysql.NewTypeProvider(schemas)

	names, found := provider.FindStructFieldNames("users")
	assert.True(t, found)
	assert.ElementsMatch(t, []string{"id", "name", "email"}, names)

	_, found = provider.FindStructFieldNames("nonexistent")
	assert.False(t, found)
}

func TestTypeProvider_FindStructFieldType(t *testing.T) {
	schemas := map[string]mysql.Schema{
		"test_table": mysql.NewSchema([]mysql.FieldSchema{
			{Name: "str_field", Type: "varchar"},
			{Name: "int_field", Type: "int"},
			{Name: "bigint_field", Type: "bigint"},
			{Name: "float_field", Type: "float"},
			{Name: "double_field", Type: "double"},
			{Name: "decimal_field", Type: "decimal"},
			{Name: "bool_field", Type: "boolean"},
			{Name: "blob_field", Type: "blob"},
			{Name: "json_field", Type: "json"},
			{Name: "datetime_field", Type: "datetime"},
			{Name: "timestamp_field", Type: "timestamp"},
			{Name: "text_field", Type: "text"},
			{Name: "enum_field", Type: "enum"},
		}),
	}
	provider := mysql.NewTypeProvider(schemas)

	tests := []struct {
		fieldName string
		wantType  *types.Type
		wantFound bool
	}{
		{"str_field", types.StringType, true},
		{"int_field", types.IntType, true},
		{"bigint_field", types.IntType, true},
		{"float_field", types.DoubleType, true},
		{"double_field", types.DoubleType, true},
		{"decimal_field", types.DoubleType, true},
		{"bool_field", types.BoolType, true},
		{"blob_field", types.BytesType, true},
		{"json_field", types.DynType, true},
		{"datetime_field", types.TimestampType, true},
		{"timestamp_field", types.TimestampType, true},
		{"text_field", types.StringType, true},
		{"enum_field", types.StringType, true},
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
	provider := mysql.NewTypeProvider(nil)
	// Close should not panic
	provider.Close()
}

// setupMySQLContainer starts a MySQL 8 container and returns a database connection.
func setupMySQLContainer(ctx context.Context, t *testing.T) (*tcmysql.MySQLContainer, *sql.DB) {
	t.Helper()

	container, err := tcmysql.Run(ctx, "mysql:8.0",
		tcmysql.WithDatabase("testdb"),
		tcmysql.WithUsername("testuser"),
		tcmysql.WithPassword("testpass"),
	)
	require.NoError(t, err, "Failed to start MySQL container")

	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "3306")
	require.NoError(t, err)

	connStr := fmt.Sprintf("testuser:testpass@tcp(%s:%s)/testdb?parseTime=true",
		host, port.Port())
	db, err := sql.Open("mysql", connStr)
	require.NoError(t, err, "Failed to connect to MySQL database")

	err = db.Ping()
	require.NoError(t, err, "Failed to ping MySQL database")

	return container, db
}

func TestLoadTableSchema_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	container, db := setupMySQLContainer(ctx, t)
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Logf("failed to close db: %v", closeErr)
		}
		if err := testcontainers.TerminateContainer(container); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	// Create test table with various MySQL types
	_, err := db.ExecContext(ctx, `
		CREATE TABLE schema_test (
			id INT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			age INT,
			score DOUBLE,
			price DECIMAL(10,2),
			is_active BOOLEAN,
			avatar BLOB,
			metadata JSON,
			created_at DATETIME,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err)

	provider, err := mysql.NewTypeProviderWithConnection(ctx, db)
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
		{"price", types.DoubleType},
		{"is_active", types.IntType}, // MySQL BOOLEAN is TINYINT(1), data_type shows "tinyint"
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
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	container, db := setupMySQLContainer(ctx, t)
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Logf("failed to close db: %v", closeErr)
		}
		if err := testcontainers.TerminateContainer(container); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	provider, err := mysql.NewTypeProviderWithConnection(ctx, db)
	require.NoError(t, err)

	err = provider.LoadTableSchema(ctx, "nonexistent_table")
	require.Error(t, err)
	assert.ErrorIs(t, err, mysql.ErrInvalidSchema)
}
