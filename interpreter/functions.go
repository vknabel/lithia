package interpreter

import (
	"fmt"
	"reflect"
	"strings"
)

type Callable interface {
	Call(arguments []*LazyRuntimeValue) (*LazyRuntimeValue, error)
}

type CurriedCallable struct {
	actual         Callable
	args           []*LazyRuntimeValue
	remainingArity int
}

type Function struct {
	name      string
	arguments []string
	body      func(*Interpreter) ([]*LazyRuntimeValue, error)
	closure   *Interpreter
}

func (f Function) String() string {
	return fmt.Sprintf("{ %s => @(%s) }", strings.Join(f.arguments, ","), strings.Join(f.closure.path, "."))
}

func (Function) RuntimeType() RuntimeType {
	return PreludeFunctionType{}.RuntimeType()
}

func (CurriedCallable) RuntimeType() RuntimeType {
	return PreludeFunctionType{}.RuntimeType()
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

func (f Function) Call(arguments []*LazyRuntimeValue) (*LazyRuntimeValue, error) {
	if len(arguments) < len(f.arguments) {
		return NewConstantRuntimeValue(CurriedCallable{
			actual:         f,
			args:           arguments,
			remainingArity: len(f.arguments) - len(arguments),
		}), nil
	}
	for i, argName := range f.arguments {
		err := f.closure.environment.Declare(argName, arguments[i])
		if err != nil {
			return nil, err
		}
	}
	return NewLazyRuntimeValue(func() (RuntimeValue, error) {
		var (
			lastValue RuntimeValue
			err       error
		)
		statements, err := f.body(f.closure)
		if err != nil {
			return nil, err
		}
		for _, statement := range statements {
			lastValue, err = statement.Evaluate()
			if err != nil {
				return nil, err
			}
		}

		if len(arguments) == len(f.arguments) {
			return lastValue, nil
		} else if function, ok := lastValue.(Callable); ok {
			lazyResult, err := function.Call(arguments[len(f.arguments)-1:])
			if err != nil {
				return nil, err
			}
			return lazyResult.Evaluate()
		} else {
			return nil, fmt.Errorf("function %s returns %s, which is not callable", f.name, lastValue)
		}
	}), nil
}
