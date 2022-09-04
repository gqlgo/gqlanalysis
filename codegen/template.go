package codegen

import (
	"bytes"
	"fmt"
	"text/template"
	"unicode"

	"github.com/vektah/gqlparser/v2/ast"
)

func NewTemplate(pass *Pass, name string) *template.Template {
	tmpl := template.New(name)
	return tmpl.Funcs(NewTemplateFuncs(tmpl, pass))
}

func NewTemplateFuncs(tmpl *template.Template, pass *Pass) template.FuncMap {
	var meta interface{}
	return template.FuncMap{
		"field":          toField,
		"fragmentspread": toFragmentSpread,
		"inlinefragment": toInlineFragment,
		"lower": func(s string) string {
			runes := []rune(s)
			runes[0] = unicode.ToLower(runes[0])
			return string(runes)
		},
		"upper": func(s string) string {
			runes := []rune(s)
			runes[0] = unicode.ToUpper(runes[0])
			return string(runes)
		},
		"cat": func(vs ...interface{}) string {
			var s string
			for _, v := range vs {
				s += fmt.Sprint(v)
			}
			return s
		},
		"typeOf": func(s string) *ast.Definition {
			return pass.Schema.Types[s]
		},
		"templateWithMeta": func(name string, _meta, data interface{}) string {
			var buf bytes.Buffer
			before := meta
			defer func() {
				meta = before
			}()
			meta = _meta
			if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
				return err.Error()
			}
			return buf.String()
		},
		"meta": func() interface{} {
			return meta
		},
	}
}

func toField(v interface{}) *ast.Field {
	f, _ := v.(*ast.Field)
	return f
}

func toFragmentSpread(v interface{}) *ast.FragmentSpread {
	f, _ := v.(*ast.FragmentSpread)
	return f
}

func toInlineFragment(v interface{}) *ast.InlineFragment {
	f, _ := v.(*ast.InlineFragment)
	return f
}
