package runtime

import (
	"fmt"

	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/resolution"
)

func (inter *Interpreter) NewPreludeEnvironment(resolvedModule resolution.ResolvedModule) (*Environment, error) {
	if inter.Prelude != nil {
		return inter.Prelude, nil
	}
	env := NewEnvironment(nil)
	inter.Prelude = env

	module, err := inter.LoadModuleIfNeeded(ast.ModuleName("prelude"), resolvedModule)
	if err != nil {
		return nil, fmt.Errorf("error loading prelude: %v\nIs $LITHIA_STDLIB set up correctly?", err)
	}
	// These declares override the ones in the prelude.
	env.Parent = &Environment{Parent: nil, Scope: module.Environment.Scope, Unexported: module.Environment.Unexported}

	return env, nil
}
