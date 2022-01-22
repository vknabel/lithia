package runtime

import "github.com/vknabel/go-lithia/ast"

func Call(function RuntimeValue, args []Evaluatable, fromExpr ast.Expr) (RuntimeValue, *RuntimeError) {
	callable, ok := function.(CallableRuntimeValue)
	if !ok {
		return nil, NewRuntimeErrorf("%s is not callable", function).CascadeCall(nil, fromExpr)
	}

	arity := callable.Arity()
	if len(args) < arity {
		return MakeCurriedCallable(callable, args), nil
	}
	intermediate, err := callable.Call(args[:arity], fromExpr)
	if err != nil {
		return nil, err.CascadeCall(callable, fromExpr)
	}
	if len(args) == arity {
		return intermediate, nil
	}
	if g, ok := intermediate.(CallableRuntimeValue); ok {
		return g.Call(args[arity:], fromExpr)
	} else {
		return nil, NewRuntimeErrorf("%s is not callable", g).CascadeCall(nil, fromExpr)
	}
}
