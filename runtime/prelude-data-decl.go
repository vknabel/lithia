package runtime

import (
	"fmt"
	"reflect"

	"github.com/vknabel/lithia/ast"
)

var _ RuntimeValue = PreludeDataDecl{}
var _ DeclRuntimeValue = PreludeDataDecl{}
var _ CallableRuntimeValue = PreludeDataDecl{}

type PreludeDataDecl struct {
	Decl ast.DeclData
}

func (d PreludeDataDecl) Lookup(member string) (Evaluatable, *RuntimeError) {
	return nil, NewRuntimeErrorf("cannot access member %s of data type %s, see https://github.com/vknabel/lithia/discussions/25", member, d.Decl.Name)
}

func (PreludeDataDecl) RuntimeType() RuntimeTypeRef {
	return PreludeAnyTypeRef
}

func (d PreludeDataDecl) String() string {
	return fmt.Sprintf("data %s", d.Decl.Name)
}

func (d PreludeDataDecl) HasInstance(value RuntimeValue) (bool, *RuntimeError) {
	if dataVal, ok := value.(DataRuntimeValue); ok {
		return reflect.DeepEqual(*dataVal.TypeDecl, d), nil
	} else {
		return false, nil
	}
}

func (d PreludeDataDecl) Arity() int {
	return len(d.Decl.Fields)
}

func (d PreludeDataDecl) Call(args []Evaluatable, fromExpr ast.Expr) (RuntimeValue, *RuntimeError) {
	if len(args) != d.Arity() {
		panic("use Call to call functions!")
	}
	if d.Arity() == 0 {
		dataVal, err := MakeDataRuntimeValueMemberwise(&d, make(map[string]Evaluatable))
		return dataVal, err.CascadeDecl(d.Decl)
	}
	members := make(map[string]Evaluatable)
	for i, field := range d.Decl.Fields {
		members[string(field.DeclName())] = args[i]
	}
	dataVal, err := MakeDataRuntimeValueMemberwise(&d, members)
	if err != nil {
		return dataVal, err.CascadeDecl(d.Decl)
	}
	return dataVal, nil
}

func (f PreludeDataDecl) Source() *ast.Source {
	return f.Decl.Meta().Source
}
