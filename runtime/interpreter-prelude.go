package runtime

import (
	"fmt"
	"os"

	"github.com/vknabel/lithia/ast"
)

func (inter *Interpreter) NewPreludeEnvironment() *Environment {
	if inter.Prelude != nil {
		return inter.Prelude
	}
	env := NewEnvironment(nil)
	inter.Prelude = env

	module, err := inter.LoadModuleIfNeeded(ast.ModuleName("prelude"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: prelude not loaded\n    %s\n", err)
	}
	// These declares override the ones in the prelude.
	env.Parent = &Environment{Parent: nil, Scope: module.Environment.Scope, Unexported: module.Environment.Unexported}

	return env
}
