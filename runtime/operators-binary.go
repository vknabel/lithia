package runtime

import "reflect"

func (ex *InterpreterContext) BinaryOperatorFunction(operator string) (func(Evaluatable, Evaluatable) (RuntimeValue, *RuntimeError), *RuntimeError) {
	switch operator {
	case "==":
		return func(lazyLeft, lazyRight Evaluatable) (RuntimeValue, *RuntimeError) {
			return ex.genericGreedyComparision(lazyLeft, lazyRight, func(left, right RuntimeValue) (bool, *RuntimeError) {
				return ex.DeepEqual(left, right)
			})
		}, nil
	case "!=":
		return func(lazyLeft, lazyRight Evaluatable) (RuntimeValue, *RuntimeError) {
			return ex.genericGreedyComparision(lazyLeft, lazyRight, func(left, right RuntimeValue) (bool, *RuntimeError) {
				result, err := ex.DeepEqual(left, right)
				if err != nil {
					return false, err
				}
				return !result, nil
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
	compare func(RuntimeValue, RuntimeValue) (bool, *RuntimeError),
) (RuntimeValue, *RuntimeError) {
	left, err := lazyLeft.Evaluate()
	if err != nil {
		return nil, err
	}
	right, err := lazyRight.Evaluate()
	if err != nil {
		return nil, err
	}
	flag, err := compare(left, right)
	if err != nil {
		return nil, err
	}
	return ex.boolToRuntimeValue(flag)
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
			return nil, ReportBinaryOperatorOnlySupportsType(
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
			return nil, ReportBinaryOperatorOnlySupportsType(
				operator,
				[]RuntimeTypeRef{PreludeInt(0).RuntimeType(), PreludeFloat(0).RuntimeType()},
				left,
			)
		}
	default:
		return nil, ReportBinaryOperatorOnlySupportsType(
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
			return nil, ReportBinaryOperatorOnlySupportsType(
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
			return nil, ReportBinaryOperatorOnlySupportsType(
				operator,
				[]RuntimeTypeRef{PreludeInt(0).RuntimeType(), PreludeFloat(0).RuntimeType()},
				left,
			)
		}
	default:
		return nil, ReportBinaryOperatorOnlySupportsType(
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
	if ok, err := boolRef.HasInstance(left); !ok || err != nil {
		if err != nil {
			return nil, err
		}
		return nil, ReportBinaryOperatorOnlySupportsType(
			operator,
			[]RuntimeTypeRef{boolRef},
			left,
		)
	}

	trueRef := MakeRuntimeTypeRef("True", "prelude")
	isLeftTrue, err := trueRef.HasInstance(left)
	if err != nil {
		return nil, NewRuntimeError(err)
	}

	isComparisionTrue, err := compare(isLeftTrue, func() (bool, *RuntimeError) {
		right, err := lazyRight.Evaluate()
		if err != nil {
			return false, nil
		}
		isRightTrue, err := trueRef.HasInstance(right)
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

func (ex *InterpreterContext) DeepEqual(left, right RuntimeValue) (bool, *RuntimeError) {
	// TODO: should be moved to RuntimeValue
	switch left := left.(type) {
	case DataRuntimeValue:
		right, ok := right.(DataRuntimeValue)
		if !ok {
			return false, nil
		}
		if len(left.Members) != len(right.Members) {
			return false, nil
		}
		ok, err := ex.DeepEqual(left.TypeDecl, right.TypeDecl)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
		for memberName, lazyLeftMemeberValue := range left.Members {
			lazyRightMemberValue, ok := right.Members[memberName]
			if !ok {
				return false, nil
			}
			leftMemberValue, err := lazyLeftMemeberValue.Evaluate()
			if err != nil {
				return false, err
			}
			rightMemberValue, err := lazyRightMemberValue.Evaluate()
			if err != nil {
				return false, err
			}

			areEqual, err := ex.DeepEqual(leftMemberValue, rightMemberValue)
			if err != nil {
				return false, err
			}
			if !areEqual {
				return false, nil
			}
		}
		return true, nil
	default:
		return reflect.DeepEqual(left, right), nil
	}
}
