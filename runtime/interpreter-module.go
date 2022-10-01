package runtime

import (
	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/resolution"
)

type FileName string

type RuntimeModule struct {
	Name        ast.ModuleName
	Environment *Environment
	Files       map[FileName]*InterpreterContext
	resolved    resolution.ResolvedModule

	Decl *ast.ContextModule

	// docs can be derived from the files
}

func (inter *Interpreter) NewModule(resolvedModule resolution.ResolvedModule) (*RuntimeModule, error) {
	env, err := inter.NewPreludeEnvironment(resolvedModule)
	if err != nil {
		return nil, err
	}
	name := resolvedModule.AbsoluteModuleName()
	module := &RuntimeModule{
		Name:        name,
		Environment: NewEnvironment(env),
		Files:       make(map[FileName]*InterpreterContext),
		Decl:        ast.MakeContextModule(name),
		resolved:    resolvedModule,
	}
	inter.Modules[name] = module
	return module, nil
}
