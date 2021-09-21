package interpreter

import "fmt"

var _ RuntimeValue = PreludeFunctionType{}
var _ DocumentedRuntimeValue = PreludeFunctionType{}

type PreludeFunctionType struct {
	docs Docs
}

func (f PreludeFunctionType) String() string {
	return "Function"
}

func (PreludeFunctionType) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Function",
		moduleName: "prelude",
	}
}

func (f PreludeFunctionType) Lookup(member string) (Evaluatable, error) {
	return nil, fmt.Errorf("function type %s has no member %s", f, member)
}

func (f PreludeFunctionType) GetDocs() Docs {
	return f.docs
}
