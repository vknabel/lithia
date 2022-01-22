package runtime

import "fmt"

var _ RuntimeValue = PreludeCurriedCallable{}
var _ CallableRuntimeValue = PreludeExternFunction{}

type PreludeCurriedCallable struct {
	actual         CallableRuntimeValue
	arguments      []Evaluatable
	remainingArity int
}

func MakeCurriedCallable(actual CallableRuntimeValue, arguments []Evaluatable) PreludeCurriedCallable {
	return PreludeCurriedCallable{
		actual:         actual,
		arguments:      arguments,
		remainingArity: actual.Arity() - len(arguments),
	}
}

func (f PreludeCurriedCallable) Lookup(member string) (Evaluatable, *RuntimeError) {
	switch member {
	case "arity":
		return NewConstantRuntimeValue(PreludeInt(f.Arity())), nil
	default:
		return nil, NewRuntimeErrorf("no such member: %s", member)
	}
}

func (PreludeCurriedCallable) RuntimeType() RuntimeTypeRef {
	return PreludeFunctionTypeRef
}

func (f PreludeCurriedCallable) String() string {
	return fmt.Sprintf("%s curried by %d", f.actual.String(), f.remainingArity)
}

func (f PreludeCurriedCallable) Arity() int {
	return f.remainingArity
}

func (f PreludeCurriedCallable) Call(args []Evaluatable) (RuntimeValue, *RuntimeError) {
	allArgs := append(f.arguments, args...)
	if len(args) < f.remainingArity {
		return MakeCurriedCallable(f.actual, allArgs), nil
	} else {
		return f.actual.Call(allArgs)
	}
}
