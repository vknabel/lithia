package interpreter

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/vknabel/go-lithia/parser"
)

type Interpreter struct {
	root *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		root: NewEnvironment(nil),
	}
}

func (interpreter *Interpreter) Interpret(script *sitter.Tree, source []byte) (*LazyRuntimeValue, error) {
	if script.RootNode().Type() != parser.TYPE_NODE_SOURCE_FILE {
		return nil, fmt.Errorf("expected source file, got node of type %s", script.RootNode().Type())
	}
	return interpreter.EvaluateSourceFile(script.RootNode(), source)
}

func (interpreter *Interpreter) EvaluateSourceFile(sourceFile *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	count := sourceFile.ChildCount()
	children := make([]*sitter.Node, count)
	for i := uint32(0); i < count; i++ {
		child := sourceFile.Child(int(i))
		children[i] = child
	}
	sort.SliceStable(children, func(i, j int) bool {
		lp := priority(children[i].Type())
		rp := priority(children[j].Type())
		return lp > rp
	})

	var lastValue RuntimeValue
	for _, child := range children {
		lazyValue, err := interpreter.EvaluateNode(child, source)
		if err != nil {
			return nil, err
		}
		if lazyValue != nil {
			lastValue, err = lazyValue.Evaluate()
			if err != nil {
				return nil, err
			}
		}
	}
	return NewConstantRuntimeValue(lastValue), nil
}

func (interpreter *Interpreter) EvaluateLetDeclaration(letDecl *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	nameNode := letDecl.ChildByFieldName("name")
	valueNode := letDecl.ChildByFieldName("value")
	if nameNode == nil || valueNode == nil {
		return nil, fmt.Errorf("let declaration must have name and value")
	}
	lazyValue, err := interpreter.EvaluateNode(valueNode, source)
	if err != nil {
		return nil, err
	}
	err = interpreter.root.Declare(nameNode.Content([]byte(source)), lazyValue)
	if err != nil {
		return nil, err
	}
	return lazyValue, nil
}
func (interpreter *Interpreter) EvaluateFunctionDeclaration(funcDecl *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	return nil, nil
}
func (interpreter *Interpreter) EvaluateDataDeclaration(dataDecl *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	name := dataDecl.ChildByFieldName("name").Content(source)
	propertiesNode := dataDecl.ChildByFieldName("properties")

	var numberOfFields int
	if propertiesNode != nil {
		numberOfFields = int(propertiesNode.ChildCount())
	} else {
		numberOfFields = 0
	}

	data := DataDeclRuntimeValue{
		name:   name,
		fields: make([]DataDeclField, numberOfFields),
	}

	for i := 0; i < numberOfFields; i++ {
		child := propertiesNode.Child(i)
		switch child.Type() {
		case parser.TYPE_NODE_DATA_PROPERTY_VALUE:
			name := child.ChildByFieldName("name").Content(source)
			data.fields[i] = DataDeclField{name: name}
		case parser.TYPE_NODE_DATA_PROPERTY_FUNCTION:
			name := child.ChildByFieldName("name").Content(source)
			parameters, error := interpreter.ParseParamterList(child.ChildByFieldName("parameters"), source)
			if error != nil {
				return nil, error
			}
			data.fields[i] = DataDeclField{name: name, params: parameters}
		default:
			return nil, fmt.Errorf("unexpected node type %s", child.Type())
		}
	}

	constantValue := NewConstantRuntimeValue((data))
	interpreter.root.Declare(name, constantValue)
	return constantValue, nil
}

func (interpreter *Interpreter) ParseParamterList(list *sitter.Node, source []byte) ([]string, error) {
	params := make([]string, list.ChildCount())
	for i := 0; i < int(list.ChildCount()); i++ {
		child := list.Child(i)
		params[i] = child.Content(source)
	}
	return params, nil
}

func (interpreter *Interpreter) EvaluateNumberLiteral(numberDecl *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	literal := numberDecl.Content(source)
	integer, err := strconv.ParseInt(literal, 10, 64)
	if err != nil {
		return nil, err
	}
	return NewConstantRuntimeValue(PreludeInt(integer)), nil
}
func (interpreter *Interpreter) EvaluateComplexInvocationExpr(expr *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	functionNode := expr.ChildByFieldName("function")
	lazyFunc, err := interpreter.EvaluateNode(functionNode, source)
	if err != nil {
		return nil, err
	}
	functionValue, err := lazyFunc.Evaluate()
	if err != nil {
		return nil, err
	}
	function, ok := reflect.ValueOf(functionValue).Interface().(Callable)
	if !ok {
		return nil, fmt.Errorf("expected callable, got %T", functionValue)
	}

	args := make([]*LazyRuntimeValue, expr.ChildCount()-1)
	for i := 0; i < int(expr.ChildCount()-1); i++ {
		child := expr.Child(i + 1)
		lazyValue, err := interpreter.EvaluateNode(child, source)
		if err != nil {
			return nil, err
		}
		args[i] = lazyValue
	}

	return function.Call(args)
}

func (interpreter *Interpreter) EvaluateIdentifier(node *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	string := node.Content(source)
	if value, ok := interpreter.root.Get(string); ok {
		return value, nil
	} else {
		return nil, fmt.Errorf("undefined identifier %s", string)
	}
}

func (interpreter *Interpreter) EvaluateNode(node *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	fmt.Printf("evaluating node of type %s\n", node.Type())

	switch node.Type() {
	case parser.TYPE_NODE_SOURCE_FILE:
		return interpreter.EvaluateSourceFile(node, source)
	// case parser.TYPE_NODE_PACKAGE_DECLARATION:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_IMPORT_DECLARATION:
	// 	return interpreter.Evaluate(node)
	case parser.TYPE_NODE_LET_DECLARATION:
		return interpreter.EvaluateLetDeclaration(node, source)
	case parser.TYPE_NODE_FUNCTION_DECLARATION:
		return interpreter.EvaluateFunctionDeclaration(node, source)
	case parser.TYPE_NODE_DATA_DECLARATION:
		return interpreter.EvaluateDataDeclaration(node, source)
	// case parser.TYPE_NODE_DATA_PROPERTY_LIST:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_DATA_PROPERTY_VALUE:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_DATA_PROPERTY_FUNCTION:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_ENUM_DECLARATION:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_ENUM_CASE_LIST:
	// 	return interpreter.Evaluate(node)
	case parser.TYPE_NODE_COMPLEX_INVOCATION_EXPRESSION:
		return interpreter.EvaluateComplexInvocationExpr(node, source)
	// case parser.TYPE_NODE_SIMPLE_INVOCATION_EXPRESSION:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_UNARY_EXPRESSION:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_BINARY_EXPRESSION:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_MEMBER_ACCESS:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_TYPE_EXPRESSION:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_TYPE_BODY:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_TYPE_CASE:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_STRING_LITERAL:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_ESCAPE_SEQUENCE:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_GROUP_LITERAL:
	// 	return interpreter.Evaluate(node)
	case parser.TYPE_NODE_NUMBER_LITERAL:
		return interpreter.EvaluateNumberLiteral(node, source)
	// case parser.TYPE_NODE_ARRAY_LITERAL:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_FUNCTION_LITERAL:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_PARAMETER_LIST:
	// 	return interpreter.Evaluate(node)
	case parser.TYPE_NODE_IDENTIFIER:
		return interpreter.EvaluateIdentifier(node, source)
	// case parser.TYPE_NODE_COMMENT:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_ENUM_CASE_REFERENCE:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_ERROR:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_UNEXPECTED:
	// 	return interpreter.Evaluate(node)
	default:
		return nil, nil
	}
}
