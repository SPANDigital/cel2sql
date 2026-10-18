// Copyright 2026 Observe, Inc. Licensed under Apache 2.0.

package dialect

import "strings"

// JSONContainerWriter is implemented by a dialect that can read a list or map
// held in a variable declared with WithJSONVariables.
//
// Such a variable is a JSON column, and reading one of its values as text is
// only right when the value is a scalar. A value that is itself a list or a map
// -- map<string, list<string>> object tags, say -- has to stay JSON so it can be
// tested for membership, iterated, or measured, and a map variable's keys have
// to be read as a set. The converter builds those reads from the dialect's
// existing JSON writers plus the methods below; a dialect that has not had them
// verified against a list- or map-valued variable does not implement this
// interface, and the converter refuses those expressions for it with
// ErrUnsupportedDialectFeature rather than emit SQL it cannot vouch for.
type JSONContainerWriter interface {
	// WriteJSONObjectKeys writes a set-returning call that yields each key of the
	// JSON object writeObject names as one text row, for use as a FROM source.
	WriteJSONObjectKeys(w *strings.Builder, writeObject func() error) error

	// WriteJSONArrayContainsText writes a test of whether the JSON array
	// writeArray names holds the string writeElem names as an element.
	WriteJSONArrayContainsText(w *strings.Builder, writeArray, writeElem func() error) error

	// WriteJSONArrayGuard writes writeBody where writeArray names a JSON array,
	// and NULL where it names anything else, including a key that is absent. CEL
	// errors on such a read, and NULL is the SQL counterpart: it propagates
	// through NOT, and AND and OR absorb it exactly as they absorb a CEL error.
	WriteJSONArrayGuard(w *strings.Builder, writeArray, writeBody func() error) error
}

// GetJSONContainerWriter returns the dialect's JSONContainerWriter, if it implements one.
func GetJSONContainerWriter(d Dialect) (JSONContainerWriter, bool) {
	w, ok := d.(JSONContainerWriter)
	return w, ok
}
