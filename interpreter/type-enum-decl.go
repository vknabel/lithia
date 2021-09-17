package interpreter

import "fmt"

var _ RuntimeValue = EnumDeclRuntimeValue{}

type EnumDeclRuntimeValue struct {
	name  string
	cases map[string]*LazyRuntimeValue
}

func (EnumDeclRuntimeValue) RuntimeType() RuntimeType {
	return PreludeAnyType{}.RuntimeType()
}

func (e EnumDeclRuntimeValue) Lookup(member string) (*LazyRuntimeValue, error) {
	if typeDecl, ok := e.cases[member]; ok {
		return typeDecl, nil
	} else {
		return nil, fmt.Errorf("enum type %s has no member %s", fmt.Sprint(e), member)
	}
}
