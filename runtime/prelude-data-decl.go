package runtime

import (
	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeValue = PreludeDataDecl{}

type PreludeDataDecl struct {
	Decl ast.DeclData
}

func (PreludeDataDecl) Lookup(member string) (Evaluatable, *RuntimeError) {
	panic("TODO: not implemented")
}

func (PreludeDataDecl) RuntimeType() RuntimeTypeRef {
	return PreludeAnyTypeRef
}

func (PreludeDataDecl) String() string {
	panic("TODO: not implemented")
}
