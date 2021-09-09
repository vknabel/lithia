package interpreter

import (
	"fmt"
)

type Callable interface {
	Call(arguments []*LazyRuntimeValue) (*LazyRuntimeValue, error)
}

type CurriedCallable struct {
	actual         Callable
	args           []*LazyRuntimeValue
	remainingArity int
}

func (c CurriedCallable) Call(arguments []*LazyRuntimeValue) (*LazyRuntimeValue, error) {
	allArgs := append(c.args, arguments...)
	if len(arguments) < c.remainingArity {
		lazy := NewLazyRuntimeValue(func() (RuntimeValue, error) {
			return CurriedCallable{
				actual:         c.actual,
				args:           allArgs,
				remainingArity: c.remainingArity - len(arguments),
			}, nil
		})
		return lazy, nil
	} else {
		return c.actual.Call(allArgs)
	}
}

func (dataDecl DataDeclRuntimeValue) Call(arguments []*LazyRuntimeValue) (*LazyRuntimeValue, error) {
	if len(arguments) < len(dataDecl.fields) {
		lazy := NewLazyRuntimeValue(func() (RuntimeValue, error) {
			return CurriedCallable{
				actual:         dataDecl,
				args:           arguments,
				remainingArity: len(dataDecl.fields) - len(arguments),
			}, nil
		})
		return lazy, nil
	} else if len(arguments) == len(dataDecl.fields) {
		members := make(map[string]*LazyRuntimeValue, len(dataDecl.fields))
		for i, field := range dataDecl.fields {
			arg := arguments[i]
			members[field.name] = arg
		}
		instance := DataRuntimeValue{
			typeValue: &dataDecl,
			members:   members,
		}
		return NewConstantRuntimeValue(instance), nil
	} else {
		// error
		return nil, fmt.Errorf("too many arguments")
	}
}
