package runtime

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeValue = PreludeDataDecl{}
var _ DeclRuntimeValue = PreludeDataDecl{}
var _ CallableRuntimeValue = PreludeDataDecl{}

type PreludeDataDecl struct {
	Decl ast.DeclData
}

func (d PreludeDataDecl) Lookup(member string) (Evaluatable, *RuntimeError) {
	panic(fmt.Sprintf("cannot access member %s of data type %s, see https://github.com/vknabel/lithia/discussions/25", member, d.Decl.Name))
}

func (PreludeDataDecl) RuntimeType() RuntimeTypeRef {
	return PreludeAnyTypeRef
}

func (d PreludeDataDecl) String() string {
	return fmt.Sprintf("data %s", d.Decl.Name)
}

func (d PreludeDataDecl) HasInstance(interpreter *Interpreter, value RuntimeValue) (bool, *RuntimeError) {
	if dataVal, ok := value.(DataRuntimeValue); ok {
		return reflect.DeepEqual(*dataVal.TypeDecl, d), nil
	} else {
		return false, nil
	}
}

func (d PreludeDataDecl) Arity() int {
	return len(d.Decl.Fields)
}

func (d PreludeDataDecl) Call(args []Evaluatable) (RuntimeValue, *RuntimeError) {
	arity := d.Arity()
	if arity == 0 && len(args) == 0 {
		dataVal, err := MakeDataRuntimeValueMemberwise(&d, make(map[string]Evaluatable))
		return dataVal, err.CascadeDecl(d.Decl)
	}
	if arity > len(args) {
		return MakeCurriedCallable(d, args), nil
	}
	members := make(map[string]Evaluatable)
	for i, field := range d.Decl.Fields {
		members[string(field.DeclName())] = args[i]
	}
	dataVal, err := MakeDataRuntimeValueMemberwise(&d, members)
	if err != nil {
		return dataVal, err.CascadeDecl(d.Decl)
	}
	if len(args) > arity {
		stringifiedArgs := make([]string, len(args[arity:]))
		for i, arg := range args {
			stringifiedArgs[i] = fmt.Sprintf("%s", arg)
		}
		return nil, NewRuntimeErrorf(
			"cannot call non-function value of type %s: %q with args %s",
			d.Decl.Name,
			dataVal,
			strings.Join(stringifiedArgs, ", "),
		).CascadeDecl(d.Decl)
	}
	return dataVal, nil
}
