package rx

import (
	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/runtime"
)

var _ runtime.ExternalDefinition = ExternalRx{}

type ExternalRx struct {
	inter *runtime.Interpreter
}

func New(inter *runtime.Interpreter) ExternalRx {
	return ExternalRx{inter}
}

func (e ExternalRx) Lookup(name string, env *runtime.Environment, decl ast.Decl) (runtime.RuntimeValue, bool) {
	switch name {
	case "Variable":
		if decl, ok := decl.(ast.DeclExternType); ok {
			return RxVariableType{decl}, true
		} else {
			panic("rx.Variable must be an extern type")
		}
	case "Future":
		if decl, ok := decl.(ast.DeclExternType); ok {
			return RxFutureType{decl, e}, true
		} else {
			panic("rx.Future must be an extern type")
		}
	default:
		return nil, false
	}
}
