package interpreter

import "reflect"

func (ex *EvaluationContext) BinaryOperatorFunction(operator string) (func(*LazyRuntimeValue, *LazyRuntimeValue) (RuntimeValue, error), error) {
	switch operator {
	case "==":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, error) {
			return ex.genericGreedyComparision(lazyLeft, lazyRight, func(left, right RuntimeValue) bool {
				return reflect.DeepEqual(left, right)
			})
		}, nil
	case "!=":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, error) {
			return ex.genericGreedyComparision(lazyLeft, lazyRight, func(left, right RuntimeValue) bool {
				return !reflect.DeepEqual(left, right)
			})
		}, nil
	case "&&":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, error) {
			return ex.lazyLogicComparision(lazyLeft, lazyRight, func(left bool, right func() (bool, error)) (bool, error) {
				if !left {
					return false, nil
				} else {
					return right()
				}
			})
		}, nil
	case "||":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, error) {
			return ex.lazyLogicComparision(lazyLeft, lazyRight, func(left bool, right func() (bool, error)) (bool, error) {
				if left {
					return true, nil
				} else {
					return right()
				}
			})
		}, nil
	case ">":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, error) {
			return ex.numericGreedyComparision(lazyLeft, lazyRight, func(left, right PreludeInt) bool {
				return left > right
			}, func(left, right PreludeFloat) bool {
				return left > right
			})
		}, nil
	case ">=":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, error) {
			return ex.numericGreedyComparision(lazyLeft, lazyRight, func(left, right PreludeInt) bool {
				return left >= right
			}, func(left, right PreludeFloat) bool {
				return left >= right
			})
		}, nil
	case "<":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, error) {
			return ex.numericGreedyComparision(lazyLeft, lazyRight, func(left, right PreludeInt) bool {
				return left < right
			}, func(left, right PreludeFloat) bool {
				return left < right
			})
		}, nil
	case "<=":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, error) {
			return ex.numericGreedyComparision(lazyLeft, lazyRight, func(left, right PreludeInt) bool {
				return left <= right
			}, func(left, right PreludeFloat) bool {
				return left <= right
			})
		}, nil
	case "+":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, error) {
			return ex.numericGreedyOperation(lazyLeft, lazyRight, func(left, right PreludeInt) PreludeInt {
				return left + right
			}, func(left, right PreludeFloat) PreludeFloat {
				return left + right
			})
		}, nil
	case "-":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, error) {
			return ex.numericGreedyOperation(lazyLeft, lazyRight, func(left, right PreludeInt) PreludeInt {
				return left - right
			}, func(left, right PreludeFloat) PreludeFloat {
				return left - right
			})
		}, nil
	case "*":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, error) {
			return ex.numericGreedyOperation(lazyLeft, lazyRight, func(left, right PreludeInt) PreludeInt {
				return left * right
			}, func(left, right PreludeFloat) PreludeFloat {
				return left * right
			})
		}, nil
	case "/":
		return func(lazyLeft, lazyRight *LazyRuntimeValue) (RuntimeValue, error) {
			return ex.numericGreedyOperation(lazyLeft, lazyRight, func(left, right PreludeInt) PreludeInt {
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
) (RuntimeValue, error) {
	left, err := lazyLeft.Evaluate()
	if err != nil {
		return nil, err
	}
	right, err := lazyRight.Evaluate()
	if err != nil {
		return nil, err
	}
	if compare(left, right) {
		return ex.environment.MakeEmptyDataRuntimeValue("True")
	} else {
		return ex.environment.MakeEmptyDataRuntimeValue("False")
	}
}

func (ex *EvaluationContext) numericGreedyComparision(
	lazyLeft, lazyRight *LazyRuntimeValue,
	compareInt func(PreludeInt, PreludeInt) bool,
	compareFloat func(PreludeFloat, PreludeFloat) bool,
) (RuntimeValue, error) {
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
			return ex.environment.boolToRuntimeValue(compareInt(left, right))
		case PreludeFloat:
			return ex.environment.boolToRuntimeValue(compareFloat(PreludeFloat(left), right))
		default:
			return nil, ex.RuntimeErrorf("expected numeric, got %d", left)
		}
	case PreludeFloat:
		right, err := lazyRight.Evaluate()
		if err != nil {
			return nil, err
		}
		switch right := right.(type) {
		case PreludeInt:
			return ex.environment.boolToRuntimeValue(compareFloat(left, PreludeFloat(right)))
		case PreludeFloat:
			return ex.environment.boolToRuntimeValue(compareFloat(left, right))
		default:
			return nil, ex.RuntimeErrorf("expected numeric, got %f", left)
		}
	default:
		return nil, ex.RuntimeErrorf("expected numeric, got %s", left)
	}
}

func (ex *EvaluationContext) numericGreedyOperation(
	lazyLeft, lazyRight *LazyRuntimeValue,
	combineInt func(PreludeInt, PreludeInt) PreludeInt,
	combineFloat func(PreludeFloat, PreludeFloat) PreludeFloat,
) (RuntimeValue, error) {
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
			return nil, ex.RuntimeErrorf("expected numeric, got %d", left)
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
			return nil, ex.RuntimeErrorf("expected numeric, got %f", left)
		}
	default:
		return nil, ex.RuntimeErrorf("expected numeric, got %s", left)
	}
}

func (ex *EvaluationContext) lazyLogicComparision(
	lazyLeft, lazyRight *LazyRuntimeValue,
	compare func(bool, func() (bool, error)) (bool, error),
) (RuntimeValue, error) {
	left, err := lazyLeft.Evaluate()
	if err != nil {
		return nil, err
	}
	boolTypeValue, err := ex.environment.LookupRuntimeValue("Bool")
	boolType := RuntimeType{
		name:       "Bool",
		modulePath: []string{"prelude"},
		typeValue:  &boolTypeValue,
	}
	if err != nil {
		return nil, err
	}
	trueTypeValue, err := ex.environment.LookupRuntimeValue("True")
	if err != nil {
		return nil, err
	}
	trueType := RuntimeType{
		name:       "True",
		modulePath: []string{"prelude"},
		typeValue:  &trueTypeValue,
	}
	if ok, err := boolType.IncludesValue(left); !ok || err != nil {
		return nil, ex.RuntimeErrorf("expected bool, got %s", left)
	}

	isLeftTrue, err := trueType.IncludesValue(left)
	if err != nil {
		return nil, err
	}

	isComparisionTrue, err := compare(isLeftTrue, func() (bool, error) {
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
		return nil, err
	}

	return ex.environment.boolToRuntimeValue(isComparisionTrue)
}

func (env *Environment) boolToRuntimeValue(value bool) (RuntimeValue, error) {
	if value {
		return env.MakeEmptyDataRuntimeValue("True")
	} else {
		return env.MakeEmptyDataRuntimeValue("False")
	}
}
