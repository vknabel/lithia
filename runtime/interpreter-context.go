package runtime

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/vknabel/lithia/ast"
)

type InterpreterContext struct {
	interpreter *Interpreter
	fileDef     *ast.SourceFile
	module      *RuntimeModule

	path        []string
	environment *Environment

	evalCache *LazyEvaluationCache
}

func (inter *Interpreter) NewInterpreterContext(fileDef *ast.SourceFile, module *RuntimeModule, node *sitter.Node, source []byte, environment *Environment) *InterpreterContext {
	if environment == nil {
		environment = NewEnvironment(module.Environment.Private())
	}
	return &InterpreterContext{
		interpreter: inter,
		fileDef:     fileDef,
		module:      module,
		path:        []string{},
		environment: environment,
		evalCache:   NewLazyEvaluationCache(inter.Context),
	}
}

func (i *InterpreterContext) NestedInterpreterContext(name string) *InterpreterContext {
	return &InterpreterContext{
		interpreter: i.interpreter,
		fileDef:     i.fileDef,
		module:      i.module,
		path:        append(i.path, name),
		environment: NewEnvironment(i.environment),
		evalCache:   NewLazyEvaluationCache(i.interpreter.Context),
	}
}

func (i *InterpreterContext) Evaluate() (RuntimeValue, *RuntimeError) {
	return i.evalCache.Evaluate(func() (RuntimeValue, *RuntimeError) {
		if len(i.fileDef.Statements) == 0 {
			return nil, nil
		}
		var result RuntimeValue
		for _, stmt := range i.fileDef.Statements {
			select {
			case <-i.interpreter.Context.Done():
				return nil, NewRuntimeError(i.interpreter.Context.Err())
			default:
			}
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
