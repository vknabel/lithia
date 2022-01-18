package runtime

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

func (PreludeCurriedCallable) Lookup(member string) (Evaluatable, *RuntimeError) {
	panic("TODO: not implemented PreludeCurriedCallable")
}

func (PreludeCurriedCallable) RuntimeType() RuntimeTypeRef {
	return PreludeFunctionTypeRef
}

func (PreludeCurriedCallable) String() string {
	panic("TODO: not implemented PreludeCurriedCallable")
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
