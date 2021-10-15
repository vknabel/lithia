package runtime

import "github.com/vknabel/go-lithia/ast"

var _ RuntimeValue = PreludeExternFunction{}

type PreludeExternFunction struct {
	Decl *ast.DeclExternFunc
	Impl func(args []Evaluatable) (RuntimeValue, error)
}

func (PreludeExternFunction) Lookup(member string) (Evaluatable, error) {
	panic("TODO: not implemented")
}

func (PreludeExternFunction) RuntimeType() RuntimeTypeRef {
	return PreludeFunctionTypeRef
}

func (PreludeExternFunction) String() string {
	panic("TODO: not implemented")
}
