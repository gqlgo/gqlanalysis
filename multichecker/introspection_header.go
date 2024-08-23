package multichecker

import (
	"fmt"
	"net/http"
	"strings"
)

// introspectionHeader confirmed to `Value` interface in `flag` package.
type introspectionHeader http.Header

func (ih introspectionHeader) String() string {
	var s string
	for k, v := range ih {
		if len(s) != 0 {
			s += ","
		}
		s += fmt.Sprintf("%v:%v", k, v[0])
	}
	return s
}

func (ih introspectionHeader) Set(args string) error {
	for _, keyAndValueString := range strings.Split(args, ",") {
		keyAndValue := strings.Split(keyAndValueString, ":")
		key := keyAndValue[0]
		value := keyAndValue[1]

		// Supports only one value.
		ih[key] = []string{value}
	}
	return nil

}
