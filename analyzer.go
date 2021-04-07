package gqlanalysis

import (
	"flag"
	"fmt"
	"reflect"
)

// Analyzer is an analyzer for .graphql file.
// It is inspired by golang.org/x/tools/go/analysis.Analyzer.
type Analyzer struct {
	Name       string
	Run        func(pass *Pass) (interface{}, error)
	Doc        string
	Flags      flag.FlagSet
	Requires   []*Analyzer
	ResultType reflect.Type
}

func (a *Analyzer) String() string {
	return a.Name
}

var _ fmt.Stringer = (*Analyzer)(nil)
