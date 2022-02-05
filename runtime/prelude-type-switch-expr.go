package runtime

import (
	"fmt"
	"strings"

	"github.com/vknabel/lithia/ast"
)

var _ RuntimeValue = PreludeTypeSwitchExpr{}
var _ CallableRuntimeValue = PreludeTypeSwitchExpr{}

type PreludeTypeSwitchExpr struct {
	Decl ast.ExprTypeSwitch

	context      *InterpreterContext
	enumDefValue Evaluatable
	caseValue    map[ast.Identifier]Evaluatable
}

func MakePreludeTypeSwitchExpr(context *InterpreterContext, decl ast.ExprTypeSwitch) PreludeTypeSwitchExpr {
	caseValue := make(map[ast.Identifier]Evaluatable)
	for _, identifier := range decl.CaseOrder {
		caseDef := decl.Cases[identifier]
		caseValue[identifier] = MakeEvaluatableExpr(context, caseDef)
	}
	return PreludeTypeSwitchExpr{decl, context, MakeEvaluatableExpr(context, decl.Type), caseValue}
}

func (f PreludeTypeSwitchExpr) Lookup(member string) (Evaluatable, *RuntimeError) {
	switch member {
	case "arity":
		return NewConstantRuntimeValue(PreludeInt(f.Arity())), nil
	default:
		return nil, NewRuntimeErrorf("no such member: %s", member)
	}
}

func (PreludeTypeSwitchExpr) RuntimeType() RuntimeTypeRef {
	return PreludeFunctionTypeRef
}

func (t PreludeTypeSwitchExpr) String() string {
	value, err := t.enumDefValue.Evaluate()
	if err != nil {
		panic(fmt.Sprintf("error: %s", err))
	}
	return fmt.Sprintf("<type %s.type %s>", strings.Join(t.context.path, "."), value.RuntimeType().Name)
}

func (t PreludeTypeSwitchExpr) Arity() int {
	return 1
}

func (t PreludeTypeSwitchExpr) Call(args []Evaluatable, fromExpr ast.Expr) (RuntimeValue, *RuntimeError) {
	if len(args) != t.Arity() {
		panic("use Call to call functions!")
	}
	primaryArg, err := args[0].Evaluate()
	if err != nil {
		return nil, err.CascadeExpr(t.Decl)
	}

	// TODO: optimization can be cached
	enumDefValue, err := t.enumDefValue.Evaluate()
	if err != nil {
		return nil, err.CascadeExpr(t.Decl)
	}
	enumDef, ok := enumDefValue.(PreludeEnumDecl)
	if !ok {
		return nil, NewRuntimeErrorf(
			"type switch requires enum type, got %s: %s",
			enumDefValue.RuntimeType(),
			enumDefValue,
		).CascadeExpr(t.Decl)
	}

	err = validateTypeSwitchAgainstEnumDecl(enumDef, t.Decl)
	if err != nil {
		return nil, err.CascadeExpr(t.Decl)
	}

	// TODO: optimization can be cached
	for _, caseIdentifier := range t.Decl.CaseOrder {
		if caseIdentifier == "Any" {
			intermediate, err := t.caseValue[caseIdentifier].Evaluate()
			if err != nil {
				return nil, err.CascadeExpr(t.Decl)
			}
			fun, ok := intermediate.(CallableRuntimeValue)
			if !ok {
				return nil, NewRuntimeErrorf("cannot call non function %T", intermediate).CascadeExpr(t.Decl)
			}
			return Call(fun, args, t.Decl.Cases[caseIdentifier])
		}
		lazyCaseValue, err := enumDef.LookupCase(string(caseIdentifier))
		if err != nil {
			return nil, err.CascadeExpr(t.Decl)
		}
		caseValue, err := lazyCaseValue.Evaluate()
		if err != nil {
			return nil, err.CascadeExpr(t.Decl)
		}
		runtimeDecl, ok := caseValue.(DeclRuntimeValue)
		if !ok {
			return nil, NewRuntimeErrorf(
				"case %s is not a type in %s",
				caseIdentifier,
				enumDefValue,
			).CascadeExpr(t.Decl)
		}

		ok, err = runtimeDecl.HasInstance(t.context.interpreter, primaryArg)
		if err != nil {
			return nil, err.CascadeExpr(t.Decl)
		}
		if !ok {
			continue
		}
		fun, err := t.caseValue[caseIdentifier].Evaluate()
		if err != nil {
			return nil, err.CascadeExpr(t.Decl)
		}
		return Call(fun, args, t.Decl.Cases[caseIdentifier])
	}
	return nil, NewRuntimeErrorf("no matching case %s", primaryArg).CascadeExpr(t.Decl)
}

func (f PreludeTypeSwitchExpr) Source() *ast.Source {
	return f.Decl.Meta().Source
}

func validateTypeSwitchAgainstEnumDecl(enumDecl PreludeEnumDecl, typeSwitch ast.ExprTypeSwitch) *RuntimeError {
	err := validateTypeSwitchAgainstEnumDeclForCount(enumDecl, typeSwitch)
	if err != nil {
		return err
	}
	err = validateTypeSwitchAgainstEnumDeclForUnknownCases(enumDecl, typeSwitch)
	if err != nil {
		return err
	}
	return nil
}

func validateTypeSwitchAgainstEnumDeclForCount(enumDecl PreludeEnumDecl, typeSwitch ast.ExprTypeSwitch) *RuntimeError {
	if len(typeSwitch.Cases) != len(enumDecl.Decl.Cases) {
		for caseIdentifier := range typeSwitch.Cases {
			if caseIdentifier == "Any" {
				return nil
			}
		}
		return ReportNonExhaustiveTypeSwitch(enumDecl, typeSwitch)
	}
	return nil
}

func validateTypeSwitchAgainstEnumDeclForUnknownCases(enumDecl PreludeEnumDecl, typeSwitch ast.ExprTypeSwitch) *RuntimeError {
	for _, caseIdentifier := range typeSwitch.CaseOrder {
		if caseIdentifier == "Any" {
			continue
		}
		_, err := enumDecl.LookupCase(string(caseIdentifier))
		if err != nil {
			return err.CascadeExpr(typeSwitch.Cases[caseIdentifier])
		}
	}
	return nil
}
