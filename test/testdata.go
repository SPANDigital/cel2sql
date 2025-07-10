package test

import (
	"github.com/cockscomb/cel2sql/pg"
)

func NewTrigramsTableSchema() pg.Schema {
	return pg.Schema{
		{
			Name:     "ngram",
			Type:     "text",
			Repeated: false,
		},
		{
			Name:     "first",
			Type:     "text",
			Repeated: false,
		},
		{
			Name:     "second",
			Type:     "text",
			Repeated: false,
		},
		{
			Name:     "third",
			Type:     "text",
			Repeated: false,
		},
		{
			Name:     "fourth",
			Type:     "text",
			Repeated: false,
		},
		{
			Name:     "fifth",
			Type:     "text",
			Repeated: false,
		},
		{
			Name:     "cell",
			Type:     "composite",
			Repeated: true, // Array of composite types
			Schema: []pg.FieldSchema{
				{
					Name:     "value",
					Type:     "text",
					Repeated: true, // Array of text values
				},
				{
					Name:     "volume_count",
					Type:     "bigint",
					Repeated: false,
				},
				{
					Name:     "volume_fraction",
					Type:     "double precision",
					Repeated: false,
				},
				{
					Name:     "page_count",
					Type:     "bigint",
					Repeated: false,
				},
				{
					Name:     "match_count",
					Type:     "bigint",
					Repeated: false,
				},
				{
					Name:     "sample",
					Type:     "composite",
					Repeated: true,
					Schema: []pg.FieldSchema{
						{
							Name:     "id",
							Type:     "text",
							Repeated: false,
						},
						{
							Name:     "text",
							Type:     "text",
							Repeated: false,
						},
						{
							Name:     "title",
							Type:     "text",
							Repeated: false,
						},
						{
							Name:     "subtitle",
							Type:     "text",
							Repeated: false,
						},
						{
							Name:     "authors",
							Type:     "text",
							Repeated: false,
						},
						{
							Name:     "url",
							Type:     "text",
							Repeated: false,
						},
					},
				},
			},
		},
	}
}

func NewWikipediaTableSchema() pg.Schema {
	return pg.Schema{
		{
			Name:     "title",
			Type:     "text",
			Repeated: false,
		},
		{
			Name:     "id",
			Type:     "bigint",
			Repeated: false,
		},
		{
			Name:     "language",
			Type:     "text",
			Repeated: false,
		},
		{
			Name:     "wp_namespace",
			Type:     "bigint",
			Repeated: false,
		},
		{
			Name:     "is_redirect",
			Type:     "boolean",
			Repeated: false,
		},
		{
			Name:     "revision_id",
			Type:     "bigint",
			Repeated: false,
		},
		{
			Name:     "contributor_ip",
			Type:     "text",
			Repeated: false,
		},
		{
			Name:     "contributor_id",
			Type:     "bigint",
			Repeated: false,
		},
		{
			Name:     "contributor_username",
			Type:     "text",
			Repeated: false,
		},
		{
			Name:     "timestamp",
			Type:     "bigint",
			Repeated: false,
		},
		{
			Name:     "is_minor",
			Type:     "boolean",
			Repeated: false,
		},
		{
			Name:     "is_bot",
			Type:     "boolean",
			Repeated: false,
		},
		{
			Name:     "reversion_id",
			Type:     "bigint",
			Repeated: false,
		},
		{
			Name:     "comment",
			Type:     "text",
			Repeated: false,
		},
		{
			Name:     "num_characters",
			Type:     "bigint",
			Repeated: false,
		},
	}
}
