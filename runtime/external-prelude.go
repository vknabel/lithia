package runtime

import (
	"fmt"

	"github.com/vknabel/go-lithia/ast"
)

var _ ExternalDefinition = ExternalPrelude{}

type ExternalPrelude struct{}

func (e ExternalPrelude) Lookup(name string, env *Environment, decl ast.Decl) (RuntimeValue, bool) {
	switch name {
	case "Int":
		if externDecl, ok := decl.(ast.DeclExternType); ok {
			return PreludePrimitiveExternType{&externDecl, func(value RuntimeValue) (bool, *RuntimeError) {
				_, ok := value.(PreludeInt)
				return ok, nil
			}}, true
		} else {
			return nil, false
		}
	case "Float":
		if externDecl, ok := decl.(ast.DeclExternType); ok {
			return PreludePrimitiveExternType{&externDecl, func(value RuntimeValue) (bool, *RuntimeError) {
				_, ok := value.(PreludeFloat)
				return ok, nil
			}}, true
		} else {
			return nil, false
		}
	case "String":
		if externDecl, ok := decl.(ast.DeclExternType); ok {
			return PreludePrimitiveExternType{&externDecl, func(value RuntimeValue) (bool, *RuntimeError) {
				_, ok := value.(PreludeString)
				return ok, nil
			}}, true
		} else {
			return nil, false
		}
	case "Char":
		if externDecl, ok := decl.(ast.DeclExternType); ok {
			return PreludePrimitiveExternType{&externDecl, func(value RuntimeValue) (bool, *RuntimeError) {
				panic("char not implemented")
			}}, true
		} else {
			return nil, false
		}
	case "Function":
		if externDecl, ok := decl.(ast.DeclExternType); ok {
			return PreludePrimitiveExternType{&externDecl, func(value RuntimeValue) (bool, *RuntimeError) {
				_, ok := value.(CallableRuntimeValue)
				return ok, nil
			}}, true
		} else {
			return nil, false
		}
	case "Module":
		if externDecl, ok := decl.(ast.DeclExternType); ok {
			return PreludePrimitiveExternType{&externDecl, func(value RuntimeValue) (bool, *RuntimeError) {
				_, ok := value.(PreludeModule)
				return ok, nil
			}}, true
		} else {
			return nil, false
		}
	case "Any":
		if externDecl, ok := decl.(ast.DeclExternType); ok {
			return PreludePrimitiveExternType{&externDecl, func(value RuntimeValue) (bool, *RuntimeError) {
				return true, nil
			}}, true
		} else {
			return nil, false
		}

	case "print":
		return builtinPrint(decl), true
	case "debug":
		return builtinDebug(decl), true

	default:
		return nil, false
	}
}

func builtinDebug(decl ast.Decl) PreludeExternFunction {
	return MakeExternFunction(
		decl,
		func(args []Evaluatable) (RuntimeValue, *RuntimeError) {
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			fmt.Printf("DEBUG: (%s: %s)\n", value.RuntimeType().Name, value)
			return value, nil
		},
	)
}

func builtinPrint(decl ast.Decl) PreludeExternFunction {
	return MakeExternFunction(
		decl,
		func(args []Evaluatable) (RuntimeValue, *RuntimeError) {
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			fmt.Println(value)
			return value, nil
		},
	)
}
