package interpreter_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	i "github.com/vknabel/go-lithia/interpreter"
)

func TestStdlib(t *testing.T) {
	pathToStdlib := "../stdlib"
	inter := i.NewInterpreter(pathToStdlib)
	mockOS := &mockExternalOS{
		calledExitCode: -1,
		env:            map[string]string{"LITHIA_TESTS": "1"},
	}
	inter.ExternalDefinitions["os"] = mockOS

	scriptData, err := os.ReadFile(filepath.Join(pathToStdlib, "stdlib-tests.lithia"))
	if err != nil {
		t.Errorf("Error reading stdlib-tests.lithia: %s", err)
		return
	}
	_, err = inter.Interpret("stdlib-tests.lithia", string(scriptData))
	if err != nil {
		t.Errorf("Error interpreting stdlib-tests.lithia: %s", err)
		return
	}
	if mockOS.calledExitCode != i.PreludeInt(0) && mockOS.calledExitCode != i.PreludeInt(-1) {
		t.Errorf("lithia tests failed with exit code %d", mockOS.calledExitCode)
		return
	}
}

type mockExternalOS struct {
	calledExitCode i.PreludeInt
	env            map[string]string
}

func (e *mockExternalOS) Lookup(name string, env *i.Environment, docs i.Docs) (i.DocumentedRuntimeValue, bool) {
	switch name {
	case "exit":
		return mockOsExit(env, func(code i.PreludeInt) {
			e.calledExitCode = code
		}), true
	case "env":
		return mockOsEnv(env, e.env), true
	default:
		return nil, false
	}
}

func mockOsExit(prelude *i.Environment, impl func(i.PreludeInt)) i.BuiltinFunction {
	return i.NewBuiltinFunction(
		"osExit",
		[]string{"code"},
		i.Docs{},
		func(args []*i.LazyRuntimeValue) (i.RuntimeValue, error) {
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			if code, ok := value.(i.PreludeInt); ok {
				impl(code)
				return code, nil
			} else {
				return nil, fmt.Errorf("%s is not an integer", value)
			}
		},
	)
}
func mockOsEnv(prelude *i.Environment, fakeOsEnv map[string]string) i.BuiltinFunction {
	return i.NewBuiltinFunction(
		"osEnv",
		[]string{"key"},
		i.Docs{},
		func(args []*i.LazyRuntimeValue) (i.RuntimeValue, error) {
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			if key, ok := value.(i.PreludeString); ok {
				if env, ok := fakeOsEnv[string(key)]; ok {
					return prelude.MakeDataRuntimeValue("Some", map[string]*i.LazyRuntimeValue{
						"value": i.NewConstantRuntimeValue(i.PreludeString(env)),
					})
				} else {
					return prelude.MakeDataRuntimeValue("Some", map[string]*i.LazyRuntimeValue{
						"value": i.NewConstantRuntimeValue(i.PreludeString(env)),
					})
				}
			} else {
				return nil, fmt.Errorf("%s is not a string", value)
			}
		},
	)
}
