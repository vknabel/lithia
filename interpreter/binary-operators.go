package interpreter

import "reflect"

func (ex *EvaluationContext) BinaryOperatorFunction(operator string) (func(*LazyRuntimeValue, *LazyRuntimeValue) (RuntimeValue, LocatableError), LocatableError) {
	switch operator {
	case "==":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, LocatableError) {
			return ex.genericGreedyComparision(lazyLeft, lazyRight, func(left, right RuntimeValue) bool {
				return reflect.DeepEqual(left, right)
			})
		}, nil
	case "!=":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, LocatableError) {
			return ex.genericGreedyComparision(lazyLeft, lazyRight, func(left, right RuntimeValue) bool {
				return !reflect.DeepEqual(left, right)
			})
		}, nil
	case "&&":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, LocatableError) {
			return ex.lazyLogicComparision(operator, lazyLeft, lazyRight, func(left bool, right func() (bool, LocatableError)) (bool, LocatableError) {
				if !left {
					return false, nil
				} else {
					return right()
				}
			})
		}, nil
	case "||":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, LocatableError) {
			return ex.lazyLogicComparision(operator, lazyLeft, lazyRight, func(left bool, right func() (bool, LocatableError)) (bool, LocatableError) {
				if left {
					return true, nil
				} else {
					return right()
				}
			})
		}, nil
	case ">":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, LocatableError) {
			return ex.numericGreedyComparision(operator, lazyLeft, lazyRight, func(left, right PreludeInt) bool {
				return left > right
			}, func(left, right PreludeFloat) bool {
				return left > right
			})
		}, nil
	case ">=":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, LocatableError) {
			return ex.numericGreedyComparision(operator, lazyLeft, lazyRight, func(left, right PreludeInt) bool {
				return left >= right
			}, func(left, right PreludeFloat) bool {
				return left >= right
			})
		}, nil
	case "<":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, LocatableError) {
			return ex.numericGreedyComparision(operator, lazyLeft, lazyRight, func(left, right PreludeInt) bool {
				return left < right
			}, func(left, right PreludeFloat) bool {
				return left < right
			})
		}, nil
	case "<=":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, LocatableError) {
			return ex.numericGreedyComparision(operator, lazyLeft, lazyRight, func(left, right PreludeInt) bool {
				return left <= right
			}, func(left, right PreludeFloat) bool {
				return left <= right
			})
		}, nil
	case "+":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, LocatableError) {
			return ex.numericGreedyOperation(operator, lazyLeft, lazyRight, func(left, right PreludeInt) PreludeInt {
				return left + right
			}, func(left, right PreludeFloat) PreludeFloat {
				return left + right
			})
		}, nil
	case "-":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, LocatableError) {
			return ex.numericGreedyOperation(operator, lazyLeft, lazyRight, func(left, right PreludeInt) PreludeInt {
				return left - right
			}, func(left, right PreludeFloat) PreludeFloat {
				return left - right
			})
		}, nil
	case "*":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, LocatableError) {
			return ex.numericGreedyOperation(operator, lazyLeft, lazyRight, func(left, right PreludeInt) PreludeInt {
				return left * right
			}, func(left, right PreludeFloat) PreludeFloat {
				return left * right
			})
		}, nil
	case "/":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, LocatableError) {
			return ex.numericGreedyOperation(operator, lazyLeft, lazyRight, func(left, right PreludeInt) PreludeInt {
				return left / right
			}, func(left, right PreludeFloat) PreludeFloat {
				return left / right
			})
		}, nil
	default:
		return nil, ex.SyntaxErrorf("unknown binary operator: %s", operator)
	}
}

func (ex *EvaluationContext) genericGreedyComparision(
	lazyLeft, lazyRight *LazyRuntimeValue,
	compare func(RuntimeValue, RuntimeValue) bool,
) (RuntimeValue, LocatableError) {
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

func (ex *EvaluationContext) numericGreedyComparision(
	operator string,
	lazyLeft, lazyRight *LazyRuntimeValue,
	compareInt func(PreludeInt, PreludeInt) bool,
	compareFloat func(PreludeFloat, PreludeFloat) bool,
) (RuntimeValue, LocatableError) {
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
			return nil, ex.RuntimeBinaryOperatorOnlySupportsType(
				operator,
				[]RuntimeType{PreludeInt(0).RuntimeType(), PreludeFloat(0).RuntimeType()},
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
			return nil, ex.RuntimeBinaryOperatorOnlySupportsType(
				operator,
				[]RuntimeType{PreludeInt(0).RuntimeType(), PreludeFloat(0).RuntimeType()},
				left,
			)
		}
	default:
		return nil, ex.RuntimeBinaryOperatorOnlySupportsType(
			operator,
			[]RuntimeType{PreludeInt(0).RuntimeType(), PreludeFloat(0).RuntimeType()},
			left,
		)
	}
}

func (ex *EvaluationContext) numericGreedyOperation(
	operator string,
	lazyLeft, lazyRight *LazyRuntimeValue,
	combineInt func(PreludeInt, PreludeInt) PreludeInt,
	combineFloat func(PreludeFloat, PreludeFloat) PreludeFloat,
) (RuntimeValue, LocatableError) {
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
			return nil, ex.RuntimeBinaryOperatorOnlySupportsType(
				operator,
				[]RuntimeType{PreludeInt(0).RuntimeType(), PreludeFloat(0).RuntimeType()},
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
			return nil, ex.RuntimeBinaryOperatorOnlySupportsType(
				operator,
				[]RuntimeType{PreludeInt(0).RuntimeType(), PreludeFloat(0).RuntimeType()},
				left,
			)
		}
	default:
		return nil, ex.RuntimeBinaryOperatorOnlySupportsType(
			operator,
			[]RuntimeType{PreludeInt(0).RuntimeType(), PreludeFloat(0).RuntimeType()},
			left,
		)
	}
}

func (ex *EvaluationContext) lazyLogicComparision(
	operator string,
	lazyLeft, lazyRight *LazyRuntimeValue,
	compare func(bool, func() (bool, LocatableError)) (bool, LocatableError),
) (RuntimeValue, LocatableError) {
	var err error
	left, err := lazyLeft.Evaluate()
	if err != nil {
		return nil, ex.LocatableErrorOrConvert(err)
	}
	boolTypeValue, err := ex.environment.LookupRuntimeValue("Bool")
	boolType := RuntimeType{
		name:       "Bool",
		moduleName: "prelude",
		typeValue:  &boolTypeValue,
	}
	if err != nil {
		return nil, ex.LocatableErrorOrConvert(err)
	}
	trueTypeValue, err := ex.environment.LookupRuntimeValue("True")
	if err != nil {
		return nil, ex.LocatableErrorOrConvert(err)
	}
	trueType := RuntimeType{
		name:       "True",
		moduleName: "prelude",
		typeValue:  &trueTypeValue,
	}
	if ok, err := boolType.IncludesValue(left); !ok || err != nil {
		return nil, ex.RuntimeBinaryOperatorOnlySupportsType(
			operator,
			[]RuntimeType{boolType},
			left,
		)
	}

	isLeftTrue, err := trueType.IncludesValue(left)
	if err != nil {
		return nil, ex.LocatableErrorOrConvert(err)
	}

	isComparisionTrue, err := compare(isLeftTrue, func() (bool, LocatableError) {
		right, err := lazyRight.Evaluate()
		if err != nil {
			return false, nil
		}
		isRightTrue, err := boolType.IncludesValue(right)
		if err != nil {
			return false, nil
		}
		return isRightTrue, nil
	})

	if err != nil {
		return nil, ex.LocatableErrorOrConvert(err)
	}

	return ex.boolToRuntimeValue(isComparisionTrue)
}

func (env *Environment) boolToRuntimeValue(value bool) (RuntimeValue, error) {
	if value {
		return env.MakeEmptyDataRuntimeValue("True")
	} else {
		return env.MakeEmptyDataRuntimeValue("False")
	}
}

func (ex *EvaluationContext) boolToRuntimeValue(value bool) (RuntimeValue, LocatableError) {
	runtimeBool, err := ex.environment.boolToRuntimeValue(value)
	return runtimeBool, ex.LocatableErrorOrConvert(err)
}