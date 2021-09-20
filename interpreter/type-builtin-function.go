package interpreter

import (
	"fmt"
	"strings"
)

var _ RuntimeValue = BuiltinFunction{}
var _ Callable = BuiltinFunction{}
var _ DocumentedRuntimeValue = BuiltinFunction{}

type BuiltinFunction struct {
	name string
	args []string
	docs Docs
	impl func(args []*LazyRuntimeValue) (RuntimeValue, error)
}

func NewBuiltinFunction(
	name string,
	args []string,
	docs Docs,
	impl func(args []*LazyRuntimeValue) (RuntimeValue, error),
) BuiltinFunction {
	if docs.name == "" {
		docs.name = name
	}
	f := BuiltinFunction{
		name: name,
		args: args,
		impl: impl,
		docs: docs,
	}
	var _ RuntimeValue = f
	var _ Callable = f
	return f
}

func (f BuiltinFunction) RuntimeType() RuntimeType {
	return PreludeFunctionType{}.RuntimeType()
}

func (f BuiltinFunction) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("function %s has no member %s", fmt.Sprint(f), member)
}

func (f BuiltinFunction) String() string {
	return fmt.Sprintf("{ %s => @(%s) }", strings.Join(f.args, ","), f.name)
}

func (f BuiltinFunction) Call(arguments []*LazyRuntimeValue) (RuntimeValue, error) {
	if len(arguments) < len(f.args) {
		return CurriedCallable{
			actual:         f,
			args:           arguments,
			remainingArity: len(f.args) - len(arguments),
		}, nil
	}
	intermediate, err := f.impl(arguments[:len(f.args)])
	if err != nil {
		return nil, err
	}
	if len(arguments) == len(f.args) {
		return intermediate, nil
	}
	if g, ok := intermediate.(Callable); ok {
		return g.Call(arguments[len(f.args):])
	} else {
		return nil, fmt.Errorf("%s is not callable", g)
	}
}

func (f BuiltinFunction) GetDocs() Docs {
	return f.docs
}
