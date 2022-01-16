package runtime

import (
	"github.com/vknabel/go-lithia/ast"
)

var _ ExternalDefinition = ExternalPrelude{}

type ExternalPrelude struct{}

func (e ExternalPrelude) Lookup(name string, env *Environment, decl *ast.Decl) (RuntimeValue, bool) {
	switch name {
	case "Int":
		if externDecl, ok := (*decl).(ast.DeclExternType); ok {
			return PreludePrimitiveExternType{&externDecl}, true
		} else {
			return nil, false
		}
	case "Float":
		if externDecl, ok := (*decl).(ast.DeclExternType); ok {
			return PreludePrimitiveExternType{&externDecl}, true
		} else {
			return nil, false
		}
	case "String":
		if externDecl, ok := (*decl).(ast.DeclExternType); ok {
			return PreludePrimitiveExternType{&externDecl}, true
		} else {
			return nil, false
		}
	case "Char":
		if externDecl, ok := (*decl).(ast.DeclExternType); ok {
			return PreludePrimitiveExternType{&externDecl}, true
		} else {
			return nil, false
		}
	case "Function":
		if externDecl, ok := (*decl).(ast.DeclExternType); ok {
			return PreludePrimitiveExternType{&externDecl}, true
		} else {
			return nil, false
		}
	case "Module":
		if externDecl, ok := (*decl).(ast.DeclExternType); ok {
			return PreludePrimitiveExternType{&externDecl}, true
		} else {
			return nil, false
		}
	case "Any":
		if externDecl, ok := (*decl).(ast.DeclExternType); ok {
			return PreludePrimitiveExternType{&externDecl}, true
		} else {
			return nil, false
		}

	// case "print":
	// 	return builtinPrint(docs), true
	// case "debug":
	// 	return builtinDebug(docs), true

	default:
		return nil, false
	}
}

// func builtinDebug(docs Docs) BuiltinFunction {
// 	return NewBuiltinFunction(
// 		"debug",
// 		[]string{"message"},
// 		docs,
// 		func(args []Evaluatable) (RuntimeValue, error) {
// 			value, err := args[0].Evaluate()
// 			if err != nil {
// 				return nil, err
// 			}
// 			fmt.Printf("DEBUG: (%s: %s)\n", value.RuntimeType().name, value)
// 			return value, nil
// 		},
// 	)
// }

// func builtinPrint(docs Docs) BuiltinFunction {
// 	return NewBuiltinFunction(
// 		"print",
// 		[]string{"message"},
// 		docs,
// 		func(args []Evaluatable) (RuntimeValue, error) {
// 			value, err := args[0].Evaluate()
// 			if err != nil {
// 				return nil, err
// 			}
// 			fmt.Println(value)
// 			return value, nil
// 		},
// 	)
// }
