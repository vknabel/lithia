package runtime

import (
	"fmt"
	"sync"
)

var _ RuntimeValue = RxVariable{}

type RxVariable struct {
	lock         *sync.RWMutex
	current      RuntimeValue
	variableType *RxVariableType
}

func MakeRxVariable(variableType *RxVariableType, current RuntimeValue) RxVariable {
	return RxVariable{
		lock:         &sync.RWMutex{},
		current:      current,
		variableType: variableType,
	}
}

func (RxVariable) RuntimeType() RuntimeTypeRef {
	return RxVariableTypeRef
}

func (v RxVariable) String() string {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return fmt.Sprintf("(%s %s)", v.RuntimeType().Name, v.current)
}

func (v RxVariable) Lookup(member string) (Evaluatable, *RuntimeError) {
	switch member {
	case "accept":
		return NewConstantRuntimeValue(MakeExternTypeMethod(
			v.variableType.Fields["accept"],
			func(args []Evaluatable) (RuntimeValue, *RuntimeError) {
				if len(args) == 1 {
				}
				return v.Accept(args[0])
			},
		)), nil
	case "current":
		return NewLazyRuntimeValue(func() (RuntimeValue, *RuntimeError) {
			return v.Current()
		}), nil
	default:
		return nil, NewRuntimeErrorf("variable %s has no member %s", fmt.Sprint(v), member)
	}
}

func (v *RxVariable) Accept(lazyValue Evaluatable) (RuntimeValue, *RuntimeError) {
	value, err := lazyValue.Evaluate()
	if err != nil {
		return nil, err
	}
	v.lock.Lock()
	defer v.lock.Unlock()
	v.current = value
	return value, nil
}

func (v *RxVariable) Current() (RuntimeValue, *RuntimeError) {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return v.current, nil
}
