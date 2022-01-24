package parser

import "github.com/vknabel/go-lithia/ast"

func (fp *FileParser) ParseGroupExpr() (*ast.ExprGroup, []SyntaxError) {
	exprNode := fp.Node.ChildByFieldName("expression")
	expr, errors := fp.ChildParserConsumingComments(exprNode).ParseExpression()
	return ast.MakeExprGroup(expr, fp.AstSource()), errors
}
