package runtime

import (
	"fmt"
	"strings"
)

// func RuntimeNonExhaustiveTypeExpression(
// 	enumDecl EnumDeclRuntimeValue,
// 	caseNames []string,
// ) *RuntimeError {
// 	missing := []string{}

// 	for caseName := range enumDecl.cases {
// 		if !contains(caseNames, caseName) {
// 			missing = append(missing, caseName)
// 		}
// 	}
// 	return ex.RuntimeErrorf("non-exhaustive type expression for %s: missing [%s]", enumDecl.name, strings.Join(missing, ", "))
// }

// func RuntimeInvalidCaseTypeExpression(
// 	enumDecl EnumDeclRuntimeValue,
// 	caseNames []string,
// ) *RuntimeError {
// 	invalids := []string{}
// 	for _, caseName := range caseNames {
// 		if _, ok := enumDecl.cases[caseName]; !ok && caseName != "Any" {
// 			invalids = append(invalids, caseName)
// 		}
// 	}
// 	return ex.RuntimeErrorf("invalid type expression for %s: invalid [%s]", enumDecl.name, strings.Join(invalids, ", "))
// }

func RuntimeCannotCallNonFunction(nonFunction RuntimeValue, args []Evaluatable) *RuntimeError {
	stringifiedArgs := []string{}
	for _, lazyArg := range args {
		arg, err := lazyArg.Evaluate()
		if err != nil {
			stringifiedArgs = append(stringifiedArgs, "<error>")
		} else {
			stringifiedArgs = append(stringifiedArgs, fmt.Sprintf("%q", arg.String()))
		}
	}
	return NewRuntimeErrorf("cannot call non-function value of type %s: %q with args %s", nonFunction.RuntimeType().Name, nonFunction.String(), strings.Join(stringifiedArgs, ", "))
}

func RuntimeBinaryOperatorOnlySupportsType(operator string, supportedTypes []RuntimeTypeRef, gotValue RuntimeValue) *RuntimeError {
	supportedTypesNames := []string{}
	for _, supported := range supportedTypes {
		supportedTypesNames = append(supportedTypesNames, string(supported.Name))
	}
	return NewRuntimeErrorf(
		"binary operator %q only supports %s; value of type %s given: %q",
		operator,
		strings.Join(supportedTypesNames, ", "),
		gotValue.RuntimeType().Name,
		gotValue.String(),
	)
}
