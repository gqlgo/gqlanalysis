package main

import (
	"embed"
	"io/fs"
	"text/template"

	"github.com/gostaticanalysis/skeleton/v2/skeleton"
	"github.com/josharian/txtarfs"
	"golang.org/x/tools/txtar"
)

//go:embed _template/*
var tmplFS embed.FS

//go:embed version.txt
var version string

var tmpl = template.Must(parseTemplate(tmplFS, "gqlskeleton", "_template"))

func parseTemplate(tmplFS embed.FS, name, prefix string) (*template.Template, error) {
	fsys, err := fs.Sub(tmplFS, prefix)
	if err != nil {
		return nil, err
	}
	ar, err := txtarfs.From(fsys)
	if err != nil {
		return nil, err
	}
	strTmpl := string(txtar.Format(ar))

	return template.New(name).Delims("@@", "@@").Funcs(skeleton.DefaultFuncMap).Parse(strTmpl)
}
