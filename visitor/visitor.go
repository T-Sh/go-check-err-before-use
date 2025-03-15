package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"strings"
)

func main() {
	v := visitor{fset: token.NewFileSet()}
	for _, filePath := range os.Args[1:] {
		if filePath == "--" { // to be able to run this like "go run example_files.go -- input.go"
			continue
		}

		f, err := parser.ParseFile(v.fset, filePath, nil, 0)
		if err != nil {
			log.Fatalf("Failed to parse file %s: %s", filePath, err)
		}

		ast.Walk(&v, f)
	}
}

type visitor struct {
	fset *token.FileSet
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	return v
}

func (v *visitor) checkAssignWithErr(node ast.Node) ast.Visitor {
	assignStmt, ok := node.(*ast.AssignStmt)
	if ok {
		for _, expr := range assignStmt.Lhs {
			if ident, ok := expr.(*ast.Ident); ok {
				if strings.HasPrefix(ident.Name, "err") {
					var buf bytes.Buffer
					printer.Fprint(&buf, v.fset, expr)
					fmt.Printf("%s | %#v\n", buf.String(), expr)
				}
			}
		}
	}

	return v
}

func (v *visitor) checkIfWithErr(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if ok {
		if binExpr, ok := ifStmt.Cond.(*ast.BinaryExpr); ok {
			if X, ok := binExpr.X.(*ast.Ident); ok {
				if strings.HasPrefix(X.Name, "err") {
					var buf bytes.Buffer
					printer.Fprint(&buf, v.fset, binExpr)
					fmt.Printf("%s | %#v\n", buf.String(), binExpr)
				}
			}
		}
	}

	return v
}
