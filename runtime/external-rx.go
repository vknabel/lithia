package runtime

import (
	"github.com/vknabel/lithia/ast"
)

var _ ExternalDefinition = ExternalRx{}

type ExternalRx struct{}

func (e ExternalRx) Lookup(name string, env *Environment, decl ast.Decl) (RuntimeValue, bool) {
	switch name {
	case "Variable":
		if decl, ok := decl.(ast.DeclExternType); ok {
			return RxVariableType{decl}, true
		} else {
			panic("rx.Variable must be an extern type")
		}
	case "Future":
		if decl, ok := decl.(ast.DeclExternType); ok {
			return RxFutureType{decl}, true
		} else {
			panic("rx.Future must be an extern type")
		}
	default:
		return nil, false
	}
}
