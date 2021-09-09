package parser

import (
	sitter "github.com/smacker/go-tree-sitter"
)

var (
	TYPE_NODE_SOURCE_FILE                   = "source_file"
	TYPE_NODE_PACKAGE_DECLARATION           = "package_declaration"
	TYPE_NODE_IMPORT_DECLARATION            = "import_declaration"
	TYPE_NODE_LET_DECLARATION               = "let_declaration"
	TYPE_NODE_FUNCTION_DECLARATION          = "function_declaration"
	TYPE_NODE_DATA_DECLARATION              = "data_declaration"
	TYPE_NODE_DATA_PROPERTY_LIST            = "data_property_list"
	TYPE_NODE_DATA_PROPERTY_VALUE           = "data_property_value"
	TYPE_NODE_DATA_PROPERTY_FUNCTION        = "data_property_function"
	TYPE_NODE_ENUM_DECLARATION              = "enum_declaration"
	TYPE_NODE_ENUM_CASE_LIST                = "enum_case_list"
	TYPE_NODE_COMPLEX_INVOCATION_EXPRESSION = "complex_invocation_expression"
	TYPE_NODE_SIMPLE_INVOCATION_EXPRESSION  = "simple_invocation_expression"
	TYPE_NODE_UNARY_EXPRESSION              = "unary_expression"
	TYPE_NODE_BINARY_EXPRESSION             = "binary_expression"
	TYPE_NODE_MEMBER_ACCESS                 = "member_access"
	TYPE_NODE_TYPE_EXPRESSION               = "type_expression"
	TYPE_NODE_TYPE_BODY                     = "type_body"
	TYPE_NODE_TYPE_CASE                     = "type_case"
	TYPE_NODE_STRING_LITERAL                = "string_literal"
	TYPE_NODE_ESCAPE_SEQUENCE               = "escape_sequence"
	TYPE_NODE_GROUP_LITERAL                 = "group_literal"
	TYPE_NODE_NUMBER_LITERAL                = "number_literal"
	TYPE_NODE_ARRAY_LITERAL                 = "array_literal"
	TYPE_NODE_FUNCTION_LITERAL              = "function_literal"
	TYPE_NODE_PARAMETER_LIST                = "parameter_list"
	TYPE_NODE_IDENTIFIER                    = "identifier"
	TYPE_NODE_COMMENT                       = "comment"

	// alias
	TYPE_NODE_ENUM_CASE_REFERENCE = "enum_case_reference"

	// errors
	TYPE_NODE_ERROR      = "ERROR"
	TYPE_NODE_UNEXPECTED = "UNEXPECTED"
)

type NodeVisitor interface {
	AcceptSourceFile(sourceFile *sitter.Node) (interface{}, error)
	AcceptPackageDeclaration(packageDeclaration *sitter.Node) (interface{}, error)
	AcceptImportDeclaration(importDeclaration *sitter.Node) (interface{}, error)
	AcceptLetDeclaration(letDeclaration *sitter.Node) (interface{}, error)
	AcceptFunctionDeclaration(functionDeclaration *sitter.Node) (interface{}, error)
	AcceptDataDeclaration(dataDeclaration *sitter.Node) (interface{}, error)
	AcceptDataPropertyList(dataPropertyList *sitter.Node) (interface{}, error)
	AcceptDataPropertyValue(dataPropertyValue *sitter.Node) (interface{}, error)
	AcceptDataPropertyFunction(dataPropertyFunction *sitter.Node) (interface{}, error)
	AcceptEnumDeclaration(enumDeclaration *sitter.Node) (interface{}, error)
	AcceptEnumCaseList(enumCaseList *sitter.Node) (interface{}, error)
	AcceptComplexInvocationExpression(complexInvocationExpression *sitter.Node) (interface{}, error)
	AcceptSimpleInvocationExpression(simpleInvocationExpression *sitter.Node) (interface{}, error)
	AcceptUnaryExpression(unaryExpression *sitter.Node) (interface{}, error)
	AcceptBinaryExpression(binaryExpression *sitter.Node) (interface{}, error)
	AcceptMemberAccess(memberAccess *sitter.Node) (interface{}, error)
	AcceptTypeExpression(typeExpression *sitter.Node) (interface{}, error)
	AcceptTypeBody(typeBody *sitter.Node) (interface{}, error)
	AcceptTypeCase(typeCase *sitter.Node) (interface{}, error)
	AcceptStringLiteral(stringLiteral *sitter.Node) (interface{}, error)
	AcceptEscapeSequence(escapeSequence *sitter.Node) (interface{}, error)
	AcceptGroupLiteral(groupLiteral *sitter.Node) (interface{}, error)
	AcceptNumberLiteral(numberLiteral *sitter.Node) (interface{}, error)
	AcceptArrayLiteral(arrayLiteral *sitter.Node) (interface{}, error)
	AcceptFunctionLiteral(functionLiteral *sitter.Node) (interface{}, error)
	AcceptParameterList(parameterList *sitter.Node) (interface{}, error)
	AcceptIdentifier(identifier *sitter.Node) (interface{}, error)
	AcceptComment(comment *sitter.Node) (interface{}, error)

	// alias
	AcceptEnumCaseReference(enumCaseReference *sitter.Node) (interface{}, error)

	// errors error
	AcceptError(error *sitter.Node) (interface{}, error)
	AcceptUnexpected(unexpected *sitter.Node) (interface{}, error)
	AcceptUnknown(unknown *sitter.Node) (interface{}, error)
}

func Accept(visitor NodeVisitor, node *sitter.Node) (interface{}, error) {
	switch node.Type() {
	case TYPE_NODE_SOURCE_FILE:
		return visitor.AcceptSourceFile(node)
	case TYPE_NODE_PACKAGE_DECLARATION:
		return visitor.AcceptPackageDeclaration(node)
	case TYPE_NODE_IMPORT_DECLARATION:
		return visitor.AcceptImportDeclaration(node)
	case TYPE_NODE_LET_DECLARATION:
		return visitor.AcceptLetDeclaration(node)
	case TYPE_NODE_FUNCTION_DECLARATION:
		return visitor.AcceptFunctionDeclaration(node)
	case TYPE_NODE_DATA_DECLARATION:
		return visitor.AcceptDataDeclaration(node)
	case TYPE_NODE_DATA_PROPERTY_LIST:
		return visitor.AcceptDataPropertyList(node)
	case TYPE_NODE_DATA_PROPERTY_VALUE:
		return visitor.AcceptDataPropertyValue(node)
	case TYPE_NODE_DATA_PROPERTY_FUNCTION:
		return visitor.AcceptDataPropertyFunction(node)
	case TYPE_NODE_ENUM_DECLARATION:
		return visitor.AcceptEnumDeclaration(node)
	case TYPE_NODE_ENUM_CASE_LIST:
		return visitor.AcceptEnumCaseList(node)
	case TYPE_NODE_COMPLEX_INVOCATION_EXPRESSION:
		return visitor.AcceptComplexInvocationExpression(node)
	case TYPE_NODE_SIMPLE_INVOCATION_EXPRESSION:
		return visitor.AcceptSimpleInvocationExpression(node)
	case TYPE_NODE_UNARY_EXPRESSION:
		return visitor.AcceptUnaryExpression(node)
	case TYPE_NODE_BINARY_EXPRESSION:
		return visitor.AcceptBinaryExpression(node)
	case TYPE_NODE_MEMBER_ACCESS:
		return visitor.AcceptMemberAccess(node)
	case TYPE_NODE_TYPE_EXPRESSION:
		return visitor.AcceptTypeExpression(node)
	case TYPE_NODE_TYPE_BODY:
		return visitor.AcceptTypeBody(node)
	case TYPE_NODE_TYPE_CASE:
		return visitor.AcceptTypeCase(node)
	case TYPE_NODE_STRING_LITERAL:
		return visitor.AcceptStringLiteral(node)
	case TYPE_NODE_ESCAPE_SEQUENCE:
		return visitor.AcceptEscapeSequence(node)
	case TYPE_NODE_GROUP_LITERAL:
		return visitor.AcceptGroupLiteral(node)
	case TYPE_NODE_NUMBER_LITERAL:
		return visitor.AcceptNumberLiteral(node)
	case TYPE_NODE_ARRAY_LITERAL:
		return visitor.AcceptArrayLiteral(node)
	case TYPE_NODE_FUNCTION_LITERAL:
		return visitor.AcceptFunctionLiteral(node)
	case TYPE_NODE_PARAMETER_LIST:
		return visitor.AcceptParameterList(node)
	case TYPE_NODE_IDENTIFIER:
		return visitor.AcceptIdentifier(node)
	case TYPE_NODE_COMMENT:
		return visitor.AcceptComment(node)

	// alias
	case TYPE_NODE_ENUM_CASE_REFERENCE:
		return visitor.AcceptEnumCaseReference(node)

	// error
	case TYPE_NODE_ERROR:
		return visitor.AcceptError(node)
	case TYPE_NODE_UNEXPECTED:
		return visitor.AcceptUnexpected(node)

	default:
		return visitor.AcceptUnknown(node)
	}
}
