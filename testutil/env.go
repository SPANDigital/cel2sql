package testutil

import (
	"fmt"

	"github.com/google/cel-go/cel"

	"github.com/spandigital/cel2sql/v3"
	dialectpkg "github.com/spandigital/cel2sql/v3/dialect"
	bigqueryDialect "github.com/spandigital/cel2sql/v3/dialect/bigquery"
	duckdbDialect "github.com/spandigital/cel2sql/v3/dialect/duckdb"
	mysqlDialect "github.com/spandigital/cel2sql/v3/dialect/mysql"
	sqliteDialect "github.com/spandigital/cel2sql/v3/dialect/sqlite"
	"github.com/spandigital/cel2sql/v3/pg"
	"github.com/spandigital/cel2sql/v3/sqltypes"
	"github.com/spandigital/cel2sql/v3/testcases"
)

// EnvResult holds both the CEL environment and convert options needed for testing.
type EnvResult struct {
	Env  *cel.Env
	Opts []cel2sql.ConvertOption
}

// NewDefaultEnv creates a basic CEL environment with standard variable types.
func NewDefaultEnv() (*EnvResult, error) {
	env, err := cel.NewEnv(
		cel.Types(
			sqltypes.Date, sqltypes.Time, sqltypes.DateTime, sqltypes.Interval, sqltypes.DatePart,
		),
		cel.Variable("name", cel.StringType),
		cel.Variable("age", cel.IntType),
		cel.Variable("adult", cel.BoolType),
		cel.Variable("height", cel.DoubleType),
		cel.Variable("string_list", cel.ListType(cel.StringType)),
		cel.Variable("string_int_map", cel.MapType(cel.StringType, cel.IntType)),
		cel.Variable("null_var", cel.NullType),
		cel.Variable("created_at", cel.TimestampType),
		cel.Variable("page", cel.MapType(cel.StringType, cel.StringType)),
		cel.Variable("salary", cel.DoubleType),
		cel.Variable("active", cel.BoolType),
		cel.Variable("data", cel.BytesType),
		cel.Variable("tags", cel.ListType(cel.StringType)),
		cel.Variable("scores", cel.ListType(cel.IntType)),
		// Cast functions
		cel.Function("bool", cel.Overload("bool_from_int", []*cel.Type{cel.IntType}, cel.BoolType)),
		cel.Function("int", cel.Overload("int_from_bool", []*cel.Type{cel.BoolType}, cel.IntType)),
	)
	if err != nil {
		return nil, err
	}
	return &EnvResult{Env: env}, nil
}

// NewTimestampEnv creates a CEL environment with timestamp-related types and functions.
func NewTimestampEnv() (*EnvResult, error) {
	env, err := cel.NewEnv(
		cel.Types(
			sqltypes.Date, sqltypes.Time, sqltypes.DateTime, sqltypes.Interval, sqltypes.DatePart,
		),
		cel.Variable("name", cel.StringType),
		cel.Variable("age", cel.IntType),
		cel.Variable("adult", cel.BoolType),
		cel.Variable("height", cel.DoubleType),
		cel.Variable("string_list", cel.ListType(cel.StringType)),
		cel.Variable("string_int_map", cel.MapType(cel.StringType, cel.IntType)),
		cel.Variable("null_var", cel.NullType),
		cel.Variable("birthday", cel.ObjectType("DATE")),
		cel.Variable("fixed_time", cel.ObjectType("TIME")),
		cel.Variable("scheduled_at", cel.ObjectType("DATETIME")),
		cel.Variable("created_at", cel.TimestampType),
		cel.Variable("page", cel.MapType(cel.StringType, cel.StringType)),
		// Date part constants
		cel.Variable("YEAR", cel.ObjectType("date_part")),
		cel.Variable("MONTH", cel.ObjectType("date_part")),
		cel.Variable("DAY", cel.ObjectType("date_part")),
		cel.Variable("HOUR", cel.ObjectType("date_part")),
		cel.Variable("MINUTE", cel.ObjectType("date_part")),
		cel.Variable("SECOND", cel.ObjectType("date_part")),
		// SQL functions
		cel.Function("date",
			cel.Overload("date_string", []*cel.Type{cel.StringType}, cel.ObjectType("DATE")),
			cel.Overload("date_int_int_int", []*cel.Type{cel.IntType, cel.IntType, cel.IntType}, cel.ObjectType("DATE"))),
		cel.Function("time", cel.Overload("time_string", []*cel.Type{cel.StringType}, cel.ObjectType("TIME"))),
		cel.Function("datetime",
			cel.Overload("datetime_string", []*cel.Type{cel.StringType}, cel.ObjectType("DATETIME")),
			cel.Overload("datetime_date_time", []*cel.Type{cel.ObjectType("DATE"), cel.ObjectType("TIME")}, cel.ObjectType("DATETIME"))),
		cel.Function("timestamp",
			cel.Overload("timestamp_datetime_string", []*cel.Type{cel.ObjectType("DATETIME"), cel.StringType}, cel.TimestampType)),
		cel.Function("interval", cel.Overload("interval_int_datepart", []*cel.Type{cel.IntType, cel.ObjectType("date_part")}, cel.ObjectType("INTERVAL"))),
		cel.Function("current_date", cel.Overload("current_date", []*cel.Type{}, cel.ObjectType("DATE"))),
		cel.Function("current_datetime", cel.Overload("current_datetime_string", []*cel.Type{cel.StringType}, cel.ObjectType("DATETIME"))),
		// Date/Time arithmetic operators
		cel.Function("_+_",
			cel.Overload("date_add_interval", []*cel.Type{cel.ObjectType("DATE"), cel.ObjectType("INTERVAL")}, cel.ObjectType("DATE")),
			cel.Overload("date_add_int", []*cel.Type{cel.ObjectType("DATE"), cel.IntType}, cel.ObjectType("DATE")),
			cel.Overload("time_add_interval", []*cel.Type{cel.ObjectType("TIME"), cel.ObjectType("INTERVAL")}, cel.ObjectType("TIME")),
			cel.Overload("datetime_add_interval", []*cel.Type{cel.ObjectType("DATETIME"), cel.ObjectType("INTERVAL")}, cel.ObjectType("DATETIME")),
			cel.Overload("timestamp_add_interval", []*cel.Type{cel.TimestampType, cel.ObjectType("INTERVAL")}, cel.TimestampType)),
		cel.Function("_-_",
			cel.Overload("date_sub_interval", []*cel.Type{cel.ObjectType("DATE"), cel.ObjectType("INTERVAL")}, cel.ObjectType("DATE")),
			cel.Overload("time_sub_interval", []*cel.Type{cel.ObjectType("TIME"), cel.ObjectType("INTERVAL")}, cel.ObjectType("TIME")),
			cel.Overload("datetime_sub_interval", []*cel.Type{cel.ObjectType("DATETIME"), cel.ObjectType("INTERVAL")}, cel.ObjectType("DATETIME")),
			cel.Overload("timestamp_sub_interval", []*cel.Type{cel.TimestampType, cel.ObjectType("INTERVAL")}, cel.TimestampType)),
		// Date/Time comparison operators
		cel.Function("_>_",
			cel.Overload("date_gt_date", []*cel.Type{cel.ObjectType("DATE"), cel.ObjectType("DATE")}, cel.BoolType)),
		// Date/Time methods
		cel.Function("getFullYear", cel.MemberOverload("date_getFullYear", []*cel.Type{cel.ObjectType("DATE")}, cel.IntType)),
		cel.Function("getMonth", cel.MemberOverload("datetime_getMonth", []*cel.Type{cel.ObjectType("DATETIME")}, cel.IntType)),
		cel.Function("getDayOfMonth", cel.MemberOverload("datetime_getDayOfMonth", []*cel.Type{cel.ObjectType("DATETIME")}, cel.IntType)),
		cel.Function("getMinutes", cel.MemberOverload("time_getMinutes", []*cel.Type{cel.ObjectType("TIME")}, cel.IntType)),
		// Cast functions
		cel.Function("bool", cel.Overload("bool_from_int", []*cel.Type{cel.IntType}, cel.BoolType)),
		cel.Function("int", cel.Overload("int_from_bool", []*cel.Type{cel.BoolType}, cel.IntType)),
	)
	if err != nil {
		return nil, err
	}
	return &EnvResult{Env: env}, nil
}

// NewJSONSchemaEnv creates a CEL environment with a JSON-enabled schema type provider.
func NewJSONSchemaEnv() (*EnvResult, error) {
	productSchema := testcases.NewProductPGSchema()
	schemas := map[string]pg.Schema{
		"product": productSchema,
	}
	provider := pg.NewTypeProvider(schemas)
	env, err := cel.NewEnv(
		cel.CustomTypeProvider(provider),
		cel.Variable("product", cel.ObjectType("product")),
	)
	if err != nil {
		return nil, err
	}
	return &EnvResult{
		Env:  env,
		Opts: []cel2sql.ConvertOption{cel2sql.WithSchemas(schemas)},
	}, nil
}

// PostgreSQLEnvFactory returns an environment factory for PostgreSQL tests.
func PostgreSQLEnvFactory() func(envSetup string) (*EnvResult, error) {
	return func(envSetup string) (*EnvResult, error) {
		switch envSetup {
		case testcases.EnvDefault:
			return NewDefaultEnv()
		case testcases.EnvWithTimestamp:
			return NewTimestampEnv()
		case testcases.EnvWithJSON:
			return NewJSONSchemaEnv()
		default:
			return nil, fmt.Errorf("unknown environment setup: %s", envSetup)
		}
	}
}

// MySQLEnvFactory returns an environment factory for MySQL tests.
// It uses the same CEL environments as PostgreSQL (CEL compilation is dialect-independent)
// but sets the MySQL dialect for SQL generation.
func MySQLEnvFactory() func(envSetup string) (*EnvResult, error) {
	return func(envSetup string) (*EnvResult, error) {
		switch envSetup {
		case testcases.EnvDefault:
			result, err := NewDefaultEnv()
			if err != nil {
				return nil, err
			}
			result.Opts = append(result.Opts, cel2sql.WithDialect(mysqlDialect.New()))
			return result, nil
		case testcases.EnvWithTimestamp:
			result, err := NewTimestampEnv()
			if err != nil {
				return nil, err
			}
			result.Opts = append(result.Opts, cel2sql.WithDialect(mysqlDialect.New()))
			return result, nil
		case testcases.EnvWithJSON:
			result, err := NewJSONSchemaEnv()
			if err != nil {
				return nil, err
			}
			result.Opts = append(result.Opts, cel2sql.WithDialect(mysqlDialect.New()))
			return result, nil
		default:
			return nil, fmt.Errorf("unknown environment setup: %s", envSetup)
		}
	}
}

// SQLiteEnvFactory returns an environment factory for SQLite tests.
// It uses the same CEL environments as PostgreSQL (CEL compilation is dialect-independent)
// but sets the SQLite dialect for SQL generation.
func SQLiteEnvFactory() func(envSetup string) (*EnvResult, error) {
	return func(envSetup string) (*EnvResult, error) {
		switch envSetup {
		case testcases.EnvDefault:
			result, err := NewDefaultEnv()
			if err != nil {
				return nil, err
			}
			result.Opts = append(result.Opts, cel2sql.WithDialect(sqliteDialect.New()))
			return result, nil
		case testcases.EnvWithTimestamp:
			result, err := NewTimestampEnv()
			if err != nil {
				return nil, err
			}
			result.Opts = append(result.Opts, cel2sql.WithDialect(sqliteDialect.New()))
			return result, nil
		case testcases.EnvWithJSON:
			result, err := NewJSONSchemaEnv()
			if err != nil {
				return nil, err
			}
			result.Opts = append(result.Opts, cel2sql.WithDialect(sqliteDialect.New()))
			return result, nil
		default:
			return nil, fmt.Errorf("unknown environment setup: %s", envSetup)
		}
	}
}

// DuckDBEnvFactory returns an environment factory for DuckDB tests.
// It uses the same CEL environments as PostgreSQL (CEL compilation is dialect-independent)
// but sets the DuckDB dialect for SQL generation.
func DuckDBEnvFactory() func(envSetup string) (*EnvResult, error) {
	return func(envSetup string) (*EnvResult, error) {
		switch envSetup {
		case testcases.EnvDefault:
			result, err := NewDefaultEnv()
			if err != nil {
				return nil, err
			}
			result.Opts = append(result.Opts, cel2sql.WithDialect(duckdbDialect.New()))
			return result, nil
		case testcases.EnvWithTimestamp:
			result, err := NewTimestampEnv()
			if err != nil {
				return nil, err
			}
			result.Opts = append(result.Opts, cel2sql.WithDialect(duckdbDialect.New()))
			return result, nil
		case testcases.EnvWithJSON:
			result, err := NewJSONSchemaEnv()
			if err != nil {
				return nil, err
			}
			result.Opts = append(result.Opts, cel2sql.WithDialect(duckdbDialect.New()))
			return result, nil
		default:
			return nil, fmt.Errorf("unknown environment setup: %s", envSetup)
		}
	}
}

// BigQueryEnvFactory returns an environment factory for BigQuery tests.
// It uses the same CEL environments as PostgreSQL (CEL compilation is dialect-independent)
// but sets the BigQuery dialect for SQL generation.
func BigQueryEnvFactory() func(envSetup string) (*EnvResult, error) {
	return func(envSetup string) (*EnvResult, error) {
		switch envSetup {
		case testcases.EnvDefault:
			result, err := NewDefaultEnv()
			if err != nil {
				return nil, err
			}
			result.Opts = append(result.Opts, cel2sql.WithDialect(bigqueryDialect.New()))
			return result, nil
		case testcases.EnvWithTimestamp:
			result, err := NewTimestampEnv()
			if err != nil {
				return nil, err
			}
			result.Opts = append(result.Opts, cel2sql.WithDialect(bigqueryDialect.New()))
			return result, nil
		case testcases.EnvWithJSON:
			result, err := NewJSONSchemaEnv()
			if err != nil {
				return nil, err
			}
			result.Opts = append(result.Opts, cel2sql.WithDialect(bigqueryDialect.New()))
			return result, nil
		default:
			return nil, fmt.Errorf("unknown environment setup: %s", envSetup)
		}
	}
}

// DialectEnvFactory returns an environment factory for the given dialect.
// This is a convenience function that maps dialect names to their env factories.
func DialectEnvFactory(d dialectpkg.Name) func(envSetup string) (*EnvResult, error) {
	switch d {
	case dialectpkg.PostgreSQL:
		return PostgreSQLEnvFactory()
	case dialectpkg.MySQL:
		return MySQLEnvFactory()
	case dialectpkg.SQLite:
		return SQLiteEnvFactory()
	case dialectpkg.DuckDB:
		return DuckDBEnvFactory()
	case dialectpkg.BigQuery:
		return BigQueryEnvFactory()
	default:
		return func(_ string) (*EnvResult, error) {
			return nil, fmt.Errorf("no environment factory for dialect %s", d)
		}
	}
}
