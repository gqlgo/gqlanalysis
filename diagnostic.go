package gqlanalysis

import "github.com/vektah/gqlparser/v2/ast"

// A Diagnostic is a message associated with a source location.
type Diagnostic struct {
	Pos     *ast.Position
	Message string
}
