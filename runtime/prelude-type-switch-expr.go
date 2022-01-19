package runtime

import (
	"fmt"

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

func (PreludeTypeSwitchExpr) Lookup(member string) (Evaluatable, *RuntimeError) {
	panic("TODO: not implemented PreludeTypeSwitchExpr")
}

func (PreludeTypeSwitchExpr) RuntimeType() RuntimeTypeRef {
	return PreludeFunctionTypeRef
}

func (t PreludeTypeSwitchExpr) String() string {
	value, err := t.enumDefValue.Evaluate()
	if err != nil {
		panic(fmt.Sprintf("error: %s", err))
	}
	return fmt.Sprintf("(type %s {})", value)
}

func (t PreludeTypeSwitchExpr) Arity() int {
	return 1
}

func (t PreludeTypeSwitchExpr) Call(args []Evaluatable) (RuntimeValue, *RuntimeError) {
	if len(args) < t.Arity() {
		return MakeCurriedCallable(t, args), nil
	}
	primaryArg, err := args[0].Evaluate()
	if err != nil {
		return nil, err.Cascade(*t.Decl.Meta().Source)
	}

	// TODO: optimization can be cached
	enumDefValue, err := t.enumDefValue.Evaluate()
	if err != nil {
		return nil, err.Cascade(*t.Decl.Meta().Source)
	}
	enumDef, ok := enumDefValue.(PreludeEnumDecl)
	if !ok {
		return nil, NewRuntimeErrorf(
			"type switch requires enum type, got %s: %s",
			enumDefValue.RuntimeType(),
			enumDefValue,
		).Cascade(*t.Decl.Meta().Source)
	}
	// TODO: optimization can be cached
	// TODO: more validation
	for _, caseIdentifier := range t.Decl.CaseOrder {
		lazyCaseValue, err := enumDef.Lookup(string(caseIdentifier))
		if err != nil {
			return nil, err.Cascade(*t.Decl.Meta().Source)
		}
		caseValue, err := lazyCaseValue.Evaluate()
		if err != nil {
			return nil, err.Cascade(*t.Decl.Meta().Source)
		}
		runtimeDecl, ok := caseValue.(DeclRuntimeValue)
		if !ok {
			return nil, NewRuntimeErrorf(
				"case %s is not a type, got: %s",
				caseIdentifier,
				enumDefValue,
			).Cascade(*t.Decl.Meta().Source)
		}

		ok, err = runtimeDecl.HasInstance(primaryArg)
		if err != nil {
			return nil, err.Cascade(*t.Decl.Meta().Source)
		}
		if !ok {
			continue
		}
		intermediate, err := t.caseValue[caseIdentifier].Evaluate()
		if err != nil {
			return nil, err.Cascade(*t.Decl.Meta().Source)
		}
		fun, ok := intermediate.(CallableRuntimeValue)
		if !ok {
			return nil, NewRuntimeErrorf("cannot call non function %T", intermediate).Cascade(*t.Decl.Meta().Source)
		}
		return fun.Call(args)
	}
	return nil, NewRuntimeErrorf("no matching case %s", primaryArg).Cascade(*t.Decl.Meta().Source)
}
