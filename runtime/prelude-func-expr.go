package runtime

import (
	"fmt"
	"strings"

	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeValue = PreludeFuncExpr{}
var _ CallableRuntimeValue = PreludeFuncExpr{}

type PreludeFuncExpr struct {
	context *InterpreterContext
	Decl    ast.ExprFunc
}

func MakePreludeFuncExpr(context *InterpreterContext, expr ast.ExprFunc) (PreludeFuncExpr, *RuntimeError) {
	if context.fileDef.Path != expr.Meta().Source.FileName {
		panic("Mixing files in function expressions!")
	}
	fx := context.NestedInterpreterContext(string(expr.Name))
	for _, decl := range expr.Declarations {
		switch decl := decl.(type) {
		case ast.DeclConstant, ast.DeclFunc:
			continue
		default:
			eval, err := MakeRuntimeValueDecl(fx, decl)
			if err != nil {
				return PreludeFuncExpr{}, err
			}
			fx.environment.DeclareExported(string(decl.DeclName()), eval)
		}
	}
	return PreludeFuncExpr{
		fx,
		expr,
	}, nil
}

func (f PreludeFuncExpr) Lookup(member string) (Evaluatable, *RuntimeError) {
	switch member {
	case "arity":
		return NewConstantRuntimeValue(PreludeInt(f.Arity())), nil
	default:
		return nil, NewRuntimeErrorf("no such member: %s", member)
	}
}

func (PreludeFuncExpr) RuntimeType() RuntimeTypeRef {
	return PreludeFunctionTypeRef
}

func (f PreludeFuncExpr) String() string {
	paramList := make([]string, len(f.Decl.Parameters))
	for i, param := range f.Decl.Parameters {
		paramList[i] = string(param.Name)
	}
	return fmt.Sprintf("<func %s.%s %s>", f.context.module.Name, strings.Join(f.context.path, "."), strings.Join(paramList, ", "))
}

func (f PreludeFuncExpr) Arity() int {
	return len(f.Decl.Parameters)
}

func (f PreludeFuncExpr) Call(args []Evaluatable, fromExpr ast.Expr) (RuntimeValue, *RuntimeError) {
	if len(args) != f.Arity() {
		panic("use Call to call functions!")
	}
	ex := f.context.NestedInterpreterContext("()")
	for _, decl := range f.Decl.Declarations {
		switch decl := decl.(type) {
		case ast.DeclConstant, ast.DeclFunc:
			eval, err := MakeRuntimeValueDecl(ex, decl)
			if err != nil {
				return nil, err
			}
			ex.environment.DeclareExported(string(decl.DeclName()), eval)
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
	return value, nil
}

func (f PreludeFuncExpr) Source() *ast.Source {
	return f.Decl.Meta().Source
}
