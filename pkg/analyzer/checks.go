package analyzer

import (
	"go/ast"
	"strings"
)

const (
	errPrefix = "err"
	okPrefix  = "ok"
)

func checkValueNameWithErr(name string) bool {
	return strings.HasPrefix(name, errPrefix) || name == okPrefix
}

func allChecks(nextStmt ast.Node) bool {
	return isIfWithErr(nextStmt) || isCallWithErr(nextStmt) || isReturnWithErr(nextStmt) || isSwitchWithErr(nextStmt) || isAssignWithErrUse(nextStmt)
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

// Checks that assigment contains err in return values with var.
// Example: var v, err = someFunc()
// Skips single err return.
// Example for skip:  var err = someFunc()
func isDeclWithErr(node ast.Node) bool {
	if declStmt, ok := node.(*ast.DeclStmt); ok {
		if genDecl, ok := declStmt.Decl.(*ast.GenDecl); ok {
			if genDecl.Tok.String() == "var" {
				if len(genDecl.Specs) == 1 {
					return false
				}
				for _, spec := range genDecl.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						for _, name := range valueSpec.Names {
							if isExpContainsErr(name) {
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

// Checks that if statement contains err check inside.
// Example: if err != nil ...
// Example: if IsError(err) ...
func isIfWithErr(node ast.Node) bool {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return false
	}

	if isExpContainsErr(ifStmt.Cond) {
		return true
	}

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

	for _, stmt := range ifStmt.Body.List {
		if allChecks(stmt) {
			return true
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

	if unaryExp, ok := expr.(*ast.UnaryExpr); ok {
		if isExpContainsErr(unaryExp.X) {
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
func isSwitchWithErr(node ast.Node) bool {
	if switchStmt, ok := node.(*ast.SwitchStmt); ok {
		return isSwitchWithTag(switchStmt) || isSwitchWithCase(switchStmt)
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
func isSwitchWithCase(switchStmt *ast.SwitchStmt) bool {
	for _, caseClauseStmt := range switchStmt.Body.List {
		if caseClause, ok := caseClauseStmt.(*ast.CaseClause); ok {
			for _, expr := range caseClause.List {
				if allChecks(expr) || isExpContainsErr(expr) {
					return true
				}
			}

			if isCaseWithBody(caseClause) {
				return true
			}
		}
	}

	return false
}

func isCaseWithBody(caseStmt *ast.CaseClause) bool {
	for _, stmt := range caseStmt.Body {
		if allChecks(stmt) {
			return true
		}
	}

	return false
}

// Checks that assigment uses err in func call.
// Example: res := errCheck(err)
func isAssignWithErrUse(node ast.Node) bool {
	assignStmt, ok := node.(*ast.AssignStmt)
	if !ok {
		return false
	}

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

	return false
}

// Checks that return stmt contains err.
// Example: return res, err
func isReturnWithErr(node ast.Node) bool {
	if returnStmt, ok := node.(*ast.ReturnStmt); ok {
		for _, result := range returnStmt.Results {
			if isExpContainsErr(result) {
				return true
			}
		}
	}

	return false
}
