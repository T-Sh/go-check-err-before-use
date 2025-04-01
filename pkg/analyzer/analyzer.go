package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "gocheckerrbeforeuse",
	Doc:      "Checks that err is checked before struct use.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

const errPrefix = "err"

func run(pass *analysis.Pass) (interface{}, error) {
	orderedInspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.BlockStmt)(nil),
	}

	orderedInspector.Preorder(nodeFilter, func(node ast.Node) {
		blockStmt := node.(*ast.BlockStmt)

		for pos, stmt := range blockStmt.List {
			if isAssignWithErr(stmt) {
				nextStmtPos := pos + 1
				if nextStmtPos >= len(blockStmt.List) {
					return
				}

				nextStmt := blockStmt.List[nextStmtPos]

				if isIfWithErr(nextStmt) || isCallWithErr(nextStmt) || isSwitch(nextStmt) || isAssignWithErrUse(nextStmt) {
					return
				}

				pass.Reportf(stmt.Pos(), "error must be checked right after receiving")
			}
		}
	},
	)

	return nil, nil //nolint:nilnil
}

func checkValueNameWithErr(name string) bool {
	return name == errPrefix
}

// Checks that assigment contains err in return values.
// Example: v, err := someFunc()
// Skips single err return.
// Example for skip:  err := someFunc()
func isAssignWithErr(node ast.Node) bool {
	if assignStmt, ok := node.(*ast.AssignStmt); ok {
		if len(assignStmt.Lhs) == 1 {
			return false
		}

		for _, expr := range assignStmt.Lhs {
			if ident, ok := expr.(*ast.Ident); ok {
				if checkValueNameWithErr(ident.Name) {
					return true
				}
			}
		}
	}

	return false
}

// Checks that if statement contains err check inside.
// Example: if err != nil ...
// Example: if IsError(err) ...
func isIfWithErr(node ast.Node) bool {
	ifStmt, ok := node.(*ast.IfStmt)
	if ok { //nolint:nestif
		if binExpr, ok := ifStmt.Cond.(*ast.BinaryExpr); ok {
			if isExpContainsErr(binExpr.X) {
				return true
			}

			if Xbin, ok := binExpr.X.(*ast.BinaryExpr); ok {
				if isExpContainsErr(Xbin.X) {
					return true
				}
			}

			if Ybin, ok := binExpr.Y.(*ast.BinaryExpr); ok {
				if isExpContainsErr(Ybin.X) {
					return true
				}
			}
		}

		if callExpr, ok := ifStmt.Cond.(*ast.CallExpr); ok {
			return isExprContainsErrInCall(callExpr)
		}
	}

	return false
}

// Checks that if inside contains err usage.
// Example: err != nil ...
func isExpContainsErr(expr ast.Expr) bool {
	if X, ok := expr.(*ast.Ident); ok {
		if checkValueNameWithErr(X.Name) {
			return true
		}
	}

	if callExpr, ok := expr.(*ast.CallExpr); ok {
		return isExprContainsErrInCall(callExpr)
	}

	return false
}

// Checks that func call uses err inside.
// Example: IsError(err) ...
func isExprContainsErrInCall(expr *ast.CallExpr) bool {
	for _, arg := range expr.Args {
		if ident, ok := arg.(*ast.Ident); ok {
			if checkValueNameWithErr(ident.Name) {
				return true
			}
		}
	}

	return false
}

// Checks that func call uses err inside.
// Example: checkError(err)
func isCallWithErr(node ast.Node) bool {
	if exprStmt, ok := node.(*ast.ExprStmt); ok {
		if X, ok := exprStmt.X.(*ast.CallExpr); ok {
			return isExpContainsErr(X)
		}
	}

	return false
}

// Checks that switch uses err variable.
// Example: switch err {...
// Example: switch { err != nil {...
func isSwitch(node ast.Node) bool {
	if switchStmt, ok := node.(*ast.SwitchStmt); ok {
		return isSwitchWithTag(switchStmt) || isSwitchWithBody(switchStmt)
	}

	return false
}

// Checks that switch uses err variable in tag.
// Example: switch err {...
func isSwitchWithTag(switchStmt *ast.SwitchStmt) bool {
	if ident, ok := switchStmt.Tag.(*ast.Ident); ok {
		if checkValueNameWithErr(ident.Name) {
			return true
		}
	}

	return false
}

// Checks that switch uses err variable in body.
// Example: switch { err != nil {...
func isSwitchWithBody(switchStmt *ast.SwitchStmt) bool {
	for _, caseClauseStmt := range switchStmt.Body.List {
		if caseClause, ok := caseClauseStmt.(*ast.CaseClause); ok {
			for _, expr := range caseClause.List {
				return isIfWithErr(expr) || isExpContainsErr(expr)
			}
		}
	}

	return false
}

// Checks that assigment uses err in func call.
// Example: res := errCheck(err)
func isAssignWithErrUse(node ast.Node) bool {
	assignStmt, ok := node.(*ast.AssignStmt)
	if ok { //nolint:nestif
		for _, stmt := range assignStmt.Rhs {
			if right, ok := stmt.(*ast.CompositeLit); ok {
				for _, elt := range right.Elts {
					if expr, ok := elt.(*ast.KeyValueExpr); ok {
						if ident, ok := expr.Value.(*ast.Ident); ok {
							if checkValueNameWithErr(ident.Name) {
								return true
							}
						}
					}
				}
			}
		}
	}

	return false
}
