package runtime

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/vknabel/go-lithia/ast"
)

type InterpreterContext struct {
	interpreter *Interpreter
	fileDef     *ast.SourceFile
	module      *Module

	path        []string
	environment *Environment
}

func (inter *Interpreter) NewInterpreterContext(fileDef *ast.SourceFile, module *Module, node *sitter.Node, source []byte, environment *Environment) *InterpreterContext {
	if environment == nil {
		environment = NewEnvironment(module.Environment.Private())
	}
	return &InterpreterContext{
		interpreter: inter,
		fileDef:     fileDef,
		module:      module,
		path:        []string{},
		environment: environment,
	}
}

func (i *InterpreterContext) NestedInterpreterContext(name string) *InterpreterContext {
	return &InterpreterContext{
		interpreter: i.interpreter,
		fileDef:     i.fileDef,
		module:      i.module,
		path:        append(i.path, name),
		environment: NewEnvironment(i.environment),
	}
}
