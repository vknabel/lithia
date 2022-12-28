package os

import (
	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/runtime"
	"github.com/vknabel/lithia/world"
)

var _ runtime.ExternalDefinition = ExternalOS{}

type ExternalOS struct {
	inter *runtime.Interpreter
}

func New(inter *runtime.Interpreter) ExternalOS {
	return ExternalOS{inter}
}

func (e ExternalOS) Lookup(name string, env *runtime.Environment, decl ast.Decl) (runtime.RuntimeValue, bool) {
	switch name {
	case "exit":
		return builtinOsExit(decl), true
	case "env":
		return builtinOsEnv(env, decl), true
	case "args":
		relevantArgs := world.Current.Args
		runtimeArgs := make([]runtime.RuntimeValue, len(relevantArgs))
		for i, arg := range relevantArgs {
			runtimeArgs[i] = runtime.PreludeString(arg)
		}
		list, err := env.MakeEagerList(runtimeArgs)
		if err != nil {
			return nil, false
		}
		return list, true
	default:
		return nil, false
	}
}

func builtinOsExit(decl ast.Decl) runtime.PreludeExternFunction {
	return runtime.MakeExternFunction(
		decl,
		func(args []runtime.Evaluatable) (runtime.RuntimeValue, *runtime.RuntimeError) {
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			if code, ok := value.(runtime.PreludeInt); ok {
				world.Current.Env.Exit(int(code))
				return value, nil
			} else {
				return nil, runtime.NewRuntimeErrorf("%s is not an int", value).CascadeDecl(decl)
			}
		},
	)
}

func builtinOsEnv(prelude *runtime.Environment, decl ast.Decl) runtime.PreludeExternFunction {
	return runtime.MakeExternFunction(
		decl,
		func(args []runtime.Evaluatable) (runtime.RuntimeValue, *runtime.RuntimeError) {
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err.CascadeDecl(decl)
			}
			if key, ok := value.(runtime.PreludeString); ok {
				if env, ok := world.Current.Env.LookupEnv(string(key)); ok && env != "" {
					value, err := prelude.MakeDataRuntimeValue("Some", map[string]runtime.Evaluatable{
						"value": runtime.NewConstantRuntimeValue(runtime.PreludeString(env)),
					})
					return value, err.CascadeDecl(decl)
				} else {
					value, err := prelude.MakeEmptyDataRuntimeValue("None")
					return value, err.CascadeDecl(decl)
				}
			} else {
				return nil, runtime.NewRuntimeErrorf("%s is not a string", value).CascadeDecl(decl)
			}
		},
	)
}
