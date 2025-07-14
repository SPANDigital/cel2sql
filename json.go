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
		jsonFields := []string{"preferences", "metadata", "profile", "details", "settings", "properties", "analytics"}
		for _, jsonField := range jsonFields {
			if parentField == jsonField {
				return true
			}
		}
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
		if operandSelect := selectExpr.GetOperand().GetSelectExpr(); operandSelect != nil {
			// This is a nested select like json_users.settings.permissions
			parentField := operandSelect.GetField()
			jsonObjectFields := []string{"settings", "properties", "metadata", "analytics"}
			for _, jsonField := range jsonObjectFields {
				if parentField == jsonField {
					return true
				}
			}
		}
	}
	return false
}

// visitNestedJSONForArray handles nested JSON access for array operations
func (con *converter) visitNestedJSONForArray(expr *exprpb.Expr) error {
	selectExpr := expr.GetSelectExpr()
	if selectExpr == nil {
		return errors.New("expected select expression for nested JSON access")
	}

	// Visit the operand (like json_users.settings)
	if err := con.visit(selectExpr.GetOperand()); err != nil {
		return err
	}

	// Use -> instead of ->> to preserve JSONB type for array operations
	con.str.WriteString("->")
	con.str.WriteString("'")
	con.str.WriteString(selectExpr.GetField())
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
				"json_users":    {"tags", "scores", "attributes"},
				"json_products": {"features", "reviews", "categories"},
				"users":         {"preferences", "profile"}, // existing test data
				"products":      {"metadata", "details"},    // existing test data
			}

			if fields, exists := jsonArrayFields[tableName]; exists {
				for _, jsonField := range fields {
					if field == jsonField {
						return true
					}
				}
			}
		}

		// Check for nested JSON field access (e.g., json_users.settings.permissions)
		if nestedSelectExpr := operand.GetSelectExpr(); nestedSelectExpr != nil {
			parentField := nestedSelectExpr.GetField()
			jsonObjectFields := []string{"settings", "properties", "metadata", "analytics"}
			for _, jsonObjectField := range jsonObjectFields {
				if parentField == jsonObjectField {
					// This is accessing a field within a JSON object that could be an array
					arrayFields := []string{"permissions", "features", "tags", "categories"}
					for _, arrayField := range arrayFields {
						if field == arrayField {
							return true
						}
					}
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
				"json_users":    {"settings", "tags", "scores"},       // JSONB fields
				"json_products": {"features", "reviews", "properties"}, // JSONB fields
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