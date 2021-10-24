package interpreter

import "fmt"

var _ RuntimeValue = TypeExpression{}
var _ Callable = TypeExpression{}

type TypeExpression struct {
	typeValue EnumDeclRuntimeValue
	caseNames []string
	cases     []Evaluatable
}

func (TypeExpression) RuntimeType() RuntimeType {
	return PreludeFunctionType{}.RuntimeType()
}
func (t TypeExpression) String() string {
	return fmt.Sprintf("{ value => type %s }", t.typeValue.name)
}

func (t TypeExpression) Lookup(member string) (Evaluatable, error) {
	return nil, fmt.Errorf("function %s has no member %s", fmt.Sprint(t), member)
}

func (typeExpr TypeExpression) Call(arguments []Evaluatable) (RuntimeValue, error) {
	if len(arguments) == 0 {
		return typeExpr, nil
	}
	lazyValueArgument := arguments[0]
	valueArgument, err := lazyValueArgument.Evaluate()
	if err != nil {
		return nil, err
	}
	for i, caseName := range typeExpr.caseNames {
		lazyCaseImpl := typeExpr.cases[i]
		caseTypeDef := typeExpr.typeValue.cases[caseName]
		if caseTypeDef == nil && caseName == "Any" {
			// always ok
		} else if caseTypeDef == nil {
			return nil, fmt.Errorf("case %s not defined", caseName)
		} else {
			caseTypeValue, err := caseTypeDef.Evaluate()
			if err != nil {
				return nil, err
			}
			ok, err := RuntimeTypeValueIncludesValue(caseTypeValue, valueArgument)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
		}

		intermediate, err := lazyCaseImpl.Evaluate()
		if err != nil {
			return nil, err
		}
		callable, ok := intermediate.(Callable)
		if !ok {
			return nil, fmt.Errorf("case %s is not callable", caseName)
		}
		return callable.Call(arguments)
	}
	return nil, fmt.Errorf("no %s has no matching case for value %s of type %s", typeExpr.typeValue.name, fmt.Sprint(valueArgument), fmt.Sprint(valueArgument.RuntimeType().name))
}

func contains(names []string, name string) bool {
	for _, n := range names {
		if n == name {
			return true
		}
	}
	return false
}
