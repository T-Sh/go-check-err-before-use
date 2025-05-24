package main

import (
	"flag"

	"github.com/T-Sh/go-check-err-before-use/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	var flagDistance = flag.Int("distance", 1, "The acceptable distance between the assignment of an error and its checking")
	flag.Parse()

	setting := analyzer.Settings{Distance: *flagDistance}
	singlechecker.Main(analyzer.NewAnalyzer(setting))
}
