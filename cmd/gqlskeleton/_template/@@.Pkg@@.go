@@ if eq .Kind "query" -@@
package @@.Pkg@@

import (
	"github.com/vektah/gqlparser/v2/ast"

	"github.com/gqlgo/gqlanalysis"
)

const doc = "@@.Pkg@@ is ..."

// Analyzer is ...
var Analyzer = &gqlanalysis.Analyzer{
	Name: "@@.Pkg@@",
	Doc:  doc,
	Run:  run,
}

func run(pass *gqlanalysis.Pass) (interface{}, error) {
	for _, q := range pass.Queries {
		for _, f := range q.Fragments {
			for _, sel := range f.SelectionSet {
				switch sel := sel.(type) {
				case *ast.Field:
					if sel.Name == "name" {
						pass.Reportf(sel.Position, "NG")
					}
				}
			}
		}
	}
	return nil, nil
}
@@ end -@@
@@ if eq .Kind "codegen" -@@
package @@.Pkg@@

import (
	"fmt"
	"io"
	"os/exec"
	"text/template"

	"github.com/gqlgo/gqlanalysis/codegen"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	flagOutput string
)

const doc = "@@.Pkg@@ is ..."

var Generator = &codegen.Generator{
	Name: "@@.Pkg@@",
	Doc:  doc,
	Run:  run,
}

func init() {
	Generator.Flags.StringVar(&flagOutput, "output", "@@.Pkg@@.kt", "output file")
}

func run(pass *codegen.Pass) (rerr error) {
	output, path := pass.CreateTemp("@@.Pkg@@.kt")

	for _, q := range pass.Queries {
		if len(q.Fragments) == 0 {
			continue
		}
		tmpl := codegen.NewTemplate(pass, "@@.Pkg@@-template")
		_, err := tmpl.Funcs(funcMap(pass, tmpl)).Parse(tmplStr)
		if err != nil {
			return err
		}
		if err := tmpl.ExecuteTemplate(output, "fragments", q.Fragments); err != nil {
			return err
		}
	}

	if err := exec.Command("ktfmt", path).Run(); err != nil {
		return err
	}

	if _, err := output.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if _, err := io.Copy(pass.Output, output); err != nil {
		return err
	}

	return nil
}

func funcMap(pass *codegen.Pass, tmpl *template.Template) template.FuncMap {
	return map[string]interface{}{
		"zero": func(typ *ast.Type) string {
			switch typ.Name() {
			case "String", "ID":
				return `""`
			case "Boolean":
				return "false"
			}

			td := pass.Schema.Types[typ.Name()]
			if td != nil && len(td.EnumValues) != 0 {
				return typ.Name() + "." + td.EnumValues[0].Name
			}

			panic(fmt.Sprintf("unexpected type: %#v", typ))
		},
	}
}

var tmplStr = `
{{define "fragments"}}
{{range .}}val {{templateWithMeta "fragment" "" .}}{{end}}
{{end}}
{{define "fragment"}}
{{- lower .Name}} = {{meta}}{{.Name}}(
	{{templateWithMeta "selectionSet" (cat .Name ".") .SelectionSet -}}
)
{{end}}
{{define "selection" -}}
	{{- with field . }}
		{{if .SelectionSet}}{{lower .Name}} = {{meta}}{{upper .Name}}(
			{{- templateWithMeta "selectionSet" meta .SelectionSet -}}
		),
		{{else}}{{.Name}} = {{zero .Definition.Type}},{{end}}
	{{- end -}}
	{{- with (fragmentspread .) }}
		{{with .Definition}}fragments = {{meta}}Fragments(
		{{lower .Name}} = {{upper .Name}}(
			{{- templateWithMeta "selectionSet" meta .SelectionSet -}}
		)){{else}}{{end}}
	{{end}}
	{{- with (inlinefragment .) }}
		{{templateWithMeta "selectionSet" meta .SelectionSet}}
	{{end}}
{{end}}
{{define "selectionSet"}}{{range .}}{{templateWithMeta "selection" meta .}}{{end}}{{end}}
`
@@ end -@@
