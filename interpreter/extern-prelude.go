package interpreter

import "fmt"

var _ ExternalDefinition = ExternalPrelude{}

type ExternalPrelude struct{}

func (e ExternalPrelude) Lookup(name string, env *Environment, docs Docs) (DocumentedRuntimeValue, bool) {
	switch name {
	case "Int":
		runtimeType := PreludeInt(0).RuntimeType()
		runtimeType.docs = docs
		return runtimeType, true
	case "Float":
		runtimeType := PreludeFloat(0).RuntimeType()
		runtimeType.docs = docs
		return runtimeType, true
	case "String":
		runtimeType := PreludeString("").RuntimeType()
		runtimeType.docs = docs
		return runtimeType, true
	case "Char":
		runtimeType := PreludeChar(0).RuntimeType()
		runtimeType.docs = docs
		return runtimeType, true
	case "Function":
		return PreludeFunctionType{docs}, true
	case "Module":
		return PreludeModuleType{docs}, true
	case "Any":
		return PreludeAnyType{docs}, true

	case "print":
		return builtinPrint(docs), true
	case "debug":
		return builtinDebug(docs), true

	default:
		return nil, false
	}
}

func builtinDebug(docs Docs) BuiltinFunction {
	return NewBuiltinFunction(
		"debug",
		[]string{"message"},
		docs,
		func(args []Evaluatable) (RuntimeValue, error) {
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			fmt.Printf("DEBUG: (%s: %s)\n", value.RuntimeType().name, value)
			return value, nil
		},
	)
}

func builtinPrint(docs Docs) BuiltinFunction {
	return NewBuiltinFunction(
		"print",
		[]string{"message"},
		docs,
		func(args []Evaluatable) (RuntimeValue, error) {
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			fmt.Println(value)
			return value, nil
		},
	)
}
