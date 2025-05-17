package main

import (
	"github.com/T-Sh/go-check-err-before-use/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	setting := analyzer.Settings{}
	singlechecker.Main(analyzer.NewAnalyzer(setting))
}
