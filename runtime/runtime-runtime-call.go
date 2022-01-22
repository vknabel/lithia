package runtime

func Call(function RuntimeValue, args []Evaluatable) (RuntimeValue, *RuntimeError) {
	callable, ok := function.(CallableRuntimeValue)
	if !ok {
		return nil, NewRuntimeErrorf("%s is not callable", function)
	}

	arity := callable.Arity()
	if len(args) < arity {
		return MakeCurriedCallable(callable, args), nil
	}
	intermediate, err := callable.Call(args[:arity])
	if err != nil {
		return nil, err
	}
	if len(args) == arity {
		return intermediate, nil
	}
	if g, ok := intermediate.(CallableRuntimeValue); ok {
		return g.Call(args[arity:])
	} else {
		return nil, NewRuntimeErrorf("%s is not callable", g)
	}
}
