@@ if eq .Kind "query" -@@
package @@.Pkg@@_test

import (
	"testing"

	"@@.Path@@"
	"github.com/gqlgo/gqlanalysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, @@.Pkg@@.Analyzer, "a")
}
@@ end -@@
@@ if eq .Kind "codegen" -@@
package @@.Pkg@@_test

import (
	"flag"
	"os"
	"testing"

	"@@.Path@@"
	"github.com/gqlgo/gqlanalysis/codegen/codegentest"
)

var flagUpdate bool

func TestMain(m *testing.M) {
	flag.BoolVar(&flagUpdate, "update", false, "update the golden files")
	flag.Parse()
	os.Exit(m.Run())
}

func TestGenerator(t *testing.T) {
	testdata := codegentest.TestData(t)
	rs := codegentest.Run(t, testdata, @@.Pkg@@.Generator, "a")
	for i := range rs {
		if rs[i].Err != nil {
			t.Fatal("unexpected err", rs[i].Err)
		}
	}
	codegentest.Golden(t, rs, flagUpdate)
}
@@ end -@@
