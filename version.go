package gqlanalysis

//go:embed version.txt
var version string

// Version returns version of gqlanalysis.
func Version() string {
	return version
}
