package multichecker

import (
	"flag"
	"net/http"
	"os"

	"github.com/gqlgo/gqlanalysis"
	"github.com/gqlgo/gqlanalysis/internal/checker"
)

var (
	flagSchema string
	flagQuery  string
)

var ih = make(introspectionHeader)

func init() {
	flag.StringVar(&flagSchema, "schema", "schema/**/*.graphql", "pattern of schema")
	flag.StringVar(&flagQuery, "query", "query/**/*.graphql", "pattern of query")
	flag.Var(ih, "introspection-header", "format key1:value1,key2:value2")
}

func Main(analyzers ...*gqlanalysis.Analyzer) {
	flag.Parse()
	checker := &checker.Checker{
		Schema:              flagSchema,
		Query:               flagQuery,
		IntrospectionHeader: http.Header(ih),
	}

	os.Exit(checker.Run(analyzers...))
}
