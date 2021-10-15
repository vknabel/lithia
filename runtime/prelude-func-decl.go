package runtime

import (
	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeValue = PreludeFuncDecl{}

type PreludeFuncDecl struct {
	Decl ast.DeclFunc
}

func (PreludeFuncDecl) Lookup(member string) (Evaluatable, error) {
	panic("TODO: not implemented")
}

func (PreludeFuncDecl) RuntimeType() RuntimeTypeRef {
	return PreludeFunctionTypeRef
}

func (PreludeFuncDecl) String() string {
	panic("TODO: not implemented")
}
