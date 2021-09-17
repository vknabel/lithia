package interpreter

import (
	"fmt"
	"os"
)

var _ ExternalDefinition = ExternalOS{}

type ExternalOS struct{}

func (e ExternalOS) Lookup(name string, env *Environment) (RuntimeValue, bool) {
	switch name {
	case "exit":
		return builtinOsExit, true
	case "env":
		return builtinOsEnv(env), true
	default:
		return nil, false
	}
}

var builtinOsExit = NewBuiltinFunction(
	"osExit",
	[]string{"code"},
	func(args []*LazyRuntimeValue) (RuntimeValue, error) {
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

func builtinOsEnv(prelude *Environment) BuiltinFunction {
	return NewBuiltinFunction(
		"osEnv",
		[]string{"key"},
		func(args []*LazyRuntimeValue) (RuntimeValue, error) {
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			if key, ok := value.(PreludeString); ok {
				if env, ok := os.LookupEnv(string(key)); ok {
					return prelude.MakeDataRuntimeValue("Some", map[string]*LazyRuntimeValue{
						"value": NewConstantRuntimeValue(PreludeString(env)),
					})
				} else {
					return prelude.MakeDataRuntimeValue("Some", map[string]*LazyRuntimeValue{
						"value": NewConstantRuntimeValue(PreludeString(env)),
					})
				}
			} else {
				return nil, fmt.Errorf("%s is not a string", value)
			}
		},
	)
}
