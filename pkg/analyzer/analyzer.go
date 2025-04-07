package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	linterName          = "gocheckerrbeforeuse"
	linterDoc           = "Checks that err is checked before struct use."
	linterReturnMessage = "error must be checked right after receiving"
)

func NewAnalyzer(settings *Settings) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     linterName,
		Doc:      linterDoc,
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

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

				pass.Reportf(stmt.Pos(), linterReturnMessage)
			}
		}
	},
	)

	return nil, nil //nolint:nilnil
}
