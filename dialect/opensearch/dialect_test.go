package opensearch_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/spandigital/cel2sql/v3/dialect"
	"github.com/spandigital/cel2sql/v3/dialect/opensearch"
)

func noopWrite() func() error {
	return func() error { return nil }
}

func writeTo(w *strings.Builder, s string) func() error {
	return func() error { w.WriteString(s); return nil }
}

func TestRegistered(t *testing.T) {
	d, err := dialect.Get(dialect.OpenSearch)
	require.NoError(t, err)
	assert.Equal(t, dialect.OpenSearch, d.Name())
}

func TestWriteArrayMembership(t *testing.T) {
	d := opensearch.New()

	tests := []struct {
		name    string
		array   string
		want    string
		wantErr bool
	}{
		{name: "literal list", array: "('a', 'b', 'c')", want: "name IN ('a', 'b', 'c')"},
		{name: "literal with parens and quotes inside strings", array: "('a)', 'it''s (x')", want: "name IN ('a)', 'it''s (x')"},
		{name: "empty list", array: "(NULL)", want: "name IN (NULL)"},
		{name: "list variable", array: "string_list", wantErr: true},
		{name: "concatenated literals", array: "(1) || (2, 3)", wantErr: true},
		{name: "nested field", array: "p.tags", wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var w strings.Builder
			err := d.WriteArrayMembership(&w, writeTo(&w, "name"), writeTo(&w, tc.array))
			if tc.wantErr {
				require.ErrorIs(t, err, dialect.ErrUnsupportedFeature)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, w.String())
		})
	}
}

func TestWriteExtract(t *testing.T) {
	d := opensearch.New()

	tests := []struct {
		part string
		want string
	}{
		{"DOW", "(DAYOFWEEK(ts) + 5) % 7"},
		{"DOY", "DAYOFYEAR(ts)"},
		{"MILLISECONDS", "FLOOR(MICROSECOND(ts) / 1000)"},
		{"YEAR", "EXTRACT(YEAR FROM ts)"},
		{"HOUR", "EXTRACT(HOUR FROM ts)"},
	}
	for _, tc := range tests {
		t.Run(tc.part, func(t *testing.T) {
			var w strings.Builder
			require.NoError(t, d.WriteExtract(&w, tc.part, writeTo(&w, "ts"), nil))
			assert.Equal(t, tc.want, w.String())
		})
	}

	t.Run("time zone is rejected", func(t *testing.T) {
		var w strings.Builder
		err := d.WriteExtract(&w, "HOUR", writeTo(&w, "ts"), writeTo(&w, "'Asia/Tokyo'"))
		require.ErrorIs(t, err, dialect.ErrUnsupportedFeature)
	})
}

func TestWriteTimestampArithmetic(t *testing.T) {
	d := opensearch.New()

	var w strings.Builder
	require.NoError(t, d.WriteTimestampArithmetic(&w, "+", writeTo(&w, "ts"), writeTo(&w, "INTERVAL 1 HOUR")))
	assert.Equal(t, "DATE_ADD(ts, INTERVAL 1 HOUR)", w.String())

	w.Reset()
	require.NoError(t, d.WriteTimestampArithmetic(&w, "-", writeTo(&w, "ts"), writeTo(&w, "INTERVAL 1 HOUR")))
	assert.Equal(t, "DATE_SUB(ts, INTERVAL 1 HOUR)", w.String())
}

func TestWriteDuration(t *testing.T) {
	d := opensearch.New()

	var w strings.Builder
	d.WriteDuration(&w, 10, "SECOND")
	assert.Equal(t, "INTERVAL 10 SECOND", w.String())

	w.Reset()
	d.WriteDuration(&w, 250, "MILLISECOND")
	assert.Equal(t, "INTERVAL 250000 MICROSECOND", w.String())
}

func TestComprehensionForms(t *testing.T) {
	d := opensearch.New()

	var w strings.Builder
	require.NoError(t, d.WriteComprehensionExists(&w, func() error {
		if err := d.WriteUnnest(&w, writeTo(&w, "p.variants")); err != nil {
			return err
		}
		w.WriteString(" AS v WHERE v.sku = 'A'")
		return nil
	}))
	assert.Equal(t, "EXISTS (SELECT * FROM p.variants AS v WHERE v.sku = 'A')", w.String())

	w.Reset()
	require.ErrorIs(t, d.WriteComprehensionNotExists(&w, writeTo(&w, "p.variants AS v WHERE NOT (v.qty > 0)")), dialect.ErrUnsupportedFeature)

	w.Reset()
	require.NoError(t, d.WriteIterVarRef(&w, "v"))
	assert.Equal(t, "v", w.String())
}

func TestJSONAccessIsDotPath(t *testing.T) {
	d := opensearch.New()

	var w strings.Builder
	require.NoError(t, d.WriteJSONFieldAccess(&w, writeTo(&w, "p.metadata"), "brand", true))
	assert.Equal(t, "p.metadata.brand", w.String())

	w.Reset()
	require.NoError(t, d.WriteJSONExistence(&w, false, "brand", writeTo(&w, "p.metadata")))
	assert.Equal(t, "p.metadata.brand IS NOT NULL", w.String())

	w.Reset()
	require.NoError(t, d.WriteJSONExtractPath(&w, []string{"specs", "color"}, writeTo(&w, "p.metadata")))
	assert.Equal(t, "p.metadata.specs.color IS NOT NULL", w.String())

	w.Reset()
	err := d.WriteJSONFieldAccess(&w, writeTo(&w, "p.metadata"), "bad'name", true)
	require.Error(t, err)
}

func TestUnsupportedFeatures(t *testing.T) {
	d := opensearch.New()
	var w strings.Builder
	noop := noopWrite()

	assert.ErrorIs(t, d.WriteBytesLiteral(&w, []byte{1}), dialect.ErrUnsupportedFeature)
	assert.ErrorIs(t, d.WriteArrayLength(&w, 1, noop), dialect.ErrUnsupportedFeature)
	assert.ErrorIs(t, d.WriteJSONArrayLength(&w, noop), dialect.ErrUnsupportedFeature)
	assert.ErrorIs(t, d.WriteListIndex(&w, noop, noop), dialect.ErrUnsupportedFeature)
	assert.ErrorIs(t, d.WriteListIndexConst(&w, noop, 0), dialect.ErrUnsupportedFeature)
	assert.ErrorIs(t, d.WriteSplit(&w, noop, noop), dialect.ErrUnsupportedFeature)
	assert.ErrorIs(t, d.WriteSplitWithLimit(&w, noop, noop, 2), dialect.ErrUnsupportedFeature)
	assert.ErrorIs(t, d.WriteJoin(&w, noop, noop), dialect.ErrUnsupportedFeature)
	assert.ErrorIs(t, d.WriteJSONArrayMembership(&w, "", noop, noop), dialect.ErrUnsupportedFeature)
	assert.ErrorIs(t, d.WriteNestedJSONArrayMembership(&w, noop, noop), dialect.ErrUnsupportedFeature)
}

func TestCapabilities(t *testing.T) {
	d := opensearch.New()
	assert.True(t, d.SupportsRegex())
	assert.False(t, d.SupportsNativeArrays())
	assert.False(t, d.SupportsJSONB())
	assert.False(t, d.SupportsIndexAnalysis())
	assert.False(t, d.SupportsSubqueryExpressions())
	assert.Equal(t, 0, d.MaxIdentifierLength())
}

func TestRegexPassThrough(t *testing.T) {
	d := opensearch.New()

	p, ci, err := d.ConvertRegex(`(?i)^user_\d+$`)
	require.NoError(t, err)
	assert.Equal(t, `(?i)^user_\d+$`, p)
	assert.False(t, ci)

	_, _, err = d.ConvertRegex(`(a+)+b`)
	require.Error(t, err)

	var w strings.Builder
	require.NoError(t, d.WriteRegexMatch(&w, writeTo(&w, "name"), `it's`, false))
	assert.Equal(t, "(name REGEXP 'it''s') = 1", w.String())
}

func TestWriteBooleanFieldPredicate(t *testing.T) {
	d := opensearch.New()
	var w strings.Builder
	require.NoError(t, d.WriteBooleanFieldPredicate(&w, writeTo(&w, "p.active")))
	assert.Equal(t, "p.active = true", w.String())
}

func TestValidateFieldName(t *testing.T) {
	d := opensearch.New()
	require.NoError(t, d.ValidateFieldName("host"))
	require.NoError(t, d.ValidateFieldName("message"))
	require.Error(t, d.ValidateFieldName("select"))
	require.Error(t, d.ValidateFieldName("timestamp"))
	require.Error(t, d.ValidateFieldName("bad-name"))
	require.Error(t, d.ValidateFieldName(""))
}
