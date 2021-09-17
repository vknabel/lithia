package interpreter

import "fmt"

var _ ExternalDefinition = ExternalPrelude{}

type ExternalPrelude struct{}

func (e ExternalPrelude) Lookup(name string, env *Environment) (RuntimeValue, bool) {
	switch name {
	case "Int":
		return PreludeInt(0).RuntimeType(), true
	case "Float":
		return PreludeFloat(0).RuntimeType(), true
	case "String":
		return PreludeString("").RuntimeType(), true
	case "Char":
		return PreludeChar(0).RuntimeType(), true
	case "Function":
		return PreludeFunctionType{}, true
	case "Module":
		return PreludeModuleType{}, true
	case "Any":
		return PreludeAnyType{}, true

	case "print":
		return builtinPrint, true
	case "debug":
		return builtinDebug, true

	default:
		return nil, false
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
