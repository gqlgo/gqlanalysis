package fschecker

import (
	"flag"
	"io/fs"
	"os"

	"github.com/gqlgo/gqlanalysis"
	"github.com/gqlgo/gqlanalysis/internal/checker"
)

var (
	flagSchema string
	flagQuery  string
)

func init() {
	flag.StringVar(&flagSchema, "schema", "schema", "pattern of schema")
	flag.StringVar(&flagQuery, "query", "query", "pattern of query")
}

func Main(fsys fs.FS, analyzers ...*gqlanalysis.Analyzer) {
	flag.Parse()
	checker := &checker.Checker{
		Fsys:   fsys,
		Schema: flagSchema,
		Query:  flagQuery,
	}
	os.Exit(checker.Run(analyzers...))
}

type Result = checker.AnalyzerResult

func RunSingle(fsys fs.FS, schema, query string, analyzer *gqlanalysis.Analyzer) ([]*Result, error) {
	checker := &checker.Checker{
		Fsys:   fsys,
		Schema: schema,
		Query:  query,
	}

	result, err := checker.RunSingle(analyzer)
	if err != nil {
		return nil, err
	}

	return result, nil
}
