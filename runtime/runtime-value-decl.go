package runtime

import (
	"fmt"

	"github.com/vknabel/lithia/ast"
)

func MakeRuntimeValueDecl(context *InterpreterContext, decl ast.Decl) (Evaluatable, *RuntimeError) {
	switch decl := decl.(type) {
	case ast.DeclConstant:
		return MakeEvaluatableExpr(context, decl.Value), nil
	case ast.DeclData:
		return NewConstantRuntimeValue(PreludeDataDecl{Decl: decl}), nil
	case ast.DeclEnum:
		return NewLazyRuntimeValue(func() (RuntimeValue, *RuntimeError) {
			enumDecl, err := MakeEnumDecl(context, decl)
			if err != nil {
				return nil, err
			}
			return enumDecl, nil
		}), nil
	case ast.DeclFunc:
		value, err := MakePreludeFuncDecl(context, decl)
		if err != nil {
			return nil, err
		}
		return NewConstantRuntimeValue(value), nil
	case ast.DeclExternFunc, ast.DeclExternType:
		definition, ok := context.interpreter.ExternalDefinitions[context.module.Name]
		if !ok {
			return nil, NewRuntimeErrorf("extern definitions not found for module %s. assumed by extern %s", context.module.Name, decl.DeclName())
		}
		value, ok := definition.Lookup(string(decl.DeclName()), context.environment, decl)
		if !ok {
			return nil, NewRuntimeErrorf("extern definition not found in module %s: extern %s", context.module.Name, decl.DeclName())
		}
		return NewConstantRuntimeValue(value), nil
	case ast.DeclModule:
		return NewConstantRuntimeValue(PreludeModule{Module: context.module}), nil
	case ast.DeclImport:
		module, err := context.interpreter.LoadModuleIfNeeded(decl.ModuleName)
		if err != nil {
			return nil, NewRuntimeError(err).CascadeDecl(decl)
		}
		return NewConstantRuntimeValue(PreludeModule{Module: module}), nil
	case ast.DeclImportMember:
		return NewLazyRuntimeValue(func() (RuntimeValue, *RuntimeError) {
			module, err := context.interpreter.LoadModuleIfNeeded(decl.ModuleName)
			if err != nil {
				return nil, NewRuntimeError(err).CascadeDecl(decl)
			}
			value, err := module.Environment.GetEvaluatedRuntimeValue(string(decl.DeclName()))
			return value, NewRuntimeError(err).CascadeDecl(decl)
		}), nil
	default:
		panic(fmt.Errorf("unknown decl: %T %s", decl, decl.DeclName()))
	}
}
