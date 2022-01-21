package runtime

import "reflect"

func (ex *InterpreterContext) BinaryOperatorFunction(operator string) (func(Evaluatable, Evaluatable) (RuntimeValue, *RuntimeError), *RuntimeError) {
	switch operator {
	case "==":
		return func(lazyLeft, lazyRight Evaluatable) (RuntimeValue, *RuntimeError) {
			return ex.genericGreedyComparision(lazyLeft, lazyRight, func(left, right RuntimeValue) bool {
				return reflect.DeepEqual(left, right)
			})
		}, nil
	case "!=":
		return func(lazyLeft, lazyRight Evaluatable) (RuntimeValue, *RuntimeError) {
			return ex.genericGreedyComparision(lazyLeft, lazyRight, func(left, right RuntimeValue) bool {

				return !reflect.DeepEqual(left, right)
			})
		}, nil
	case "&&":
		return func(lazyLeft, lazyRight Evaluatable) (RuntimeValue, *RuntimeError) {
			return ex.lazyLogicComparision(operator, lazyLeft, lazyRight, func(left bool, right func() (bool, *RuntimeError)) (bool, *RuntimeError) {
				if !left {
					return false, nil
				} else {
					return right()
				}
			})
		}, nil
	case "||":
		return func(lazyLeft, lazyRight Evaluatable) (RuntimeValue, *RuntimeError) {
			return ex.lazyLogicComparision(operator, lazyLeft, lazyRight, func(left bool, right func() (bool, *RuntimeError)) (bool, *RuntimeError) {
				if left {
					return true, nil
				} else {
					return right()
				}
			})
		}, nil
	case ">":
		return func(lazyLeft, lazyRight Evaluatable) (RuntimeValue, *RuntimeError) {
			return ex.numericGreedyComparision(operator, lazyLeft, lazyRight, func(left, right PreludeInt) bool {
				return left > right
			}, func(left, right PreludeFloat) bool {
				return left > right
			})
		}, nil
	case ">=":
		return func(lazyLeft, lazyRight Evaluatable) (RuntimeValue, *RuntimeError) {
			return ex.numericGreedyComparision(operator, lazyLeft, lazyRight, func(left, right PreludeInt) bool {
				return left >= right
			}, func(left, right PreludeFloat) bool {
				return left >= right
			})
		}, nil
	case "<":
		return func(lazyLeft, lazyRight Evaluatable) (RuntimeValue, *RuntimeError) {
			return ex.numericGreedyComparision(operator, lazyLeft, lazyRight, func(left, right PreludeInt) bool {
				return left < right
			}, func(left, right PreludeFloat) bool {
				return left < right
			})
		}, nil
	case "<=":
		return func(lazyLeft, lazyRight Evaluatable) (RuntimeValue, *RuntimeError) {
			return ex.numericGreedyComparision(operator, lazyLeft, lazyRight, func(left, right PreludeInt) bool {
				return left <= right
			}, func(left, right PreludeFloat) bool {
				return left <= right
			})
		}, nil
	case "+":
		return func(lazyLeft, lazyRight Evaluatable) (RuntimeValue, *RuntimeError) {
			return ex.numericGreedyOperation(operator, lazyLeft, lazyRight, func(left, right PreludeInt) PreludeInt {
				return left + right
			}, func(left, right PreludeFloat) PreludeFloat {
				return left + right
			})
		}, nil
	case "-":
		return func(lazyLeft, lazyRight Evaluatable) (RuntimeValue, *RuntimeError) {
			return ex.numericGreedyOperation(operator, lazyLeft, lazyRight, func(left, right PreludeInt) PreludeInt {
				return left - right
			}, func(left, right PreludeFloat) PreludeFloat {
				return left - right
			})
		}, nil
	case "*":
		return func(lazyLeft, lazyRight Evaluatable) (RuntimeValue, *RuntimeError) {
			return ex.numericGreedyOperation(operator, lazyLeft, lazyRight, func(left, right PreludeInt) PreludeInt {
				return left * right
			}, func(left, right PreludeFloat) PreludeFloat {
				return left * right
			})
		}, nil
	case "/":
		return func(lazyLeft, lazyRight Evaluatable) (RuntimeValue, *RuntimeError) {
			return ex.numericGreedyOperation(operator, lazyLeft, lazyRight, func(left, right PreludeInt) PreludeInt {
				return left / right
			}, func(left, right PreludeFloat) PreludeFloat {
				return left / right
			})
		}, nil
	default:
		return nil, NewRuntimeErrorf("unknown binary operator: %s", operator)
	}
}

func (ex *InterpreterContext) genericGreedyComparision(
	lazyLeft, lazyRight Evaluatable,
	compare func(RuntimeValue, RuntimeValue) bool,
) (RuntimeValue, *RuntimeError) {
	left, err := lazyLeft.Evaluate()
	if err != nil {
		return nil, err
	}
	right, err := lazyRight.Evaluate()
	if err != nil {
		return nil, err
	}
	return ex.boolToRuntimeValue(compare(left, right))
}

func (ex *InterpreterContext) numericGreedyComparision(
	operator string,
	lazyLeft, lazyRight Evaluatable,
	compareInt func(PreludeInt, PreludeInt) bool,
	compareFloat func(PreludeFloat, PreludeFloat) bool,
) (RuntimeValue, *RuntimeError) {
	left, err := lazyLeft.Evaluate()
	if err != nil {
		return nil, err
	}
	switch left := left.(type) {
	case PreludeInt:
		right, err := lazyRight.Evaluate()
		if err != nil {
			return nil, err
		}
		switch right := right.(type) {
		case PreludeInt:
			return ex.boolToRuntimeValue(compareInt(left, right))
		case PreludeFloat:
			return ex.boolToRuntimeValue(compareFloat(PreludeFloat(left), right))
		default:
			return nil, RuntimeBinaryOperatorOnlySupportsType(
				operator,
				[]RuntimeTypeRef{PreludeInt(0).RuntimeType(), PreludeFloat(0).RuntimeType()},
				left,
			)
		}
	case PreludeFloat:
		right, err := lazyRight.Evaluate()
		if err != nil {
			return nil, err
		}
		switch right := right.(type) {
		case PreludeInt:
			return ex.boolToRuntimeValue(compareFloat(left, PreludeFloat(right)))
		case PreludeFloat:
			return ex.boolToRuntimeValue(compareFloat(left, right))
		default:
			return nil, RuntimeBinaryOperatorOnlySupportsType(
				operator,
				[]RuntimeTypeRef{PreludeInt(0).RuntimeType(), PreludeFloat(0).RuntimeType()},
				left,
			)
		}
	default:
		return nil, RuntimeBinaryOperatorOnlySupportsType(
			operator,
			[]RuntimeTypeRef{PreludeInt(0).RuntimeType(), PreludeFloat(0).RuntimeType()},
			left,
		)
	}
}

func (ex *InterpreterContext) numericGreedyOperation(
	operator string,
	lazyLeft, lazyRight Evaluatable,
	combineInt func(PreludeInt, PreludeInt) PreludeInt,
	combineFloat func(PreludeFloat, PreludeFloat) PreludeFloat,
) (RuntimeValue, *RuntimeError) {
	left, err := lazyLeft.Evaluate()
	if err != nil {
		return nil, err
	}
	switch left := left.(type) {
	case PreludeInt:
		right, err := lazyRight.Evaluate()
		if err != nil {
			return nil, err
		}
		switch right := right.(type) {
		case PreludeInt:
			return combineInt(left, right), nil
		case PreludeFloat:
			return combineFloat(PreludeFloat(left), right), nil
		default:
			return nil, RuntimeBinaryOperatorOnlySupportsType(
				operator,
				[]RuntimeTypeRef{PreludeInt(0).RuntimeType(), PreludeFloat(0).RuntimeType()},
				left,
			)
		}
	case PreludeFloat:
		right, err := lazyRight.Evaluate()
		if err != nil {
			return nil, err
		}
		switch right := right.(type) {
		case PreludeInt:
			return combineFloat(left, PreludeFloat(right)), nil
		case PreludeFloat:
			return combineFloat(left, right), nil
		default:
			return nil, RuntimeBinaryOperatorOnlySupportsType(
				operator,
				[]RuntimeTypeRef{PreludeInt(0).RuntimeType(), PreludeFloat(0).RuntimeType()},
				left,
			)
		}
	default:
		return nil, RuntimeBinaryOperatorOnlySupportsType(
			operator,
			[]RuntimeTypeRef{PreludeInt(0).RuntimeType(), PreludeFloat(0).RuntimeType()},
			left,
		)
	}
}

func (ex *InterpreterContext) lazyLogicComparision(
	operator string,
	lazyLeft, lazyRight Evaluatable,
	compare func(bool, func() (bool, *RuntimeError)) (bool, *RuntimeError),
) (RuntimeValue, *RuntimeError) {
	var err *RuntimeError
	left, err := lazyLeft.Evaluate()
	if err != nil {
		return nil, NewRuntimeError(err)
	}
	boolRef := MakeRuntimeTypeRef("Bool", "prelude")
	// trueTypeValue, err := ex.Environment.GetEvaluatedRuntimeValue("True")
	// if err != nil {
	// 	return nil, err
	// }
	// // trueType := RuntimeType{
	// // 	name:       "True",
	// // 	moduleName: "prelude",
	// // 	typeValue:  &trueTypeValue,
	// // }
	if ok, err := boolRef.HasInstance(ex.interpreter, left); !ok || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, RuntimeBinaryOperatorOnlySupportsType(
			operator,
			[]RuntimeTypeRef{boolRef},
			left,
		)
	}

	trueRef := MakeRuntimeTypeRef("True", "prelude")
	isLeftTrue, err := trueRef.HasInstance(ex.interpreter, left)
	if err != nil {
		return nil, NewRuntimeError(err)
	}

	isComparisionTrue, err := compare(isLeftTrue, func() (bool, *RuntimeError) {
		right, err := lazyRight.Evaluate()
		if err != nil {
			return false, nil
		}
		isRightTrue, err := boolRef.HasInstance(ex.interpreter, right)
		if err != nil {
			return false, nil
		}
		return isRightTrue, nil
	})

	if err != nil {
		return nil, err
	}

	return ex.boolToRuntimeValue(isComparisionTrue)
}

func (env *Environment) boolToRuntimeValue(value bool) (RuntimeValue, *RuntimeError) {
	if value {
		return env.MakeEmptyDataRuntimeValue("True")
	} else {
		return env.MakeEmptyDataRuntimeValue("False")
	}
}

func (ex *InterpreterContext) boolToRuntimeValue(value bool) (RuntimeValue, *RuntimeError) {
	return ex.environment.boolToRuntimeValue(value)
}
