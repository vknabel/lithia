package interpreter

import (
	"fmt"
	"strings"
)

var _ RuntimeValue = ExternFunctionDecl{}

type ExternFunctionDecl struct {
	name   string
	params map[string]*LazyRuntimeValue
	docs   DocString
}

func (e ExternFunctionDecl) String() string {
	params := make([]string, 0, len(e.params))
	for key := range e.params {
		params = append(params, key)
	}
	return fmt.Sprintf("extern %s %s", e.name, strings.Join(params, ", "))
}

func (ExternFunctionDecl) RuntimeType() RuntimeType {
	return PreludeAnyType{}.RuntimeType()
}

func (e ExternFunctionDecl) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("extern type %s has no member %s", fmt.Sprint(e), member)
}
