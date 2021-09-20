package interpreter

import (
	"fmt"
)

var _ RuntimeValue = PreludeModuleType{}
var _ DocumentedRuntimeValue = PreludeModuleType{}

type PreludeModuleType struct {
	docs Docs
}

func (PreludeModuleType) String() string {
	return "Module"
}

func (PreludeModuleType) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Module",
		moduleName: "prelude",
	}
}

func (a PreludeModuleType) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("module type %s has no member %s", a, member)
}

func (m PreludeModuleType) GetDocs() Docs {
	return m.docs
}
