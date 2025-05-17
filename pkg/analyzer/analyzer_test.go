package analyzer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/T-Sh/go-check-err-before-use/pkg/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestPositive(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata")
	analysistest.Run(t, testdata, analyzer.NewAnalyzer(analyzer.Settings{}), "positive")
}

func TestNegative(t *testing.T) {
	t.Parallel()
	t.Skip("skipping negative cases")

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata")
	analysistest.Run(t, testdata, analyzer.NewAnalyzer(analyzer.Settings{}), "negative")
}
