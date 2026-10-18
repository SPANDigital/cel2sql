// Copyright 2026 Observe, Inc. Licensed under Apache 2.0.
// List- and map-valued JSON variable support added by Observe, Inc.

package cel2sql

import (
	"fmt"

	"github.com/google/cel-go/common/operators"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"github.com/observeinc/cel2sql/v3/dialect"
)

// A variable declared with WithJSONVariables is a JSONB column, and a read of one
// of its values, v.k or v["k"], is written as v->>'k': text. That is the right
// reading when the value is a string, the map<string, string> shape the option
// was written for, and the wrong one when the value is a list or a map. Object
// tags, map<string, list<string>>, are the case in hand: v->>'env' is the text
// '["prod","dev"]', which ANY and UNNEST reject, and a test of the map's own keys
// ("env" in v, has(v.env), v.exists(k, ...)) fell through to SQL that does not
// parse at all.
//
// The CEL checker already knows each value's type, so the forms below are chosen
// by it rather than guessed from the expression's shape. For PostgreSQL, with
// GUARD(a, b) standing for CASE WHEN jsonb_typeof(a) = 'array' THEN b END:
//
//	"k" in v, has(v.k)              v ? 'k'
//	x in v.k, list(string)          GUARD(v->'k', (v->'k') ? (x))
//	v.k.exists(e, p), list(string)  GUARD(v->'k', EXISTS (SELECT 1 FROM jsonb_array_elements_text(v->'k') AS e WHERE p))
//	v.exists(k, p)                  EXISTS (SELECT 1 FROM jsonb_object_keys(v) AS k WHERE p)
//	size(v.k), list                 GUARD(v->'k', COALESCE(jsonb_array_length(v->'k'), 0))
//	size(v)                         (SELECT COUNT(*) FROM jsonb_object_keys(v))
//
// The guard is what keeps an absent key faithful. CEL errors on v.k when v has
// no key k, so an in-memory filter cannot say the row matches or that it does
// not; the SQL counterpart of that error is NULL, which a WHERE clause treats as
// not matching, NOT leaves NULL, and AND and OR absorb exactly as CEL absorbs an
// error (false && error is false, true || error is true). Without the guard an
// absent list would read as empty, and !("prod" in v.env) would match every row
// that has no env at all. The same guard keeps a value that is not an array from
// reaching jsonb_array_elements_text, which raises an error on one.
//
// Any other use of a list or map value -- comparing it, indexing into it, passing
// it to a function -- is refused, because no single SQL form is right for it and
// the text reading is always wrong. A scalar value is untouched: v.k == "x" is
// still v->>'k' = 'x' in every dialect.
//
// Only a dialect that implements dialect.JSONContainerWriter gets the forms
// above; the others refuse them with ErrUnsupportedDialectFeature.

// jsonVariableValue is one read of a JSON variable's value, v.k or v["k"].
type jsonVariableValue struct {
	variable *exprpb.Expr // the JSON variable's identifier
	key      string
	typ      *exprpb.Type // the value's CEL type
}

// isJSONVariableIdent reports whether expr is a bare JSON variable, v.
func (con *converter) isJSONVariableIdent(expr *exprpb.Expr) bool {
	ident := expr.GetIdentExpr()
	return ident != nil && con.isJSONVariable(ident.GetName())
}

// jsonVariableValueOf recognizes v.k and v["k"] on a JSON variable. ok is false for
// any other expression. A recognized read whose key the dialect cannot render
// returns ok with the error, so the caller reports it rather than falling back.
func (con *converter) jsonVariableValueOf(expr *exprpb.Expr) (jsonVariableValue, bool, error) {
	if sel := expr.GetSelectExpr(); sel != nil && !sel.GetTestOnly() && con.isJSONVariableIdent(sel.GetOperand()) {
		if err := con.dialect.ValidateFieldName(sel.GetField()); err != nil {
			return jsonVariableValue{}, true, fmt.Errorf("%w: %w", ErrInvalidFieldName, err)
		}
		return jsonVariableValue{variable: sel.GetOperand(), key: sel.GetField(), typ: con.getType(expr)}, true, nil
	}
	if call := expr.GetCallExpr(); call != nil && call.GetFunction() == operators.Index &&
		len(call.GetArgs()) == 2 && con.isJSONVariableIdent(call.GetArgs()[0]) {
		key, err := con.extractJSONVariableKey(call.GetArgs()[1])
		if err != nil {
			return jsonVariableValue{}, true, err
		}
		return jsonVariableValue{variable: call.GetArgs()[0], key: key, typ: con.getType(expr)}, true, nil
	}
	return jsonVariableValue{}, false, nil
}

func isContainerType(typ *exprpb.Type) bool {
	return isListType(typ) || isMapType(typ)
}

func isStringListType(typ *exprpb.Type) bool {
	return isListType(typ) && typ.GetListType().GetElemType().GetPrimitive() == exprpb.Type_STRING
}

// jsonContainerWriter returns the dialect's JSONContainerWriter, or the error that
// refuses a list- or map-valued read for a dialect without one.
func (con *converter) jsonContainerWriter(form string) (dialect.JSONContainerWriter, error) {
	w, ok := dialect.GetJSONContainerWriter(con.dialect)
	if !ok {
		return nil, fmt.Errorf("%w: %s on a JSON variable is not supported by %s dialect",
			ErrUnsupportedDialectFeature, form, con.dialect.Name())
	}
	return w, nil
}

// refuseJSONContainerValue is the error for a list- or map-valued read in a
// position none of the forms above covers.
func refuseJSONContainerValue(v jsonVariableValue) error {
	return newConversionErrorf(errMsgUnsupportedType,
		"JSON variable key %q holds a list or map; it can only be tested with in, has(), size(), exists() or all()", v.key)
}

// writeJSONVariableContainer writes v->'k', the value kept as JSON so it can be
// iterated or tested rather than read as text.
func (con *converter) writeJSONVariableContainer(v jsonVariableValue) error {
	return con.dialect.WriteJSONFieldAccess(&con.str, func() error { return con.visit(v.variable) }, v.key, false)
}

// visitJSONVariableIn writes `lhs in rhs` when rhs is a JSON variable or a list
// value of one. handled is false when rhs is neither, so the caller converts the
// expression as before.
func (con *converter) visitJSONVariableIn(lhs, rhs *exprpb.Expr) (handled bool, err error) {
	if con.isJSONVariableIdent(rhs) {
		// "k" in v tests the map's keys. The key is named in the SQL, as for v.k, so
		// it must be a literal.
		if _, err := con.jsonContainerWriter("in"); err != nil {
			return true, err
		}
		key, err := con.extractJSONVariableKey(lhs)
		if err != nil {
			return true, err
		}
		return true, con.dialect.WriteJSONExistence(&con.str, true, key, func() error { return con.visit(rhs) })
	}

	v, ok, err := con.jsonVariableValueOf(rhs)
	if !ok || err != nil {
		return ok, err
	}
	if !isContainerType(v.typ) {
		return false, nil
	}
	if !isStringListType(v.typ) {
		return true, refuseJSONContainerValue(v)
	}
	w, err := con.jsonContainerWriter("in")
	if err != nil {
		return true, err
	}
	writeArray := func() error { return con.writeJSONVariableContainer(v) }
	return true, w.WriteJSONArrayGuard(&con.str, writeArray, func() error {
		return w.WriteJSONArrayContainsText(&con.str, writeArray, func() error { return con.visit(lhs) })
	})
}

// visitJSONVariableComprehension writes a comprehension over a list value of a
// JSON variable inside the array guard, so the whole comprehension is NULL when
// the key is absent rather than true for all() and false for exists(). handled is
// false for any other range, which the caller converts as before.
func (con *converter) visitJSONVariableComprehension(iterRange *exprpb.Expr, writeComprehension func() error) (handled bool, err error) {
	v, ok, err := con.jsonVariableValueOf(iterRange)
	if !ok || err != nil {
		return ok, err
	}
	if !isStringListType(v.typ) {
		// Refused, or a map or scalar value, by writeJSONVariableComprehensionSource.
		return false, nil
	}
	w, err := con.jsonContainerWriter("a comprehension")
	if err != nil {
		return true, err
	}
	return true, w.WriteJSONArrayGuard(&con.str, func() error { return con.writeJSONVariableContainer(v) }, writeComprehension)
}

// writeJSONVariableComprehensionSource writes the FROM source of a comprehension
// over a JSON variable's keys or one of its list values. handled is false for any
// other range.
func (con *converter) writeJSONVariableComprehensionSource(iterRange *exprpb.Expr) (handled bool, err error) {
	if con.isJSONVariableIdent(iterRange) {
		w, err := con.jsonContainerWriter("a comprehension")
		if err != nil {
			return true, err
		}
		return true, w.WriteJSONObjectKeys(&con.str, func() error { return con.visit(iterRange) })
	}

	v, ok, err := con.jsonVariableValueOf(iterRange)
	if !ok || err != nil {
		return ok, err
	}
	if !isContainerType(v.typ) {
		return false, nil
	}
	if !isStringListType(v.typ) {
		return true, refuseJSONContainerValue(v)
	}
	if _, err := con.jsonContainerWriter("a comprehension"); err != nil {
		return true, err
	}
	return true, con.dialect.WriteJSONArrayElements(&con.str, true, true, func() error {
		return con.writeJSONVariableContainer(v)
	})
}

// visitJSONVariableHas writes has(v.k) on a JSON variable. handled is false when
// the operand is not one.
func (con *converter) visitJSONVariableHas(operand *exprpb.Expr, key string) (handled bool, err error) {
	if !con.isJSONVariableIdent(operand) {
		return false, nil
	}
	if _, err := con.jsonContainerWriter("has()"); err != nil {
		return true, err
	}
	return true, con.dialect.WriteJSONExistence(&con.str, true, key, func() error { return con.visit(operand) })
}

// visitJSONVariableSize writes size() of a JSON variable or of one of its list
// values. handled is false for any other argument.
func (con *converter) visitJSONVariableSize(arg *exprpb.Expr) (handled bool, err error) {
	if con.isJSONVariableIdent(arg) {
		w, err := con.jsonContainerWriter("size()")
		if err != nil {
			return true, err
		}
		con.str.WriteString("(SELECT COUNT(*) FROM ")
		if err := w.WriteJSONObjectKeys(&con.str, func() error { return con.visit(arg) }); err != nil {
			return true, err
		}
		con.str.WriteString(")")
		return true, nil
	}

	v, ok, err := con.jsonVariableValueOf(arg)
	if !ok || err != nil {
		return ok, err
	}
	if !isListType(v.typ) {
		// A string value keeps its LENGTH(v->>'k'); a map value has no form here.
		if isMapType(v.typ) {
			return true, refuseJSONContainerValue(v)
		}
		return false, nil
	}
	w, err := con.jsonContainerWriter("size()")
	if err != nil {
		return true, err
	}
	writeArray := func() error { return con.writeJSONVariableContainer(v) }
	return true, w.WriteJSONArrayGuard(&con.str, writeArray, func() error {
		return con.dialect.WriteJSONArrayLength(&con.str, writeArray)
	})
}

// checkJSONVariableScalarRead refuses a list- or map-valued read reached as an
// ordinary value, where it would otherwise be written as text. The forms above
// write their reads themselves and never visit them, so reaching here means
// nothing claimed the value.
func (con *converter) checkJSONVariableScalarRead(expr *exprpb.Expr) error {
	v, ok, err := con.jsonVariableValueOf(expr)
	if !ok || err != nil {
		return err
	}
	if isContainerType(v.typ) {
		return refuseJSONContainerValue(v)
	}
	return nil
}
