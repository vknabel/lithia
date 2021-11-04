package runtime

import (
	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeValue = PreludeEnumDecl{}

type PreludeEnumDecl struct {
	Decl ast.DeclEnum
}

func (PreludeEnumDecl) Lookup(member string) (Evaluatable, *RuntimeError) {
	panic("TODO: not implemented")
}

func (PreludeEnumDecl) RuntimeType() RuntimeTypeRef {
	return PreludeAnyTypeRef
}

func (PreludeEnumDecl) String() string {
	panic("TODO: not implemented")
}
