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
	defaultDistance     = 1
)

func NewAnalyzer(settings Settings) *analysis.Analyzer {
	if settings.Distance < 1 {
		settings.Distance = defaultDistance
	}

	return &analysis.Analyzer{
		Name: linterName,
		Doc:  linterDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			run(pass, settings.Distance)

			return nil, nil //nolint:nilnil
		},
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func run(pass *analysis.Pass, distance int) {
	orderedInspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.BlockStmt)(nil),
	}

	orderedInspector.Preorder(nodeFilter, func(node ast.Node) {
		blockStmt := node.(*ast.BlockStmt)

		for pos, stmt := range blockStmt.List {
			if isAssignWithErr(stmt) || isDeclWithErr(stmt) {
				for idx := 1; idx <= distance; idx++ {
					nextStmtPos := pos + idx
					if nextStmtPos >= len(blockStmt.List) {
						return
					}

					nextStmt := blockStmt.List[nextStmtPos]

					if allChecks(nextStmt) {
						return
					}
				}

				pass.Reportf(stmt.Pos(), linterReturnMessage)
			}
		}
	},
	)
}
