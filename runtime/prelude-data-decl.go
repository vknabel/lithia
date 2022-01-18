package runtime

import (
	"reflect"

	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeValue = PreludeDataDecl{}
var _ DeclRuntimeValue = PreludeDataDecl{}

type PreludeDataDecl struct {
	Decl ast.DeclData
}

func (PreludeDataDecl) Lookup(member string) (Evaluatable, *RuntimeError) {
	panic("TODO: data not implemented")
}

func (PreludeDataDecl) RuntimeType() RuntimeTypeRef {
	return PreludeAnyTypeRef
}

func (PreludeDataDecl) String() string {
	panic("TODO: data not implemented")
}

func (d PreludeDataDecl) HasInstance(value RuntimeValue) (bool, *RuntimeError) {
	if dataVal, ok := value.(DataRuntimeValue); ok {
		return reflect.DeepEqual(*dataVal.TypeDecl, d), nil
	} else {
		return false, nil
	}
}
