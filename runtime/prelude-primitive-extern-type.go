package runtime

import "github.com/vknabel/lithia/ast"

var _ RuntimeValue = PreludePrimitiveExternType{}
var _ DeclRuntimeValue = PreludePrimitiveExternType{}

type PreludePrimitiveExternType struct {
	*ast.DeclExternType
	hasInstance func(value RuntimeValue) (bool, *RuntimeError)
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

func (t PreludePrimitiveExternType) HasInstance(interpreter *Interpreter, value RuntimeValue) (bool, *RuntimeError) {
	return t.hasInstance(value)
}
