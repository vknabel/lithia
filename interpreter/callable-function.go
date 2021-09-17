package interpreter

import (
	"fmt"
	"strings"
)

type Function struct {
	name      string
	arguments []string
	body      func(*EvaluationContext) ([]*LazyRuntimeValue, error)
	parent    *EvaluationContext
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
