package interpreter

import (
	"fmt"
	"os"
)

func (inter *Interpreter) NewPreludeEnvironment() *Environment {
	if inter.prelude != nil {
		return inter.prelude
	}
	env := NewEnvironment(nil)
	inter.prelude = env

	env.Declare("Int", NewConstantRuntimeValue(PreludeInt(0)))
	env.Declare("Float", NewConstantRuntimeValue(PreludeFloat(0.0)))
	env.Declare("String", NewConstantRuntimeValue(PreludeString("")))
	env.Declare("Rune", NewConstantRuntimeValue(PreludeRune('r')))
	env.Declare("Function", NewConstantRuntimeValue(PreludeFunctionType{}))
	env.Declare("Variable", NewConstantRuntimeValue(PreludeVariableType{}))
	env.Declare("Module", NewConstantRuntimeValue(PreludeModuleType{}))
	env.Declare("Any", NewConstantRuntimeValue(PreludeAnyType{}))

	env.Declare("print", NewConstantRuntimeValue(builtinPrint))
	env.Declare("debug", NewConstantRuntimeValue(builtinDebug))
	env.Declare("osExit", NewConstantRuntimeValue(builtinOsExit))

	module, err := inter.LoadModuleIfNeeded(ModuleName("prelude"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: prelude not loaded\n    %s\n", err)
	}
	// These declares override the ones in the prelude.
	env.Parent = &Environment{Parent: nil, Scope: module.environment.Scope, Unexported: module.environment.Unexported}
	env.Declare("osEnv", NewConstantRuntimeValue(builtinOsEnv(env)))

	return env
}

var builtinDebug = NewBuiltinFunction(
	"debug",
	[]string{"message"},
	func(args []*LazyRuntimeValue) (RuntimeValue, error) {
		value, err := args[0].Evaluate()
		if err != nil {
			return nil, err
		}
		fmt.Printf("DEBUG: (%s: %s)\n", value.RuntimeType().name, value)
		return value, nil
	},
)

var builtinPrint = NewBuiltinFunction(
	"print",
	[]string{"message"},
	func(args []*LazyRuntimeValue) (RuntimeValue, error) {
		value, err := args[0].Evaluate()
		if err != nil {
			return nil, err
		}
		fmt.Println(value)
		return value, nil
	},
)

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
