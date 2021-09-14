package interpreter

import (
	"sort"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/vknabel/go-lithia/parser"
)

func (interpreter *ExecutionContext) AcceptSourceFile(sourceFile *sitter.Node) (interface{}, error) {
	count := sourceFile.ChildCount()
	children := make([]*sitter.Node, count)
	for i := uint32(0); i < count; i++ {
		child := sourceFile.Child(int(i))
		children = append(children, child)
	}
	sort.Slice(children, func(i, j int) bool {
		lp := priority(children[i].Type())
		rp := priority(children[j].Type())
		return lp < rp
	})

	for _, child := range children {
		_, err := parser.Accept(interpreter, child)
		if err != nil {
			return nil, err
		}

	}

	// now run all sorted children

	return nil, nil
}

func priority(nodeType string) int {
	switch nodeType {
	case parser.TYPE_NODE_PACKAGE_DECLARATION:
		return 19
	case parser.TYPE_NODE_IMPORT_DECLARATION:
		return 17
	case parser.TYPE_NODE_DATA_DECLARATION:
		return 15
	case parser.TYPE_NODE_ENUM_DECLARATION:
		return 13
	case parser.TYPE_NODE_FUNCTION_DECLARATION:
		return 7
	case parser.TYPE_NODE_LET_DECLARATION:
		return 3
	default:
		return 0
	}
}

func (i *ExecutionContext) AcceptPackageDeclaration(packageDeclaration *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptImportDeclaration(importDeclaration *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptLetDeclaration(letDeclaration *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptFunctionDeclaration(functionDeclaration *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptDataDeclaration(dataDeclaration *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptDataPropertyList(dataPropertyList *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptDataPropertyValue(dataPropertyValue *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptDataPropertyFunction(dataPropertyFunction *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptEnumDeclaration(enumDeclaration *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptEnumCaseList(enumCaseList *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptComplexInvocationExpression(complexInvocationExpression *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptSimpleInvocationExpression(simpleInvocationExpression *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptUnaryExpression(unaryExpression *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptBinaryExpression(binaryExpression *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptMemberAccess(memberAccess *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptTypeExpression(typeExpression *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptTypeBody(typeBody *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptTypeCase(typeCase *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptStringLiteral(stringLiteral *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptEscapeSequence(escapeSequence *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptGroupLiteral(groupLiteral *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptNumberLiteral(numberLiteral *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptArrayLiteral(arrayLiteral *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptFunctionLiteral(functionLiteral *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptParameterList(parameterList *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptIdentifier(identifier *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptComment(comment *sitter.Node) (interface{}, error) {
	return nil, nil
}

// alias
func (i *ExecutionContext) AcceptEnumCaseReference(enumCaseReference *sitter.Node) (interface{}, error) {
	return nil, nil
}

// errors
func (i *ExecutionContext) AcceptError(error *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptMissing(error *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptUnexpected(unexpected *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *ExecutionContext) AcceptUnknown(unknown *sitter.Node) (interface{}, error) {
	return nil, nil
}
