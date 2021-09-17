package interpreter

import (
	"reflect"
	"sync"
)

type RuntimeValue interface {
	RuntimeType() RuntimeType
	Lookup(name string) (*LazyRuntimeValue, error)
}

type RuntimeType struct {
	name       string
	modulePath []string
	typeValue  *RuntimeValue
}

type LazyRuntimeValue struct {
	once  *sync.Once
	value RuntimeValue
	err   error
	eval  func() (RuntimeValue, error)
}

func NewLazyRuntimeValue(eval func() (RuntimeValue, error)) *LazyRuntimeValue {
	return &LazyRuntimeValue{
		once:  &sync.Once{},
		eval:  eval,
		value: nil,
	}
}

func NewConstantRuntimeValue(value RuntimeValue) *LazyRuntimeValue {
	return &LazyRuntimeValue{
		once:  &sync.Once{},
		eval:  func() (RuntimeValue, error) { return value, nil },
		value: value,
	}
}

func (l *LazyRuntimeValue) Evaluate() (RuntimeValue, error) {
	l.once.Do(func() {
		l.value, l.err = l.eval()
	})
	return l.value, l.err
}

func RuntimeTypeValueIncludesValue(t RuntimeValue, v RuntimeValue) (bool, error) {
	if _, ok := t.(PreludeAnyType); ok {
		return true, nil
	}
	if enumDecl, ok := t.(EnumDeclRuntimeValue); ok {
		for _, lazyValue := range enumDecl.cases {
			value, err := lazyValue.Evaluate()
			if err != nil {
				return false, err
			}
			ok, err := value.RuntimeType().IncludesValue(v)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
		}
		return false, nil
	} else {
		if _, ok := t.(RuntimeType); ok {
			t = *t.(RuntimeType).typeValue
		}
		if v.RuntimeType().typeValue == nil {
			return false, nil
		}
		valueType := *v.RuntimeType().typeValue
		if _, ok := valueType.(RuntimeType); ok {
			valueType = *valueType.(RuntimeType).typeValue
		}
		return reflect.DeepEqual(valueType, t), nil
	}
}

func (t RuntimeType) IncludesValue(v RuntimeValue) (bool, error) {
	return RuntimeTypeValueIncludesValue(t, v)
}
