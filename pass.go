package gqlanalysis

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
)

// A Pass provides information to the Run function that applies a specific analyzer.
// The Run function should not call any of the Pass functions concurrently.
type Pass struct {
	Analyzer *Analyzer

	Schema   *ast.Schema
	Queries  []*ast.QueryDocument
	Comments []*Comment

	Report   func(*Diagnostic)
	ResultOf map[*Analyzer]interface{}
}

// Reportf reports a diagnostic with a format.
func (pass *Pass) Reportf(pos *ast.Position, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	pass.Report(&Diagnostic{Pos: pos, Message: msg})
}
