package interpreter

import (
	"fmt"
	"sync"
)

var _ RuntimeValue = PreludeVariableType{}
var _ Callable = &PreludeVariableType{}

type PreludeVariableType struct{}

func (PreludeVariableType) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Variable",
		modulePath: []string{"prelude"},
	}
}

func (v PreludeVariableType) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("variable type %s has no member %s", v, member)
}

func (v PreludeVariableType) String() string {
	return "(extern Variable)"
}

func (v PreludeVariableType) Call(arguments []*LazyRuntimeValue) (RuntimeValue, error) {
	var _ Callable = v
	if len(arguments) != 1 {
		return nil, fmt.Errorf("too many arguments for variable type %s", v)
	}
	value, err := arguments[0].Evaluate()
	if err != nil {
		return nil, err
	}
	return &RuntimeVariable{current: value, lock: &sync.RWMutex{}}, nil
}
