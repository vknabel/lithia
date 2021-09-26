package parser

import (
	"fmt"

	"github.com/vknabel/go-lithia/ast"
)

func (fp *FileParser) ParseFunctionExpr() (*ast.ExprFunc, []SyntaxError) {
	parametersNode := fp.Node.ChildByFieldName("parameters")
	bodyNode := fp.Node.ChildByFieldName("body")

	errors := []SyntaxError{}
	params, paramsErrors := fp.ChildParser(parametersNode).ParseParameterDeclarationList()
	if paramsErrors != nil {
		errors = append(errors, paramsErrors...)
	}
	if params == nil {
		return nil, errors
	}
	function := ast.MakeExprFunc(params, fp.AstSource())
	bodyParser := fp.ChildParser(bodyNode)
	for i := 0; i < int(bodyNode.NamedChildCount()); i++ {
		child := bodyNode.NamedChild(i)
		if bodyParser.ParseChildCommentIfNeeded(child) {
			continue
		}

		childParser := bodyParser.ChildParserConsumingComments(child)
		decls, declErrors := childParser.ParseDeclsIfGiven()
		if declErrors != nil {
			errors = append(errors, declErrors...)
		}
		if decls != nil {
			for _, decl := range decls {
				function.AddDecl(&decl)
			}
			continue
		}

		expr, exprErrors := childParser.ParseExpressionIfGiven()
		if exprErrors != nil {
			errors = append(errors, exprErrors...)
		}
		if expr != nil {
			function.AddExpr(&expr)
			continue
		}

		errors = append(errors, fmt.Errorf("unknown child type"))
	}

	return function, errors
}
