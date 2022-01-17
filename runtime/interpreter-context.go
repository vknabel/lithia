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

	evalCache *LazyEvaluationCache
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
		evalCache:   NewLazyEvaluationCache(),
	}
}

func (i *InterpreterContext) NestedInterpreterContext(name string) *InterpreterContext {
	return &InterpreterContext{
		interpreter: i.interpreter,
		fileDef:     i.fileDef,
		module:      i.module,
		path:        append(i.path, name),
		environment: NewEnvironment(i.environment),
		evalCache:   NewLazyEvaluationCache(),
	}
}

func (i *InterpreterContext) Evaluate() (RuntimeValue, *RuntimeError) {
	return i.evalCache.Evaluate(func() (RuntimeValue, *RuntimeError) {
		if len(i.fileDef.Statements) == 0 {
			return nil, nil
		}
		var result RuntimeValue
		for _, stmt := range i.fileDef.Statements {
			expr := MakeEvaluatableExpr(i, stmt)
			value, err := expr.Evaluate()
			if err != nil {
				return nil, err
			}
			result = value
		}
		return result, nil
	})
}
