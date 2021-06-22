package codegen

import (
	"flag"
	"io"
	"os"

	"github.com/gqlgo/gqlanalysis"
	"github.com/vektah/gqlparser/v2/ast"
)

// A Generator describes a code generator function and its options.
type Generator struct {
	Name     string
	Doc      string
	Flags    flag.FlagSet
	Run      func(*Pass) error
	Requires []*gqlanalysis.Analyzer
	Output   io.Writer
}

// ToAnalyzer converts the generator to an analyzer.
func (g *Generator) ToAnalyzer() *gqlanalysis.Analyzer {
	requires := make([]*gqlanalysis.Analyzer, len(g.Requires))
	requiresMap := make(map[*gqlanalysis.Analyzer]*gqlanalysis.Analyzer)
	for i := range requires {
		_a := g.Requires[i]
		a := *_a // copy
		a.Run = func(pass *gqlanalysis.Pass) (interface{}, error) {
			pass.Report = func(*gqlanalysis.Diagnostic) {}
			return g.Requires[i].Run(pass)
		}
		requires[i] = &a
		requiresMap[&a] = _a
	}

	return &gqlanalysis.Analyzer{
		Name: g.Name,
		Doc:  g.Doc,
		Run: func(pass *gqlanalysis.Pass) (interface{}, error) {
			var output io.Writer = os.Stdout
			if g.Output != nil {
				output = g.Output
			}

			resultOf := make(map[*gqlanalysis.Analyzer]interface{})
			for k, v := range pass.ResultOf {
				a, ok := requiresMap[k]
				if ok {
					resultOf[a] = v
				}
			}

			gpass := &Pass{
				Generator: g,
				Schema:    pass.Schema,
				Queries:   pass.Queries,
				Comments:  pass.Comments,
				ResultOf:  resultOf,
				Output:    output,
			}

			return nil, g.Run(gpass)
		},
		Requires: requires,
	}
}

type Pass struct {
	Generator *Generator

	Schema   *ast.Schema
	Queries  []*ast.QueryDocument
	Comments []*gqlanalysis.Comment
	ResultOf map[*gqlanalysis.Analyzer]interface{}

	Output io.Writer
}
