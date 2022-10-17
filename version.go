package gqlanalysis

import _ "embed"

//go:embed versionfile
var version string

// Version returns version of gqlanalysis.
func Version() string {
	return version
}
