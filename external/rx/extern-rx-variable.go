package rx

import (
	"fmt"
	"sync"

	"github.com/vknabel/lithia/runtime"
)

var _ runtime.RuntimeValue = RxVariable{}

type RxVariable struct {
	lock         *sync.RWMutex
	current      *runtime.RuntimeValue
	variableType *RxVariableType
}

func MakeRxVariable(variableType *RxVariableType, current runtime.RuntimeValue) RxVariable {
	return RxVariable{
		lock:         &sync.RWMutex{},
		current:      &current,
		variableType: variableType,
	}
}

func (RxVariable) RuntimeType() runtime.RuntimeTypeRef {
	return RxVariableTypeRef
}

func (v RxVariable) String() string {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return fmt.Sprintf("(%s %s)", v.RuntimeType().Name, *v.current)
}

func (v RxVariable) Lookup(member string) (runtime.Evaluatable, *runtime.RuntimeError) {
	switch member {
	case "accept":
		return runtime.NewConstantRuntimeValue(runtime.MakeExternTypeMethod(
			v.variableType.Fields["accept"],
			func(args []runtime.Evaluatable) (runtime.RuntimeValue, *runtime.RuntimeError) {
				return v.Accept(args[0])
			},
		)), nil
	case "current":
		return runtime.NewLazyRuntimeValue(func() (runtime.RuntimeValue, *runtime.RuntimeError) {
			return v.Current()
		}), nil
	default:
		return nil, runtime.NewRuntimeErrorf("variable %s has no member %s", fmt.Sprint(v), member)
	}
}

func (v *RxVariable) Accept(lazyValue runtime.Evaluatable) (runtime.RuntimeValue, *runtime.RuntimeError) {
	value, err := lazyValue.Evaluate()
	if err != nil {
		return nil, err
	}
	if eagerEvaluatable, ok := value.(runtime.EagerEvaluatableRuntimeValue); ok {
		err = eagerEvaluatable.EagerEvaluate()
		if err != nil {
			return nil, err
		}
	}
	v.lock.Lock()
	defer v.lock.Unlock()
	*v.current = value
	return value, nil
}

func (v *RxVariable) Current() (runtime.RuntimeValue, *runtime.RuntimeError) {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return *v.current, nil
}
