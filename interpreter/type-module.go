package interpreter

import "fmt"

var _ RuntimeValue = RuntimeModule{}

type RuntimeModule struct {
	module *Module
	docs   DocString
}

func (m RuntimeModule) String() string {
	return fmt.Sprintf("(import %s)", m.module.name)
}

func (RuntimeModule) RuntimeType() RuntimeType {
	return PreludeModuleType{}.RuntimeType()
}

func (m RuntimeModule) Lookup(member string) (Evaluatable, error) {
	if lazy, ok := m.module.environment.Get(member); ok {
		return lazy, nil
	} else {
		return nil, fmt.Errorf("module %s has no member %s", m.module.name, member)
	}
}
