package runtime

import (
	"fmt"

	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeValue = PreludeFuncExpr{}
var _ CallableRuntimeValue = PreludeFuncExpr{}

type PreludeFuncExpr struct {
	context *InterpreterContext
	Decl    ast.ExprFunc
}

func MakePreludeFuncExpr(context *InterpreterContext, expr ast.ExprFunc) PreludeFuncExpr {
	fx := context.NestedInterpreterContext(string(expr.Name))
	for _, decl := range expr.Declarations {
		switch decl := decl.(type) {
		case ast.DeclConstant, ast.DeclFunc:
			continue
		default:
			fx.environment.DeclareExported(string(decl.DeclName()), MakeRuntimeValueDecl(fx, decl))
		}
	}
	return PreludeFuncExpr{
		fx,
		expr,
	}
}

func (PreludeFuncExpr) Lookup(member string) (Evaluatable, *RuntimeError) {
	panic("TODO: not implemented PreludeFuncExpr")
}

func (PreludeFuncExpr) RuntimeType() RuntimeTypeRef {
	return PreludeFunctionTypeRef
}

func (f PreludeFuncExpr) String() string {
	return fmt.Sprintf("%s", f.Decl)
	// panic("TODO: not implemented PreludeFuncExpr")
}

func (f PreludeFuncExpr) Arity() int {
	return len(f.Decl.Parameters)
}

func (f PreludeFuncExpr) Call(args []Evaluatable) (RuntimeValue, *RuntimeError) {
	arity := f.Arity()
	if arity > len(args) {
		return MakeCurriedCallable(f, args), nil
	}

	ex := f.context.NestedInterpreterContext("()")
	for _, decl := range f.Decl.Declarations {
		switch decl := decl.(type) {
		case ast.DeclConstant, ast.DeclFunc:
			ex.environment.DeclareExported(string(decl.DeclName()), MakeRuntimeValueDecl(ex, decl))
		default:
			continue
		}
	}
	for i, param := range f.Decl.Parameters {
		ex.environment.DeclareUnexported(string(param.Name), args[i])
	}
	var value RuntimeValue
	for _, expr := range f.Decl.Expressions {
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
