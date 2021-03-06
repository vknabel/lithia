package parser

import (
	"github.com/vknabel/lithia/ast"
)

func (fp *FileParser) ParseExprArray() (*ast.ExprArray, []SyntaxError) {
	numberOfElements := int(fp.Node.NamedChildCount())
	elements := make([]ast.Expr, 0, numberOfElements)
	for i := 0; i < numberOfElements; i++ {
		elementNode := fp.Node.NamedChild(i)
		if elementNode.Type() == TYPE_NODE_COMMENT {
			fp.Comments = append(fp.Comments, elementNode.Content(fp.Source))
			continue
		}
		expr, errs := fp.ChildParserConsumingComments(elementNode).ParseExpression()
		if len(errs) > 0 {
			return nil, errs
		}
		elements = append(elements, expr)
	}
	return ast.MakeExprArray(elements, fp.AstSource()), nil
}
