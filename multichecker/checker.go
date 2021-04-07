package multichecker

import (
	"flag"
	"fmt"
	"os"

	"github.com/gqlgo/gqlanalysis"
	"github.com/gqlgo/gqlanalysis/internal/checker"
)

var (
	flagSchema string
	flagQuery  string
)

func init() {
	flag.StringVar(&flagSchema, "schema", "schema/**/*.graphql", "pattern of schema")
	flag.StringVar(&flagQuery, "query", "query/**/*.graphql", "pattern of query")
}

func Main(analyzers ...*gqlanalysis.Analyzer) {
	flag.Parse()
	checker := &checker.Checker{
		Schema: flagSchema,
		Query:  flagQuery,
	}
	if err := checker.Run(analyzers...); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
