package interpreter

import (
	"fmt"
	"os"
)

var _ ExternalDefinition = ExternalFS{}

type ExternalFS struct{}

func (e ExternalFS) Lookup(name string, env *Environment, docs Docs) (DocumentedRuntimeValue, bool) {
	switch name {
	case "writeString":
		return builtinFsWrite(env, docs), true
	case "readString":
		return builtinFsRead(env, docs), true
	case "exists":
		return builtinFsExists(env, docs), true
	case "delete":
		return builtinFsDelete(env, docs), true
	case "deleteAll":
		return builtinFsDeleteAll(env, docs), true
	default:
		return nil, false
	}
}

func builtinFsWrite(env *Environment, docs Docs) BuiltinFunction {
	return NewBuiltinFunction(
		"writeString",
		[]string{"toPath", "contents"},
		docs,
		func(args []Evaluatable) (RuntimeValue, error) {
			toPathValue, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			contentsValue, err := args[1].Evaluate()
			if err != nil {
				return nil, err
			}
			toPath, ok := toPathValue.(PreludeString)
			if !ok {
				return nil, fmt.Errorf("%s is not a string", toPathValue)
			}
			contents, ok := contentsValue.(PreludeString)
			if !ok {
				return nil, fmt.Errorf("%s is not a string", contentsValue)
			}
			writeError := os.WriteFile(string(toPath), []byte(string(contents)), 0644)
			if writeError != nil {
				return env.MakeDataRuntimeValue("Failure", map[string]Evaluatable{
					"error": NewConstantRuntimeValue(PreludeString(writeError.Error())),
				})
			} else {
				return env.MakeDataRuntimeValue("Success", map[string]Evaluatable{
					"value": NewConstantRuntimeValue(toPath),
				})
			}
		},
	)
}

func builtinFsRead(env *Environment, docs Docs) BuiltinFunction {
	return NewBuiltinFunction(
		"readString",
		[]string{"fromPath"},
		docs,
		func(args []Evaluatable) (RuntimeValue, error) {
			fromPathValue, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			fromPath, ok := fromPathValue.(PreludeString)
			if !ok {
				return nil, fmt.Errorf("%s is not a string", fromPath)
			}
			bytes, writeError := os.ReadFile(string(fromPath))
			if writeError != nil {
				return env.MakeDataRuntimeValue("Failure", map[string]Evaluatable{
					"error": NewConstantRuntimeValue(PreludeString(writeError.Error())),
				})
			} else {
				return env.MakeDataRuntimeValue("Success", map[string]Evaluatable{
					"value": NewConstantRuntimeValue(PreludeString(string(bytes))),
				})
			}
		},
	)
}

func builtinFsExists(env *Environment, docs Docs) BuiltinFunction {
	return NewBuiltinFunction(
		"exists",
		[]string{"atPath"},
		docs,
		func(args []Evaluatable) (RuntimeValue, error) {
			atPathValue, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			atPath, ok := atPathValue.(PreludeString)
			if !ok {
				return nil, fmt.Errorf("%s is not a string", atPath)
			}
			_, writeError := os.Stat(string(atPath))
			if os.IsNotExist(writeError) {
				return env.MakeEmptyDataRuntimeValue("False")
			} else {
				return env.MakeEmptyDataRuntimeValue("True")
			}
		},
	)
}

func builtinFsDelete(env *Environment, docs Docs) BuiltinFunction {
	return NewBuiltinFunction(
		"delete",
		[]string{"atPath"},
		docs,
		func(args []Evaluatable) (RuntimeValue, error) {
			atPathValue, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			atPath, ok := atPathValue.(PreludeString)
			if !ok {
				return nil, fmt.Errorf("%s is not a string", atPath)
			}
			writeError := os.Remove(string(atPath))
			if writeError != nil {
				return env.MakeDataRuntimeValue("Failure", map[string]Evaluatable{
					"error": NewConstantRuntimeValue(PreludeString(writeError.Error())),
				})
			} else {
				return env.MakeDataRuntimeValue("Success", map[string]Evaluatable{
					"value": NewConstantRuntimeValue(atPath),
				})
			}
		},
	)
}

func builtinFsDeleteAll(env *Environment, docs Docs) BuiltinFunction {
	return NewBuiltinFunction(
		"deleteAll",
		[]string{"atPath"},
		docs,
		func(args []Evaluatable) (RuntimeValue, error) {
			atPathValue, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			atPath, ok := atPathValue.(PreludeString)
			if !ok {
				return nil, fmt.Errorf("%s is not a string", atPath)
			}
			writeError := os.Remove(string(atPath))
			if writeError != nil {
				return env.MakeDataRuntimeValue("Failure", map[string]Evaluatable{
					"error": NewConstantRuntimeValue(PreludeString(writeError.Error())),
				})
			} else {
				return env.MakeDataRuntimeValue("Success", map[string]Evaluatable{
					"value": NewConstantRuntimeValue(atPath),
				})
			}
		},
	)
}
