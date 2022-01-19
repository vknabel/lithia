package runtime

import (
	"fmt"

	"github.com/vknabel/go-lithia/ast"
)

func MakeRuntimeValueDecl(context *InterpreterContext, decl ast.Decl) Evaluatable {
	switch decl := decl.(type) {
	case ast.DeclConstant:
		return MakeEvaluatableExpr(context, decl.Value)
	case ast.DeclData:
		return NewConstantRuntimeValue(PreludeDataDecl{Decl: decl})
	case ast.DeclEnum:
		return NewConstantRuntimeValue(MakeEnumDecl(context, decl))
	case ast.DeclFunc:
		return NewConstantRuntimeValue(MakePreludeFuncDecl(context, decl))
	case ast.DeclExternFunc, ast.DeclExternType:
		definition, ok := context.interpreter.ExternalDefinitions[context.module.Name]
		if !ok {
			panic(fmt.Sprintf("extern definitions not found for module %s. assumed by extern %s", context.module.Name, decl.DeclName()))
		}
		value, ok := definition.Lookup(string(decl.DeclName()), context.environment, decl)
		if !ok {
			panic(fmt.Sprintf("extern definition not found in module %s: extern %s", context.module.Name, decl.DeclName()))
		}
		return NewConstantRuntimeValue(value)
	case ast.DeclModule:
		return NewConstantRuntimeValue(PreludeModule{Module: context.module})
	case ast.DeclImport:
		return NewLazyRuntimeValue(func() (RuntimeValue, *RuntimeError) {
			module, err := context.interpreter.LoadModuleIfNeeded(decl.ModuleName)
			if err != nil {
				return nil, NewRuntimeError(err).Cascade(*decl.Meta().Source)
			}
			return PreludeModule{Module: module}, nil
		})
	case ast.DeclImportMember:
		return NewLazyRuntimeValue(func() (RuntimeValue, *RuntimeError) {
			module, err := context.interpreter.LoadModuleIfNeeded(decl.ModuleName)
			if err != nil {
				return nil, NewRuntimeError(err).Cascade(*decl.Meta().Source)
			}
			value, err := module.Environment.GetEvaluatedRuntimeValue(string(decl.DeclName()))
			return value, NewRuntimeError(err).Cascade(*decl.Meta().Source)
		})
	default:
		panic(fmt.Sprintf("unknown decl: %T %s", decl, decl.DeclName()))
	}
}
