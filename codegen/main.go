package codegen

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gqlgo/gqlanalysis/multichecker"
)

// Main is the main function for a code generation command for a generator.
// It is a wrapper of multichecker.Main.
// See github.com/gqlgo/gqlanalysis/singlechecker.
func Main(g *Generator) {

	g.Flags.Parse(os.Args[1:])
	programName := os.Args[0]
	os.Args = make([]string, g.Flags.NArg()+1)
	os.Args[0] = programName
	copy(os.Args[1:], g.Flags.Args())
	flag.CommandLine.SetOutput(ioutil.Discard)

	a := g.ToAnalyzer()
	g.Flags.Usage = func() {
		paras := strings.Split(g.Doc, "\n\n")
		fmt.Fprintf(os.Stderr, "%s: %s\n\n", g.Name, paras[0])
		fmt.Fprintf(os.Stderr, "Usage: %s [-flag] [package]\n\n", g.Name)
		if len(paras) > 1 {
			fmt.Fprintln(os.Stderr, strings.Join(paras[1:], "\n\n"))
		}
		fmt.Fprintln(os.Stderr, "\nFlags:")
		g.Flags.PrintDefaults()
	}

	if g.Flags.NArg() == 0 {
		g.Flags.Usage()
		os.Exit(1)
	}

	multichecker.Main(a)
}
