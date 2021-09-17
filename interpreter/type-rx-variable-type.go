package interpreter

import (
	"fmt"
	"sync"
)

var _ RuntimeValue = RxVariableType{}
var _ Callable = &RxVariableType{}

type RxVariableType struct{}

func (v RxVariableType) String() string {
	return v.RuntimeType().String()
}

func (RxVariableType) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Variable",
		moduleName: "rx",
	}
}

func (v RxVariableType) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("variable type %s has no member %s", v, member)
}

func (v RxVariableType) Call(arguments []*LazyRuntimeValue) (RuntimeValue, error) {
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
