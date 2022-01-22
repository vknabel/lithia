package runtime

import "github.com/vknabel/go-lithia/ast"

type FileName string

// TODO: Rename to Runtime Module
type Module struct {
	Name        ast.ModuleName
	Environment *Environment
	Files       map[FileName]*InterpreterContext

	Decl *ast.ContextModule

	// docs can be derived from the files
}

func (inter *Interpreter) NewModule(name ast.ModuleName) *Module {
	module := &Module{
		Name:        name,
		Environment: NewEnvironment(inter.NewPreludeEnvironment()),
		Files:       make(map[FileName]*InterpreterContext),
		Decl:        ast.MakeContextModule(name),
	}
	inter.Modules[name] = module
	return module
}
