package interpreter

import (
	"fmt"
	"strings"
)

func (ex *EvaluationContext) RuntimeErrorf(format string, args ...interface{}) LocatableError {
	return ex.LocatableErrorf("runtime error", format, args...)
}

func (ex *EvaluationContext) RuntimeNonExhaustiveTypeExpression(
	enumDecl EnumDeclRuntimeValue,
	caseNames []string,
) LocatableError {
	missing := []string{}

	for caseName := range enumDecl.cases {
		if !contains(caseNames, caseName) {
			missing = append(missing, caseName)
		}
	}
	return ex.RuntimeErrorf("non-exhaustive type expression for %s: missing [%s]", enumDecl.name, strings.Join(missing, ", "))
}

func (ex *EvaluationContext) RuntimeInvalidCaseTypeExpression(
	enumDecl EnumDeclRuntimeValue,
	caseNames []string,
) LocatableError {
	invalids := []string{}
	for _, caseName := range caseNames {
		if _, ok := enumDecl.cases[caseName]; !ok && caseName != "Any" {
			invalids = append(invalids, caseName)
		}
	}
	return ex.RuntimeErrorf("invalid type expression for %s: invalid [%s]", enumDecl.name, strings.Join(invalids, ", "))
}

func (ex *EvaluationContext) RuntimeCannotCallNonFunction(nonFunction RuntimeValue, args []Evaluatable) LocatableError {
	stringifiedArgs := []string{}
	for _, lazyArg := range args {
		arg, err := lazyArg.Evaluate()
		if err != nil {
			stringifiedArgs = append(stringifiedArgs, "<error>")
		} else {
			stringifiedArgs = append(stringifiedArgs, fmt.Sprintf("%q", arg.String()))
		}
	}
	return ex.RuntimeErrorf("cannot call non-function value of type %s: %q with args %s", nonFunction.RuntimeType().name, nonFunction.String(), strings.Join(stringifiedArgs, ", "))
}

func (ex *EvaluationContext) RuntimeBinaryOperatorOnlySupportsType(operator string, supportedTypes []RuntimeType, gotValue RuntimeValue) LocatableError {
	supportedTypesNames := []string{}
	for _, supported := range supportedTypes {
		supportedTypesNames = append(supportedTypesNames, supported.name)
	}
	return ex.RuntimeErrorf(
		"binary operator %q only supports %s; value of type %s given: %q",
		operator,
		strings.Join(supportedTypesNames, ", "),
		gotValue.RuntimeType().name,
		gotValue.String(),
	)
}
