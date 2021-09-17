package interpreter

import "fmt"

var _ RuntimeValue = PreludeModuleType{}

type PreludeModuleType struct{}

func (PreludeModuleType) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Module",
		modulePath: []string{"prelude"},
	}
}

func (a PreludeModuleType) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("module type %s has no member %s", a, member)
}
