package main

import "flag"

// Kind represents kind of skeleton codes.
// Kind implements flag.Value.
type Kind string

var _ flag.Value = (*Kind)(nil)

const (
	KindQuery   Kind = "query"
	KindCodegen Kind = "codegen"
)

func (k Kind) String() string {
	switch k {
	case KindCodegen:
		return "codegen"
	default:
		return "query"
	}
}

// "codegen" -> KindCodegen otherwise KindQuery.
func (k *Kind) Set(s string) error {
	switch s {
	case "codegen":
		*k = KindCodegen
	default:
		*k = KindQuery
	}
	return nil
}
