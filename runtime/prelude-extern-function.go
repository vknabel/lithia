package runtime

import (
	"fmt"
	"strings"

	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeValue = PreludeExternFunction{}
var _ CallableRuntimeValue = PreludeExternFunction{}

type PreludeExternFunction struct {
	Decl ast.DeclExternFunc
	Impl func(args []Evaluatable) (RuntimeValue, *RuntimeError)
}

func MakeExternFunction(
	decl ast.Decl,
	impl func(args []Evaluatable) (RuntimeValue, *RuntimeError),
) PreludeExternFunction {
	externDecl, ok := decl.(ast.DeclExternFunc)
	if !ok {
		panic(fmt.Errorf("extern func declaration requires func definition: %T %s", decl, decl.DeclName()))
	}
	return PreludeExternFunction{
		Decl: externDecl,
		Impl: impl,
	}
}

func (f PreludeExternFunction) Lookup(member string) (Evaluatable, *RuntimeError) {
	switch member {
	case "arity":
		return NewConstantRuntimeValue(PreludeInt(f.Arity())), nil
	default:
		return nil, NewRuntimeErrorf("no such member: %s", member)
	}
}

func (PreludeExternFunction) RuntimeType() RuntimeTypeRef {
	return PreludeFunctionTypeRef
}

func (f PreludeExternFunction) String() string {
	argNames := make([]string, len(f.Decl.Parameters))
	for i, param := range f.Decl.Parameters {
		argNames[i] = string(param.Name)
	}
	return fmt.Sprintf("(extern %s %s)", f.Decl.Name, strings.Join(argNames, ", "))
}

func (f PreludeExternFunction) Arity() int {
	return len(f.Decl.Parameters)
}

func (f PreludeExternFunction) Call(args []Evaluatable) (RuntimeValue, *RuntimeError) {
	if len(args) < len(f.Decl.Parameters) {
		return MakeCurriedCallable(f, args), nil
	}
	intermediate, err := f.Impl(args[:len(f.Decl.Parameters)])
	if err != nil {
		return nil, err
	}
	if len(args) == len(f.Decl.Parameters) {
		return intermediate, nil
	}
	if g, ok := intermediate.(CallableRuntimeValue); ok {
		return g.Call(args[len(f.Decl.Parameters):])
	} else {
		return nil, NewRuntimeErrorf("%s is not callable", g)
	}
}
