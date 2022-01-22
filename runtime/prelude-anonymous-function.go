package runtime

import (
	"fmt"
	"strings"
)

var _ RuntimeValue = PreludeAnonymousFunction{}
var _ CallableRuntimeValue = PreludeAnonymousFunction{}

type PreludeAnonymousFunction struct {
	Name   string
	Params []string
	Impl   func(args []Evaluatable) (RuntimeValue, *RuntimeError)
}

func MakeAnonymousFunction(
	name string,
	params []string,
	impl func(args []Evaluatable) (RuntimeValue, *RuntimeError),
) PreludeAnonymousFunction {
	return PreludeAnonymousFunction{
		Name:   name,
		Params: params,
		Impl:   impl,
	}
}

func (f PreludeAnonymousFunction) Lookup(member string) (Evaluatable, *RuntimeError) {
	switch member {
	case "arity":
		return NewConstantRuntimeValue(PreludeInt(f.Arity())), nil
	default:
		return nil, NewRuntimeErrorf("no such member: %s", member)
	}
}

func (PreludeAnonymousFunction) RuntimeType() RuntimeTypeRef {
	return PreludeFunctionTypeRef
}

func (f PreludeAnonymousFunction) String() string {
	return fmt.Sprintf("<extern %s %s>", f.Name, strings.Join(f.Params, ", "))
}

func (f PreludeAnonymousFunction) Arity() int {
	return len(f.Params)
}

func (f PreludeAnonymousFunction) Call(args []Evaluatable) (RuntimeValue, *RuntimeError) {
	if len(args) != f.Arity() {
		panic("use Call to call functions!")
	}
	return f.Impl(args)
}
