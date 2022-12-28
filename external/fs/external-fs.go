package fs

import (
	"os"

	"github.com/vknabel/lithia/ast"
	. "github.com/vknabel/lithia/runtime"
	"github.com/vknabel/lithia/world"
)

var _ ExternalDefinition = ExternalFS{}

type ExternalFS struct {
	inter *Interpreter
}

func New(inter *Interpreter) ExternalFS {
	return ExternalFS{inter}
}

func (e ExternalFS) Lookup(name string, env *Environment, decl ast.Decl) (RuntimeValue, bool) {
	switch name {
	case "writeString":
		return builtinFsWrite(env, decl), true
	case "readString":
		return builtinFsRead(env, decl), true
	case "exists":
		return builtinFsExists(env, decl), true
	case "delete":
		return builtinFsDelete(env, decl), true
	case "deleteAll":
		return builtinFsDeleteAll(env, decl), true
	default:
		return nil, false
	}
}

func builtinFsWrite(env *Environment, decl ast.Decl) PreludeExternFunction {
	return MakeExternFunction(
		decl,
		func(args []Evaluatable) (RuntimeValue, *RuntimeError) {
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
				return nil, NewRuntimeErrorf("%s is not a string", toPathValue)
			}
			contents, ok := contentsValue.(PreludeString)
			if !ok {
				return nil, NewRuntimeErrorf("%s is not a string", contentsValue)
			}
			writeError := world.Current.FS.WriteFile(string(toPath), []byte(string(contents)), 0644)
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

func builtinFsRead(env *Environment, decl ast.Decl) PreludeExternFunction {
	return MakeExternFunction(
		decl,
		func(args []Evaluatable) (RuntimeValue, *RuntimeError) {
			fromPathValue, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			fromPath, ok := fromPathValue.(PreludeString)
			if !ok {
				return nil, NewRuntimeErrorf("%s is not a string", fromPath)
			}
			bytes, writeError := world.Current.FS.ReadFile(string(fromPath))
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

func builtinFsExists(env *Environment, decl ast.Decl) PreludeExternFunction {
	return MakeExternFunction(
		decl,
		func(args []Evaluatable) (RuntimeValue, *RuntimeError) {
			atPathValue, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			atPath, ok := atPathValue.(PreludeString)
			if !ok {
				return nil, NewRuntimeErrorf("%s is not a string", atPath)
			}
			_, writeError := world.Current.FS.Stat(string(atPath))
			if os.IsNotExist(writeError) {
				return env.MakeEmptyDataRuntimeValue("False")
			} else {
				return env.MakeEmptyDataRuntimeValue("True")
			}
		},
	)
}

func builtinFsDelete(env *Environment, decl ast.Decl) PreludeExternFunction {
	return MakeExternFunction(
		decl,
		func(args []Evaluatable) (RuntimeValue, *RuntimeError) {
			atPathValue, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			atPath, ok := atPathValue.(PreludeString)
			if !ok {
				return nil, NewRuntimeErrorf("%s is not a string", atPath)
			}
			writeError := world.Current.FS.Remove(string(atPath))
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

func builtinFsDeleteAll(env *Environment, decl ast.Decl) PreludeExternFunction {
	return MakeExternFunction(
		decl,
		func(args []Evaluatable) (RuntimeValue, *RuntimeError) {
			atPathValue, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			atPath, ok := atPathValue.(PreludeString)
			if !ok {
				return nil, NewRuntimeErrorf("%s is not a string", atPath)
			}
			writeError := world.Current.FS.Remove(string(atPath))
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
