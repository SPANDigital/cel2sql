package cel2sql_test

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v2"
	"github.com/spandigital/cel2sql/v2/sqltypes"
)

func TestConvert(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Types(
			// Custom abstract types
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
		cel.Variable("page", cel.MapType(cel.StringType, cel.StringType)), // simplified version
		cel.Variable("trigram", cel.MapType(cel.StringType, cel.DynType)), // simplified version
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
	require.NoError(t, err)
	type args struct {
		source string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "startsWith",
			args:    args{source: `name.startsWith("a")`},
			want:    "name LIKE 'a%'",
			wantErr: false,
		},
		{
			name:    "endsWith",
			args:    args{source: `name.endsWith("z")`},
			want:    "name LIKE '%z'",
			wantErr: false,
		},
		{
			name:    "matches",
			args:    args{source: `name.matches("a+")`},
			want:    "name ~ 'a+'",
			wantErr: false,
		},
		{
			name:    "matches_function_style",
			args:    args{source: `matches(name, "^[0-9]+$")`},
			want:    "name ~ '^[0-9]+$'",
			wantErr: false,
		},
		{
			name:    "matches_with_word_boundary",
			args:    args{source: `name.matches("\\btest\\b")`},
			want:    "name ~ '\\ytest\\y'",
			wantErr: false,
		},
		{
			name:    "matches_with_digit_class",
			args:    args{source: `name.matches("\\d{3}-\\d{4}")`},
			want:    "name ~ '[[:digit:]]{3}-[[:digit:]]{4}'",
			wantErr: false,
		},
		{
			name:    "matches_with_word_class",
			args:    args{source: `name.matches("\\w+@\\w+\\.\\w+")`},
			want:    "name ~ '[[:alnum:]_]+@[[:alnum:]_]+\\.[[:alnum:]_]+'",
			wantErr: false,
		},
		{
			name:    "matches_complex_pattern",
			args:    args{source: `name.matches(".*pattern.*")`},
			want:    "name ~ '.*pattern.*'",
			wantErr: false,
		},
		{
			name:    "contains",
			args:    args{source: `name.contains("abc")`},
			want:    "POSITION('abc' IN name) > 0",
			wantErr: false,
		},
		{
			name:    "&&",
			args:    args{source: `name.startsWith("a") && name.endsWith("z")`},
			want:    "name LIKE 'a%' AND name LIKE '%z'",
			wantErr: false,
		},
		{
			name:    "||",
			args:    args{source: `name.startsWith("a") || name.endsWith("z")`},
			want:    "name LIKE 'a%' OR name LIKE '%z'",
			wantErr: false,
		},
		{
			name:    "()",
			args:    args{source: `age >= 10 && (name.startsWith("a") || name.endsWith("z"))`},
			want:    "age >= 10 AND (name LIKE 'a%' OR name LIKE '%z')",
			wantErr: false,
		},
		{
			name:    "IF",
			args:    args{source: `name == "a" ? "a" : "b"`},
			want:    "CASE WHEN name = 'a' THEN 'a' ELSE 'b' END",
			wantErr: false,
		},
		{
			name:    "==",
			args:    args{source: `name == "a"`},
			want:    "name = 'a'",
			wantErr: false,
		},
		{
			name:    "!=",
			args:    args{source: `age != 20`},
			want:    "age != 20",
			wantErr: false,
		},
		{
			name:    "IS NULL",
			args:    args{source: `null_var == null`},
			want:    "null_var IS NULL",
			wantErr: false,
		},
		{
			name:    "IS NOT TRUE",
			args:    args{source: `adult != true`},
			want:    "adult IS NOT TRUE",
			wantErr: false,
		},
		{
			name:    "<",
			args:    args{source: `age < 20`},
			want:    "age < 20",
			wantErr: false,
		},
		{
			name:    ">=",
			args:    args{source: `height >= 1.6180339887`},
			want:    "height >= 1.6180339887",
			wantErr: false,
		},
		{
			name:    "NOT",
			args:    args{source: `!adult`},
			want:    "NOT adult",
			wantErr: false,
		},
		{
			name:    "-",
			args:    args{source: `-1`},
			want:    "-1",
			wantErr: false,
		},
		{
			name:    "list",
			args:    args{source: `[1, 2, 3][0] == 1`},
			want:    "ARRAY[1, 2, 3][1] = 1", // PostgreSQL arrays are 1-indexed
			wantErr: false,
		},
		{
			name:    "list_var",
			args:    args{source: `string_list[0] == "a"`},
			want:    "string_list[1] = 'a'", // PostgreSQL arrays are 1-indexed
			wantErr: false,
		},
		{
			name:    "array_index_overflow",
			args:    args{source: `string_list[9223372036854775807]`}, // math.MaxInt64
			want:    "",
			wantErr: true,
		},
		{
			name:    "array_index_negative",
			args:    args{source: `string_list[-1]`},
			want:    "",
			wantErr: true,
		},
		{
			name:    "map",
			args:    args{source: `{"one": 1, "two": 2, "three": 3}["one"] == 1`},
			want:    "ROW(1, 2, 3).one = 1",
			wantErr: false,
		},
		{
			name:    "map_var",
			args:    args{source: `string_int_map["one"] == 1`},
			want:    "string_int_map.one = 1",
			wantErr: false,
		},
		{
			name:    "invalidFieldType",
			args:    args{source: `{1: 1}[1]`},
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalidFieldName",
			args:    args{source: `{"on e": 1}["on e"]`},
			want:    "",
			wantErr: true,
		},
		{
			name:    "add",
			args:    args{source: `1 + 2 == 3`},
			want:    "1 + 2 = 3",
			wantErr: false,
		},
		{
			name:    "concatString",
			args:    args{source: `"a" + "b" == "ab"`},
			want:    "'a' || 'b' = 'ab'",
			wantErr: false,
		},
		{
			name:    "concatList",
			args:    args{source: `1 in [1] + [2, 3]`},
			want:    "1 = ANY(ARRAY[1] || ARRAY[2, 3])", // PostgreSQL array concatenation and membership
			wantErr: false,
		},
		{
			name:    "modulo",
			args:    args{source: `5 % 3 == 2`},
			want:    "MOD(5, 3) = 2",
			wantErr: false,
		},
		{
			name:    "date",
			args:    args{source: `birthday > date(2000, 1, 1) + 1`},
			want:    "birthday > DATE(2000, 1, 1) + 1",
			wantErr: false,
		},
		{
			name:    "time",
			args:    args{source: `fixed_time == time("18:00:00")`},
			want:    "fixed_time = TIME('18:00:00')",
			wantErr: false,
		},
		{
			name:    "datetime",
			args:    args{source: `scheduled_at != datetime(date("2021-09-01"), fixed_time)`},
			want:    "scheduled_at != DATETIME(DATE('2021-09-01'), fixed_time)",
			wantErr: false,
		},
		{
			name:    "timestamp",
			args:    args{source: `created_at - duration("60m") <= timestamp(datetime("2021-09-01 18:00:00"), "Asia/Tokyo")`},
			want:    "created_at - INTERVAL 1 HOUR <= DATETIME('2021-09-01 18:00:00') AT TIME ZONE 'Asia/Tokyo'",
			wantErr: false,
		},
		{
			name:    "duration_second",
			args:    args{source: `duration("10s")`},
			want:    "INTERVAL 10 SECOND",
			wantErr: false,
		},
		{
			name:    "duration_minute",
			args:    args{source: `duration("1h1m")`},
			want:    "INTERVAL 61 MINUTE",
			wantErr: false,
		},
		{
			name:    "duration_hour",
			args:    args{source: `duration("60m")`},
			want:    "INTERVAL 1 HOUR",
			wantErr: false,
		},
		{
			name:    "interval",
			args:    args{source: `interval(1, MONTH)`},
			want:    "INTERVAL 1 MONTH",
			wantErr: false,
		},
		{
			name:    "date_add",
			args:    args{source: `date("2021-09-01") + interval(1, DAY)`},
			want:    "DATE('2021-09-01') + INTERVAL 1 DAY",
			wantErr: false,
		},
		{
			name:    "date_sub",
			args:    args{source: `current_date() - interval(1, DAY)`},
			want:    "CURRENT_DATE() - INTERVAL 1 DAY",
			wantErr: false,
		},
		{
			name:    "time_add",
			args:    args{source: `time("09:00:00") + interval(1, MINUTE)`},
			want:    "TIME('09:00:00') + INTERVAL 1 MINUTE",
			wantErr: false,
		},
		{
			name:    "time_sub",
			args:    args{source: `time("09:00:00") - interval(1, MINUTE)`},
			want:    "TIME('09:00:00') - INTERVAL 1 MINUTE",
			wantErr: false,
		},
		{
			name:    "datetime_add",
			args:    args{source: `datetime("2021-09-01 18:00:00") + interval(1, MINUTE)`},
			want:    "DATETIME('2021-09-01 18:00:00') + INTERVAL 1 MINUTE",
			wantErr: false,
		},
		{
			name:    "datetime_sub",
			args:    args{source: `current_datetime("Asia/Tokyo") - interval(1, MINUTE)`},
			want:    "CURRENT_DATETIME('Asia/Tokyo') - INTERVAL 1 MINUTE",
			wantErr: false,
		},
		{
			name:    "timestamp_add",
			args:    args{source: `duration("1h") + timestamp("2021-09-01T18:00:00Z")`},
			want:    "CAST('2021-09-01T18:00:00Z' AS TIMESTAMP WITH TIME ZONE) + INTERVAL 1 HOUR",
			wantErr: false,
		},
		{
			name:    "timestamp_sub",
			args:    args{source: `created_at - interval(1, HOUR)`},
			want:    "created_at - INTERVAL 1 HOUR",
			wantErr: false,
		},
		{
			name:    "timestamp_getSeconds",
			args:    args{source: `created_at.getSeconds()`},
			want:    "EXTRACT(SECOND FROM created_at)",
			wantErr: false,
		},
		{
			name:    "\"timestamp_getHours_withTimezone",
			args:    args{source: `created_at.getHours("Asia/Tokyo")`},
			want:    "EXTRACT(HOUR FROM created_at AT 'Asia/Tokyo')",
			wantErr: false,
		},
		{
			name:    "date_getFullYear",
			args:    args{source: `birthday.getFullYear()`},
			want:    "EXTRACT(YEAR FROM birthday)",
			wantErr: false,
		},
		{
			name:    "datetime_getMonth",
			args:    args{source: `scheduled_at.getMonth()`},
			want:    "EXTRACT(MONTH FROM scheduled_at) - 1",
			wantErr: false,
		},
		{
			name:    "datetime_getDayOfMonth",
			args:    args{source: `scheduled_at.getDayOfMonth()`},
			want:    "EXTRACT(DAY FROM scheduled_at) - 1",
			wantErr: false,
		},
		{
			name:    "time_getMinutes",
			args:    args{source: `fixed_time.getMinutes()`},
			want:    "EXTRACT(MINUTE FROM fixed_time)",
			wantErr: false,
		},
		{
			name:    "fieldSelect",
			args:    args{source: `page.title == "test"`},
			want:    "page.title = 'test'",
			wantErr: false,
		},
		{
			name:    "fieldSelect_startsWith",
			args:    args{source: `page.title.startsWith("test")`},
			want:    "page.title LIKE 'test%'",
			wantErr: false,
		},
		{
			name:    "fieldSelect_add",
			args:    args{source: `trigram.cell[0].page_count + 1`},
			want:    "trigram.cell[1].page_count + 1", // PostgreSQL 1-indexed arrays
			wantErr: false,
		},
		{
			name:    "fieldSelect_concatString",
			args:    args{source: `trigram.cell[0].sample[0].title + "test"`},
			want:    "trigram.cell[1].sample[1].title || 'test'", // PostgreSQL syntax
			wantErr: false,
		},
		{
			name:    "fieldSelect_in",
			args:    args{source: `"test" in trigram.cell[0].value`},
			want:    "'test' = ANY(trigram.cell[1].value)", // PostgreSQL array membership
			wantErr: false,
		},
		{
			name:    "cast_bool",
			args:    args{source: `bool(0) == false`},
			want:    "CAST(0 AS BOOLEAN) IS FALSE",
			wantErr: false,
		},
		{
			name:    "cast_bytes",
			args:    args{source: `bytes("test")`},
			want:    "CAST('test' AS BYTEA)",
			wantErr: false,
		},
		{
			name:    "cast_int",
			args:    args{source: `int(true) == 1`},
			want:    "CAST(TRUE AS BIGINT) = 1",
			wantErr: false,
		},
		{
			name:    "cast_string",
			args:    args{source: `string(true) == "true"`},
			want:    "CAST(TRUE AS TEXT) = 'true'",
			wantErr: false,
		},
		{
			name:    "cast_string_from_timestamp",
			args:    args{source: `string(created_at)`},
			want:    "CAST(created_at AS TEXT)",
			wantErr: false,
		},
		{
			name:    "cast_int_epoch",
			args:    args{source: `int(created_at)`},
			want:    "EXTRACT(EPOCH FROM created_at)::bigint",
			wantErr: false,
		},
		{
			name:    "size_string",
			args:    args{source: `size("test")`},
			want:    "LENGTH('test')",
			wantErr: false,
		},
		{
			name:    "size_bytes",
			args:    args{source: `size(bytes("test"))`},
			want:    "LENGTH(CAST('test' AS BYTEA))",
			wantErr: false,
		},
		{
			name:    "size_list",
			args:    args{source: `size(string_list)`},
			want:    "ARRAY_LENGTH(string_list, 1)",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.args.source)
			require.Empty(t, issues)

			got, err := cel2sql.Convert(ast)
			if !tt.wantErr && assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestNullByteRejection tests that string literals with null bytes are rejected
func TestNullByteRejection(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("name", cel.StringType),
	)
	require.NoError(t, err)

	tests := []struct {
		name    string
		expr    string
		wantErr string
	}{
		{
			name:    "null byte at start",
			expr:    `name == "\x00test"`,
			wantErr: "string literals cannot contain null bytes",
		},
		{
			name:    "null byte in middle",
			expr:    `name == "test\x00value"`,
			wantErr: "string literals cannot contain null bytes",
		},
		{
			name:    "null byte at end",
			expr:    `name == "test\x00"`,
			wantErr: "string literals cannot contain null bytes",
		},
		{
			name:    "only null byte",
			expr:    `name == "\x00"`,
			wantErr: "string literals cannot contain null bytes",
		},
		{
			name:    "multiple null bytes",
			expr:    `name == "\x00\x00\x00"`,
			wantErr: "string literals cannot contain null bytes",
		},
		{
			name:    "valid string without null byte",
			expr:    `name == "valid"`,
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expr)
			require.NoError(t, issues.Err())

			got, err := cel2sql.Convert(ast)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "name = 'valid'", got)
			}
		})
	}
}

// TestNullByteInLikePatterns tests that LIKE patterns with null bytes are rejected
func TestNullByteInLikePatterns(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("name", cel.StringType),
	)
	require.NoError(t, err)

	tests := []struct {
		name    string
		expr    string
		wantErr string
	}{
		{
			name:    "startsWith with null byte at start",
			expr:    `name.startsWith("\x00test")`,
			wantErr: "LIKE patterns cannot contain null bytes",
		},
		{
			name:    "startsWith with null byte in middle",
			expr:    `name.startsWith("test\x00value")`,
			wantErr: "LIKE patterns cannot contain null bytes",
		},
		{
			name:    "startsWith with only null byte",
			expr:    `name.startsWith("\x00")`,
			wantErr: "LIKE patterns cannot contain null bytes",
		},
		{
			name:    "endsWith with null byte at end",
			expr:    `name.endsWith("test\x00")`,
			wantErr: "LIKE patterns cannot contain null bytes",
		},
		{
			name:    "endsWith with null byte in middle",
			expr:    `name.endsWith("\x00value")`,
			wantErr: "LIKE patterns cannot contain null bytes",
		},
		{
			name:    "endsWith with only null byte",
			expr:    `name.endsWith("\x00")`,
			wantErr: "LIKE patterns cannot contain null bytes",
		},
		{
			name:    "startsWith valid pattern",
			expr:    `name.startsWith("valid")`,
			wantErr: "",
		},
		{
			name:    "endsWith valid pattern",
			expr:    `name.endsWith("valid")`,
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expr)
			require.NoError(t, issues.Err())

			got, err := cel2sql.Convert(ast)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
				// Just verify it doesn't error - actual SQL format tested elsewhere
				assert.NotEmpty(t, got)
			}
		})
	}
}

// TestNullByteInRegexPatterns tests that regex patterns with null bytes are rejected
func TestNullByteInRegexPatterns(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Variable("email", cel.StringType),
	)
	require.NoError(t, err)

	tests := []struct {
		name    string
		expr    string
		wantErr string
	}{
		{
			name:    "matches with only null byte",
			expr:    `email.matches("\x00")`,
			wantErr: "regex patterns cannot contain null bytes",
		},
		{
			name:    "matches with null byte at start",
			expr:    `email.matches("\x00test")`,
			wantErr: "regex patterns cannot contain null bytes",
		},
		{
			name:    "matches with null byte at end",
			expr:    `email.matches("test\x00")`,
			wantErr: "regex patterns cannot contain null bytes",
		},
		{
			name:    "matches with null byte in middle",
			expr:    `email.matches("te\x00st")`,
			wantErr: "regex patterns cannot contain null bytes",
		},
		{
			name:    "matches function style with null byte",
			expr:    `matches(email, "\x00")`,
			wantErr: "regex patterns cannot contain null bytes",
		},
		{
			name:    "matches with multiple null bytes",
			expr:    `email.matches("\x00\x00\x00")`,
			wantErr: "regex patterns cannot contain null bytes",
		},
		{
			name:    "matches with valid pattern",
			expr:    `email.matches(r"^[a-z]+@.*\.com$")`,
			wantErr: "",
		},
		{
			name:    "matches with valid complex pattern",
			expr:    `email.matches(r".*")`,
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, issues := env.Compile(tt.expr)
			require.NoError(t, issues.Err())

			got, err := cel2sql.Convert(ast)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
				// Just verify it doesn't error - actual SQL format tested elsewhere
				assert.NotEmpty(t, got)
			}
		})
	}
}
