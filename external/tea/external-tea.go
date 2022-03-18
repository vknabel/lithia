package tea

import (
	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/runtime"
)

var _ runtime.ExternalDefinition = ExternalTea{}

type ExternalTea struct{}

func (e ExternalTea) Lookup(name string, env *runtime.Environment, decl ast.Decl) (runtime.RuntimeValue, bool) {
	switch name {
	case "Program":
		if decl, ok := decl.(ast.DeclExternType); ok {
			return TeaProgramType{decl}, true
		} else {
			panic("tea.Program must be an extern type")
		}
	case "Model":
		if decl, ok := decl.(ast.DeclExternType); ok {
			return TeaModelType{decl, env}, true
		} else {
			panic("tea.Model must be an extern type")
		}
	default:
		return nil, false
	}
}
