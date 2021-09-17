package interpreter

import (
	"fmt"
	"reflect"
	"sync"
)

type RuntimeValue interface {
	RuntimeType() RuntimeType
	Lookup(name string) (*LazyRuntimeValue, error)
	String() string
}

var _ RuntimeValue = RuntimeType{}

type RuntimeType struct {
	name       string
	moduleName ModuleName
	typeValue  *RuntimeValue
}

func (t RuntimeType) String() string {
	return fmt.Sprintf("%s.%s", t.moduleName, t.name)
}

func (t RuntimeType) RuntimeTypeValue() RuntimeValue {
	if t.typeValue == nil {
		return t.RuntimeType()
	}
	return *t.typeValue
}

type LazyRuntimeValue struct {
	once  *sync.Once
	value RuntimeValue
	err   LocatableError
	eval  func() (RuntimeValue, LocatableError)
}

func NewLazyRuntimeValue(eval func() (RuntimeValue, LocatableError)) *LazyRuntimeValue {
	return &LazyRuntimeValue{
		once:  &sync.Once{},
		eval:  eval,
		value: nil,
	}
}

func NewConstantRuntimeValue(value RuntimeValue) *LazyRuntimeValue {
	return &LazyRuntimeValue{
		once:  &sync.Once{},
		eval:  func() (RuntimeValue, LocatableError) { return value, nil },
		value: value,
	}
}

func (l *LazyRuntimeValue) Evaluate() (RuntimeValue, LocatableError) {
	l.once.Do(func() {
		l.value, l.err = l.eval()
	})
	return l.value, l.err
}

func RuntimeTypeValueIncludesValue(t RuntimeValue, v RuntimeValue) (bool, LocatableError) {
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
			t = t.(RuntimeType).RuntimeTypeValue()
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

func (t RuntimeType) IncludesValue(v RuntimeValue) (bool, LocatableError) {
	return RuntimeTypeValueIncludesValue(t, v)
}
