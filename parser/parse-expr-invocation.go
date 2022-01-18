package parser

import (
	"github.com/vknabel/go-lithia/ast"
)

func (fp *FileParser) ParseInvocationExpr() (*ast.ExprInvocation, []SyntaxError) {
	errors := []SyntaxError{}
	functionParser := fp.ChildParser(fp.Node.ChildByFieldName("function"))
	functionExpr, functionErrors := functionParser.ParseExpression()

	if functionErrors != nil {
		errors = append(errors, functionErrors...)
	}
	if functionExpr == nil {
		return nil, errors
	}

	function := ast.MakeExprInvocation(functionExpr, functionParser.AstSource())
	for i := 1; i < int(fp.Node.NamedChildCount()); i++ {
		child := fp.Node.NamedChild(i)
		if fp.ParseChildCommentIfNeeded(child) {
			continue
		}
		childParser := fp.ChildParser(child)
		expr, childErrs := childParser.ParseExpression()
		if childErrs != nil {
			errors = append(errors, childErrs...)
		}
		if expr != nil {
			function.AddArgument(expr)
		}
	}
	return function, errors
}
