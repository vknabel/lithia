package interpreter

import (
	"fmt"
	"strings"
)

type Callable interface {
	Call(arguments []*LazyRuntimeValue) (RuntimeValue, error)
	String() string
}

type CurriedCallable struct {
	actual         Callable
	args           []*LazyRuntimeValue
	remainingArity int
}

type Function struct {
	name      string
	arguments []string
	body      func(*ExecutionContext) ([]*LazyRuntimeValue, error)
	parent    *ExecutionContext
}

func (f Function) String() string {
	return fmt.Sprintf("{ %s => @(%s) }", strings.Join(f.arguments, ","), f.name)
}

func (Function) RuntimeType() RuntimeType {
	return PreludeFunctionType{}.RuntimeType()
}

func (f Function) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("function %s has no member %s", f, member)
}

func (CurriedCallable) RuntimeType() RuntimeType {
	return PreludeFunctionType{}.RuntimeType()
}

func (f CurriedCallable) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("function %s has no member %s", f, member)
}

func (f CurriedCallable) String() string {
	return fmt.Sprintf("{ -%d => @(%s) }", len(f.args), f.actual)
}

func (c CurriedCallable) Call(arguments []*LazyRuntimeValue) (RuntimeValue, error) {
	allArgs := append(c.args, arguments...)
	if len(arguments) < c.remainingArity {
		lazy := NewLazyRuntimeValue(func() (RuntimeValue, error) {
			return CurriedCallable{
				actual:         c.actual,
				args:           allArgs,
				remainingArity: c.remainingArity - len(arguments),
			}, nil
		})
		return lazy.Evaluate()
	} else {
		return c.actual.Call(allArgs)
	}
}

func (dataDecl DataDeclRuntimeValue) Call(arguments []*LazyRuntimeValue) (RuntimeValue, error) {
	if len(arguments) < len(dataDecl.fields) {
		lazy := NewLazyRuntimeValue(func() (RuntimeValue, error) {
			return CurriedCallable{
				actual:         dataDecl,
				args:           arguments,
				remainingArity: len(dataDecl.fields) - len(arguments),
			}, nil
		})
		return lazy.Evaluate()
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
		return instance, nil
	} else {
		// error
		return nil, fmt.Errorf("too many arguments")
	}
}

func (typeExpr TypeExpression) Call(arguments []*LazyRuntimeValue) (RuntimeValue, error) {
	if len(arguments) == 0 {
		return typeExpr, nil
	}
	lazyValueArgument := arguments[0]
	valueArgument, err := lazyValueArgument.Evaluate()
	if err != nil {
		return nil, err
	}
	for caseName, lazyCaseImpl := range typeExpr.cases {
		caseTypeValue, err := typeExpr.typeValue.cases[caseName].Evaluate()
		if err != nil {
			return nil, err
		}
		ok, err := RuntimeTypeValueIncludesValue(caseTypeValue, valueArgument)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}

		intermediate, err := lazyCaseImpl.Evaluate()
		if err != nil {
			return nil, err
		}
		callable, ok := intermediate.(Callable)
		if !ok {
			return nil, fmt.Errorf("case %s is not callable", caseName)
		}
		return callable.Call(arguments)
	}
	return nil, fmt.Errorf("no %s has no matching case for value %s of type %s", typeExpr.typeValue.name, fmt.Sprint(valueArgument), fmt.Sprint(valueArgument.RuntimeType().name))
}

func (f Function) Call(arguments []*LazyRuntimeValue) (RuntimeValue, error) {
	closure := f.parent.NestedExecutionContext(f.name)
	if len(arguments) < len(f.arguments) {
		return CurriedCallable{
			actual:         f,
			args:           arguments,
			remainingArity: len(f.arguments) - len(arguments),
		}, nil
	}
	for i, argName := range f.arguments {
		err := closure.environment.Declare(argName, arguments[i])
		if err != nil {
			return nil, err
		}
	}

	var (
		lastValue RuntimeValue
		err       error
	)
	statements, err := f.body(closure)
	if err != nil {
		return nil, err
	}
	for _, statement := range statements {
		if statement == nil {
			continue
		}
		lastValue, err = statement.Evaluate()
		if err != nil {
			return nil, err
		}
	}

	if len(arguments) == len(f.arguments) {
		return lastValue, nil
	} else if function, ok := lastValue.(Callable); ok {
		lazyResult, err := function.Call(arguments[len(f.arguments):])
		if err != nil {
			return nil, err
		}
		return lazyResult, nil
	} else {
		return nil, fmt.Errorf("function %s returns %s, which is not callable", f.name, lastValue)
	}

}
