package fragment_test

import (
	"flag"
	"os"
	"testing"

	"github.com/gqlgo/gqlanalysis/codegen/_example/fragment"
	"github.com/gqlgo/gqlanalysis/codegen/codegentest"
)

var flagUpdate bool

func TestMain(m *testing.M) {
	flag.BoolVar(&flagUpdate, "update", false, "update the golden files")
	flag.Parse()
	os.Exit(m.Run())
}

func Test(t *testing.T) {
	testdata := codegentest.TestData(t)
	rs := codegentest.Run(t, testdata, fragment.Generator, "a")
	for i := range rs {
		if rs[i].Err != nil {
			t.Fatal("unexpected err", rs[i].Err)
		}
	}
	codegentest.Golden(t, rs, flagUpdate)
}
