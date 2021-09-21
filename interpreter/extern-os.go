package interpreter

import (
	"fmt"
	"os"
)

var _ ExternalDefinition = ExternalOS{}

type ExternalOS struct{}

func (e ExternalOS) Lookup(name string, env *Environment, docs Docs) (DocumentedRuntimeValue, bool) {
	switch name {
	case "exit":
		return builtinOsExit(docs), true
	case "env":
		return builtinOsEnv(env, docs), true
	default:
		return nil, false
	}
}

func builtinOsExit(docs Docs) BuiltinFunction {
	return NewBuiltinFunction(
		"exit",
		[]string{"code"},
		docs,
		func(args []Evaluatable) (RuntimeValue, error) {
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			if code, ok := value.(PreludeInt); ok {
				os.Exit(int(code))
				return value, nil
			} else {
				return nil, fmt.Errorf("%s is not an int", value)
			}
		},
	)
}

func builtinOsEnv(prelude *Environment, docs Docs) BuiltinFunction {
	return NewBuiltinFunction(
		"env",
		[]string{"key"},
		docs,
		func(args []Evaluatable) (RuntimeValue, error) {
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			if key, ok := value.(PreludeString); ok {
				if env, ok := os.LookupEnv(string(key)); ok && env != "" {
					return prelude.MakeDataRuntimeValue("Some", map[string]Evaluatable{
						"value": NewConstantRuntimeValue(PreludeString(env)),
					})
				} else {
					return prelude.MakeEmptyDataRuntimeValue("None")
				}
			} else {
				return nil, fmt.Errorf("%s is not a string", value)
			}
		},
	)
}
