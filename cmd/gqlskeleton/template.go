package main

import (
	"embed"
	"text/template"

	"github.com/gostaticanalysis/skeletonkit"
)

//go:embed _template/*
var tmplFS embed.FS

var tmpl = template.Must(skeletonkit.ParseTemplate(tmplFS, "gqlskeleton", "_template"))
