package runtime

import (
	"fmt"
	"strings"

	"github.com/vknabel/go-lithia/ast"
)

type EvaluatableExpr struct {
	Context *InterpreterContext
	Expr    ast.Expr
	cache   *LazyEvaluationCache
}

func MakeEvaluatableExpr(context *InterpreterContext, expr ast.Expr) EvaluatableExpr {
	copy := *context
	if context.fileDef.Path != expr.Meta().Source.FileName {
		panic("Mixing files in declared evaluatable expr!")
	}
	return EvaluatableExpr{&copy, expr, NewLazyEvaluationCache()}
}

func (e EvaluatableExpr) Evaluate() (RuntimeValue, *RuntimeError) {
	if e.Context.fileDef.Path != e.Expr.Meta().Source.FileName {
		panic("Mixing files in declared evaluatable expr!")
	}
	value, err := e.cache.Evaluate(func() (value RuntimeValue, error *RuntimeError) {
		if e.Context.fileDef.Path != e.Expr.Meta().Source.FileName {
			panic("Mixing files in declared evaluatable expr!")
		}
		defer func() {
			if err := recover(); err != nil {
				error = NewRuntimeError(fmt.Errorf("panic: %q", err))
			}
		}()

		if e.Expr == nil {
			panic("cannot evaluate nil expr")
		}
		switch expr := e.Expr.(type) {
		case *ast.ExprArray:
			return e.EvaluateExprArray(*expr)
		case *ast.ExprFloat:
			return e.EvaluateExprFloat(*expr)
		case *ast.ExprFunc:
			return e.EvaluateExprFunc(*expr)
		case *ast.ExprGroup:
			return e.EvaluateExprGroup(*expr)
		case *ast.ExprIdentifier:
			return e.EvaluateExprIdentifier(*expr)
		case *ast.ExprInt:
			return e.EvaluateExprInt(*expr)
		case *ast.ExprInvocation:
			return e.EvaluateExprInvocation(*expr)
		case *ast.ExprMemberAccess:
			return e.EvaluateExprMemberAccess(*expr)
		case *ast.ExprOperatorBinary:
			return e.EvaluateExprOperatorBinary(*expr)
		case *ast.ExprOperatorUnary:
			return e.EvaluateExprOperatorUnary(*expr)
		case *ast.ExprString:
			return e.EvaluateExprString(*expr)
		case *ast.ExprTypeSwitch:
			return e.EvaluateExprTypeSwitch(*expr)
		default:
			panic(fmt.Sprintf("unknown expr: %T %s", expr, expr))
		}
	})

	return value, err.Cascade(*e.Expr.Meta().Source)
}

func (e EvaluatableExpr) EvaluateExprArray(expr ast.ExprArray) (RuntimeValue, *RuntimeError) {
	evaluatables := make([]Evaluatable, len(expr.Elements))
	for i, element := range expr.Elements {
		evaluatables[i] = MakeEvaluatableExpr(e.Context, element)
	}
	list, err := e.Context.environment.MakeList(evaluatables)
	return list, err.Cascade(*expr.Meta().Source)
}

func (e EvaluatableExpr) EvaluateExprFloat(expr ast.ExprFloat) (RuntimeValue, *RuntimeError) {
	return PreludeFloat(expr.Literal), nil
}

func (e EvaluatableExpr) EvaluateExprFunc(expr ast.ExprFunc) (RuntimeValue, *RuntimeError) {
	if e.Context.fileDef.Path != expr.Meta().Source.FileName {
		panic("Mixing files in function expressions!")
	}
	return MakePreludeFuncExpr(e.Context, expr), nil
}

func (e EvaluatableExpr) EvaluateExprGroup(expr ast.ExprGroup) (RuntimeValue, *RuntimeError) {
	var inner ast.Expr = expr.Expr
	return MakeEvaluatableExpr(e.Context, inner).Evaluate()
}

func (e EvaluatableExpr) EvaluateExprIdentifier(expr ast.ExprIdentifier) (RuntimeValue, *RuntimeError) {
	if unevaluated, ok := e.Context.environment.GetPrivate(string(expr.Name)); ok {
		value, err := unevaluated.Evaluate()
		if err != nil {
			return nil, err
		}
		if fun, ok := value.(CallableRuntimeValue); ok {
			if fun.Arity() == 0 {
				return fun.Call(nil)
			}
		}
		return value, nil
	} else {
		return nil, NewRuntimeErrorf("undeclared %s in %s", expr.Name, strings.Join(e.Context.path, "."))
	}
}

func (e EvaluatableExpr) EvaluateExprInt(expr ast.ExprInt) (RuntimeValue, *RuntimeError) {
	return PreludeInt(expr.Literal), nil
}

func (e EvaluatableExpr) EvaluateExprInvocation(expr ast.ExprInvocation) (RuntimeValue, *RuntimeError) {
	function, err := MakeEvaluatableExpr(e.Context, expr.Function).Evaluate()
	if err != nil {
		return nil, err
	}
	callable, ok := function.(CallableRuntimeValue)
	if !ok {
		return nil, NewRuntimeErrorf("cannot call %T %s", function, function)
	}
	args := make([]Evaluatable, len(expr.Arguments))
	for i, argExpr := range expr.Arguments {
		args[i] = MakeEvaluatableExpr(e.Context, *argExpr)
	}
	return callable.Call(args)
}

func (e EvaluatableExpr) EvaluateExprMemberAccess(expr ast.ExprMemberAccess) (RuntimeValue, *RuntimeError) {
	evaluatableTargetExpr := MakeEvaluatableExpr(e.Context, expr.Target)
	target, err := evaluatableTargetExpr.Evaluate()
	if err != nil {
		return nil, err
	}

	var evaluatableTarget Evaluatable
	for _, identifier := range expr.AccessPath {
		evaluatableTarget, err = target.Lookup(string(identifier))
		if err != nil {
			return nil, err
		}
		target, err = evaluatableTarget.Evaluate()
		if err != nil {
			return nil, err
		}
	}
	return target, nil
}

func (e EvaluatableExpr) EvaluateExprOperatorBinary(expr ast.ExprOperatorBinary) (RuntimeValue, *RuntimeError) {
	impl, err := e.Context.BinaryOperatorFunction(string(expr.Operator))
	if err != nil {
		return nil, err
	}
	leftEvalExpr := MakeEvaluatableExpr(e.Context, expr.Left)
	rightEvalExpr := MakeEvaluatableExpr(e.Context, expr.Right)

	return impl(leftEvalExpr, rightEvalExpr)
}

func (e EvaluatableExpr) EvaluateExprOperatorUnary(expr ast.ExprOperatorUnary) (RuntimeValue, *RuntimeError) {
	panic("TODO: not implemented EvaluateExprOperatorUnary")
}

func (e EvaluatableExpr) EvaluateExprString(expr ast.ExprString) (RuntimeValue, *RuntimeError) {
	return PreludeString(expr.Literal), nil
}

func (e EvaluatableExpr) EvaluateExprTypeSwitch(expr ast.ExprTypeSwitch) (RuntimeValue, *RuntimeError) {
	return MakePreludeTypeSwitchExpr(e.Context, expr), nil
}
