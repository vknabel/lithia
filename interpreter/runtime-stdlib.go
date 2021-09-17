package interpreter

import (
	"fmt"
	"os"
)

func (inter *Interpreter) NewPreludeEnvironment() *Environment {
	if inter.prelude != nil {
		return inter.prelude
	}
	env := NewEnvironment(nil)
	inter.prelude = env

	module, err := inter.LoadModuleIfNeeded(ModuleName("prelude"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: prelude not loaded\n    %s\n", err)
	}
	// These declares override the ones in the prelude.
	env.Parent = &Environment{Parent: nil, Scope: module.environment.Scope, Unexported: module.environment.Unexported}

	return env
}
