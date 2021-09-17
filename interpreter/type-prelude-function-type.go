package interpreter

import "fmt"

var _ RuntimeValue = PreludeFunctionType{}

type PreludeFunctionType struct{}

func (f PreludeFunctionType) String() string {
	return "Function"
}

func (PreludeFunctionType) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Function",
		moduleName: "prelude",
	}
}

func (f PreludeFunctionType) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("function type %s has no member %s", f, member)
}
