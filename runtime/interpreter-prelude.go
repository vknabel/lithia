package runtime

import (
	"fmt"

	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/resolution"
	"github.com/vknabel/lithia/world"
)

func (inter *Interpreter) NewPreludeEnvironment(resolvedModule resolution.ResolvedModule) *Environment {
	if inter.Prelude != nil {
		return inter.Prelude
	}
	env := NewEnvironment(nil)
	inter.Prelude = env

	module, err := inter.LoadModuleIfNeeded(ast.ModuleName("prelude"), resolvedModule)
	if err != nil {
		fmt.Fprintf(world.Current.Stderr, "error: prelude not loaded\n    %s\n", err)
	}
	// These declares override the ones in the prelude.
	env.Parent = &Environment{Parent: nil, Scope: module.Environment.Scope, Unexported: module.Environment.Unexported}

	return env
}
