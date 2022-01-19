package parser

import (
	"github.com/vknabel/go-lithia/ast"
)

func (fp *FileParser) ParseExpressionIfGiven() (ast.Expr, []SyntaxError) {
	switch fp.Node.Type() {
	case TYPE_NODE_COMPLEX_INVOCATION_EXPRESSION:
		expr, errs := fp.ParseInvocationExpr()
		if expr == nil {
			return nil, errs
		} else {
			return expr, errs
		}
	case TYPE_NODE_SIMPLE_INVOCATION_EXPRESSION:
		expr, errs := fp.ParseInvocationExpr()
		if expr == nil {
			return nil, errs
		} else {
			return expr, errs
		}
	case TYPE_NODE_UNARY_EXPRESSION:
		expr, errs := fp.ParseUnaryExpr()
		if expr == nil {
			return nil, errs
		} else {
			return expr, errs
		}
	case TYPE_NODE_BINARY_EXPRESSION:
		expr, errs := fp.ParseBinaryExpr()
		if expr == nil {
			return nil, errs
		} else {
			return expr, errs
		}
	case TYPE_NODE_MEMBER_ACCESS:
		expr, errs := fp.ParseExprMemberAccess()
		if expr == nil {
			return nil, errs
		} else {
			return expr, errs
		}
	case TYPE_NODE_TYPE_EXPRESSION:
		expr, errs := fp.ParseExprTypeSwitch()
		if expr == nil {
			return nil, errs
		} else {
			return expr, errs
		}
	case TYPE_NODE_TYPE_BODY:
		panic("not implemented")
	case TYPE_NODE_TYPE_CASE:
		panic("not implemented")
	case TYPE_NODE_STRING_LITERAL:
		expr, errs := fp.ParseExprString()
		if expr == nil {
			return nil, errs
		} else {
			return expr, errs
		}
	// case TYPE_NODE_ESCAPE_SEQUENCE:
	// 	panic("not implemented")
	case TYPE_NODE_GROUP_LITERAL:
		expr, errs := fp.ParseGroupExpr()
		if expr == nil {
			return nil, errs
		} else {
			return expr, errs
		}
	case TYPE_NODE_NUMBER_LITERAL:
		expr, errs := fp.ParseIntExpr()
		if expr == nil {
			return nil, errs
		} else {
			return expr, errs
		}
	case TYPE_NODE_ARRAY_LITERAL:
		expr, errs := fp.ParseExprArray()
		if expr == nil {
			return nil, errs
		} else {
			return expr, errs
		}
	case TYPE_NODE_FUNCTION_LITERAL:
		expr, errs := fp.ParseFunctionExpr()
		if expr == nil {
			return nil, errs
		} else {
			return expr, errs
		}
	case TYPE_NODE_PARAMETER_LIST:
		panic("not implemented")
	case TYPE_NODE_IDENTIFIER:
		expr, errs := fp.parseExprIdentifier()
		if expr == nil {
			return nil, errs
		} else {
			return expr, errs
		}

	default:
		return nil, nil
	}
}

func (fp *FileParser) ParseExpression() (ast.Expr, []SyntaxError) {
	expr, errs := fp.ParseExpressionIfGiven()
	if expr != nil || len(errs) > 0 {
		return expr, errs
	}
	return nil, []SyntaxError{fp.SyntaxErrorf("expected expression, got %s", fp.Node.Type())}
}
