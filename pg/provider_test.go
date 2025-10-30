// Package pg provides PostgreSQL type provider for CEL type system integration.
package pg_test

import (
	"testing"

	"github.com/google/cel-go/common/types"
	"github.com/stretchr/testify/assert"

	"github.com/spandigital/cel2sql/v3/pg"
	"github.com/spandigital/cel2sql/v3/test"
)

func Test_typeProvider_FindStructType(t *testing.T) {
	typeProvider := pg.NewTypeProvider(map[string]pg.Schema{
		"trigrams":  test.NewTrigramsTableSchema(),
		"wikipedia": test.NewWikipediaTableSchema(),
	})

	type args struct {
		structType string
	}
	tests := []struct {
		name      string
		args      args
		wantFound bool
	}{
		{
			name:      "trigrams",
			args:      args{structType: "trigrams"},
			wantFound: true,
		},
		{
			name:      "trigrams.cell",
			args:      args{structType: "trigrams.cell"},
			wantFound: true,
		},
		{
			name:      "trigrams.cell.value",
			args:      args{structType: "trigrams.cell.value"},
			wantFound: false, // value is a primitive field, not a composite type
		},
		{
			name:      "not_exists",
			args:      args{structType: "not_exists"},
			wantFound: false,
		},
		{
			name:      "trigrams.cell.not_exists",
			args:      args{structType: "trigrams.cell.not_exists"},
			wantFound: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotFound := typeProvider.FindStructType(tt.args.structType)
			assert.Equal(t, tt.wantFound, gotFound)
			if gotFound {
				assert.NotNil(t, got)
				assert.Equal(t, tt.args.structType, got.TypeName())
			}
		})
	}
}

func Test_typeProvider_FindStructFieldNames(t *testing.T) {
	typeProvider := pg.NewTypeProvider(map[string]pg.Schema{
		"trigrams":  test.NewTrigramsTableSchema(),
		"wikipedia": test.NewWikipediaTableSchema(),
	})

	type args struct {
		structType string
	}
	tests := []struct {
		name           string
		args           args
		wantFieldNames []string
		wantFound      bool
	}{
		{
			name: "wikipedia",
			args: args{structType: "wikipedia"},
			wantFieldNames: []string{
				"title", "id", "language", "wp_namespace", "is_redirect",
				"revision_id", "contributor_ip", "contributor_id", "contributor_username",
				"timestamp", "is_minor", "is_bot", "reversion_id", "comment", "num_characters",
			},
			wantFound: true,
		},
		{
			name: "trigrams",
			args: args{structType: "trigrams"},
			wantFieldNames: []string{
				"ngram", "first", "second", "third", "fourth", "fifth", "cell",
			},
			wantFound: true,
		},
		{
			name: "trigrams.cell",
			args: args{structType: "trigrams.cell"},
			wantFieldNames: []string{
				"value", "volume_count", "volume_fraction", "page_count", "match_count", "sample",
			},
			wantFound: true,
		},
		{
			name: "trigrams.cell.sample",
			args: args{structType: "trigrams.cell.sample"},
			wantFieldNames: []string{
				"id", "text", "title", "subtitle", "authors", "url",
			},
			wantFound: true,
		},
		{
			name:      "not_exists",
			args:      args{structType: "not_exists"},
			wantFound: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotFound := typeProvider.FindStructFieldNames(tt.args.structType)
			assert.Equal(t, tt.wantFound, gotFound)
			if gotFound {
				assert.ElementsMatch(t, tt.wantFieldNames, got)
			}
		})
	}
}

func Test_typeProvider_FindStructFieldType(t *testing.T) {
	typeProvider := pg.NewTypeProvider(map[string]pg.Schema{
		"trigrams":  test.NewTrigramsTableSchema(),
		"wikipedia": test.NewWikipediaTableSchema(),
	})

	type args struct {
		structType string
		fieldName  string
	}
	tests := []struct {
		name      string
		args      args
		wantType  *types.Type
		wantFound bool
	}{
		{
			name: "wikipedia.title",
			args: args{
				structType: "wikipedia",
				fieldName:  "title",
			},
			wantType:  types.StringType,
			wantFound: true,
		},
		{
			name: "wikipedia.id",
			args: args{
				structType: "wikipedia",
				fieldName:  "id",
			},
			wantType:  types.IntType,
			wantFound: true,
		},
		{
			name: "wikipedia.is_redirect",
			args: args{
				structType: "wikipedia",
				fieldName:  "is_redirect",
			},
			wantType:  types.BoolType,
			wantFound: true,
		},
		{
			name: "trigrams.cell",
			args: args{
				structType: "trigrams",
				fieldName:  "cell",
			},
			wantType:  types.NewListType(types.NewObjectType("trigrams.cell")),
			wantFound: true,
		},
		{
			name: "trigrams.cell.value",
			args: args{
				structType: "trigrams.cell",
				fieldName:  "value",
			},
			wantType:  types.NewListType(types.StringType),
			wantFound: true,
		},
		{
			name: "trigrams.cell.sample",
			args: args{
				structType: "trigrams.cell",
				fieldName:  "sample",
			},
			wantType:  types.NewListType(types.NewObjectType("trigrams.cell.sample")),
			wantFound: true,
		},
		{
			name: "trigrams.cell.sample.id",
			args: args{
				structType: "trigrams.cell.sample",
				fieldName:  "id",
			},
			wantType:  types.StringType,
			wantFound: true,
		},
		{
			name: "not_exists_struct",
			args: args{
				structType: "not_exists",
				fieldName:  "",
			},
			wantFound: false,
		},
		{
			name: "not_exists_field",
			args: args{
				structType: "wikipedia",
				fieldName:  "not_exists",
			},
			wantFound: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotFound := typeProvider.FindStructFieldType(tt.args.structType, tt.args.fieldName)
			assert.Equal(t, tt.wantFound, gotFound)
			if gotFound {
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantType, got.Type)
			}
		})
	}
}

func Test_typeProvider_PostgreSQLTypes(t *testing.T) {
	tests := []struct {
		name      string
		pgType    string
		repeated  bool
		wantType  *types.Type
		wantFound bool
	}{
		{
			name:      "uuid",
			pgType:    "uuid",
			wantType:  types.BytesType,
			wantFound: true,
		},
		{
			name:      "uuid array",
			pgType:    "uuid",
			repeated:  true,
			wantType:  types.NewListType(types.BytesType),
			wantFound: true,
		},
		{
			name:      "inet",
			pgType:    "inet",
			wantType:  types.StringType,
			wantFound: true,
		},
		{
			name:      "cidr",
			pgType:    "cidr",
			wantType:  types.StringType,
			wantFound: true,
		},
		{
			name:      "macaddr",
			pgType:    "macaddr",
			wantType:  types.StringType,
			wantFound: true,
		},
		{
			name:      "macaddr8",
			pgType:    "macaddr8",
			wantType:  types.StringType,
			wantFound: true,
		},
		{
			name:      "xml",
			pgType:    "xml",
			wantType:  types.StringType,
			wantFound: true,
		},
		{
			name:      "money",
			pgType:    "money",
			wantType:  types.DoubleType,
			wantFound: true,
		},
		{
			name:      "tsvector",
			pgType:    "tsvector",
			wantType:  types.StringType,
			wantFound: true,
		},
		{
			name:      "tsquery",
			pgType:    "tsquery",
			wantType:  types.StringType,
			wantFound: true,
		},
		{
			name:      "unknown_type returns false",
			pgType:    "unknown_custom_type",
			wantFound: false,
		},
		{
			name:      "point returns false",
			pgType:    "point",
			wantFound: false,
		},
		{
			name:      "polygon returns false",
			pgType:    "polygon",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := pg.NewSchema([]pg.FieldSchema{
				{Name: "test_field", Type: tt.pgType, Repeated: tt.repeated},
			})
			typeProvider := pg.NewTypeProvider(map[string]pg.Schema{
				"test_table": schema,
			})

			got, gotFound := typeProvider.FindStructFieldType("test_table", "test_field")
			assert.Equal(t, tt.wantFound, gotFound)
			if tt.wantFound {
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantType, got.Type)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

// Benchmark tests to verify O(1) performance improvement
func BenchmarkFieldLookup_Small(b *testing.B) {
	// 10 fields - small schema
	fields := make([]pg.FieldSchema, 10)
	for i := 0; i < 10; i++ {
		fields[i] = pg.FieldSchema{
			Name: "field_" + string(rune('a'+i)),
			Type: "text",
		}
	}
	schema := pg.NewSchema(fields)
	typeProvider := pg.NewTypeProvider(map[string]pg.Schema{
		"test_table": schema,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Lookup last field (worst case for O(n), same as O(1) for indexed)
		_, _ = typeProvider.FindStructFieldType("test_table", "field_j")
	}
}

func BenchmarkFieldLookup_Medium(b *testing.B) {
	// 100 fields - medium schema
	fields := make([]pg.FieldSchema, 100)
	for i := 0; i < 100; i++ {
		fields[i] = pg.FieldSchema{
			Name: "field_" + string(rune('0'+i%10)) + string(rune('0'+i/10)),
			Type: "text",
		}
	}
	schema := pg.NewSchema(fields)
	typeProvider := pg.NewTypeProvider(map[string]pg.Schema{
		"test_table": schema,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Lookup last field (worst case for O(n))
		_, _ = typeProvider.FindStructFieldType("test_table", "field_99")
	}
}

func BenchmarkFieldLookup_Large(b *testing.B) {
	// 1000 fields - large schema (real-world worst case)
	fields := make([]pg.FieldSchema, 1000)
	for i := 0; i < 1000; i++ {
		fields[i] = pg.FieldSchema{
			Name: "field_" + string(rune('0'+i%10)) + string(rune('0'+(i/10)%10)) + string(rune('0'+i/100)),
			Type: "text",
		}
	}
	schema := pg.NewSchema(fields)
	typeProvider := pg.NewTypeProvider(map[string]pg.Schema{
		"test_table": schema,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Lookup last field (worst case for O(n))
		_, _ = typeProvider.FindStructFieldType("test_table", "field_999")
	}
}

func BenchmarkFieldNames_Small(b *testing.B) {
	fields := make([]pg.FieldSchema, 10)
	for i := 0; i < 10; i++ {
		fields[i] = pg.FieldSchema{
			Name: "field_" + string(rune('a'+i)),
			Type: "text",
		}
	}
	schema := pg.NewSchema(fields)
	typeProvider := pg.NewTypeProvider(map[string]pg.Schema{
		"test_table": schema,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = typeProvider.FindStructFieldNames("test_table")
	}
}

func BenchmarkFieldNames_Large(b *testing.B) {
	fields := make([]pg.FieldSchema, 1000)
	for i := 0; i < 1000; i++ {
		fields[i] = pg.FieldSchema{
			Name: "field_" + string(rune('0'+i%10)) + string(rune('0'+(i/10)%10)) + string(rune('0'+i/100)),
			Type: "text",
		}
	}
	schema := pg.NewSchema(fields)
	typeProvider := pg.NewTypeProvider(map[string]pg.Schema{
		"test_table": schema,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = typeProvider.FindStructFieldNames("test_table")
	}
}
