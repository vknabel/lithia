package interpreter

import (
	"sort"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/vknabel/go-lithia/parser"
)

func (interpreter *Interpreter) AcceptSourceFile(sourceFile *sitter.Node) (interface{}, error) {
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

func (i *Interpreter) AcceptPackageDeclaration(packageDeclaration *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptImportDeclaration(importDeclaration *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptLetDeclaration(letDeclaration *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptFunctionDeclaration(functionDeclaration *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptDataDeclaration(dataDeclaration *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptDataPropertyList(dataPropertyList *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptDataPropertyValue(dataPropertyValue *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptDataPropertyFunction(dataPropertyFunction *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptEnumDeclaration(enumDeclaration *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptEnumCaseList(enumCaseList *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptComplexInvocationExpression(complexInvocationExpression *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptSimpleInvocationExpression(simpleInvocationExpression *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptUnaryExpression(unaryExpression *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptBinaryExpression(binaryExpression *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptMemberAccess(memberAccess *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptTypeExpression(typeExpression *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptTypeBody(typeBody *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptTypeCase(typeCase *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptStringLiteral(stringLiteral *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptEscapeSequence(escapeSequence *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptGroupLiteral(groupLiteral *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptNumberLiteral(numberLiteral *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptArrayLiteral(arrayLiteral *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptFunctionLiteral(functionLiteral *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptParameterList(parameterList *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptIdentifier(identifier *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptComment(comment *sitter.Node) (interface{}, error) {
	return nil, nil
}

// alias
func (i *Interpreter) AcceptEnumCaseReference(enumCaseReference *sitter.Node) (interface{}, error) {
	return nil, nil
}

// errors
func (i *Interpreter) AcceptError(error *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptUnexpected(unexpected *sitter.Node) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) AcceptUnknown(unknown *sitter.Node) (interface{}, error) {
	return nil, nil
}
