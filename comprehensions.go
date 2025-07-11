package cel2sql

import (
	"fmt"

	"github.com/google/cel-go/common/operators"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// ComprehensionType represents the type of comprehension operation
type ComprehensionType int

const (
	ComprehensionAll ComprehensionType = iota
	ComprehensionExists
	ComprehensionExistsOne
	ComprehensionMap
	ComprehensionFilter
	ComprehensionTransformList
	ComprehensionTransformMap
	ComprehensionTransformMapEntry
	ComprehensionUnknown
)

// String returns a string representation of the comprehension type
func (ct ComprehensionType) String() string {
	switch ct {
	case ComprehensionAll:
		return "all"
	case ComprehensionExists:
		return "exists"
	case ComprehensionExistsOne:
		return "exists_one"
	case ComprehensionMap:
		return "map"
	case ComprehensionFilter:
		return "filter"
	case ComprehensionTransformList:
		return "transformList"
	case ComprehensionTransformMap:
		return "transformMap"
	case ComprehensionTransformMapEntry:
		return "transformMapEntry"
	default:
		return "unknown"
	}
}

// ComprehensionInfo holds metadata about a comprehension operation
type ComprehensionInfo struct {
	Type      ComprehensionType
	IterVar   string
	IndexVar  string // for two-variable comprehensions
	AccuVar   string
	HasFilter bool
	IsTwoVar  bool
	Transform *exprpb.Expr // transform expression for map/transformList
	Predicate *exprpb.Expr // predicate expression for filtering
	Filter    *exprpb.Expr // filter expression for map with filter
}

// identifyComprehension analyzes the AST to determine if an expression
// is a comprehension macro and extracts its metadata
func (con *converter) identifyComprehension(expr *exprpb.Expr) (*ComprehensionInfo, error) {
	comprehension := expr.GetComprehensionExpr()
	if comprehension == nil {
		return nil, fmt.Errorf("expression is not a comprehension")
	}

	// Analyze the comprehension structure to determine its type
	return con.analyzeComprehensionPattern(comprehension)
}

// analyzeComprehensionPattern examines the comprehension AST structure to identify
// which CEL macro it represents by pattern matching the characteristic expressions
func (con *converter) analyzeComprehensionPattern(comp *exprpb.Expr_Comprehension) (*ComprehensionInfo, error) {
	info := &ComprehensionInfo{
		IterVar: comp.GetIterVar(),
		AccuVar: comp.GetAccuVar(),
	}

	// Check for two-variable comprehensions
	if comp.GetIterVar2() != "" {
		info.IsTwoVar = true
		info.IndexVar = comp.GetIterVar2()
	}

	// Check accumulator initialization to determine type
	accuInit := comp.GetAccuInit()

	// All: accuInit = true, step = accu && predicate, result = accu
	if con.isBoolTrue(accuInit) {
		if con.isLogicalAndStep(comp.GetLoopStep(), comp.GetAccuVar()) {
			info.Type = ComprehensionAll
			info.Predicate = con.extractPredicateFromAndStep(comp.GetLoopStep(), comp.GetAccuVar())
			return info, nil
		}
	}

	// Exists: accuInit = false, step = accu || predicate, result = accu
	if con.isBoolFalse(accuInit) {
		if con.isLogicalOrStep(comp.GetLoopStep(), comp.GetAccuVar()) {
			info.Type = ComprehensionExists
			info.Predicate = con.extractPredicateFromOrStep(comp.GetLoopStep(), comp.GetAccuVar())
			return info, nil
		}
	}

	// ExistsOne: accuInit = 0, step = conditional(predicate, accu + 1, accu), result = accu == 1
	if con.isIntZero(accuInit) {
		if con.isConditionalCountStep(comp.GetLoopStep(), comp.GetAccuVar()) && con.isEqualsOneResult(comp.GetResult(), comp.GetAccuVar()) {
			info.Type = ComprehensionExistsOne
			info.Predicate = con.extractPredicateFromConditionalStep(comp.GetLoopStep())
			return info, nil
		}
	}

	// Map: accuInit = [], step = accu + [transform], result = accu
	if con.isEmptyList(accuInit) {
		if con.isListAppendStep(comp.GetLoopStep(), comp.GetAccuVar()) {
			info.Type = ComprehensionMap
			info.Transform = con.extractTransformFromAppendStep(comp.GetLoopStep(), comp.GetAccuVar())
			return info, nil
		}
		// Map with filter: step = conditional(filter, accu + [transform], accu)
		if con.isConditionalAppendStep(comp.GetLoopStep(), comp.GetAccuVar()) {
			info.Type = ComprehensionMap
			info.HasFilter = true
			filter, transform := con.extractFilterAndTransformFromConditionalStep(comp.GetLoopStep(), comp.GetAccuVar())
			info.Filter = filter
			info.Transform = transform
			return info, nil
		}
	}

	// Filter: accuInit = [], step = conditional(predicate, accu + [iterVar], accu), result = accu
	if con.isEmptyList(accuInit) {
		if con.isConditionalFilterStep(comp.GetLoopStep(), comp.GetAccuVar(), comp.GetIterVar()) {
			info.Type = ComprehensionFilter
			info.Predicate = con.extractPredicateFromConditionalStep(comp.GetLoopStep())
			return info, nil
		}
	}

	// If we can't identify the pattern, mark as unknown for now
	info.Type = ComprehensionUnknown
	return info, fmt.Errorf("unrecognized comprehension pattern for %s", comp.String())
}

// Helper functions to identify patterns in comprehension expressions

func (con *converter) isBoolTrue(expr *exprpb.Expr) bool {
	if constant := expr.GetConstExpr(); constant != nil {
		if boolVal := constant.GetBoolValue(); boolVal {
			return true
		}
	}
	return false
}

func (con *converter) isBoolFalse(expr *exprpb.Expr) bool {
	if constant := expr.GetConstExpr(); constant != nil {
		if boolVal := constant.GetBoolValue(); !boolVal {
			return true
		}
	}
	return false
}

func (con *converter) isIntZero(expr *exprpb.Expr) bool {
	if constant := expr.GetConstExpr(); constant != nil {
		if intVal := constant.GetInt64Value(); intVal == 0 {
			return true
		}
	}
	return false
}

func (con *converter) isEmptyList(expr *exprpb.Expr) bool {
	if listExpr := expr.GetListExpr(); listExpr != nil {
		return len(listExpr.Elements) == 0
	}
	return false
}

func (con *converter) isLogicalAndStep(step *exprpb.Expr, accuVar string) bool {
	if call := step.GetCallExpr(); call != nil {
		return call.Function == operators.LogicalAnd && con.hasAccuReference(call.Args, accuVar)
	}
	return false
}

func (con *converter) isLogicalOrStep(step *exprpb.Expr, accuVar string) bool {
	if call := step.GetCallExpr(); call != nil {
		return call.Function == operators.LogicalOr && con.hasAccuReference(call.Args, accuVar)
	}
	return false
}

func (con *converter) isListAppendStep(step *exprpb.Expr, accuVar string) bool {
	if call := step.GetCallExpr(); call != nil {
		return call.Function == operators.Add && con.hasAccuReference(call.Args, accuVar) && con.hasListConstruction(call.Args)
	}
	return false
}

func (con *converter) isConditionalCountStep(step *exprpb.Expr, accuVar string) bool {
	if call := step.GetCallExpr(); call != nil {
		return call.Function == operators.Conditional && len(call.Args) == 3
	}
	return false
}

func (con *converter) isConditionalAppendStep(step *exprpb.Expr, accuVar string) bool {
	if call := step.GetCallExpr(); call != nil {
		return call.Function == operators.Conditional && len(call.Args) == 3
	}
	return false
}

func (con *converter) isConditionalFilterStep(step *exprpb.Expr, accuVar, iterVar string) bool {
	if call := step.GetCallExpr(); call != nil {
		return call.Function == operators.Conditional && len(call.Args) == 3
	}
	return false
}

func (con *converter) isEqualsOneResult(result *exprpb.Expr, accuVar string) bool {
	if call := result.GetCallExpr(); call != nil {
		return call.Function == operators.Equals
	}
	return false
}

func (con *converter) hasAccuReference(args []*exprpb.Expr, accuVar string) bool {
	for _, arg := range args {
		if ident := arg.GetIdentExpr(); ident != nil && ident.Name == accuVar {
			return true
		}
	}
	return false
}

func (con *converter) hasListConstruction(args []*exprpb.Expr) bool {
	for _, arg := range args {
		if arg.GetListExpr() != nil {
			return true
		}
	}
	return false
}

// Extraction functions to get predicate and transform expressions

func (con *converter) extractPredicateFromAndStep(step *exprpb.Expr, accuVar string) *exprpb.Expr {
	if call := step.GetCallExpr(); call != nil && len(call.Args) == 2 {
		// In AND step: accu && predicate, return the non-accu argument
		for _, arg := range call.Args {
			if ident := arg.GetIdentExpr(); ident == nil || ident.Name != accuVar {
				return arg
			}
		}
	}
	return nil
}

func (con *converter) extractPredicateFromOrStep(step *exprpb.Expr, accuVar string) *exprpb.Expr {
	if call := step.GetCallExpr(); call != nil && len(call.Args) == 2 {
		// In OR step: accu || predicate, return the non-accu argument
		for _, arg := range call.Args {
			if ident := arg.GetIdentExpr(); ident == nil || ident.Name != accuVar {
				return arg
			}
		}
	}
	return nil
}

func (con *converter) extractTransformFromAppendStep(step *exprpb.Expr, accuVar string) *exprpb.Expr {
	if call := step.GetCallExpr(); call != nil && len(call.Args) == 2 {
		// In append step: accu + [transform], find the list and extract its first element
		for _, arg := range call.Args {
			if listExpr := arg.GetListExpr(); listExpr != nil && len(listExpr.Elements) == 1 {
				return listExpr.Elements[0]
			}
		}
	}
	return nil
}

func (con *converter) extractPredicateFromConditionalStep(step *exprpb.Expr) *exprpb.Expr {
	if call := step.GetCallExpr(); call != nil && len(call.Args) == 3 {
		// In conditional step: conditional(predicate, then, else), return predicate
		return call.Args[0]
	}
	return nil
}

func (con *converter) extractFilterAndTransformFromConditionalStep(step *exprpb.Expr, accuVar string) (*exprpb.Expr, *exprpb.Expr) {
	if call := step.GetCallExpr(); call != nil && len(call.Args) == 3 {
		filter := call.Args[0]
		thenExpr := call.Args[1]

		// Extract transform from then expression: accu + [transform]
		if thenCall := thenExpr.GetCallExpr(); thenCall != nil && thenCall.Function == operators.Add {
			for _, arg := range thenCall.Args {
				if listExpr := arg.GetListExpr(); listExpr != nil && len(listExpr.Elements) == 1 {
					return filter, listExpr.Elements[0]
				}
			}
		}
	}
	return nil, nil
}
