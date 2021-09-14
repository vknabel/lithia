package interpreter

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/vknabel/go-lithia/parser"
)

type Interpreter struct {
	importRoot              string
	parser                  *parser.Parser
	cachedExecutionContexts map[string]*ExecutionContext
}

func NewInterpreter(importRoot string) *Interpreter {
	return &Interpreter{
		importRoot:              importRoot,
		parser:                  parser.NewParser(),
		cachedExecutionContexts: make(map[string]*ExecutionContext),
	}
}

type ExecutionContext struct {
	interpreter   *Interpreter
	fileName      string
	path          []string
	environment   *Environment
	functionCount int

	node   *sitter.Node
	source []byte
}

func (inter *Interpreter) NewExecutionContext(fileName string, node *sitter.Node, source []byte) *ExecutionContext {
	return &ExecutionContext{
		interpreter:   inter,
		fileName:      fileName,
		path:          []string{},
		environment:   NewEnvironment(NewPreludeEnvironment()),
		functionCount: 0,

		node:   node,
		source: source,
	}
}

func (i *ExecutionContext) NestedExecutionContext(name string) *ExecutionContext {
	return &ExecutionContext{
		fileName:      i.fileName,
		path:          append(i.path, name),
		environment:   NewEnvironment(i.environment),
		functionCount: 0,
		node:          i.node,
		source:        i.source,
	}
}

func (i *ExecutionContext) ChildNodeExecutionContext(childNode *sitter.Node) *ExecutionContext {
	return &ExecutionContext{
		fileName:      i.fileName,
		path:          i.path,
		environment:   i.environment,
		functionCount: i.functionCount,
		node:          childNode,
		source:        i.source,
	}
}

func (inter *Interpreter) Interpret(fileName string, script string) (RuntimeValue, error) {
	ex, err := inter.LoadContext(fileName, script)
	if err != nil {
		return nil, err
	}
	lazyValue, err := ex.EvaluateNode()
	if err != nil {
		return nil, err
	}
	return lazyValue.Evaluate()
}

func (inter *Interpreter) LoadContext(fileName string, script string) (*ExecutionContext, error) {
	tree, error := inter.parser.Parse(script)
	if error != nil {
		return nil, error
	}
	ex := inter.NewExecutionContext(fileName, tree.RootNode(), []byte(script))
	if context, ok := inter.cachedExecutionContexts[fileName]; ok {
		ex.functionCount = context.functionCount
		ex.environment = context.environment
	}
	inter.cachedExecutionContexts[fileName] = ex
	return ex, nil
}

func (inter *Interpreter) LoadContextIfNeeded(fileName string, script string) (*ExecutionContext, error) {
	if context, ok := inter.cachedExecutionContexts[fileName]; ok {
		return context, nil
	}
	tree, error := inter.parser.Parse(script)
	if error != nil {
		return nil, error
	}
	ex := inter.NewExecutionContext(fileName, tree.RootNode(), []byte(script))
	inter.cachedExecutionContexts[fileName] = ex
	return ex, nil
}

func (ex *ExecutionContext) EvaluateSourceFile() (*LazyRuntimeValue, error) {
	count := ex.node.ChildCount()
	children := make([]*sitter.Node, count)
	for i := uint32(0); i < count; i++ {
		child := ex.node.Child(int(i))
		children[i] = child
	}
	sort.SliceStable(children, func(i, j int) bool {
		lp := priority(children[i].Type())
		rp := priority(children[j].Type())
		return lp > rp
	})

	return NewLazyRuntimeValue(func() (RuntimeValue, error) {
		var lastValue RuntimeValue
		for _, child := range children {
			lazyValue, err := ex.ChildNodeExecutionContext(child).EvaluateNode()
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
		return lastValue, nil
	}), nil
}

func (ex *ExecutionContext) EvaluateImport() (*LazyRuntimeValue, error) {
	importModuleNode := ex.node.ChildByFieldName("name")
	modulePath := make([]string, importModuleNode.NamedChildCount())
	for i := 0; i < int(importModuleNode.NamedChildCount()); i++ {
		modulePath[i] = importModuleNode.NamedChild(i).Content(ex.source)
	}

	childModulePath := path.Join(path.Dir(ex.fileName), path.Join(modulePath...))

	if _, err := os.Stat(childModulePath); os.IsNotExist(err) {
		// rootModulePath := path.Join(ex.interpreter.importRoot, childModulePath)
		return nil, ex.SyntaxErrorf("root module import not implemented")
	}
	matches, err := filepath.Glob(filepath.Join(childModulePath, "*.lithia"))
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, ex.SyntaxErrorf("root module import not implemented")
	}
	for _, match := range matches {
		if strings.HasSuffix(match, ".lithia") {
			scriptData, err := os.ReadFile(match)
			if err != nil {
				return nil, err
			}
			context, err := ex.interpreter.LoadContextIfNeeded(match, string(scriptData))
			if err != nil {
				return nil, err
			}
			return context.EvaluateNode()
		}
	}
	return nil, ex.SyntaxErrorf("child module import not implemented")
}

func (ex *ExecutionContext) EvaluateLetDeclaration() (*LazyRuntimeValue, error) {
	nameNode := ex.node.ChildByFieldName("name")
	valueNode := ex.node.ChildByFieldName("value")
	if nameNode == nil || valueNode == nil {
		return nil, ex.SyntaxErrorf("let declaration must have name and value")
	}
	lazyValue, err := ex.ChildNodeExecutionContext(valueNode).EvaluateNode()
	if err != nil {
		return nil, err
	}
	err = ex.environment.Declare(nameNode.Content([]byte(ex.source)), lazyValue)
	if err != nil {
		return nil, err
	}
	return lazyValue, nil
}
func (ex *ExecutionContext) EvaluateFunctionDeclaration() (*LazyRuntimeValue, error) {
	name := ex.node.ChildByFieldName("name").Content(ex.source)
	functionNode := ex.node.ChildByFieldName("function")
	function, err := ex.ChildNodeExecutionContext(functionNode).ParseFunctionLiteral(name)

	if err != nil {
		return nil, err
	}
	impl := NewConstantRuntimeValue(function)
	err = ex.environment.Declare(name, impl)
	if err != nil {
		return nil, err
	}
	return impl, nil
}
func (ex *ExecutionContext) EvaluateDataDeclaration() (*LazyRuntimeValue, error) {
	name := ex.node.ChildByFieldName("name").Content(ex.source)
	propertiesNode := ex.node.ChildByFieldName("properties")

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
			name := child.ChildByFieldName("name").Content(ex.source)
			data.fields[i] = DataDeclField{name: name}
		case parser.TYPE_NODE_DATA_PROPERTY_FUNCTION:
			name := child.ChildByFieldName("name").Content(ex.source)
			parameters, error := ex.ChildNodeExecutionContext(child.ChildByFieldName("parameters")).ParseParamterList()
			if error != nil {
				return nil, error
			}
			data.fields[i] = DataDeclField{name: name, params: parameters}
		default:
			return nil, ex.ChildNodeExecutionContext(child).SyntaxErrorf("unexpected node type %s", child.Type())
		}
	}

	constantValue := NewConstantRuntimeValue(data)
	ex.environment.Declare(name, constantValue)
	return constantValue, nil
}

func (ex *ExecutionContext) EvaluateExternDeclaration() (*LazyRuntimeValue, error) {
	return ex.ChildNodeExecutionContext(ex.node.ChildByFieldName("name")).EvaluateIdentifier()
}

func (ex *ExecutionContext) ParseParamterList() ([]string, error) {
	params := make([]string, ex.node.ChildCount())
	for i := 0; i < int(ex.node.ChildCount()); i++ {
		child := ex.node.Child(i)
		params[i] = child.Content(ex.source)
	}
	return params, nil
}

func (ex *ExecutionContext) ParseStatementList() ([]*LazyRuntimeValue, error) {
	stmts := make([]*LazyRuntimeValue, ex.node.NamedChildCount())
	for i := 0; i < int(ex.node.NamedChildCount()); i++ {
		child := ex.node.NamedChild(i)
		stmt, err := ex.ChildNodeExecutionContext(child).EvaluateNode()
		if err != nil {
			return nil, err
		}
		stmts[i] = stmt
	}
	return stmts, nil
}

func (ex *ExecutionContext) EvaluateNumberLiteral() (*LazyRuntimeValue, error) {
	literal := ex.node.Content(ex.source)
	integer, err := strconv.ParseInt(literal, 10, 64)
	if err != nil {
		return nil, err
	}
	return NewConstantRuntimeValue(PreludeInt(integer)), nil
}
func (ex *ExecutionContext) EvaluateComplexInvocationExpr() (*LazyRuntimeValue, error) {
	functionNode := ex.node.ChildByFieldName("function")
	lazyFunc, err := ex.ChildNodeExecutionContext(functionNode).EvaluateNode()
	if err != nil {
		return nil, err
	}
	return NewLazyRuntimeValue(func() (RuntimeValue, error) {
		functionValue, err := lazyFunc.Evaluate()
		if err != nil {
			return nil, err
		}
		function, ok := reflect.ValueOf(functionValue).Interface().(Callable)
		if !ok {
			return nil, ex.RuntimeErrorf("expected callable, got %T", functionValue)
		}

		args := make([]*LazyRuntimeValue, ex.node.NamedChildCount()-1)
		for i := 0; i < int(ex.node.NamedChildCount()-1); i++ {
			child := ex.node.NamedChild(i + 1)
			lazyValue, err := ex.ChildNodeExecutionContext(child).EvaluateNode()
			if err != nil {
				return nil, err
			}
			args[i] = lazyValue
		}

		return function.Call(args)
	}), nil
}

func (ex *ExecutionContext) EvaluateSimpleInvocation() (*LazyRuntimeValue, error) {
	return ex.EvaluateComplexInvocationExpr()
}

func (ex *ExecutionContext) EvaluateEnumDeclaration() (*LazyRuntimeValue, error) {
	name := ex.node.ChildByFieldName("name").Content(ex.source)
	casesNode := ex.node.ChildByFieldName("cases")
	if casesNode == nil {
		enum := NewConstantRuntimeValue(EnumDeclRuntimeValue{name: name, cases: make(map[string]*LazyRuntimeValue)})
		ex.environment.Declare(name, enum)
		return enum, nil
	}
	caseCount := int(casesNode.ChildCount())
	cases := make(map[string]*LazyRuntimeValue)
	for i := 0; i < caseCount; i++ {
		child := casesNode.Child(i)
		switch child.Type() {
		case parser.TYPE_NODE_ENUM_CASE_REFERENCE:
			caseName := child.Content(ex.source)
			lookedUp, _ := ex.environment.Get(caseName)
			if lookedUp == nil {
				return nil, ex.ChildNodeExecutionContext(child).RuntimeErrorf("undefined enum case %s", caseName)
			}
			cases[caseName] = lookedUp
		case parser.TYPE_NODE_DATA_DECLARATION:
			caseName := child.ChildByFieldName("name").Content(ex.source)
			runtimeValue, err := ex.ChildNodeExecutionContext(child).EvaluateDataDeclaration()
			if err != nil {
				return nil, err
			}
			cases[caseName] = runtimeValue
		case parser.TYPE_NODE_ENUM_DECLARATION:
			caseName := child.ChildByFieldName("name").Content(ex.source)
			runtimeValue, err := ex.ChildNodeExecutionContext(child).EvaluateEnumDeclaration()
			if err != nil {
				return nil, err
			}
			cases[caseName] = runtimeValue
		default:
			return nil, ex.ChildNodeExecutionContext(child).SyntaxErrorf("unexpected node type %s", child.Type())
		}
	}
	constantValue := NewConstantRuntimeValue(EnumDeclRuntimeValue{
		name:  name,
		cases: cases,
	})
	ex.environment.Declare(name, constantValue)
	return constantValue, nil
}

func (ex *ExecutionContext) EvaluateIdentifier() (*LazyRuntimeValue, error) {
	string := ex.node.Content(ex.source)
	return NewLazyRuntimeValue(func() (RuntimeValue, error) {
		if value, ok := ex.environment.Get(string); ok {
			return value.Evaluate()
		} else {
			return nil, ex.RuntimeErrorf("undefined identifier %s", string)
		}
	}), nil
}

func (ex *ExecutionContext) EvaluateMemberAccess() (*LazyRuntimeValue, error) {
	if ex.node.NamedChildCount() < 2 {
		return nil, ex.SyntaxErrorf("expected at least 2 children, got %d", ex.node.NamedChildCount())
	}
	literalNode := ex.node.NamedChild(0)
	lazyObject, err := ex.ChildNodeExecutionContext(literalNode).EvaluateNode()
	if err != nil {
		return nil, err
	}

	keyPath := make([]string, ex.node.NamedChildCount()-1)
	for i := 1; i < int(ex.node.NamedChildCount()); i++ {
		child := ex.node.NamedChild(i)
		keyPath[i-1] = child.Content(ex.source)
	}
	lazyResult := NewLazyRuntimeValue(func() (RuntimeValue, error) {
		object, err := lazyObject.Evaluate()
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(keyPath); i++ {
			lazyObject, err = object.Lookup(keyPath[i])
			if err != nil {
				return nil, err
			}
			object, err = lazyObject.Evaluate()
			if err != nil {
				return nil, err
			}
		}
		return object, nil
	})
	return lazyResult, nil
}

func (ex *ExecutionContext) EvaluateTypeExpression() (*LazyRuntimeValue, error) {
	typeNode := ex.node.ChildByFieldName("type")
	bodyNode := ex.node.ChildByFieldName("body")
	if typeNode == nil || bodyNode == nil {
		return nil, ex.SyntaxErrorf("expected type and body")
	}
	lazyTypeValue, err := ex.ChildNodeExecutionContext(typeNode).EvaluateNode()
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
			return nil, ex.SyntaxErrorf("expected label and body")
		}
		if err != nil {
			return nil, err
		}
		lazyBody, err := ex.ChildNodeExecutionContext(bodyNode).EvaluateNode()
		if err != nil {
			return nil, err
		}
		typeCases[labelNode.Content(ex.source)] = lazyBody
	}
	lazyCheckedTypeExpression := NewLazyRuntimeValue(func() (RuntimeValue, error) {
		typeValue, err := lazyTypeValue.Evaluate()
		if err != nil {
			return nil, err
		}
		enumDecl, ok := typeValue.(EnumDeclRuntimeValue)
		if !ok {
			return nil, ex.RuntimeErrorf("expected enum type, got %s", typeValue)
		}
		typeExpression := TypeExpression{typeValue: enumDecl, cases: typeCases}
		if len(enumDecl.cases) != len(typeCases) {
			return nil, ex.RuntimeErrorf("expected %s cases, got %s", casesToString(enumDecl.cases), casesToString(typeCases))
		}
		for label := range typeCases {
			if _, ok := enumDecl.cases[label]; !ok {
				return nil, ex.RuntimeErrorf("undefined enum case %s, expected %s", label, casesToString(enumDecl.cases))
			}
		}
		return typeExpression, nil
	})
	return lazyCheckedTypeExpression, nil
}
func casesToString(cases map[string]*LazyRuntimeValue) string {
	var labels []string
	for label := range cases {
		labels = append(labels, label)
	}
	return fmt.Sprintf("[%s]", strings.Join(labels, ", "))
}

func (ex *ExecutionContext) EvaluateGroup() (*LazyRuntimeValue, error) {
	expressionNode := ex.node.ChildByFieldName("expression")
	if expressionNode == nil {
		return nil, ex.SyntaxErrorf("expected expression")
	}
	return ex.ChildNodeExecutionContext(expressionNode).EvaluateNode()
}

func (ex *ExecutionContext) EvaluateStringLiteral() (*LazyRuntimeValue, error) {
	string, err := strconv.Unquote(ex.node.Content(ex.source))
	if err != nil {
		return nil, err
	}
	return NewConstantRuntimeValue(PreludeString(string)), nil
}

func (ex *ExecutionContext) EvaluateBinaryExpression() (*LazyRuntimeValue, error) {
	return nil, ex.SyntaxErrorf("unimplemented")
}

func (ex *ExecutionContext) EvaluateUnaryExpression() (*LazyRuntimeValue, error) {
	return nil, ex.SyntaxErrorf("unimplemented")
}

func (ex *ExecutionContext) ParseFunctionLiteral(name string) (Function, error) {
	parametersNode := ex.node.ChildByFieldName("parameters")
	bodyNode := ex.node.ChildByFieldName("body")

	// TODO: both nodes are optional!
	var (
		params []string
		err    error
	)
	if parametersNode != nil {
		params, err = ex.ChildNodeExecutionContext(parametersNode).ParseParamterList()
		if err != nil {
			return Function{}, err
		}
	} else {
		params = []string{}
	}

	if name == "" {
		name = fmt.Sprintf("func#%d", ex.functionCount)
		ex.functionCount += 1
	}
	return Function{
		name:      name,
		arguments: params,
		parent:    ex,
		body: func(i *ExecutionContext) ([]*LazyRuntimeValue, error) {
			var stmts []*LazyRuntimeValue
			if bodyNode != nil {
				stmts, err = i.ChildNodeExecutionContext(bodyNode).ParseStatementList()
				if err != nil {
					return stmts, err
				}
			} else {
				return nil, ex.RuntimeErrorf("empty functions not implemented yet")
			}
			return stmts, nil
		},
	}, nil
}

func (ex *ExecutionContext) EvaluateFunctionLiteral() (*LazyRuntimeValue, error) {
	function, err := ex.ParseFunctionLiteral("")
	if err != nil {
		return nil, err
	}
	return NewConstantRuntimeValue(function), nil
}

func (ex *ExecutionContext) EvaluateArrayLiteral() (*LazyRuntimeValue, error) {
	return nil, ex.SyntaxErrorf("unimplemented")
}

func (ex *ExecutionContext) EvaluateNode() (*LazyRuntimeValue, error) {
	switch ex.node.Type() {
	case parser.TYPE_NODE_SOURCE_FILE:
		return ex.EvaluateSourceFile()
	// case parser.TYPE_NODE_PACKAGE_DECLARATION:
	// 	return interpreter.Evaluate(node)
	case parser.TYPE_NODE_IMPORT_DECLARATION:
		return ex.EvaluateImport()
	case parser.TYPE_NODE_LET_DECLARATION:
		return ex.EvaluateLetDeclaration()
	case parser.TYPE_NODE_FUNCTION_DECLARATION:
		return ex.EvaluateFunctionDeclaration()
	case parser.TYPE_NODE_DATA_DECLARATION:
		return ex.EvaluateDataDeclaration()
	case parser.TYPE_NODE_EXTERN_DECLARATION:
		return ex.EvaluateExternDeclaration()
	case parser.TYPE_NODE_ENUM_DECLARATION:
		return ex.EvaluateEnumDeclaration()
	case parser.TYPE_NODE_COMPLEX_INVOCATION_EXPRESSION:
		return ex.EvaluateComplexInvocationExpr()
	case parser.TYPE_NODE_SIMPLE_INVOCATION_EXPRESSION:
		return ex.EvaluateSimpleInvocation()
	case parser.TYPE_NODE_UNARY_EXPRESSION:
		return ex.EvaluateUnaryExpression()
	case parser.TYPE_NODE_BINARY_EXPRESSION:
		return ex.EvaluateBinaryExpression()
	case parser.TYPE_NODE_MEMBER_ACCESS:
		return ex.EvaluateMemberAccess()
	case parser.TYPE_NODE_TYPE_EXPRESSION:
		return ex.EvaluateTypeExpression()
	case parser.TYPE_NODE_STRING_LITERAL:
		return ex.EvaluateStringLiteral()
	// case parser.TYPE_NODE_ESCAPE_SEQUENCE:
	// 	return interpreter.Evaluate(node)
	case parser.TYPE_NODE_GROUP_LITERAL:
		return ex.EvaluateGroup()
	case parser.TYPE_NODE_NUMBER_LITERAL:
		return ex.EvaluateNumberLiteral()
	case parser.TYPE_NODE_ARRAY_LITERAL:
		return ex.EvaluateArrayLiteral()
	case parser.TYPE_NODE_FUNCTION_LITERAL:
		return ex.EvaluateFunctionLiteral()
	// case parser.TYPE_NODE_PARAMETER_LIST:
	// 	return interpreter.Evaluate(node)
	case parser.TYPE_NODE_IDENTIFIER:
		return ex.EvaluateIdentifier()
	case parser.TYPE_NODE_COMMENT:
		return nil, nil
	// case parser.TYPE_NODE_ENUM_CASE_REFERENCE:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_ERROR:
	// 	return interpreter.Evaluate(node)
	// case parser.TYPE_NODE_UNEXPECTED:
	// 	return interpreter.Evaluate(node)
	default:
		return nil, ex.SyntaxErrorf("unimplemented node type %s", ex.node.Type())
	}
}
