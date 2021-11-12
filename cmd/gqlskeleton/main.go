package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/gostaticanalysis/skeletonkit"
	"github.com/gqlgo/gqlanalysis"
	"golang.org/x/mod/module"
)

const (
	ExitSuccess = 0
	ExitError   = 1
)

func main() {
	os.Exit(run(gqlanalysis.Version(), os.Args[1:]))
}

func run(version string, args []string) int {
	if len(args) > 0 && args[0] == "-v" {
		fmt.Println("gqlskeleton", version)
		return ExitSuccess
	}

	var info Info
	flags, err := parseFlag(args, &info)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return ExitError
	}

	info.Path = flags.Arg(0)
	if prefix := os.Getenv("GQLSKELETON_PREFIX"); prefix != "" {
		info.Path = path.Join(prefix, info.Path)
	}
	// allow package name only
	if module.CheckImportPath(info.Path) != nil {
		flags.Usage()
		return ExitError
	}

	if info.Pkg == "" {
		info.Pkg = path.Base(info.Path)
	}

	if err := generate(&info); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return ExitError
	}

	return ExitSuccess
}

func parseFlag(args []string, info *Info) (*flag.FlagSet, error) {
	flags := flag.NewFlagSet("gqlskeleton", flag.ContinueOnError)
	flags.SetOutput(os.Stdout)
	flags.Usage = func() {
		fmt.Fprintln(os.Stderr, "gqlskeleton [-kind,-cmd] example.com/path")
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
	fsys, err := skeletonkit.ExecuteTemplate(tmpl, info)
	if err != nil {
		return err
	}

	dst := filepath.Join(".", info.Pkg)
	if err := skeletonkit.CreateDir(skeletonkit.DefaultPrompt, dst, fsys); err != nil {
		return err
	}

	return nil
}
