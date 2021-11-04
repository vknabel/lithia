package runtime

import (
	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeValue = PreludeTypeSwitchExpr{}

type PreludeTypeSwitchExpr struct {
	Environment *Environment
	Decl        ast.ExprTypeSwitch
}

func MakePreludeTypeSwitchExpr(environment *Environment, decl ast.ExprTypeSwitch) PreludeTypeSwitchExpr {
	return PreludeTypeSwitchExpr{environment, decl}
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
