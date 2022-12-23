package runtime_test

import (
	"path/filepath"
	"testing"

	l "github.com/vknabel/lithia"
	"github.com/vknabel/lithia/ast"
	r "github.com/vknabel/lithia/runtime"
	"github.com/vknabel/lithia/world"
)

func TestStdlib(t *testing.T) {
	pathToStdlib := "../stdlib"
	inter := l.NewDefaultInterpreter(pathToStdlib, "../stdlib")
	mockOS := &mockExternalOS{
		calledExitCode: -1,
		env:            map[string]string{"LITHIA_TESTS": "1"},
	}
	inter.ExternalDefinitions["os"] = mockOS

	scriptData, err := world.Current.FS.ReadFile(filepath.Join(pathToStdlib, "stdlib-tests.lithia"))
	if err != nil {
		t.Errorf("error reading stdlib-tests.lithia: %s", err)
		return
	}
	_, err = inter.Interpret("stdlib-tests.lithia", string(scriptData))
	if err != nil {
		t.Errorf("error interpreting stdlib-tests.lithia: %s", err)
		return
	}
	if mockOS.calledExitCode != r.PreludeInt(0) && mockOS.calledExitCode != r.PreludeInt(-1) {
		t.Errorf("lithia tests failed with exit code %d", mockOS.calledExitCode)
		return
	}
}

type mockExternalOS struct {
	calledExitCode r.PreludeInt
	env            map[string]string
	args           []string
}

func (e *mockExternalOS) Lookup(name string, env *r.Environment, decl ast.Decl) (r.RuntimeValue, bool) {
	switch name {
	case "exit":
		return mockOsExit(env, decl, func(code r.PreludeInt) {
			e.calledExitCode = code
		}), true
	case "env":
		return mockOsEnv(env, decl, e.env), true
	case "args":
		preludeStrings := make([]r.RuntimeValue, len(e.args))
		for i, arg := range e.args {
			preludeStrings[i] = r.PreludeString(arg)
		}
		argsList, err := env.MakeEagerList(preludeStrings)
		if err != nil {
			return nil, false
		}
		return argsList, true
	default:
		return nil, false
	}
}

func mockOsExit(prelude *r.Environment, decl ast.Decl, impl func(r.PreludeInt)) r.PreludeExternFunction {
	return r.MakeExternFunction(
		decl,
		func(args []r.Evaluatable) (r.RuntimeValue, *r.RuntimeError) {
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			if code, ok := value.(r.PreludeInt); ok {
				impl(code)
				return code, nil
			} else {
				return nil, r.NewRuntimeErrorf("%s is not an integer", value)
			}
		},
	)
}
func mockOsEnv(prelude *r.Environment, decl ast.Decl, fakeOsEnv map[string]string) r.PreludeExternFunction {
	return r.MakeExternFunction(
		decl,
		func(args []r.Evaluatable) (r.RuntimeValue, *r.RuntimeError) {
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			if key, ok := value.(r.PreludeString); ok {
				if env, ok := fakeOsEnv[string(key)]; ok {
					return prelude.MakeDataRuntimeValue("Some", map[string]r.Evaluatable{
						"value": r.NewConstantRuntimeValue(r.PreludeString(env)),
					})
				} else {
					return prelude.MakeDataRuntimeValue("Some", map[string]r.Evaluatable{
						"value": r.NewConstantRuntimeValue(r.PreludeString(env)),
					})
				}
			} else {
				return nil, r.NewRuntimeErrorf("%s is not a string", value)
			}
		},
	)
}
