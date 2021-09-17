package interpreter

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/vknabel/go-lithia/parser"
)

type Interpreter struct {
	importRoots []string
	parser      *parser.Parser
	modules     map[ModuleName]*Module
	prelude     *Environment
}

func defaultImportRootPaths() []string {
	roots := []string{}
	if path, ok := os.LookupEnv("LITHIA_LOCALS"); ok {
		roots = append(roots, path)
	}
	if path, ok := os.LookupEnv("LITHIA_PACKAGES"); ok {
		roots = append(roots, path)
	}
	if path, ok := os.LookupEnv("LITHIA_STDLIB"); ok {
		roots = append(roots, path)
	} else {
		roots = append(roots, "/usr/local/opt/lithia/stdlib")
	}
	return roots
}

func NewInterpreter(importRoots ...string) *Interpreter {
	importRoots = append(importRoots, defaultImportRootPaths()...)
	absoluteImportRoots := make([]string, len(importRoots))
	for i, root := range importRoots {
		absolute, err := filepath.Abs(root)
		if err == nil {
			absoluteImportRoots[i] = absolute
		} else {
			absoluteImportRoots[i] = root
		}
	}
	return &Interpreter{
		importRoots: absoluteImportRoots,
		parser:      parser.NewParser(),
		modules:     make(map[ModuleName]*Module),
	}
}

type ModuleName string
type ModuleFile struct {
	name       string
	parentPath string
	module     ModuleName
}

type Module struct {
	name              ModuleName
	environment       *Environment
	executionContexts map[ModuleFile]*EvaluationContext
}

type EvaluationContext struct {
	interpreter   *Interpreter
	file          string
	module        *Module
	path          []string
	environment   *Environment
	functionCount int

	node   *sitter.Node
	source []byte
}

func (inter *Interpreter) NewEvaluationContext(file string, module *Module, node *sitter.Node, source []byte, environment *Environment) *EvaluationContext {
	if environment == nil {
		environment = NewEnvironment(inter.NewPreludeEnvironment())
	}
	return &EvaluationContext{
		interpreter:   inter,
		file:          file,
		module:        module,
		path:          []string{},
		environment:   environment,
		functionCount: 0,

		node:   node,
		source: source,
	}
}

func (i *EvaluationContext) NestedExecutionContext(name string) *EvaluationContext {
	return &EvaluationContext{
		interpreter:   i.interpreter,
		file:          i.file,
		module:        i.module,
		path:          append(i.path, name),
		environment:   NewEnvironment(i.environment),
		functionCount: 0,
		node:          i.node,
		source:        i.source,
	}
}

func (i *EvaluationContext) ChildNodeExecutionContext(childNode *sitter.Node) *EvaluationContext {
	return &EvaluationContext{
		interpreter:   i.interpreter,
		file:          i.file,
		module:        i.module,
		path:          i.path,
		environment:   i.environment,
		functionCount: i.functionCount,
		node:          childNode,
		source:        i.source,
	}
}

func (inter *Interpreter) NewModule(name ModuleName) *Module {
	module := &Module{
		name:              name,
		environment:       NewEnvironment(inter.NewPreludeEnvironment()),
		executionContexts: make(map[ModuleFile]*EvaluationContext),
	}
	inter.modules[name] = module
	return module
}
func (inter *Interpreter) Interpret(fileName string, script string) (RuntimeValue, error) {
	moduleName := ModuleName(strings.ReplaceAll(filepath.Base(fileName), ".", "_"))
	module := inter.NewModule(moduleName)
	ex, err := inter.LoadFileIntoModule(module, fileName, script)
	if err != nil {
		return nil, err
	}
	lazyValue, err := ex.EvaluateNode()
	if err != nil {
		return nil, err
	}
	return lazyValue.Evaluate()
}
func (inter *Interpreter) InterpretEmbed(fileName string, script string) (RuntimeValue, error) {
	moduleName := ModuleName(strings.ReplaceAll(filepath.Base(fileName), ".", "_"))
	module := inter.modules[moduleName]
	if module == nil {
		module = inter.NewModule(moduleName)
	}
	ex, err := inter.EmbedFileIntoModule(module, fileName, script)
	if err != nil {
		return nil, err
	}
	lazyValue, err := ex.EvaluateNode()
	if err != nil {
		return nil, err
	}
	return lazyValue.Evaluate()
}

func (inter *Interpreter) LoadFileIntoModule(module *Module, fileName string, script string) (*EvaluationContext, error) {
	tree, err := inter.parser.Parse(script)
	if err != nil {
		return nil, inter.SyntaxParsingError(fileName, script, tree)
	}
	ex := inter.NewEvaluationContext(fileName, module, tree.RootNode(), []byte(script), module.environment.Private())
	return ex, nil
}

func (inter *Interpreter) EmbedFileIntoModule(module *Module, fileName string, script string) (*EvaluationContext, error) {
	tree, err := inter.parser.Parse(script)
	if err != nil {
		return nil, inter.SyntaxParsingError(fileName, script, tree)
	}
	ex := inter.NewEvaluationContext(fileName, module, tree.RootNode(), []byte(script), module.environment)
	return ex, nil
}

func (ex *EvaluationContext) EvaluateSourceFile() (*LazyRuntimeValue, error) {
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

func priority(nodeType string) int {
	switch nodeType {
	case parser.TYPE_NODE_MODULE_DECLARATION:
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

func (ex *EvaluationContext) EvaluateModule() (*LazyRuntimeValue, error) {
	internalName := ex.node.ChildByFieldName("name").Content(ex.source)
	runtimeModule := NewConstantRuntimeValue(RuntimeModule{module: ex.module})
	ex.environment.DeclareUnexported(internalName, runtimeModule)
	return runtimeModule, nil
}

func (inter *Interpreter) LoadModuleIfNeeded(moduleName ModuleName) (*Module, error) {
	if module, ok := inter.modules[moduleName]; ok {
		return module, nil
	}
	for _, root := range inter.importRoots {
		relativeModulePath := strings.ReplaceAll(string(moduleName), ".", string(filepath.Separator))
		modulePath := filepath.Join(root, relativeModulePath)
		matches, err := filepath.Glob(filepath.Join(modulePath, "*.lithia"))
		if err != nil {
			continue
		}
		if len(matches) == 0 {
			continue
		}

		module := &Module{
			name:              moduleName,
			environment:       NewEnvironment(inter.NewPreludeEnvironment()),
			executionContexts: make(map[ModuleFile]*EvaluationContext),
		}
		inter.modules[moduleName] = module

		err = inter.LoadFilesIntoModule(module, matches)
		if err != nil {
			return module, err
		}

		return module, nil
	}
	return nil, fmt.Errorf("module %s not found", moduleName)
}

func (inter *Interpreter) LoadFilesIntoModule(module *Module, files []string) error {
	for _, file := range files {
		scriptData, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		childContext, err := inter.LoadFileIntoModule(module, file, string(scriptData))
		if err != nil {
			return err
		}
		source, err := childContext.EvaluateNode()
		if err != nil {
			return err
		}
		_, err = source.Evaluate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (ex *EvaluationContext) EvaluateImport() (*LazyRuntimeValue, error) {
	importModuleNode := ex.node.ChildByFieldName("name")
	modulePath := make([]string, importModuleNode.NamedChildCount())
	for i := 0; i < int(importModuleNode.NamedChildCount()); i++ {
		modulePath[i] = importModuleNode.NamedChild(i).Content(ex.source)
	}
	importMember := modulePath[len(modulePath)-1]
	moduleName := ModuleName(strings.Join(modulePath, "."))
	module, err := ex.interpreter.LoadModuleIfNeeded(moduleName)
	if err != nil {
		return nil, ex.SyntaxErrorf("error importing module %s: %s", moduleName, err.Error())
	}
	runtimeModule := NewConstantRuntimeValue(RuntimeModule{module: module})
	err = ex.environment.DeclareUnexported(importMember, runtimeModule)
	if err != nil {
		return nil, err
	}
	return runtimeModule, nil
}

func (ex *EvaluationContext) EvaluateLetDeclaration() (*LazyRuntimeValue, error) {
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
func (ex *EvaluationContext) EvaluateFunctionDeclaration() (*LazyRuntimeValue, error) {
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
func (ex *EvaluationContext) EvaluateDataDeclaration() (*LazyRuntimeValue, error) {
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

func (ex *EvaluationContext) EvaluateExternDeclaration() (*LazyRuntimeValue, error) {
	return ex.ChildNodeExecutionContext(ex.node.ChildByFieldName("name")).EvaluateIdentifier()
}

func (ex *EvaluationContext) ParseParamterList() ([]string, error) {
	params := make([]string, ex.node.ChildCount())
	for i := 0; i < int(ex.node.ChildCount()); i++ {
		child := ex.node.Child(i)
		params[i] = child.Content(ex.source)
	}
	return params, nil
}

func (ex *EvaluationContext) ParseStatementList() ([]*LazyRuntimeValue, error) {
	stmts := make([]*LazyRuntimeValue, 0, ex.node.NamedChildCount())
	for i := 0; i < int(ex.node.NamedChildCount()); i++ {
		child := ex.node.NamedChild(i)
		stmt, err := ex.ChildNodeExecutionContext(child).EvaluateNode()
		if err != nil {
			return nil, err
		}
		if stmt == nil {
			// TODO: comments are evil, but they shouldn't
			continue
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

func (ex *EvaluationContext) EvaluateNumberLiteral() (*LazyRuntimeValue, error) {
	literal := ex.node.Content(ex.source)
	integer, err := strconv.ParseInt(literal, 10, 64)
	if err != nil {
		return nil, err
	}
	return NewConstantRuntimeValue(PreludeInt(integer)), nil
}
func (ex *EvaluationContext) EvaluateComplexInvocationExpr() (*LazyRuntimeValue, error) {
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

func (ex *EvaluationContext) EvaluateSimpleInvocation() (*LazyRuntimeValue, error) {
	return ex.EvaluateComplexInvocationExpr()
}

func (ex *EvaluationContext) EvaluateEnumDeclaration() (*LazyRuntimeValue, error) {
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

func (ex *EvaluationContext) EvaluateIdentifier() (*LazyRuntimeValue, error) {
	content := ex.node.Content(ex.source)
	return NewLazyRuntimeValue(func() (RuntimeValue, error) {
		if lazyValue, ok := ex.environment.Get(content); ok {
			value, err := lazyValue.Evaluate()
			if err != nil {
				return nil, err
			}
			switch value := value.(type) {
			case DataDeclRuntimeValue:
				if len(value.fields) == 0 {
					return DataRuntimeValue{
						typeValue: &value,
						members:   make(map[string]*LazyRuntimeValue),
					}, nil
				} else {
					return value, nil
				}
			case Function:
				if len(value.arguments) == 0 {
					return value.Call(nil)
				} else {
					return value, nil
				}
			default:
				return value, nil
			}
		} else {
			return nil, ex.RuntimeErrorf("undefined identifier %s", content)
		}
	}), nil
}

func (ex *EvaluationContext) EvaluateMemberAccess() (*LazyRuntimeValue, error) {
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

func (ex *EvaluationContext) EvaluateTypeExpression() (*LazyRuntimeValue, error) {
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

func (ex *EvaluationContext) EvaluateGroup() (*LazyRuntimeValue, error) {
	expressionNode := ex.node.ChildByFieldName("expression")
	if expressionNode == nil {
		return nil, ex.SyntaxErrorf("expected expression")
	}
	return ex.ChildNodeExecutionContext(expressionNode).EvaluateNode()
}

func (ex *EvaluationContext) EvaluateStringLiteral() (*LazyRuntimeValue, error) {
	string, err := strconv.Unquote(ex.node.Content(ex.source))
	if err != nil {
		return nil, err
	}
	return NewConstantRuntimeValue(PreludeString(string)), nil
}

func (ex *EvaluationContext) EvaluateBinaryExpression() (*LazyRuntimeValue, error) {
	if ex.node.NamedChildCount() != 2 {
		return nil, ex.SyntaxErrorf("expected 2 children, got %d", ex.node.NamedChildCount())
	}
	lazyLeft, err := ex.ChildNodeExecutionContext(ex.node.NamedChild(0)).EvaluateNode()
	if err != nil {
		return nil, err
	}
	lazyRight, err := ex.ChildNodeExecutionContext(ex.node.NamedChild(1)).EvaluateNode()
	if err != nil {
		return nil, err
	}
	operator := ex.node.ChildByFieldName("operator").Content(ex.source)

	impl, err := ex.BinaryOperatorFunction(operator)
	if err != nil {
		return nil, err
	}
	return NewLazyRuntimeValue(func() (RuntimeValue, error) {
		return impl(lazyLeft, lazyRight)
	}), nil
}

func (ex *EvaluationContext) EvaluateUnaryExpression() (*LazyRuntimeValue, error) {
	return nil, ex.SyntaxErrorf("unimplemented")
}

func (ex *EvaluationContext) ParseFunctionLiteral(name string) (Function, error) {
	parametersNode := ex.node.ChildByFieldName("parameters")
	bodyNode := ex.node.ChildByFieldName("body")

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
		body: func(i *EvaluationContext) ([]*LazyRuntimeValue, error) {
			if bodyNode == nil {
				return []*LazyRuntimeValue{
					NewConstantRuntimeValue(PreludeString("TODO: WHY IS THIS NIL?" + name)),
				}, nil
			}
			stmts, err := i.ChildNodeExecutionContext(bodyNode).ParseStatementList()
			if err != nil {
				return nil, err
			}
			return stmts, nil
		},
	}, nil
}

func (ex *EvaluationContext) EvaluateFunctionLiteral() (*LazyRuntimeValue, error) {
	function, err := ex.ParseFunctionLiteral("")
	if err != nil {
		return nil, err
	}
	return NewConstantRuntimeValue(function), nil
}

func (ex *EvaluationContext) EvaluateArrayLiteral() (*LazyRuntimeValue, error) {
	numberOfElements := int(ex.node.NamedChildCount())
	elements := make([]*LazyRuntimeValue, numberOfElements)
	for i := 0; i < numberOfElements; i++ {
		elementNode := ex.node.NamedChild(i)
		lazyElement, err := ex.ChildNodeExecutionContext(elementNode).EvaluateNode()
		if err != nil {
			return nil, err
		}
		elements[i] = lazyElement
	}
	return NewLazyRuntimeValue(func() (RuntimeValue, error) {
		var (
			cons       DataDeclRuntimeValue
			runtimeNil DataDeclRuntimeValue
		)
		if lazyConsValue, ok := ex.environment.Get("Cons"); ok {
			consValue, err := lazyConsValue.Evaluate()
			if err != nil {
				return nil, err
			}
			cons = consValue.(DataDeclRuntimeValue)
		}
		if lazyNilValue, ok := ex.environment.Get("Nil"); ok {
			nilValue, err := lazyNilValue.Evaluate()
			if err != nil {
				return nil, err
			}
			runtimeNil = nilValue.(DataDeclRuntimeValue)
		}
		return SliceToList(cons, runtimeNil, elements), nil
	}), nil
}

func SliceToList(consDecl DataDeclRuntimeValue, nilDecl DataDeclRuntimeValue, slice []*LazyRuntimeValue) RuntimeValue {
	if len(slice) == 0 {
		return DataRuntimeValue{
			typeValue: &nilDecl,
			members:   make(map[string]*LazyRuntimeValue),
		}
	} else {
		return DataRuntimeValue{
			typeValue: &consDecl,
			members: map[string]*LazyRuntimeValue{
				"head": slice[0],
				"tail": NewLazyRuntimeValue(func() (RuntimeValue, error) {
					return SliceToList(consDecl, nilDecl, slice[1:]), nil
				}),
			},
		}
	}
}

func (ex *EvaluationContext) EvaluateNode() (*LazyRuntimeValue, error) {
	switch ex.node.Type() {
	case parser.TYPE_NODE_SOURCE_FILE:
		return ex.EvaluateSourceFile()
	case parser.TYPE_NODE_MODULE_DECLARATION:
		return ex.EvaluateModule()
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
