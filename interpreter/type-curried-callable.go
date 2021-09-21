package interpreter

import "fmt"

var _ RuntimeValue = CurriedCallable{}
var _ Callable = CurriedCallable{}

type CurriedCallable struct {
	actual         Callable
	args           []Evaluatable
	remainingArity int
}

func (CurriedCallable) RuntimeType() RuntimeType {
	return PreludeFunctionType{}.RuntimeType()
}

func (f CurriedCallable) Lookup(member string) (Evaluatable, error) {
	return nil, fmt.Errorf("function %s has no member %s", f, member)
}

func (f CurriedCallable) String() string {
	return fmt.Sprintf("{ -%d => @(%s) }", len(f.args), f.actual)
}

func (c CurriedCallable) Call(arguments []Evaluatable) (RuntimeValue, error) {
	allArgs := append(c.args, arguments...)
	if len(arguments) < c.remainingArity {
		lazy := NewLazyRuntimeValue(func() (RuntimeValue, LocatableError) {
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
