package interpreter

import (
	"fmt"
	"reflect"
	"strings"
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

type DataDeclRuntimeValue struct {
	name   string
	fields []DataDeclField
}

func (d DataDeclRuntimeValue) RuntimeType() RuntimeType {
	return PreludeAnyType{}.RuntimeType()
}

func (f DataDeclRuntimeValue) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("data type %s has no member %s", f, member)
}

type DataDeclField struct {
	name   string
	params []string
}

type DataRuntimeValue struct {
	typeValue *DataDeclRuntimeValue
	members   map[string]*LazyRuntimeValue
}

func (d DataRuntimeValue) RuntimeType() RuntimeType {
	var typeValue RuntimeValue = *d.typeValue
	return RuntimeType{
		name:      d.typeValue.name,
		typeValue: &typeValue,
	}
}

func (d DataRuntimeValue) String() string {
	params := make([]string, 0)
	for _, arg := range d.members {
		value, err := arg.Evaluate()
		if err != nil {
			params = append(params, err.Error())
		} else {
			params = append(params, fmt.Sprint(value))
		}
	}

	return fmt.Sprintf("%s %s", d.typeValue.name, strings.Join(params, ", "))
}

type EnumDeclRuntimeValue struct {
	name  string
	cases map[string]*LazyRuntimeValue
}

func (EnumDeclRuntimeValue) RuntimeType() RuntimeType {
	return PreludeAnyType{}.RuntimeType()
}

func (e EnumDeclRuntimeValue) Lookup(member string) (*LazyRuntimeValue, error) {
	if typeDecl, ok := e.cases[member]; ok {
		return typeDecl, nil
	} else {
		return nil, fmt.Errorf("enum type %s has no member %s", fmt.Sprint(e), member)
	}
}

type TypeExpression struct {
	typeValue EnumDeclRuntimeValue
	cases     map[string]*LazyRuntimeValue
}

func (TypeExpression) RuntimeType() RuntimeType {
	return PreludeFunctionType{}.RuntimeType()
}

func (t TypeExpression) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("function %s has no member %s", fmt.Sprint(t), member)
}

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
