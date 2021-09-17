package interpreter

import (
	"fmt"
)

var _ RuntimeValue = PreludeAnyType{}

type PreludeAnyType struct{}

func (PreludeAnyType) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Any",
		modulePath: []string{"prelude"},
	}
}

func (a PreludeAnyType) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("any type %s has no member %s", a, member)
}

func (t RuntimeType) RuntimeType() RuntimeType {
	typeValue := t.typeValue
	if typeValue == nil {
		return PreludeAnyType{}.RuntimeType()
	} else {
		return (*typeValue).RuntimeType()
	}
}

func (t RuntimeType) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("type %s has no member %s", t.name, member)
}
