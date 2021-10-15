package runtime

var _ RuntimeValue = PreludeCurriedCallable{}

type PreludeCurriedCallable struct {
	Actual         CallableRuntimeValue
	Arguments      []Evaluatable
	RemainingArity int
}

func (PreludeCurriedCallable) Lookup(member string) (Evaluatable, error) {
	panic("TODO: not implemented")
}

func (PreludeCurriedCallable) RuntimeType() RuntimeTypeRef {
	return PreludeFunctionTypeRef
}

func (PreludeCurriedCallable) String() string {
	panic("TODO: not implemented")
}
