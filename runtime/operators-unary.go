package runtime

func (ex *InterpreterContext) unaryOperatorFunction(operator string) (func(Evaluatable) (RuntimeValue, *RuntimeError), *RuntimeError) {
	switch operator {
	case "!":
		return func(expr Evaluatable) (RuntimeValue, *RuntimeError) {
			value, err := expr.Evaluate()
			if err != nil {
				return nil, err
			}
			flag, err := ex.boolFromRuntimeValue(value)
			if err != nil {
				return nil, err
			}
			return ex.boolToRuntimeValue(!flag)
		}, nil
	default:
		return nil, NewRuntimeErrorf("unknown binary operator: %s", operator)
	}
}

func (ex *InterpreterContext) boolFromRuntimeValue(value RuntimeValue) (bool, *RuntimeError) {
	trueRef := MakeRuntimeTypeRef("True", "prelude")
	isTrue, err := trueRef.HasInstance(ex.interpreter, value)
	if err != nil {
		return false, NewRuntimeError(err)
	} else if isTrue {
		return true, nil
	}
	falseRef := MakeRuntimeTypeRef("True", "prelude")
	_, err = falseRef.HasInstance(ex.interpreter, value)
	if err != nil {
		return false, NewRuntimeError(err)
	} else {
		return false, nil
	}
}
