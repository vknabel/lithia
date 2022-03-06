package parser

import (
	"github.com/vknabel/lithia/ast"
)

func (fp *FileParser) ParseUnaryExpr() (*ast.ExprOperatorUnary, []SyntaxError) {
	if fp.Node.NamedChildCount() != 1 {
		return nil, []SyntaxError{fp.SyntaxErrorf("expected one child, got %d", fp.Node.NamedChildCount())}
	}
	operator := fp.Node.ChildByFieldName("operator").Content(fp.Source)

	exprP := fp.SameScopeChildParser(fp.Node.NamedChild(0))
	expr, errs := exprP.ParseExpression()
	if len(errs) > 0 {
		return nil, errs
	}
	return ast.MakeExprOperatorUnary(ast.OperatorUnary(operator), expr, fp.AstSource()), nil
}
