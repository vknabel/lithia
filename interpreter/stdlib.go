package interpreter

import (
	"fmt"
	"os"
	"strings"
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

	if module, err := inter.LoadModule(ModuleName{name: "prelude"}, "."); err == nil {
		// These declares override the ones in the prelude.
		env.Parent = &Environment{Parent: nil, Scope: module.environment.Scope, Unexported: module.environment.Unexported}
		env.Declare("osEnv", NewConstantRuntimeValue(builtinOsEnv(env)))
	}

	return env
}

func NewBuiltinFunction(
	name string,
	args []string,
	impl func(args []*LazyRuntimeValue) (RuntimeValue, error),
) BuiltinFunction {
	f := BuiltinFunction{
		name: name,
		args: args,
		impl: impl,
	}
	var _ RuntimeValue = f
	var _ Callable = f
	return f
}

type BuiltinFunction struct {
	name string
	args []string
	impl func(args []*LazyRuntimeValue) (RuntimeValue, error)
}

func (f BuiltinFunction) RuntimeType() RuntimeType {
	return PreludeFunctionType{}.RuntimeType()
}

func (f BuiltinFunction) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("function %s has no member %s", fmt.Sprint(f), member)
}

func (f BuiltinFunction) String() string {
	return fmt.Sprintf("{ %s => @(%s) }", strings.Join(f.args, ","), f.name)
}

func (f BuiltinFunction) Call(arguments []*LazyRuntimeValue) (RuntimeValue, error) {
	if len(arguments) < len(f.args) {
		return CurriedCallable{
			actual:         f,
			args:           arguments,
			remainingArity: len(f.args) - len(arguments),
		}, nil
	}
	intermediate, err := f.impl(arguments[:len(f.args)])
	if err != nil {
		return nil, err
	}
	if len(arguments) == len(f.args) {
		return intermediate, nil
	}
	if g, ok := intermediate.(Callable); ok {
		return g.Call(arguments[len(f.args):])
	} else {
		return nil, fmt.Errorf("%s is not callable", g)
	}
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
					if lazySome, ok := prelude.Get("Some"); ok {
						someValue, err := lazySome.Evaluate()
						if err != nil {
							return nil, err
						}
						if someType, ok := someValue.(DataDeclRuntimeValue); ok {
							members := make(map[string]*LazyRuntimeValue)
							members["value"] = NewConstantRuntimeValue(PreludeString(env))
							return DataRuntimeValue{typeValue: &someType, members: members}, nil
						} else {
							return nil, fmt.Errorf("%s is not a data type", someValue)
						}
					} else {
						return nil, fmt.Errorf("prelude not found")
					}
				}
				if lazyNone, ok := prelude.Get("None"); ok {
					noneValue, err := lazyNone.Evaluate()
					if err != nil {
						return nil, err
					}
					if noneType, ok := noneValue.(DataDeclRuntimeValue); ok {
						members := make(map[string]*LazyRuntimeValue)
						return DataRuntimeValue{typeValue: &noneType, members: members}, nil
					} else {
						return nil, fmt.Errorf("%s is not a data type", noneType)
					}
				} else {
					return nil, fmt.Errorf("prelude not found")
				}
			} else {
				return nil, fmt.Errorf("%s is not a string", value)
			}
		},
	)
}
