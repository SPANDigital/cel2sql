// Package pg provides PostgreSQL type provider for CEL type system integration.
package pg_test

import (
	"testing"

	"github.com/google/cel-go/common/types"
	"github.com/stretchr/testify/assert"

	"github.com/spandigital/cel2sql/pg"
	"github.com/spandigital/cel2sql/test"
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
