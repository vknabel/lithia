package interpreter

import (
	"reflect"
	"sort"
	"strconv"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/vknabel/go-lithia/parser"
)

type Interpreter struct {
	path        []string
	environment *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		path:        []string{},
		environment: NewEnvironment(nil),
	}
}

func (i *Interpreter) ChildInterpreter(name string) *Interpreter {
	return &Interpreter{
		path:        append(i.path, name),
		environment: NewEnvironment(i.environment),
	}
}

func (interpreter *Interpreter) Interpret(script *sitter.Tree, source []byte) (*LazyRuntimeValue, error) {
	if script.RootNode().Type() != parser.TYPE_NODE_SOURCE_FILE {
		return nil, SyntaxErrorf(script.RootNode(), source, "expected source file, got node of type %s", script.RootNode().Type())
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
		return nil, SyntaxErrorf(letDecl, source, "let declaration must have name and value")
	}
	lazyValue, err := interpreter.EvaluateNode(valueNode, source)
	if err != nil {
		return nil, err
	}
	err = interpreter.environment.Declare(nameNode.Content([]byte(source)), lazyValue)
	if err != nil {
		return nil, err
	}
	return lazyValue, nil
}
func (interpreter *Interpreter) EvaluateFunctionDeclaration(funcDecl *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	name := funcDecl.ChildByFieldName("name").Content(source)
	functionNode := funcDecl.ChildByFieldName("function")
	function, err := interpreter.ParseFunctionLiteral(functionNode, source, interpreter.environment)
	function.name = name

	if err != nil {
		return nil, err
	}
	impl := NewConstantRuntimeValue(function)
	err = interpreter.environment.Declare(name, impl)
	if err != nil {
		return nil, err
	}
	return impl, nil
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
			return nil, SyntaxErrorf(child, source, "unexpected node type %s", child.Type())
		}
	}

	constantValue := NewConstantRuntimeValue(data)
	interpreter.environment.Declare(name, constantValue)
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

func (interpreter *Interpreter) ParseStatementList(list *sitter.Node, source []byte) ([]*LazyRuntimeValue, error) {
	stmts := make([]*LazyRuntimeValue, list.ChildCount())
	for i := 0; i < int(list.ChildCount()); i++ {
		child := list.Child(i)
		stmt, err := interpreter.EvaluateNode(child, source)
		if err != nil {
			return nil, err
		}
		stmts[i] = stmt
	}
	return stmts, nil
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
		return nil, RuntimeErrorf(expr, source, "expected callable, got %T", functionValue)
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

func (interpreter *Interpreter) EvaluateSimpleInvocation(expr *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	return interpreter.EvaluateComplexInvocationExpr(expr, source)
}

func (interpreter *Interpreter) EvaluateEnumDeclaration(enumDecl *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	name := enumDecl.ChildByFieldName("name").Content(source)
	casesNode := enumDecl.ChildByFieldName("cases")
	if casesNode == nil {
		enum := NewConstantRuntimeValue(EnumDeclRuntimeValue{name: name, cases: make(map[string]*LazyRuntimeValue)})
		interpreter.environment.Declare(name, enum)
		return enum, nil
	}
	caseCount := int(casesNode.ChildCount())
	cases := make(map[string]*LazyRuntimeValue)
	for i := 0; i < caseCount; i++ {
		child := casesNode.Child(i)
		switch child.Type() {
		case parser.TYPE_NODE_ENUM_CASE_REFERENCE:
			caseName := child.Content(source)
			lookedUp, _ := interpreter.environment.Get(caseName)
			if lookedUp == nil {
				return nil, RuntimeErrorf(child, source, "undefined enum case %s", caseName)
			}
			cases[caseName] = lookedUp
		case parser.TYPE_NODE_DATA_DECLARATION:
			caseName := child.ChildByFieldName("name").Content(source)
			runtimeValue, err := interpreter.EvaluateDataDeclaration(child, source)
			if err != nil {
				return nil, err
			}
			cases[caseName] = runtimeValue
		case parser.TYPE_NODE_ENUM_DECLARATION:
			caseName := child.ChildByFieldName("name").Content(source)
			runtimeValue, err := interpreter.EvaluateEnumDeclaration(child, source)
			if err != nil {
				return nil, err
			}
			cases[caseName] = runtimeValue
		default:
			return nil, SyntaxErrorf(child, source, "unexpected node type %s", child.Type())
		}
	}
	constantValue := NewConstantRuntimeValue(EnumDeclRuntimeValue{
		name:  name,
		cases: cases,
	})
	interpreter.environment.Declare(name, constantValue)
	return constantValue, nil
}

func (interpreter *Interpreter) EvaluateIdentifier(node *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	string := node.Content(source)
	return NewLazyRuntimeValue(func() (RuntimeValue, error) {
		if value, ok := interpreter.environment.Get(string); ok {
			return value.Evaluate()
		} else {
			return nil, RuntimeErrorf(node, source, "undefined identifier %s", string)
		}
	}), nil
}

func (interpreter *Interpreter) EvaluateMemberAccess(node *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	if node.NamedChildCount() < 2 {
		return nil, SyntaxErrorf(node, source, "expected at least 2 children, got %d", node.ChildCount())
	}
	literalNode := node.Child(0)
	lazyObject, err := interpreter.EvaluateNode(literalNode, source)
	if err != nil {
		return nil, err
	}

	keyPath := make([]string, node.NamedChildCount()-1)
	for i := 1; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)
		keyPath[i-1] = child.Content(source)
	}
	lazyResult := NewLazyRuntimeValue(func() (RuntimeValue, error) {
		objectValue, err := lazyObject.Evaluate()
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(keyPath); i++ {
			object, ok := objectValue.(MemberAccessable)
			if !ok {
				return nil, RuntimeErrorf(node, source, "cannot access %s of %s", keyPath[i], objectValue)
			}
			objectValue, err = object.Lookup(keyPath[i])
			if err != nil {
				return nil, err
			}
		}
		return objectValue, nil
	})
	return lazyResult, nil
}

func (interpreter *Interpreter) EvaluateTypeExpression(node *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	typeNode := node.ChildByFieldName("type")
	bodyNode := node.ChildByFieldName("body")
	if typeNode == nil || bodyNode == nil {
		return nil, SyntaxErrorf(node, source, "expected type and body")
	}
	lazyTypeValue, err := interpreter.EvaluateNode(typeNode, source)
	if err != nil {
		return nil, err
	}

	caseCount := int(bodyNode.NamedChildCount())
	typeCases := make(map[string]*LazyRuntimeValue, caseCount)
	for i := 0; i < caseCount; i++ {
		typeCaseNode := bodyNode.NamedChild(i)
		labelNode := typeCaseNode.ChildByFieldName("label")
		bodyNode := typeCaseNode.ChildByFieldName("body")
		if labelNode == nil || bodyNode == nil {
			return nil, SyntaxErrorf(typeCaseNode, source, "expected label and body")
		}
		if err != nil {
			return nil, err
		}
		lazyBody, err := interpreter.EvaluateNode(bodyNode, source)
		if err != nil {
			return nil, err
		}
		typeCases[labelNode.Content(source)] = lazyBody
	}
	lazyCheckedTypeExpression := NewLazyRuntimeValue(func() (RuntimeValue, error) {
		typeValue, err := lazyTypeValue.Evaluate()
		if err != nil {
			return nil, err
		}
		enumDecl, ok := typeValue.(EnumDeclRuntimeValue)
		if !ok {
			return nil, RuntimeErrorf(node, source, "expected enum type, got %s", typeValue)
		}
		typeExpression := TypeExpression{typeValue: enumDecl, cases: typeCases}
		if len(enumDecl.cases) != len(typeCases) {
			return nil, RuntimeErrorf(node, source, "expected %d cases, got %d", len(enumDecl.cases), len(typeCases))
		}
		for label := range typeCases {
			if _, ok := enumDecl.cases[label]; !ok {
				return nil, RuntimeErrorf(node, source, "undefined enum case %s", label)
			}
		}
		return typeExpression, nil
	})
	return lazyCheckedTypeExpression, nil
}

func (interpreter *Interpreter) EvaluateGroup(node *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	expressionNode := node.ChildByFieldName("expression")
	if expressionNode == nil {
		return nil, SyntaxErrorf(node, source, "expected expression")
	}
	return interpreter.EvaluateNode(expressionNode, source)
}

func (interpreter *Interpreter) EvaluateStringLiteral(node *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	string, err := strconv.Unquote(node.Content(source))
	if err != nil {
		return nil, err
	}
	return NewConstantRuntimeValue(PreludeString(string)), nil
}

func (interpreter *Interpreter) EvaluateBinaryExpression(node *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	return nil, SyntaxErrorf(node, source, "unimplemented")
}

func (interpreter *Interpreter) EvaluateUnaryExpression(node *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	return nil, SyntaxErrorf(node, source, "unimplemented")
}

func (interpreter *Interpreter) ParseFunctionLiteral(node *sitter.Node, source []byte, env *Environment) (Function, error) {
	parametersNode := node.ChildByFieldName("parameters")
	bodyNode := node.ChildByFieldName("body")

	closure := interpreter.ChildInterpreter("#func")
	// TODO: both nodes are optional!
	var (
		params []string
		err    error
	)
	if parametersNode != nil {
		params, err = closure.ParseParamterList(parametersNode, source)
		if err != nil {
			return Function{}, err
		}
	} else {
		params = []string{}
	}

	return Function{
		arguments: params,
		closure:   closure,
		body: func(i *Interpreter) ([]*LazyRuntimeValue, error) {
			var stmts []*LazyRuntimeValue
			if bodyNode != nil {
				stmts, err = i.ParseStatementList(bodyNode, source)
				if err != nil {
					return stmts, err
				}
			} else {
				return nil, RuntimeErrorf(node, source, "empty functions not implemented yet")
			}
			return stmts, nil
		},
	}, nil
}

func (interpreter *Interpreter) EvaluateFunctionLiteral(node *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	function, err := interpreter.ParseFunctionLiteral(node, source, interpreter.environment)
	if err != nil {
		return nil, err
	}
	return NewConstantRuntimeValue(function), nil
}

func (interpreter *Interpreter) EvaluateArrayLiteral(node *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
	return nil, SyntaxErrorf(node, source, "unimplemented")
}

func (interpreter *Interpreter) EvaluateNode(node *sitter.Node, source []byte) (*LazyRuntimeValue, error) {
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
	case parser.TYPE_NODE_ENUM_DECLARATION:
		return interpreter.EvaluateEnumDeclaration(node, source)
	case parser.TYPE_NODE_COMPLEX_INVOCATION_EXPRESSION:
		return interpreter.EvaluateComplexInvocationExpr(node, source)
	case parser.TYPE_NODE_SIMPLE_INVOCATION_EXPRESSION:
		return interpreter.EvaluateSimpleInvocation(node, source)
	case parser.TYPE_NODE_UNARY_EXPRESSION:
		return interpreter.EvaluateUnaryExpression(node, source)
	case parser.TYPE_NODE_BINARY_EXPRESSION:
		return interpreter.EvaluateBinaryExpression(node, source)
	case parser.TYPE_NODE_MEMBER_ACCESS:
		return interpreter.EvaluateMemberAccess(node, source)
	case parser.TYPE_NODE_TYPE_EXPRESSION:
		return interpreter.EvaluateTypeExpression(node, source)
	case parser.TYPE_NODE_STRING_LITERAL:
		return interpreter.EvaluateStringLiteral(node, source)
	// case parser.TYPE_NODE_ESCAPE_SEQUENCE:
	// 	return interpreter.Evaluate(node)
	case parser.TYPE_NODE_GROUP_LITERAL:
		return interpreter.EvaluateGroup(node, source)
	case parser.TYPE_NODE_NUMBER_LITERAL:
		return interpreter.EvaluateNumberLiteral(node, source)
	case parser.TYPE_NODE_ARRAY_LITERAL:
		return interpreter.EvaluateArrayLiteral(node, source)
	case parser.TYPE_NODE_FUNCTION_LITERAL:
		return interpreter.EvaluateFunctionLiteral(node, source)
	// case parser.TYPE_NODE_PARAMETER_LIST:
	// 	return interpreter.Evaluate(node)
	case parser.TYPE_NODE_IDENTIFIER:
		return interpreter.EvaluateIdentifier(node, source)
	case parser.TYPE_NODE_COMMENT:
		return nil, nil
	// case parser.TYPE_NODE_ENUM_CASE_REFERENCE:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_ERROR:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_UNEXPECTED:
	// 	return interpreter.Evaluate(node)
	default:
		return nil, SyntaxErrorf(node, source, "unimplemented node type %s", node.Type())
	}
}
