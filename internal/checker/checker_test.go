package checker_test

import (
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/josharian/txtarfs"
	"github.com/vektah/gqlparser/v2/ast"
	"golang.org/x/tools/txtar"

	"github.com/gqlgo/gqlanalysis"
	"github.com/gqlgo/gqlanalysis/internal/checker"
)

type fileTransport struct {
	fileName string
}

func (t *fileTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	f, err := os.Open(t.fileName)
	if err != nil {
		return nil, err
	}
	w := httptest.NewRecorder()
	if _, err := io.Copy(w, f); err != nil {
		return nil, err
	}
	return w.Result(), nil
}

func fsys(s string) fs.FS {
	a := txtar.Parse([]byte(strings.TrimSpace(s)))
	return txtarfs.As(a)
}

func TestChecker_Run_Introspection(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		fsys     fs.FS
		schema   string
		query    string
		testdata string
	}{
		"introspection": {
			fsys: fsys(`
-- query/q.graphql --
query GetA() {
    a { # check
        id
	name
    }
}
			`),
			schema:   "http://example.com",
			query:    "query/*.graphql",
			testdata: "testdata/introspection/schema.json",
		},
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			checker := &checker.Checker{
				Fsys:   tt.fsys,
				Schema: tt.schema,
				Query:  tt.query,
				HTTPClient: &http.Client{
					Transport: &fileTransport{fileName: tt.testdata},
				},
			}
			a := &gqlanalysis.Analyzer{
				Name:       name,
				Doc:        name,
				ResultType: reflect.TypeOf(false),
				Run: func(pass *gqlanalysis.Pass) (interface{}, error) {
					if len(pass.Comments) == 0 {
						return false, nil
					}
					for _, q := range pass.Queries {
						for _, op := range q.Operations {
							for _, sel := range op.SelectionSet {
								if pass.Comments[0].Pos.Line != sel.GetPosition().Line {
									continue
								}
								field, _ := sel.(*ast.Field)
								if field != nil && field.Definition != nil {
									return true, nil
								}
							}
						}
					}
					return false, nil
				},
			}
			results, err := checker.RunSingle(a)
			if err != nil {
				t.Fatal("unexpected error:", err)
			}
			got := results[0].Result.(bool)
			if !got {
				t.Error("does not get expected schema from instrospection")
			}
		})
	}
}

func TestChecker_Run_Glob(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		fsys   fs.FS
		schema string
		query  string
	}{
		"glob": {
			fsys: fsys(`
-- schema/models/a.gql --
type A {
    id: ID!
    name: String!
}
-- schema/schema.gql --
schema {
    query: Query
}
-- schema/query.gql --
type Query {
    a: A!
}
-- query/q.gql --
query GetA() {
    a { # check
        id
	name
    }
}
			`),
			schema: "schema/**/*.gql",
			query:  "query/*.gql",
		},
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			checker := &checker.Checker{
				Fsys:   tt.fsys,
				Schema: tt.schema,
				Query:  tt.query,
			}
			a := &gqlanalysis.Analyzer{
				Name:       name,
				Doc:        name,
				ResultType: reflect.TypeOf(false),
				Run: func(pass *gqlanalysis.Pass) (interface{}, error) {
					if len(pass.Comments) == 0 {
						return false, nil
					}
					for _, q := range pass.Queries {
						for _, op := range q.Operations {
							for _, sel := range op.SelectionSet {
								if pass.Comments[0].Pos.Line != sel.GetPosition().Line {
									continue
								}
								field, _ := sel.(*ast.Field)
								if field != nil && field.Definition != nil {
									return true, nil
								}
							}
						}
					}
					return false, nil
				},
			}
			results, err := checker.RunSingle(a)
			if err != nil {
				t.Fatal("unexpected error:", err)
			}
			got := results[0].Result.(bool)
			if !got {
				t.Error("does not get expected schema")
			}
		})
	}
}
