package runtime

import (
	"os"

	"github.com/vknabel/go-lithia/ast"
)

var _ ExternalDefinition = ExternalOS{}

type ExternalOS struct{}

func (e ExternalOS) Lookup(name string, env *Environment, decl ast.Decl) (RuntimeValue, bool) {
	switch name {
	case "exit":
		return builtinOsExit(decl), true
	case "env":
		return builtinOsEnv(env, decl), true
	default:
		return nil, false
	}
}

func builtinOsExit(decl ast.Decl) PreludeExternFunction {
	return MakeExternFunction(
		decl,
		func(args []Evaluatable) (RuntimeValue, *RuntimeError) {
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			if code, ok := value.(PreludeInt); ok {
				os.Exit(int(code))
				return value, nil
			} else {
				return nil, NewRuntimeErrorf("%s is not an int", value).CascadeDecl(decl)
			}
		},
	)
}

func builtinOsEnv(prelude *Environment, decl ast.Decl) PreludeExternFunction {
	return MakeExternFunction(
		decl,
		func(args []Evaluatable) (RuntimeValue, *RuntimeError) {
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err.CascadeDecl(decl)
			}
			if key, ok := value.(PreludeString); ok {
				if env, ok := os.LookupEnv(string(key)); ok && env != "" {
					value, err := prelude.MakeDataRuntimeValue("Some", map[string]Evaluatable{
						"value": NewConstantRuntimeValue(PreludeString(env)),
					})
					return value, err.CascadeDecl(decl)
				} else {
					value, err := prelude.MakeEmptyDataRuntimeValue("None")
					return value, err.CascadeDecl(decl)
				}
			} else {
				return nil, NewRuntimeErrorf("%s is not a string", value).CascadeDecl(decl)
			}
		},
	)
}
