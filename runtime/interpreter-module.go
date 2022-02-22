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

func (inter *Interpreter) NewModule(resolvedModule resolution.ResolvedModule) *RuntimeModule {
	name := resolvedModule.AbsoluteModuleName()
	module := &RuntimeModule{
		Name:        name,
		Environment: NewEnvironment(inter.NewPreludeEnvironment(resolvedModule)),
		Files:       make(map[FileName]*InterpreterContext),
		Decl:        ast.MakeContextModule(name),
		resolved:    resolvedModule,
	}
	inter.Modules[name] = module
	return module
}
