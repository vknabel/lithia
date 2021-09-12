package interpreter

import (
	"fmt"
	"strings"
	"sync"
)

type RuntimeValue interface{}

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

type DataDeclRuntimeValue struct {
	name   string
	fields []DataDeclField
}

type DataDeclField struct {
	name   string
	params []string
}

type DataRuntimeValue struct {
	typeValue *DataDeclRuntimeValue
	members   map[string]*LazyRuntimeValue
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

type TypeExpression struct {
	typeValue EnumDeclRuntimeValue
	cases     map[string]*LazyRuntimeValue
}
