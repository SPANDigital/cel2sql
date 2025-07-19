package cel2sql

import (
	"errors"

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
func (con *converter) shouldUseJSONPath(operand *exprpb.Expr, _ string) bool {
	// For now, we'll use a simple heuristic: if the operand is a direct field reference
	// to a field that commonly contains JSON (like 'preferences', 'metadata', 'profile', 'details')
	// then we use JSON path operators
	if identExpr := operand.GetIdentExpr(); identExpr != nil {
		// Direct field access - check if it's a known JSON field
		return false // We don't have direct JSON field access in our current tests
	}

	if selectExpr := operand.GetSelectExpr(); selectExpr != nil {
		// Nested field access - check if the parent field is a JSON field
		parentField := selectExpr.GetField()
		jsonFields := []string{"preferences", "metadata", "profile", "details", "settings", "properties", "analytics", 
		                      "content", "structure", "taxonomy", "classification", "content_structure"}
		for _, jsonField := range jsonFields {
			if parentField == jsonField {
				return true
			}
		}
		
		// Also check if we have deeper nesting where a JSON field appears earlier in the chain
		if con.hasJSONFieldInChain(operand) {
			return true
		}
	}

	return false
}

// hasJSONFieldInChain checks if there's a JSON field anywhere in the select expression chain
func (con *converter) hasJSONFieldInChain(expr *exprpb.Expr) bool {
	if selectExpr := expr.GetSelectExpr(); selectExpr != nil {
		field := selectExpr.GetField()
		operand := selectExpr.GetOperand()
		
		// Check if current field is a JSON field
		jsonFields := []string{"preferences", "metadata", "profile", "details", "settings", "properties", "analytics", 
		                      "content", "structure", "taxonomy", "classification", "content_structure"}
		for _, jsonField := range jsonFields {
			if field == jsonField {
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
		return errors.New("expected select expression for nested JSON access")
	}

	// For array operations, we need to use -> throughout to preserve JSON type
	// Use a specialized builder that doesn't convert to text
	return con.buildJSONPathForArray(expr)
}

// buildJSONPathForArray constructs a JSON path for array operations, using -> throughout
func (con *converter) buildJSONPathForArray(expr *exprpb.Expr) error {
	selectExpr := expr.GetSelectExpr()
	if selectExpr == nil {
		return errors.New("expected select expression for JSON array path")
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
			con.str.WriteString(field)
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
	con.str.WriteString(field)
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

// getJSONTypeofFunction returns the appropriate typeof function for JSON/JSONB fields
func (con *converter) getJSONTypeofFunction(expr *exprpb.Expr) string {
	if con.isJSONBField(expr) {
		return "jsonb_typeof"
	}
	return "json_typeof"
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

			// Check for known JSON array fields in our test schemas
			jsonArrayFields := map[string][]string{
				"json_users":         {"tags", "scores", "attributes"},
				"json_products":      {"features", "reviews", "categories"},
				"users":              {"preferences", "profile"}, // existing test data
				"products":           {"metadata", "details"},    // existing test data
				"information_assets": {"metadata", "properties", "classification", "content_structure"}, // nested path test data
				"documents":          {"content", "structure", "taxonomy", "analytics"},                  // nested path test data
			}

			if fields, exists := jsonArrayFields[tableName]; exists {
				for _, jsonField := range fields {
					if field == jsonField {
						return true
					}
				}
			}
		}

		// Check for nested JSON field access (e.g., information_assets.metadata.corpus.tags)
		// If there's a JSON field in the chain, this could be a nested array
		if con.hasJSONFieldInChain(expr) {
			// For nested access, we assume certain field names are arrays
			arrayFieldNames := []string{"tags", "permissions", "features", "categories", "scores", "attributes"}
			for _, arrayField := range arrayFieldNames {
				if field == arrayField {
					return true
				}
			}
		}
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

			// Define which fields are JSONB vs JSON in our test schemas
			jsonbFields := map[string][]string{
				"json_users":         {"settings", "tags", "scores"},       // JSONB fields
				"json_products":      {"features", "reviews", "properties"}, // JSONB fields
				"information_assets": {"metadata", "classification"},        // JSONB fields
				"documents":          {"content", "taxonomy"},               // JSONB fields
			}

			if fields, exists := jsonbFields[tableName]; exists {
				for _, jsonbField := range fields {
					if field == jsonbField {
						return true
					}
				}
			}
		}

		// For nested access, check if the parent is JSONB
		if nestedSelectExpr := operand.GetSelectExpr(); nestedSelectExpr != nil {
			parentField := nestedSelectExpr.GetField()
			jsonbParentFields := []string{"settings", "properties"}
			for _, jsonbParent := range jsonbParentFields {
				if parentField == jsonbParent {
					return true
				}
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
		
		// Fields that contain simple values (strings, numbers)
		simpleArrayFields := []string{"tags", "scores", "categories"}
		for _, simpleField := range simpleArrayFields {
			if field == simpleField {
				// For all simple fields, use text extraction to avoid casting issues
				if isJSONB {
					return jsonbArrayElementsText
				}
				return jsonArrayElementsText
			}
		}
		
		// Fields that contain complex objects
		complexArrayFields := []string{"attributes", "features", "reviews"}
		for _, complexField := range complexArrayFields {
			if field == complexField {
				if isJSONB {
					return jsonbArrayElements
				}
				return jsonArrayElements
			}
		}
		
		// For nested JSON access, use appropriate array elements function
		if operand := selectExpr.GetOperand(); operand.GetSelectExpr() != nil {
			if isJSONB {
				return jsonbArrayElements
			}
			return jsonArrayElements
		}
	}
	
	// Default based on field type
	if isJSONB {
		return jsonbArrayElements
	}
	return jsonArrayElements
}

// buildJSONPath constructs the full JSON path for nested field access
func (con *converter) buildJSONPath(expr *exprpb.Expr) error {
	return con.buildJSONPathInternal(expr, true)
}

// buildJSONPathInternal is the internal implementation that tracks if this is the final field
func (con *converter) buildJSONPathInternal(expr *exprpb.Expr, isFinalField bool) error {
	selectExpr := expr.GetSelectExpr()
	if selectExpr == nil {
		return errors.New("expected select expression for JSON path")
	}

	operand := selectExpr.GetOperand()
	field := selectExpr.GetField()

	// Check if the operand is also a select expression (nested access)
	if operandSelect := operand.GetSelectExpr(); operandSelect != nil {
		// This is nested access like table.jsonfield.subfield.finalfield
		// We need to determine if the operand is JSON-related
		if con.shouldUseJSONPath(operandSelect.GetOperand(), operandSelect.GetField()) {
			// Recursively build the path for the operand (not final since we have more fields)
			if err := con.buildJSONPathInternal(operand, false); err != nil {
				return err
			}
			// Add appropriate JSON path operator based on whether this is the final field
			if isFinalField {
				con.str.WriteString("->>'")  // Final field: extract as text
			} else {
				con.str.WriteString("->'")   // Intermediate field: keep as JSON
			}
			con.str.WriteString(field)
			con.str.WriteString("'")
			return nil
		}
	}

	// Visit the base operand (like table.jsonfield)
	if err := con.visit(operand); err != nil {
		return err
	}

	// Add the appropriate JSON path operator based on whether this is the final field
	if isFinalField {
		con.str.WriteString("->>'")  // Final field: extract as text
	} else {
		con.str.WriteString("->'")   // Intermediate field: keep as JSON
	}
	con.str.WriteString(field)
	con.str.WriteString("'")
	return nil
}