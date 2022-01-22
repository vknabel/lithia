package runtime

import (
	"fmt"
	"strings"

	"github.com/vknabel/go-lithia/ast"
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
	cases := make([]string, 0, len(t.Decl.CaseOrder))
	for i, identifier := range t.Decl.CaseOrder {
		cases[i] = string(identifier)
	}
	return fmt.Sprintf("type %s { %s })", value, strings.Join(cases, ", "))
}

func (t PreludeTypeSwitchExpr) Arity() int {
	return 1
}

func (t PreludeTypeSwitchExpr) Call(args []Evaluatable) (RuntimeValue, *RuntimeError) {
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
	// TODO: optimization can be cached
	// TODO: more validation
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
			return Call(fun, args)
		}
		lazyCaseValue, err := enumDef.Lookup(string(caseIdentifier))
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
		return Call(fun, args)
	}
	return nil, NewRuntimeErrorf("no matching case %s", primaryArg).CascadeExpr(t.Decl)
}
