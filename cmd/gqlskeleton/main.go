package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/gostaticanalysis/skeleton/v2/skeleton"
	"github.com/gqlgo/gqlanalysis"
	"golang.org/x/mod/module"
)

func main() {
	os.Exit(s.Run(gqlanalysis.Version(), os.Args[1:]))
}

func run(version string, args []string) int {
	if len(args) > 0 && args[0] == "-v" {
		fmt.Fprintln(s.Output, "gqlskeleton", version)
		return skeleton.ExitSuccess
	}

	var info Info
	flags, err := parseFlag(args, &info)
	if err != nil {
		fmt.Fprintln(s.ErrOutput, "Error:", err)
		return skeleton.ExitError
	}

	info.Path = flags.Arg(0)
	if prefix := os.Getenv("GQLSKELETON_PREFIX"); prefix != "" {
		info.Path = path.Join(prefix, info.Path)
	}
	// allow package name only
	if module.CheckImportPath(info.Path) != nil {
		flags.Usage()
		return skeleton.ExitError
	}

	if info.Pkg == "" {
		info.Pkg = path.Base(info.Path)
	}

	if err := s.run(&info); err != nil {
		fmt.Fprintln(s.ErrOutput, "Error:", err)
		return skeleton.ExitError
	}

	return skeleton.ExitSuccess
}

func parseFlag(args []string, info *Info) (*flag.FlagSet, error) {
	flags := flag.NewFlagSet("gqlskeleton", flag.ContinueOnError)
	flags.SetOutput(s.ErrOutput)
	flags.Usage = func() {
		fmt.Fprintln(s.ErrOutput, "gqlskeleton [-kind,-cmd] example.com/path")
		flags.PrintDefaults()
	}
	flags.Var(&info.Kind, "kind", "[query,codegen]")

	flags.BoolVar(&info.Cmd, "cmd", true, "create main file")
	flags.StringVar(&info.Pkg, "pkg", "", "package name")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	if info.Kind == "" {
		info.Kind = KindQuery
	}

	return flags, nil
}

func generate(info *Info) error {
	g := &skeleton.Generator{
		Template: tmpl,
	}

	fsys, err := g.Run(info)
	if err != nil {
		return err
	}

	prompt := &skeleton.Prompt{
		Output:    os.Stdout,
		ErrOutput: os.Stderr,
		Input:     os.Stdin,
	}

	dst := filepath.Join(".", info.Pkg)
	if err := skeleton.CreateDir(prompt, dst, fsys); err != nil {
		return err
	}

	return nil
}
