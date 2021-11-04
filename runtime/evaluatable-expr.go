package runtime

import (
	"fmt"

	"github.com/vknabel/go-lithia/ast"
)

type EvaluationContext struct {
	*Environment
}

type EvaluatableExpr struct {
	Context *EvaluationContext
	Expr    ast.Expr
	cache   *LazyEvaluationCache
}

func MakeEvaluatableExpr(context *EvaluationContext, expr ast.Expr) EvaluatableExpr {
	return EvaluatableExpr{context, expr, NewLazyEvaluationCache()}
}

func (e EvaluatableExpr) Evaluate() (RuntimeValue, *RuntimeError) {
	return e.cache.Evaluate(func() (RuntimeValue, *RuntimeError) {
		if e.Expr == nil {
			panic("cannot evaluate nil expr")
		}
		switch expr := e.Expr.(type) {
		case ast.ExprArray:
			return e.EvaluateExprArray(expr)
		case ast.ExprFloat:
			return e.EvaluateExprFloat(expr)
		case ast.ExprFunc:
			return e.EvaluateExprFunc(expr)
		case ast.ExprGroup:
			return e.EvaluateExprGroup(expr)
		case ast.ExprIdentifier:
			return e.EvaluateExprIdentifier(expr)
		case ast.ExprInt:
			return e.EvaluateExprInt(expr)
		case ast.ExprInvocation:
			return e.EvaluateExprInvocation(expr)
		case ast.ExprMemberAccess:
			return e.EvaluateExprMemberAccess(expr)
		case ast.ExprOperatorBinary:
			return e.EvaluateExprOperatorBinary(expr)
		case ast.ExprOperatorUnary:
			return e.EvaluateExprOperatorUnary(expr)
		case ast.ExprString:
			return e.EvaluateExprString(expr)
		case ast.ExprTypeSwitch:
			return e.EvaluateExprTypeSwitch(expr)
		default:
			panic(fmt.Sprintf("unknown expr: %s", expr))
		}
	})
}

func (e EvaluatableExpr) EvaluateExprArray(expr ast.ExprArray) (RuntimeValue, *RuntimeError) {
	panic("not implemented EvaluateExprArray")
}

func (e EvaluatableExpr) EvaluateExprFloat(expr ast.ExprFloat) (RuntimeValue, *RuntimeError) {
	return PreludeFloat(expr.Literal), nil
}

func (e EvaluatableExpr) EvaluateExprFunc(expr ast.ExprFunc) (RuntimeValue, *RuntimeError) {
	return MakePreludeFuncExpr(e.Context.Environment, expr), nil
}

func (e EvaluatableExpr) EvaluateExprGroup(expr ast.ExprGroup) (RuntimeValue, *RuntimeError) {
	var inner ast.Expr = expr
	return MakeEvaluatableExpr(e.Context, inner).Evaluate()
}

func (e EvaluatableExpr) EvaluateExprIdentifier(expr ast.ExprIdentifier) (RuntimeValue, *RuntimeError) {
	if unevaluated, ok := e.Context.GetPrivte(string(expr.Name)); ok {
		return unevaluated.Evaluate()
	} else {
		return nil, NewRuntimeError(fmt.Errorf("undeclared %s", expr.Name), *expr.Meta().Source)
	}
}

func (e EvaluatableExpr) EvaluateExprInt(expr ast.ExprInt) (RuntimeValue, *RuntimeError) {
	return PreludeInt(expr.Literal), nil
}

func (e EvaluatableExpr) EvaluateExprInvocation(expr ast.ExprInvocation) (RuntimeValue, *RuntimeError) {
	panic("not implemented EvaluateExprInvocation")
}

func (e EvaluatableExpr) EvaluateExprMemberAccess(expr ast.ExprMemberAccess) (RuntimeValue, *RuntimeError) {
	panic("not implemented EvaluateExprMemberAccess")
}

func (e EvaluatableExpr) EvaluateExprOperatorBinary(expr ast.ExprOperatorBinary) (RuntimeValue, *RuntimeError) {
	panic("not implemented EvaluateExprOperatorBinary")
}

func (e EvaluatableExpr) EvaluateExprOperatorUnary(expr ast.ExprOperatorUnary) (RuntimeValue, *RuntimeError) {
	panic("not implemented EvaluateExprOperatorUnary")
}

func (e EvaluatableExpr) EvaluateExprString(expr ast.ExprString) (RuntimeValue, *RuntimeError) {
	return PreludeString(expr.Literal), nil
}

func (e EvaluatableExpr) EvaluateExprTypeSwitch(expr ast.ExprTypeSwitch) (RuntimeValue, *RuntimeError) {
	return MakePreludeTypeSwitchExpr(e.Context.Environment, expr), nil
}
