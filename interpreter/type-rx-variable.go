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

func (v *RuntimeVariable) String() string {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return fmt.Sprintf("(%s %s)", v.RuntimeType().name, v.current)
}

func (*RuntimeVariable) RuntimeType() RuntimeType {
	return RxVariableType{}.RuntimeType()
}

func (v *RuntimeVariable) Lookup(member string) (Evaluatable, error) {
	switch member {
	case "accept":
		return NewConstantRuntimeValue(NewBuiltinFunction(
			"accept",
			[]string{"value"},
			Docs{
				docs: "Evalutes and sets the current value",
			},
			func(args []Evaluatable) (RuntimeValue, error) {
				return v.Accept(args[0])
			},
		)), nil
	case "current":
		return NewLazyRuntimeValue(func() (RuntimeValue, LocatableError) {
			return v.Current()
		}), nil
	default:
		return nil, fmt.Errorf("variable %s has no member %s", fmt.Sprint(v), member)
	}
}

func (v *RuntimeVariable) Accept(lazyValue Evaluatable) (RuntimeValue, error) {
	value, err := lazyValue.Evaluate()
	if err != nil {
		return nil, err
	}
	v.lock.Lock()
	defer v.lock.Unlock()
	v.current = value
	return value, nil
}

func (v *RuntimeVariable) Current() (RuntimeValue, LocatableError) {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return v.current, nil
}
