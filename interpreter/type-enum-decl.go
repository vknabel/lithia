package interpreter

import (
	"fmt"
	"strings"
)

var _ RuntimeValue = EnumDeclRuntimeValue{}

type EnumDeclRuntimeValue struct {
	name  string
	cases map[string]Evaluatable
	docs  DocString
}

func (e EnumDeclRuntimeValue) String() string {
	cases := make([]string, 0, len(e.cases))
	for key := range e.cases {
		cases = append(cases, key)
	}
	return fmt.Sprintf("enum %s { %s }", e.name, strings.Join(cases, ", "))
}

func (EnumDeclRuntimeValue) RuntimeType() RuntimeType {
	return PreludeAnyType{}.RuntimeType()
}

func (e EnumDeclRuntimeValue) Lookup(member string) (Evaluatable, error) {
	if typeDecl, ok := e.cases[member]; ok {
		return typeDecl, nil
	} else {
		return nil, fmt.Errorf("enum type %s has no member %s", fmt.Sprint(e), member)
	}
}
