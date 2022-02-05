package runtime

import (
	"fmt"
	"strings"

	"github.com/vknabel/lithia/ast"
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

func ReportCannotCallNonFunction(nonFunction RuntimeValue, args []Evaluatable) *RuntimeError {
	stringifiedArgs := []string{}
	for _, lazyArg := range args {
		arg, err := lazyArg.Evaluate()
		if err != nil {
			stringifiedArgs = append(stringifiedArgs, "<error>")
		} else {
			stringifiedArgs = append(stringifiedArgs, fmt.Sprintf("%q", arg.String()))
		}
	}
	return NewTypeErrorf("cannot call non-function value of type %s: %q with args %s", nonFunction.RuntimeType().Name, nonFunction.String(), strings.Join(stringifiedArgs, ", "))
}

func ReportBinaryOperatorOnlySupportsType(operator string, supportedTypes []RuntimeTypeRef, gotValue RuntimeValue) *RuntimeError {
	supportedTypesNames := []string{}
	for _, supported := range supportedTypes {
		supportedTypesNames = append(supportedTypesNames, string(supported.Name))
	}
	return NewTypeErrorf(
		"binary operator %q only supports %s; value given: %q",
		operator,
		strings.Join(supportedTypesNames, ", "),
		gotValue.String(),
	)
}

func ReportNonExhaustiveTypeSwitch(
	enumDecl PreludeEnumDecl,
	typeSwitch ast.ExprTypeSwitch,
) *RuntimeError {
	missing := []string{}

	declCaseNames := make([]string, len(enumDecl.Decl.Cases))
	for i, caseDecl := range enumDecl.Decl.Cases {
		declCaseNames[i] = string(caseDecl.Name)
	}

	switchCaseNames := make([]string, len(typeSwitch.CaseOrder))
	for i, caseDecl := range typeSwitch.CaseOrder {
		switchCaseNames[i] = string(caseDecl)
	}

	for _, caseName := range declCaseNames {
		if !contains(switchCaseNames, caseName) {
			missing = append(missing, caseName)
		}
	}

	return NewTypeErrorf("non-exhaustive type expression for %s: missing [%s]", enumDecl.Decl.Name, strings.Join(missing, ", "))
}

func contains(names []string, name string) bool {
	for _, n := range names {
		if n == name {
			return true
		}
	}
	return false
}
