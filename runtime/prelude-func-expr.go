package runtime

import (
	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeValue = PreludeFuncExpr{}

type PreludeFuncExpr struct {
	Environment *Environment
	Decl        ast.ExprFunc
}

func MakePreludeFuncExpr(env *Environment, expr ast.ExprFunc) PreludeFuncExpr {
	return PreludeFuncExpr{env, expr}
}

func (PreludeFuncExpr) Lookup(member string) (Evaluatable, *RuntimeError) {
	panic("TODO: not implemented")
}

func (PreludeFuncExpr) RuntimeType() RuntimeTypeRef {
	return PreludeFunctionTypeRef
}

func (PreludeFuncExpr) String() string {
	panic("TODO: not implemented")
}
