package runtime

import "github.com/vknabel/lithia/ast"

type ExternalDefinition interface {
	Lookup(name string, env *Environment, decl ast.Decl) (RuntimeValue, bool)
}
