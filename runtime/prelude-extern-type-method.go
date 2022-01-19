package runtime

import (
	"fmt"
	"strings"

	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeValue = PreludeExternTypeMethod{}
var _ CallableRuntimeValue = PreludeExternTypeMethod{}

type PreludeExternTypeMethod struct {
	Decl ast.DeclField
	Impl func(args []Evaluatable) (RuntimeValue, *RuntimeError)
}

func MakeExternTypeMethod(
	decl ast.Decl,
	impl func(args []Evaluatable) (RuntimeValue, *RuntimeError),
) PreludeExternTypeMethod {
	externDecl, ok := decl.(ast.DeclField)
	if !ok {
		panic(fmt.Errorf("extern func declaration requires func definition: %T %s", decl, decl.DeclName()))
	}
	if len(externDecl.Parameters) == 0 {
		panic(fmt.Errorf("extern func declaration requires at least one param: %T %s", decl, decl.DeclName()))
	}
	return PreludeExternTypeMethod{
		Decl: externDecl,
		Impl: impl,
	}
}

func (PreludeExternTypeMethod) Lookup(member string) (Evaluatable, *RuntimeError) {
	panic("TODO: not implemented PreludeExternFunction")
}

func (PreludeExternTypeMethod) RuntimeType() RuntimeTypeRef {
	return PreludeFunctionTypeRef
}

func (f PreludeExternTypeMethod) String() string {
	argNames := make([]string, len(f.Decl.Parameters))
	for i, param := range f.Decl.Parameters {
		argNames[i] = string(param.Name)
	}
	return fmt.Sprintf("(extern %s %s)", f.Decl.Name, strings.Join(argNames, ", "))
}

func (f PreludeExternTypeMethod) Arity() int {
	return len(f.Decl.Parameters)
}

func (f PreludeExternTypeMethod) Call(args []Evaluatable) (RuntimeValue, *RuntimeError) {
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
