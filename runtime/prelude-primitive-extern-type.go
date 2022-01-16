package runtime

import "github.com/vknabel/go-lithia/ast"

var _ RuntimeValue = PreludePrimitiveExternType{}

type PreludePrimitiveExternType struct {
	*ast.DeclExternType
}

func (t PreludePrimitiveExternType) String() string {
	return string(t.DeclExternType.Name)
}

func (t PreludePrimitiveExternType) Lookup(name string) (Evaluatable, *RuntimeError) {
	return nil, NewRuntimeErrorf("type %s has no member %s", t.DeclExternType.Name, name)
}

func (t PreludePrimitiveExternType) RuntimeType() RuntimeTypeRef {
	return PreludeAnyTypeRef
}
