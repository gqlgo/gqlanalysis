package codegentest

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/gqlgo/gqlanalysis/analysistest"
	"github.com/gqlgo/gqlanalysis/codegen"
)

var TestData = analysistest.TestData

// A Result holds the result of applying a generator to a package.
type Result struct {
	Dir    string
	Pass   *codegen.Pass
	Err    error
	Output *bytes.Buffer
}

func Run(t testing.TB, testdata string, g *codegen.Generator, dirs ...string) []*Result {

	var buf bytes.Buffer
	_g := *g
	g = &_g
	if g.Output != nil {
		g.Output = io.MultiWriter(g.Output, &buf)
	} else {
		g.Output = &buf
	}

	a := g.ToAnalyzer()
	rs := analysistest.Run(t, testdata, a, dirs...)
	results := make([]*Result, len(rs))
	for i := range rs {
		gpass := &codegen.Pass{
			Generator: g,
			Schema:    rs[i].Pass.Schema,
			Queries:   rs[i].Pass.Queries,
			Comments:  rs[i].Pass.Comments,
			ResultOf:  rs[i].Pass.ResultOf,
			Output:    g.Output,
		}
		results[i] = &Result{
			Dir:    dirs[i],
			Pass:   gpass,
			Err:    rs[i].Err,
			Output: &buf,
		}
	}

	return results
}

// Golden compares the results with golden files.
// Golden creates read a golden file which name is codegen.Generator.Name + ".golden".
// The golden file is stored in same directory of the package.
// If Golden cannot find a golden file or the result of Generator test is not same with the golden,
// Golden reports error via *testing.T.
// If update is true, golden files would be updated.
//
// 	var flagUpdate bool
//
// 	func TestMain(m *testing.M) {
// 		flag.BoolVar(&flagUpdate, "update", false, "update the golden files")
// 		flag.Parse()
// 		os.Exit(m.Run())
// 	}
//
// 	func TestGenerator(t *testing.T) {
// 		rs := codegentest.Run(t, codegentest.TestData(), example.Generator, "example")
// 		codegentest.Golden(t, rs, flagUpdate)
// 	}
func Golden(t testing.TB, results []*Result, update bool) {
	t.Helper()
	for _, r := range results {
		golden(t, r, update)
	}
}

func golden(t testing.TB, r *Result, update bool) {
	t.Helper()

	got := r.Output.String()
	r.Output = bytes.NewBufferString(got)

	fname := fmt.Sprintf("%s.golden", r.Pass.Generator.Name)
	fpath := filepath.Join(r.Dir, fname)
	if !update {
		gf, err := os.ReadFile(fpath)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}

		if diff := cmp.Diff(string(gf), got); diff != "" {
			gname := r.Pass.Generator.Name
			t.Errorf("%s's output is different from the golden file(%s):\n%s", gname, fpath, diff)
		}
		return
	}

	newGolden, err := os.Create(fpath)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if _, err := io.Copy(newGolden, strings.NewReader(got)); err != nil {
		t.Fatal("unexpected error:", err)
	}
	if err := newGolden.Close(); err != nil {
		t.Fatal("unexpected error:", err)
	}
}
