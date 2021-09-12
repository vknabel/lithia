package interpreter

import (
	"fmt"
	"reflect"
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

func (typeExpr TypeExpression) Call(arguments []*LazyRuntimeValue) (*LazyRuntimeValue, error) {
	if len(arguments) == 0 {
		return NewConstantRuntimeValue(typeExpr), nil
	}
	lazyValueArgument := arguments[0]
	valueArgument, err := lazyValueArgument.Evaluate()
	if err != nil {
		return nil, err
	}
	for caseName := range typeExpr.cases {
		data, ok := valueArgument.(DataRuntimeValue)
		if !ok {
			continue
		}
		caseTypeValue, err := typeExpr.typeValue.cases[caseName].Evaluate()
		if err != nil {
			return nil, err
		}
		if caseDataType, ok := caseTypeValue.(DataDeclRuntimeValue); !ok || !reflect.DeepEqual(*data.typeValue, caseDataType) {
			continue
		}
		switch specificDataType := caseTypeValue.(type) {
		case DataDeclRuntimeValue:
			reflect.DeepEqual(*data.typeValue, specificDataType)
		case EnumDeclRuntimeValue:
			return nil, fmt.Errorf("cannot use enum type as data type")
		default:
			return nil, fmt.Errorf("unexpected type %T", caseTypeValue)
		}

		intermediate, err := typeExpr.cases[caseName].Evaluate()
		if err != nil {
			return nil, err
		}
		callable, ok := intermediate.(Callable)
		if !ok {
			return nil, fmt.Errorf("case %s is not callable", caseName)
		}
		return callable.Call(arguments)
	}
	return nil, fmt.Errorf("no matching case")
}
