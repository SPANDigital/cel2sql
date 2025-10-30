package cel2sql

import (
	"context"
	"fmt"
	"log/slog"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// Constants for PostgreSQL JSON functions
const (
	jsonArrayElements      = "json_array_elements"
	jsonbArrayElements     = "jsonb_array_elements"
	jsonArrayElementsText  = "json_array_elements_text"
	jsonbArrayElementsText = "jsonb_array_elements_text"
)

// shouldUseJSONPath determines if we should use JSON path operators for field access
// This function checks if the operand represents a JSON/JSONB field using schema information
func (con *converter) shouldUseJSONPath(operand *exprpb.Expr, _ string) bool {
	// Check if the operand is a direct table.column access where column is JSON
	if selectExpr := operand.GetSelectExpr(); selectExpr != nil {
		// For obj.metadata, check if metadata is a JSON column in obj table
		if tableName, fieldName, ok := con.getTableAndFieldFromSelectChain(operand); ok {
			// Use schema information to determine if this field is JSON
			isJSON := con.isFieldJSON(tableName, fieldName)
			con.logger.LogAttrs(context.Background(), slog.LevelDebug,
				"JSON path detection",
				slog.String("table", tableName),
				slog.String("field", fieldName),
				slog.Bool("is_json", isJSON),
			)
			return isJSON
		}

		// Check if there's a JSON field somewhere in the operand chain
		if con.hasJSONFieldInChain(operand) {
			con.logger.Debug("JSON field detected in select chain")
			return true
		}
	}

	return false
}

// hasJSONFieldInChain checks if there's a JSON field anywhere in the select expression chain
func (con *converter) hasJSONFieldInChain(expr *exprpb.Expr) bool {
	if selectExpr := expr.GetSelectExpr(); selectExpr != nil {
		operand := selectExpr.GetOperand()

		// Check if this is a table.field access where field is JSON
		if identExpr := operand.GetIdentExpr(); identExpr != nil {
			tableName := identExpr.GetName()
			field := selectExpr.GetField()
			if con.isFieldJSON(tableName, field) {
				return true
			}
		}

		// Recursively check the operand
		return con.hasJSONFieldInChain(operand)
	}

	return false
}

// isJSONTextExtraction checks if an expression represents a JSON field extraction that returns text
// This is used to determine if we need numeric casting for comparisons
func (con *converter) isJSONTextExtraction(expr *exprpb.Expr) bool {
	// Check if this is a select expression that would use JSON path operators
	if selectExpr := expr.GetSelectExpr(); selectExpr != nil {
		operand := selectExpr.GetOperand()
		field := selectExpr.GetField()

		// If this would trigger JSON path generation, it's a text extraction
		return con.shouldUseJSONPath(operand, field)
	}

	return false
}

// needsNumericCasting checks if an identifier represents a numeric iteration variable from JSON
func (con *converter) needsNumericCasting(identName string) bool {
	// Common iteration variable names that come from numeric JSON arrays
	numericIterationVars := []string{"score", "value", "num", "amount", "count", "level"}

	for _, numericVar := range numericIterationVars {
		if identName == numericVar {
			return true
		}
	}

	return false
}

// isNumericJSONField checks if a JSON field name typically contains numeric values
func (con *converter) isNumericJSONField(fieldName string) bool {
	numericFields := []string{"level", "score", "value", "count", "amount", "price", "rating", "age", "size", "capacity", "megapixels", "cores", "threads", "ram", "storage", "vram", "weight", "frequency", "helpful"}

	for _, numericField := range numericFields {
		if fieldName == numericField {
			return true
		}
	}

	return false
}

// isNestedJSONAccess checks if this is nested JSON field access like settings.permissions
func (con *converter) isNestedJSONAccess(expr *exprpb.Expr) bool {
	if selectExpr := expr.GetSelectExpr(); selectExpr != nil {
		// Check if there's a JSON field somewhere in the chain
		return con.hasJSONFieldInChain(selectExpr.GetOperand())
	}
	return false
}

// visitNestedJSONForArray handles nested JSON access for array operations
func (con *converter) visitNestedJSONForArray(expr *exprpb.Expr) error {
	selectExpr := expr.GetSelectExpr()
	if selectExpr == nil {
		return fmt.Errorf("%w: expected select expression for nested JSON access", ErrInvalidJSONPath)
	}

	// For array operations, we need to use -> throughout to preserve JSON type
	// Use a specialized builder that doesn't convert to text
	return con.buildJSONPathForArray(expr)
}

// buildJSONPathForArray constructs a JSON path for array operations, using -> throughout
func (con *converter) buildJSONPathForArray(expr *exprpb.Expr) error {
	selectExpr := expr.GetSelectExpr()
	if selectExpr == nil {
		return fmt.Errorf("%w: expected select expression for JSON array path", ErrInvalidJSONPath)
	}

	operand := selectExpr.GetOperand()
	field := selectExpr.GetField()

	// Check if the operand is also a select expression (nested access)
	if operandSelect := operand.GetSelectExpr(); operandSelect != nil {
		// This is nested access - recursively build the path for the operand
		if con.hasJSONFieldInChain(operand) {
			if err := con.buildJSONPathForArray(operand); err != nil {
				return err
			}
			// Add intermediate JSON path operator (always -> for arrays)
			con.str.WriteString("->'")
			con.str.WriteString(escapeJSONFieldName(field))
			con.str.WriteString("'")
			return nil
		}
	}

	// Check if this is the base table.jsonfield access
	if operandIdent := operand.GetIdentExpr(); operandIdent != nil {
		// This is table.jsonfield - use normal table.field syntax for the base
		tableName := operandIdent.GetName()
		con.str.WriteString(tableName)
		con.str.WriteString(".")
		con.str.WriteString(field)
		return nil
	}

	// For other cases, visit the operand and add JSON operator
	if err := con.visit(operand); err != nil {
		return err
	}
	con.str.WriteString("->'")
	con.str.WriteString(escapeJSONFieldName(field))
	con.str.WriteString("'")
	return nil
}

// isJSONObjectFieldAccess determines if this is a JSON object field access in comprehensions
func (con *converter) isJSONObjectFieldAccess(expr *exprpb.Expr) bool {
	if selectExpr := expr.GetSelectExpr(); selectExpr != nil {
		operand := selectExpr.GetOperand()

		// Check if the operand is an identifier that could be a comprehension variable
		if identExpr := operand.GetIdentExpr(); identExpr != nil {
			// Common comprehension variable names that access JSON objects
			jsonObjectVars := []string{"attr", "item", "element", "obj", "feature", "review"}
			identName := identExpr.GetName()

			for _, jsonVar := range jsonObjectVars {
				if identName == jsonVar {
					return true
				}
			}
		}
	}
	return false
}

// isJSONArrayField determines if the expression refers to a JSON/JSONB array field
func (con *converter) isJSONArrayField(expr *exprpb.Expr) bool {
	// Check if this is a field selection on a JSON field
	if selectExpr := expr.GetSelectExpr(); selectExpr != nil {
		// Get the operand (the table/object being accessed)
		operand := selectExpr.GetOperand()
		field := selectExpr.GetField()

		// Check if the operand is an identifier (table name)
		if identExpr := operand.GetIdentExpr(); identExpr != nil {
			tableName := identExpr.GetName()

			// Use schema information to check if this is an array field
			return con.isFieldArray(tableName, field) && con.isFieldJSON(tableName, field)
		}

		// For nested JSON access, we can't determine array status without schema
		// Return false - schema must be provided for JSON array detection
	}

	return false
}

// isJSONBField determines if the expression refers to a JSONB field (vs JSON field)
func (con *converter) isJSONBField(expr *exprpb.Expr) bool {
	// Check if this is a field selection on a JSONB field
	if selectExpr := expr.GetSelectExpr(); selectExpr != nil {
		operand := selectExpr.GetOperand()
		field := selectExpr.GetField()

		// Check if the operand is an identifier (table name)
		if identExpr := operand.GetIdentExpr(); identExpr != nil {
			tableName := identExpr.GetName()
			// Use schema information to check if this is a JSONB field
			return con.isFieldJSONB(tableName, field)
		}

		// For nested access, check if the parent is JSONB using schema
		if nestedSelectExpr := operand.GetSelectExpr(); nestedSelectExpr != nil {
			if parentIdentExpr := nestedSelectExpr.GetOperand().GetIdentExpr(); parentIdentExpr != nil {
				tableName := parentIdentExpr.GetName()
				parentField := nestedSelectExpr.GetField()
				return con.isFieldJSONB(tableName, parentField)
			}
		}
	}
	return false
}

// getJSONArrayFunction returns the appropriate PostgreSQL function for JSON array operations
func (con *converter) getJSONArrayFunction(expr *exprpb.Expr) string {
	// Determine if this is JSON or JSONB based on the field
	isJSONB := con.isJSONBField(expr)

	if selectExpr := expr.GetSelectExpr(); selectExpr != nil {
		field := selectExpr.GetField()
		operand := selectExpr.GetOperand()

		// Use schema information to determine element type
		if identExpr := operand.GetIdentExpr(); identExpr != nil {
			tableName := identExpr.GetName()
			elementType := con.getFieldElementType(tableName, field)

			if elementType != "" {
				// Determine if this is a simple type (text, numbers) or complex (json/jsonb)
				simpleTypes := map[string]bool{
					"text":    true,
					"varchar": true,
					"integer": true,
					"bigint":  true,
					"numeric": true,
					"decimal": true,
					"real":    true,
					"double":  true,
				}

				if simpleTypes[elementType] {
					// Simple types: use text extraction
					if isJSONB {
						return jsonbArrayElementsText
					}
					return jsonArrayElementsText
				}

				// Complex types (json, jsonb, composite): use object extraction
				if isJSONB {
					return jsonbArrayElements
				}
				return jsonArrayElements
			}
		}
	}

	// Default based on field type when schema not available
	if isJSONB {
		return jsonbArrayElements
	}
	return jsonArrayElements
}

// buildJSONPath constructs the full JSON path for nested field access
func (con *converter) buildJSONPath(expr *exprpb.Expr) error {
	con.logger.Debug("building JSON path for nested access")
	return con.buildJSONPathInternal(expr, true)
}

// buildJSONPathInternal is the internal implementation that tracks if this is the final field
func (con *converter) buildJSONPathInternal(expr *exprpb.Expr, isFinalField bool) error {
	selectExpr := expr.GetSelectExpr()
	if selectExpr == nil {
		return fmt.Errorf("%w: expected select expression for JSON path", ErrInvalidJSONPath)
	}

	operand := selectExpr.GetOperand()
	field := selectExpr.GetField()

	// Check if the operand is also a select expression (nested access)
	if operandSelect := operand.GetSelectExpr(); operandSelect != nil {
		// Check if this is a direct table.column access (e.g., obj.metadata)
		// If so, we should NOT apply JSON operators to this level
		if tableName, columnName, ok := con.getTableAndFieldFromSelectChain(operand); ok {
			// This is table.column where column is JSON/JSONB
			// Render as table.column without JSON operators
			con.str.WriteString(tableName)
			con.str.WriteString(".")
			con.str.WriteString(columnName)
			// Now add JSON operator for the current field
			if isFinalField {
				con.str.WriteString("->>'") // Final field: extract as text
			} else {
				con.str.WriteString("->'") // Intermediate field: keep as JSON
			}
			con.str.WriteString(escapeJSONFieldName(field))
			con.str.WriteString("'")
			return nil
		}

		// This is deeper nesting like table.jsonfield.subfield.finalfield
		// We need to determine if the operand is JSON-related
		if con.shouldUseJSONPath(operandSelect.GetOperand(), operandSelect.GetField()) {
			// Recursively build the path for the operand (not final since we have more fields)
			if err := con.buildJSONPathInternal(operand, false); err != nil {
				return err
			}
			// Add appropriate JSON path operator based on whether this is the final field
			if isFinalField {
				con.str.WriteString("->>'") // Final field: extract as text
			} else {
				con.str.WriteString("->'") // Intermediate field: keep as JSON
			}
			con.str.WriteString(escapeJSONFieldName(field))
			con.str.WriteString("'")
			return nil
		}
	}

	// Visit the base operand (like table.jsonfield)
	if err := con.visit(operand); err != nil {
		return err
	}

	// Add the appropriate JSON path operator based on whether this is the final field
	operator := "->>"
	if !isFinalField {
		operator = "->"
	}

	con.logger.LogAttrs(context.Background(), slog.LevelDebug,
		"JSON path operator selection",
		slog.String("field", field),
		slog.String("operator", operator),
		slog.Bool("is_final", isFinalField),
	)

	con.str.WriteString(operator)
	con.str.WriteString("'")
	con.str.WriteString(escapeJSONFieldName(field))
	con.str.WriteString("'")
	return nil
}
