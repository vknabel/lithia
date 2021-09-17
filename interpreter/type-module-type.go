package interpreter

import (
	"fmt"
)

var _ RuntimeValue = PreludeModuleType{}

type PreludeModuleType struct{}

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
