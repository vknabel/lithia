package runtime

import (
	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeValue = PreludeFuncDecl{}
var _ CallableRuntimeValue = PreludeFuncDecl{}

type PreludeFuncDecl struct {
	context *InterpreterContext
	Decl    ast.DeclFunc
}

func MakePreludeFuncDecl(context *InterpreterContext, decl ast.DeclFunc) PreludeFuncDecl {
	fx := context.NestedInterpreterContext(string(decl.DeclName()))
	for _, decl := range decl.Impl.Declarations {
		fx.environment.DeclareExported(string(decl.DeclName()), MakeRuntimeValueDecl(fx, decl))
	}
	return PreludeFuncDecl{
		fx,
		decl,
	}
}

func (PreludeFuncDecl) Lookup(member string) (Evaluatable, *RuntimeError) {
	panic("TODO: not implemented")
}

func (PreludeFuncDecl) RuntimeType() RuntimeTypeRef {
	return PreludeFunctionTypeRef
}

func (PreludeFuncDecl) String() string {
	panic("TODO: not implemented")
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
