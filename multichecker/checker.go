package multichecker

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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

func Main(analyzers ...*gqlanalysis.Analyzer) {
	flag.Parse()
	dir := flag.Arg(0)
	checker := &checker.Checker{
		Schema: filepath.Join(dir, flagSchema),
		Query:  filepath.Join(dir, flagQuery),
	}
	if err := checker.Run(analyzers...); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
