package runtime

import (
	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeValue = PreludeTypeSwitchExpr{}

type PreludeTypeSwitchExpr struct {
	Decl ast.ExprTypeSwitch
}

func (PreludeTypeSwitchExpr) Lookup(member string) (Evaluatable, error) {
	panic("TODO: not implemented")
}

func (PreludeTypeSwitchExpr) RuntimeType() RuntimeTypeRef {
	return PreludeFunctionTypeRef
}

func (PreludeTypeSwitchExpr) String() string {
	panic("TODO: not implemented")
}
