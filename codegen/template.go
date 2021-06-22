package codegen

import (
	"text/template"

	"github.com/vektah/gqlparser/v2/ast"
)

var TemplateFuncs = template.FuncMap{
	"field":          toField,
	"fragmentspread": toFragmentSpread,
	"inlinefragment": toInlineFragment,
}

func NewTemplate(name string) *template.Template {
	return template.New(name).Funcs(TemplateFuncs)
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
