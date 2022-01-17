package runtime

import (
	"fmt"

	"github.com/vknabel/go-lithia/ast"
)

func MakeRuntimeValueDecl(context *InterpreterContext, decl ast.Decl) Evaluatable {
	switch decl := decl.(type) {
	case ast.DeclConstant:
		return MakeEvaluatableExpr(context, *decl.Value)
	case ast.DeclData:
		return NewConstantRuntimeValue(PreludeDataDecl{Decl: decl})
	case ast.DeclEnum:
		return NewConstantRuntimeValue(PreludeEnumDecl{Decl: decl})
	case ast.DeclFunc:
		return NewConstantRuntimeValue(PreludeFuncDecl{
			Environment: context.environment.Private(),
			Decl:        decl,
		})
	default:
		panic(fmt.Sprintf("unknown decl: %s", decl))
	}
}
