package runtime

import (
	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeValue = PreludeFuncDecl{}

type PreludeFuncDecl struct {
	Environment *Environment
	Decl        ast.DeclFunc
}

func MakePreludeFuncDecl(env *Environment, decl ast.DeclFunc) PreludeFuncDecl {
	return PreludeFuncDecl{env, decl}
}

func (PreludeFuncDecl) Lookup(member string) (Evaluatable, *RuntimeError) {
	panic("TODO: not implemented")
}

func (PreludeFuncDecl) RuntimeType() RuntimeTypeRef {
	return PreludeFunctionTypeRef
}

func (PreludeFuncDecl) String() string {
	panic("TODO: not implemented")
}
