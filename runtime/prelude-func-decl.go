package runtime

import (
	"fmt"
	"strings"

	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeValue = PreludeFuncDecl{}
var _ CallableRuntimeValue = PreludeFuncDecl{}

type PreludeFuncDecl struct {
	context *InterpreterContext
	Decl    ast.DeclFunc
}

func MakePreludeFuncDecl(context *InterpreterContext, decl ast.DeclFunc) (PreludeFuncDecl, *RuntimeError) {
	if context.fileDef.Path != decl.Meta().FileName {
		panic("Mixing files in declared functions!")
	}
	fx := context.NestedInterpreterContext(string(decl.DeclName()))
	for _, decl := range decl.Impl.Declarations {
		switch decl := decl.(type) {
		case ast.DeclConstant, ast.DeclFunc:
			continue
		default:
			declValue, err := MakeRuntimeValueDecl(fx, decl)
			if err != nil {
				return PreludeFuncDecl{}, err
			}
			fx.environment.DeclareExported(string(decl.DeclName()), declValue)
		}
	}

	return PreludeFuncDecl{
		fx,
		decl,
	}, nil
}

func (f PreludeFuncDecl) Lookup(member string) (Evaluatable, *RuntimeError) {
	switch member {
	case "arity":
		return NewConstantRuntimeValue(PreludeInt(f.Arity())), nil
	default:
		return nil, NewRuntimeErrorf("no such member: %s", member)
	}
}

func (PreludeFuncDecl) RuntimeType() RuntimeTypeRef {
	return PreludeFunctionTypeRef
}

func (f PreludeFuncDecl) String() string {
	paramList := make([]string, len(f.Decl.Impl.Declarations))
	for i, param := range f.Decl.Impl.Parameters {
		paramList[i] = string(param.Name)
	}

	return fmt.Sprintf("<func %s %s>", f.Decl.Name, strings.Join(paramList, ", "))
}

func (f PreludeFuncDecl) Arity() int {
	return len(f.Decl.Impl.Parameters)
}

func (f PreludeFuncDecl) Call(args []Evaluatable) (RuntimeValue, *RuntimeError) {
	arity := f.Arity()
	if arity > len(args) {
		return MakeCurriedCallable(f, args), nil
	}

	ex := f.context.NestedInterpreterContext("()")
	for _, decl := range f.Decl.Impl.Declarations {
		switch decl := decl.(type) {
		case ast.DeclConstant, ast.DeclFunc:
			declValue, err := MakeRuntimeValueDecl(ex, decl)
			if err != nil {
				return nil, err
			}
			ex.environment.DeclareExported(string(decl.DeclName()), declValue)
		default:
			continue
		}
	}
	for i, param := range f.Decl.Impl.Parameters {
		ex.environment.DeclareUnexported(string(param.Name), args[i])
	}
	var value RuntimeValue
	for _, expr := range f.Decl.Impl.Expressions {
		var err *RuntimeError
		value, err = MakeEvaluatableExpr(ex, expr).Evaluate()
		if err != nil {
			return nil, err
		}
	}

	if arity == len(args) {
		return value, nil
	}
	intermediate, ok := value.(CallableRuntimeValue)
	if !ok {
		return nil, NewRuntimeErrorf("cannot call %T %s", intermediate, intermediate)
	}
	remainingArgs := args[arity:]
	return intermediate.Call(remainingArgs)
}
