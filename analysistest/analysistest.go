package analysistest

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gqlgo/gqlanalysis"
	"github.com/gqlgo/gqlanalysis/internal/checker"
)

// TestData returns absolute path of testdata.
// If TestData cannot get the path, the test will be failed.
var TestData = func(t testing.TB) string {
	t.Helper()
	testdata, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal(err)
	}
	return testdata
}

// Result is a result of an analyzer.
type Result = checker.AnalyzerResult

type pos struct {
	file string
	line int
}

// Run runs tests for an analyzer.
// The given directries which are in testdata, have "schema" and "query" directory.
// Run applies the analyzer for schemas and queries which are hold on these directories.
// Run checks whether expected diagnostics is reported and unexpected ones are not by the analyzer.
// An expected diagnostic can be written with a comment in a .graphql file.
// The comment begin "want" and a Go's regular expression folows it.
// For example, if the analyzer must report "NG" as a diagnostic, it can test with such as following.
//
//	query Q {
//	    a { # want "NG"
//	         name
//	    }
//	}
func Run(t testing.TB, testdata string, a *gqlanalysis.Analyzer, dirs ...string) []*Result {
	t.Helper()
	var results []*Result

	for _, dir := range dirs {
		c := &checker.Checker{
			Schema: filepath.Join(testdata, dir, "schema", "**", "*.graphql"),
			Query:  filepath.Join(testdata, dir, "query", "**", "*.graphql"),
		}
		rs, err := c.RunSingle(a)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		} else {
			results = append(results, rs...)
		}
	}

	for _, r := range results {
		line2cmnt := wantComments(t, r.Pass.Comments)
		for _, d := range r.Diagnostics {
			p := pos{file: d.Pos.Src.Name, line: d.Pos.Line}
			regexps := line2cmnt[p]
			if len(regexps) == 0 {
				t.Errorf("unxpected diagnostic in %s:%d: %s", d.Pos.Src.Name, d.Pos.Line, d.Message)
				continue
			}

			if !regexps[0].MatchString(d.Message) {
				t.Errorf("diagnostic %q does not match %s in %s:%d", d.Message, regexps[0], d.Pos.Src.Name, d.Pos.Line)
			}

			line2cmnt[p] = regexps[1:]
			if len(line2cmnt[p]) == 0 {
				delete(line2cmnt, p)
			}
		}

		for p, regexps := range line2cmnt {
			for _, reg := range regexps {
				t.Errorf("diagnostic is not reported which match with %s in %s:%d", reg, p.file, p.line)
			}
		}
	}

	return results
}

var wantRegexp = regexp.MustCompile(`"([^"]+)"`)

func wantComments(t testing.TB, comments []*gqlanalysis.Comment) map[pos][]*regexp.Regexp {
	t.Helper()
	line2cmnt := make(map[pos][]*regexp.Regexp)
	for _, cmnt := range comments {
		line := strings.TrimLeft(cmnt.Value, "# ")
		if !strings.HasPrefix(line, "want ") {
			continue
		}

		wants := line[len("want"):]
		for _, submatches := range wantRegexp.FindAllStringSubmatchIndex(wants, -1) {
			want := wantRegexp.ExpandString(nil, "$1", wants, submatches)
			re, err := regexp.Compile(string(want))
			if err != nil {
				t.Error("want comment regexp compile:", err)
				continue
			}

			p := pos{file: cmnt.Pos.Src.Name, line: cmnt.Pos.Line}
			line2cmnt[p] = append(line2cmnt[p], re)
		}
	}

	return line2cmnt
}
