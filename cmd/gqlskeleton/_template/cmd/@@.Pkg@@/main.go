@@ if .Cmd -@@
@@ if eq .Kind "query" -@@
package main

import (
	"@@.Path@@"
	"github.com/gqlgo/gqlanalysis/multichecker"
)

func main() { multichecker.Main(@@.Pkg@@.Analyzer) }
@@ end -@@
@@ if eq .Kind "codegen" -@@
package main

import (
	"@@.Path@@"
	"github.com/gqlgo/gqlanalysis/codegen"
)

func main() {
	codegen.Main(@@.Pkg@@.Generator)
}
@@ end -@@
@@end@@
