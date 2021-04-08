//go:build go1.12
// +build go1.12

package fschecker_test

import (
	"errors"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gqlgo/gqlanalysis"
	"github.com/gqlgo/gqlanalysis/fschecker"
)

func main() {
	var (
		report   bool
		hasError bool
	)
	flag.BoolVar(&report, "report", false, "report diagnositc")
	flag.BoolVar(&hasError, "error", false, "return an error")

	a := &gqlanalysis.Analyzer{
		Name: "lint",
		Doc:  "test linter",
		Run: func(pass *gqlanalysis.Pass) (interface{}, error) {
			if report {
				pos := pass.Queries[0].Position
				pass.Reportf(pos, "NG")
			}
			if hasError {
				return nil, errors.New("error")
			}
			return nil, nil
		},
	}
	fschecker.Main(os.DirFS("."), a)
}

func TestExitCode(t *testing.T) {

	if os.Getenv("FSCHECKER_CHILD") == "1" {
		// child process

		// replace [progname -test.run=TestExitCode -- ...]
		//      by [progname ...]
		os.Args = os.Args[2:]
		os.Args[0] = "lint"
		main()
		panic("unreachable")
	}

	cases := []struct {
		args string
		want int
	}{
		{"", 0},
		{"-error", 1},
		{"-report", 2},
	}

	for _, tt := range cases {
		schema := filepath.Join("testdata", "schema", "**", "**.graphql")
		query := filepath.Join("testdata", "query", "**", "**.graphql")
		args := []string{"-test.run=TestExitCode", "--", "-schema", schema, "-query", query}
		args = append(args, strings.Split(tt.args, " ")...)
		cmd := exec.Command(os.Args[0], args...)
		cmd.Env = append(os.Environ(), "FSCHECKER_CHILD=1")
		out, err := cmd.CombinedOutput()
		if len(out) > 0 {
			t.Logf("%s: out=<<%s>>", tt.args, out)
		}
		var exitcode int
		if err, ok := err.(*exec.ExitError); ok {
			exitcode = err.ExitCode() // requires go1.12
		}
		if exitcode != tt.want {
			t.Errorf("%s: exited %d, want %d", tt.args, exitcode, tt.want)
		}
	}
}
