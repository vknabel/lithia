package interpreter

import (
	"fmt"
	"sync"
)

var _ RuntimeValue = &RuntimeVariable{}

type RuntimeVariable struct {
	lock    *sync.RWMutex
	current RuntimeValue
}

func (*RuntimeVariable) RuntimeType() RuntimeType {
	return PreludeVariableType{}.RuntimeType()
}

func (v *RuntimeVariable) Lookup(member string) (*LazyRuntimeValue, error) {
	switch member {
	case "accept":
		return NewConstantRuntimeValue(NewBuiltinFunction(
			"accept",
			[]string{"value"},
			func(args []*LazyRuntimeValue) (RuntimeValue, error) {
				return v.Accept(args[0])
			},
		)), nil
	case "current":
		return NewLazyRuntimeValue(func() (RuntimeValue, error) {
			return v.Current()
		}), nil
	default:
		return nil, fmt.Errorf("variable %s has no member %s", fmt.Sprint(v), member)
	}
}

func (v *RuntimeVariable) Accept(lazyValue *LazyRuntimeValue) (RuntimeValue, error) {
	value, err := lazyValue.Evaluate()
	if err != nil {
		return nil, err
	}
	v.lock.Lock()
	defer v.lock.Unlock()
	v.current = value
	return value, nil
}

func (v *RuntimeVariable) Current() (RuntimeValue, error) {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return v.current, nil
}
