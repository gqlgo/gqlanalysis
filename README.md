# gqlanalysis

[![pkg.go.dev][gopkg-badge]][gopkg]

`gqlanalysis` defines the interface between a modular static analysis for GraphQL in Go.
`gqlanalysis` is inspired by [go/analysis](https://golang.org/x/tools/go/analysis).

`gqlanalysis` makes easy to develop static analysis tools for GraphQL in Go.

## How to use

### Analyzer

The primary type in the API is Analyzer.
An Analyzer statically describes an analysis function: its name, documentation, flags, relationship to other analyzers, and of course, its logic.

```go
package lackid

var Analyzer = &gqlanalysis.Analyzer{
	Name: "lackid",
	Doc:  "lackid finds a selection for a type which has id field but the selection does not have id",
	Run:  run,
	...
}

func run(pass *gqlanalysis.Pass) (interface{}, error) {
	...
}
```

### Driver

An analysis driver is a program that runs a set of analyses and prints the diagnostics that they report.
The driver program must import the list of Analyzers it needs.

A typical driver can be created with multichecker package.

```go
package main

import (
        "github.com/gqlgo/gqlanalysis/multichecker"
        "github.com/gqlgo/lackid"
        "github.com/gqlgo/myanalyzer"
)

func main() {
        multichecker.Main(
		lackid.Analyzer,
		myanalyzer.Analyzer,
	)
}
```

### Pass

A Pass describes a single unit of work: the application of a particular Analyzer to given GraphQL's schema and query files.
The Pass provides information to the Analyzer's Run function about schemas and queries being analyzed, and provides operations to the Run function for reporting diagnostics and other information back to the driver.

```go
type Pass struct {
        Analyzer *Analyzer

        Schema   *ast.Schema
        Queries  []*ast.QueryDocument
        Comments []*Comment

        Report   func(*Diagnostic)
        ResultOf map[*Analyzer]interface{}
}
```

### Diagnostic

A Diagnostic is a message associated with a source location.
Pass can report a diagnostic via Report field or Reportf method.

```go
type Diagnostic struct {
        Pos     *ast.Position
        Message string
}
```

## Implementations of Analyzer

* [gqlgo/lackid](https://github.com/gqlgo/lackid) - Detect lack of id in GraphQL query

## Author

[![Appify Technologies, Inc.](appify-logo.png)](http://github.com/appify-technologies)

<!-- links -->
[gopkg]: https://pkg.go.dev/github.com/gqlgo/gqlanalysis
[gopkg-badge]: https://pkg.go.dev/badge/github.com/gqlgo/gqlanalysis?status.svg
