package parser

import (
	"github.com/vknabel/go-lithia/ast"
)

func (fp *FileParser) ParseBinaryExpr() (*ast.ExprOperatorBinary, []SyntaxError) {
	if fp.Node.NamedChildCount() != 2 {
		return nil, []SyntaxError{fp.SyntaxErrorf("expected 2 children, got %d", fp.Node.NamedChildCount())}
	}
	operator := fp.Node.ChildByFieldName("operator").Content(fp.Source)

	leftP := fp.ChildParser(fp.Node.NamedChild(0))
	left, lerrs := leftP.ParseExpression()
	if len(lerrs) > 0 {
		return nil, lerrs
	}
	rightP := fp.ChildParser(fp.Node.NamedChild(1))
	right, rerrs := rightP.ParseExpression()
	if len(rerrs) > 0 {
		return nil, rerrs
	}
	return ast.MakeExprOperatorBinary(ast.OperatorBinary(operator), left, right, fp.AstSource()), nil
}
