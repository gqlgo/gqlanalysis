package main

import (
	"github.com/gqlgo/gqlanalysis/codegen"
	"github.com/gqlgo/gqlanalysis/codegen/_example/fragment"
)

func main() {
	codegen.Main(fragment.Generator)
}
