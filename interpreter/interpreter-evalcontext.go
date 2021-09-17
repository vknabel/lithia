package interpreter

import sitter "github.com/smacker/go-tree-sitter"

type EvaluationContext struct {
	interpreter   *Interpreter
	file          string
	module        *Module
	path          []string
	environment   *Environment
	functionCount int

	node   *sitter.Node
	source []byte
}

func (inter *Interpreter) NewEvaluationContext(file string, module *Module, node *sitter.Node, source []byte, environment *Environment) *EvaluationContext {
	if environment == nil {
		environment = NewEnvironment(inter.NewPreludeEnvironment())
	}
	return &EvaluationContext{
		interpreter:   inter,
		file:          file,
		module:        module,
		path:          []string{},
		environment:   environment,
		functionCount: 0,

		node:   node,
		source: source,
	}
}

func (i *EvaluationContext) NestedExecutionContext(name string) *EvaluationContext {
	return &EvaluationContext{
		interpreter:   i.interpreter,
		file:          i.file,
		module:        i.module,
		path:          append(i.path, name),
		environment:   NewEnvironment(i.environment),
		functionCount: 0,
		node:          i.node,
		source:        i.source,
	}
}

func (i *EvaluationContext) ChildNodeExecutionContext(childNode *sitter.Node) *EvaluationContext {
	return &EvaluationContext{
		interpreter:   i.interpreter,
		file:          i.file,
		module:        i.module,
		path:          i.path,
		environment:   i.environment,
		functionCount: i.functionCount,
		node:          childNode,
		source:        i.source,
	}
}
