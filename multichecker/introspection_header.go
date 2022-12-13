package multichecker

import (
	"fmt"
	"net/http"
	"strings"
)

// introspectionHeader is conformance to `Value` in flag package.
type introspectionHeader http.Header

func (ih introspectionHeader) String() string {
	var s string
	for k, v := range ih {
		if len(s) != 0 {
			s += ","
		}
		s += fmt.Sprintf("%v=%v", k, v)
	}
	return s
}

func (ih introspectionHeader) Set(args string) error {
	keyAndValueList := strings.Split(args, ",")
	for _, keyAndValueString := range keyAndValueList {
		keyAndValue := strings.Split(keyAndValueString, "=")
		key := keyAndValue[0]
		value := keyAndValue[1]
		// Supports only one value.
		ih[key][0] = value
	}
	return nil

}
