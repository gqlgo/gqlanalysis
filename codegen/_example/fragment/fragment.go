package fragment

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"text/template"

	"github.com/gqlgo/gqlanalysis/codegen"
)

var (
	flagOutput string
)

var Generator = &codegen.Generator{
	Name: "fragment",
	Doc:  "example of codegen",
	Run:  run,
}

func init() {
	Generator.Flags.StringVar(&flagOutput, "output", "", "output dir")
}

func run(pass *codegen.Pass) error {
	var buf bytes.Buffer
	for _, q := range pass.Queries {
		if len(q.Fragments) == 0 {
			continue
		}
		if err := tmpl.ExecuteTemplate(&buf, "fragments", q.Fragments); err != nil {
			return err
		}
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	fmt.Print(string(src))

	if _, err := io.Copy(pass.Output, bytes.NewReader(src)); err != nil {
		return err
	}
	return nil
}

var tmpl = template.Must(codegen.NewTemplate("fragment-template").Funcs(map[string]interface{}{
	"type": func(n string) string {
		switch n {
		case "String", "ID":
			return "string"
		default:
			return n
		}
	},
}).Parse(`
{{define "fragments"}}
package fragments
{{range .}}{{template "fragment" .}}{{end}}
{{end}}

{{define "fragment"}}
type {{.Name}} struct {
	{{template "selectionSet" .SelectionSet}}
}
{{end}}

{{define "selection" -}}
	{{- with field . }} {{.Name}} {{type .Definition.Type.Name}} {{end -}}
	{{- with (fragmentspread .) }}
		{{template "selectionSet" .Definition.SelectionSet}}
	{{end}}
	{{- with (inlinefragment .) }}
		{{template "selectionSet" .SelectionSet}}
	{{end}}
{{end}}

{{define "selectionSet"}}{{range .}}{{template "selection" .}}{{end}}{{end}}
`))
