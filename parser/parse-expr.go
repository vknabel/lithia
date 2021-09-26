package parser

import (
	"fmt"

	"github.com/vknabel/go-lithia/ast"
)

func (fp *FileParser) ParseExpressionIfGiven() (ast.Expr, []SyntaxError) {
	switch fp.Node.Type() {
	case TYPE_NODE_COMPLEX_INVOCATION_EXPRESSION:
		panic("not implemented")
	case TYPE_NODE_SIMPLE_INVOCATION_EXPRESSION:
		panic("not implemented")
	case TYPE_NODE_UNARY_EXPRESSION:
		panic("not implemented")
	case TYPE_NODE_BINARY_EXPRESSION:
		panic("not implemented")
	case TYPE_NODE_MEMBER_ACCESS:
		panic("not implemented")
	case TYPE_NODE_TYPE_EXPRESSION:
		panic("not implemented")
	case TYPE_NODE_TYPE_BODY:
		panic("not implemented")
	case TYPE_NODE_TYPE_CASE:
		panic("not implemented")
	case TYPE_NODE_STRING_LITERAL:
		panic("not implemented")
	case TYPE_NODE_ESCAPE_SEQUENCE:
		panic("not implemented")
	case TYPE_NODE_GROUP_LITERAL:
		panic("not implemented")
	case TYPE_NODE_NUMBER_LITERAL:
		panic("not implemented")
	case TYPE_NODE_ARRAY_LITERAL:
		panic("not implemented")
	case TYPE_NODE_FUNCTION_LITERAL:
		expr, errs := fp.ParseFunctionExpr()
		return expr, errs
	case TYPE_NODE_PARAMETER_LIST:
		panic("not implemented")
	case TYPE_NODE_IDENTIFIER:
		panic("not implemented")

	default:
		return nil, nil
	}
}

func (fp *FileParser) ParseExpression() (*ast.Expr, []SyntaxError) {
	expr, errs := fp.ParseExpressionIfGiven()
	if expr != nil || errs != nil {
		return &expr, errs
	}
	return nil, []SyntaxError{fmt.Errorf("expected expression, got %s", fp.Node.Type())}
}
