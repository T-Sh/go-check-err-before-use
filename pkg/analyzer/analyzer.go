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

func isAssignWithErr(node ast.Node) bool {
	assignStmt, ok := node.(*ast.AssignStmt)
	if ok {
		for _, expr := range assignStmt.Lhs {
			if ident, ok := expr.(*ast.Ident); ok {
				if ident.Name == errPrefix {
					return true
				}
			}
		}
	}

	return false
}

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

func isExpContainsErr(expr ast.Expr) bool {
	if X, ok := expr.(*ast.Ident); ok {
		if X.Name == errPrefix {
			return true
		}
	}

	if callExpr, ok := expr.(*ast.CallExpr); ok {
		return isExprContainsErrInCall(callExpr)
	}

	return false
}

func isExprContainsErrInCall(expr *ast.CallExpr) bool {
	for _, arg := range expr.Args {
		if ident, ok := arg.(*ast.Ident); ok {
			if ident.Name == errPrefix {
				return true
			}
		}
	}

	return false
}

func isCallWithErr(node ast.Node) bool {
	if exprStmt, ok := node.(*ast.ExprStmt); ok {
		if X, ok := exprStmt.X.(*ast.CallExpr); ok {
			return isExpContainsErr(X)
		}
	}

	return false
}

func isSwitch(node ast.Node) bool {
	if switchStmt, ok := node.(*ast.SwitchStmt); ok {
		return isSwitchWithTag(switchStmt) || isSwitchWithBody(switchStmt)
	}

	return false
}

func isSwitchWithTag(switchStmt *ast.SwitchStmt) bool {
	if ident, ok := switchStmt.Tag.(*ast.Ident); ok {
		if ident.Name == errPrefix {
			return true
		}
	}

	return false
}

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

func isAssignWithErrUse(node ast.Node) bool {
	assignStmt, ok := node.(*ast.AssignStmt)
	if ok { //nolint:nestif
		for _, stmt := range assignStmt.Rhs {
			if right, ok := stmt.(*ast.CompositeLit); ok {
				for _, elt := range right.Elts {
					if expr, ok := elt.(*ast.KeyValueExpr); ok {
						if ident, ok := expr.Value.(*ast.Ident); ok {
							if ident.Name == errPrefix {
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
