package interpreter

type ModuleName string
type FileName string

type Module struct {
	name              ModuleName
	environment       *Environment
	executionContexts map[FileName]*EvaluationContext
	docs              DocString
}

func (inter *Interpreter) NewModule(name ModuleName) *Module {
	module := &Module{
		name:              name,
		environment:       NewEnvironment(inter.NewPreludeEnvironment()),
		executionContexts: make(map[FileName]*EvaluationContext),
	}
	inter.modules[name] = module
	return module
}
