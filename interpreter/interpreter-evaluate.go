package interpreter

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/vknabel/go-lithia/parser"
)

func (ex *EvaluationContext) EvaluateImport() (*LazyRuntimeValue, LocatableError) {
	importModuleNode := ex.node.ChildByFieldName("name")
	membersNode := ex.node.ChildByFieldName("members")
	modulePath := make([]string, importModuleNode.NamedChildCount())
	for i := 0; i < int(importModuleNode.NamedChildCount()); i++ {
		modulePath[i] = importModuleNode.NamedChild(i).Content(ex.source)
	}
	importMember := modulePath[len(modulePath)-1]
	moduleName := ModuleName(strings.Join(modulePath, "."))
	module, err := ex.interpreter.LoadModuleIfNeeded(moduleName)
	if err != nil {
		return nil, ex.LocatableErrorOrConvert(err)
	}
	runtimeModule := NewConstantRuntimeValue(RuntimeModule{module: module})
	err = ex.environment.DeclareUnexported(importMember, runtimeModule)
	if err != nil {
		return nil, ex.LocatableErrorOrConvert(err)
	}

	if membersNode != nil {
		for i := 0; i < int(membersNode.NamedChildCount()); i++ {
			childNode := membersNode.NamedChild(i)
			if childNode.Type() == parser.TYPE_NODE_IDENTIFIER {
				memberName := membersNode.NamedChild(i).Content(ex.source)
				member, ok := module.environment.Scope[memberName]
				if !ok {
					return nil, ex.ChildNodeExecutionContext(childNode).RuntimeErrorf("%s is not a member of %s", memberName, moduleName)
				}
				err := ex.environment.DeclareUnexported(memberName, member)
				if err != nil {
					return nil, ex.LocatableErrorOrConvert(err)
				}
			}
		}
	}

	return runtimeModule, nil
}

func (ex *EvaluationContext) EvaluateLetDeclaration() (*LazyRuntimeValue, LocatableError) {
	nameNode := ex.node.ChildByFieldName("name")
	valueNode := ex.node.ChildByFieldName("value")
	if nameNode == nil || valueNode == nil {
		return nil, ex.SyntaxErrorf("let declaration must have name and value")
	}
	var err error
	lazyValue, err := ex.ChildNodeExecutionContext(valueNode).EvaluateNode()
	if err != nil {
		return nil, ex.LocatableErrorOrConvert(err)
	}
	err = ex.environment.Declare(nameNode.Content([]byte(ex.source)), lazyValue)
	if err != nil {
		return nil, ex.LocatableErrorOrConvert(err)
	}
	return lazyValue, nil
}
func (ex *EvaluationContext) EvaluateFunctionDeclaration(docs DocString) (*LazyRuntimeValue, LocatableError) {
	var err error
	name := ex.node.ChildByFieldName("name").Content(ex.source)
	functionNode := ex.node.ChildByFieldName("function")
	function, err := ex.ChildNodeExecutionContext(functionNode).ParseFunctionLiteral(name, docs)

	if err != nil {
		return nil, ex.LocatableErrorOrConvert(err)
	}
	impl := NewConstantRuntimeValue(function)
	err = ex.environment.Declare(name, impl)
	if err != nil {
		return nil, ex.LocatableErrorOrConvert(err)
	}
	return impl, nil
}
func (ex *EvaluationContext) EvaluateDataDeclaration(docs DocString) (*LazyRuntimeValue, LocatableError) {
	name := ex.node.ChildByFieldName("name").Content(ex.source)
	propertiesNode := ex.node.ChildByFieldName("properties")

	var numberOfFields int
	if propertiesNode != nil {
		numberOfFields = int(propertiesNode.ChildCount())
	} else {
		numberOfFields = 0
	}

	data := NewDataDecl(name, make([]DataDeclField, 0, numberOfFields), docs)

	comments := make([]string, 0)
	for i := 0; i < int(ex.node.ChildCount()); i++ {
		child := ex.node.Child(i)
		if child.Type() == parser.TYPE_NODE_COMMENT {
			comments = append(comments, child.Content(ex.source))
		}
	}
	for i := 0; i < numberOfFields; i++ {
		child := propertiesNode.Child(i)
		switch child.Type() {
		case parser.TYPE_NODE_DATA_PROPERTY_VALUE:
			name := child.ChildByFieldName("name").Content(ex.source)
			fieldDecl := DataDeclField{
				name: name,
				docs: MakeDocString(comments),
			}
			comments = make([]string, 0)
			data.fields = append(data.fields, fieldDecl)
		case parser.TYPE_NODE_DATA_PROPERTY_FUNCTION:
			name := child.ChildByFieldName("name").Content(ex.source)
			var err error
			parameters, err := ex.ChildNodeExecutionContext(child.ChildByFieldName("parameters")).ParseParamterList()
			if err != nil {
				return nil, ex.LocatableErrorOrConvert(err)
			}
			fieldDecl := DataDeclField{
				name:   name,
				params: parameters,
				docs:   MakeDocString(comments),
			}
			comments = make([]string, 0)
			data.fields = append(data.fields, fieldDecl)
		case parser.TYPE_NODE_COMMENT:
			comments = append(comments, child.Content(ex.source))
		default:
			return nil, ex.ChildNodeExecutionContext(child).SyntaxErrorf("unexpected node type %s", child.Type())
		}
	}

	constantValue := NewConstantRuntimeValue(data)
	ex.environment.Declare(name, constantValue)
	return constantValue, nil
}

func (ex *EvaluationContext) EvaluateExternDeclaration(docs DocString) (*LazyRuntimeValue, LocatableError) {
	name := ex.node.ChildByFieldName("name").Content(ex.source)
	externalDef, ok := ex.interpreter.ExternalDefinitions[ex.module.name]
	if !ok {
		return nil, ex.SyntaxErrorf("no external declarations allowed in module %s", ex.module.name)
	}
	runtimeValue, ok := externalDef.Lookup(name, ex.environment, Docs{
		name: name,
		docs: docs,
	})
	if !ok {
		return nil, ex.SyntaxErrorf("unknown external declaration %s in module %s", name, ex.module.name)
	}
	constantValue := NewConstantRuntimeValue(runtimeValue)
	ex.environment.Declare(name, constantValue)
	return constantValue, nil
}

func (ex *EvaluationContext) ParseParamterList() ([]string, error) {
	params := make([]string, ex.node.ChildCount())
	for i := 0; i < int(ex.node.ChildCount()); i++ {
		child := ex.node.Child(i)
		params[i] = child.Content(ex.source)
	}
	return params, nil
}

func (ex *EvaluationContext) ParseStatementList() ([]*LazyRuntimeValue, LocatableError) {
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

func (ex *EvaluationContext) EvaluateNumberLiteral() (*LazyRuntimeValue, LocatableError) {
	literal := ex.node.Content(ex.source)
	integer, err := strconv.ParseInt(literal, 10, 64)
	if err != nil {
		return nil, ex.LocatableErrorOrConvert(err)
	}
	return NewConstantRuntimeValue(PreludeInt(integer)), nil
}
func (ex *EvaluationContext) EvaluateComplexInvocationExpr() (*LazyRuntimeValue, LocatableError) {
	functionNode := ex.node.ChildByFieldName("function")
	lazyFunc, err := ex.ChildNodeExecutionContext(functionNode).EvaluateNode()
	if err != nil {
		return nil, err
	}
	return NewLazyRuntimeValue(func() (RuntimeValue, LocatableError) {
		var err error
		functionValue, err := lazyFunc.Evaluate()
		if err != nil {
			return nil, ex.LocatableErrorOrConvert(err)
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

		result, err := function.Call(args)
		return result, ex.LocatableErrorOrConvert(err)
	}), nil
}

func (ex *EvaluationContext) EvaluateSimpleInvocation() (*LazyRuntimeValue, LocatableError) {
	return ex.EvaluateComplexInvocationExpr()
}

func (ex *EvaluationContext) EvaluateEnumDeclaration(docs DocString) (*LazyRuntimeValue, LocatableError) {
	name := ex.node.ChildByFieldName("name").Content(ex.source)
	casesNode := ex.node.ChildByFieldName("cases")
	if casesNode == nil {
		enum := NewConstantRuntimeValue(EnumDeclRuntimeValue{
			name:  name,
			cases: make(map[string]*LazyRuntimeValue),
			docs:  docs,
		})
		ex.environment.Declare(name, enum)
		return enum, nil
	}
	comments := make([]string, 0)
	for i := 0; i < int(ex.node.ChildCount()); i++ {
		child := ex.node.Child(i)
		if child.Type() == parser.TYPE_NODE_COMMENT {
			comments = append(comments, child.Content(ex.source))
		}
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
			if caseName == "" {
				fmt.Println("caseName:", caseName, len(comments))
			}
			runtimeValue, err := ex.ChildNodeExecutionContext(child).EvaluateDataDeclaration(MakeDocString(comments))
			if err != nil {
				return nil, err
			}
			cases[caseName] = runtimeValue
			comments = make([]string, 0)
		case parser.TYPE_NODE_ENUM_DECLARATION:
			caseName := child.ChildByFieldName("name").Content(ex.source)
			runtimeValue, err := ex.ChildNodeExecutionContext(child).EvaluateEnumDeclaration(MakeDocString(comments))
			if err != nil {
				return nil, err
			}
			cases[caseName] = runtimeValue
			comments = make([]string, 0)
		case parser.TYPE_NODE_COMMENT:
			comments = append(comments, child.Content(ex.source))
		default:
			return nil, ex.ChildNodeExecutionContext(child).SyntaxErrorf("unexpected node type %s", child.Type())
		}
	}
	constantValue := NewConstantRuntimeValue(EnumDeclRuntimeValue{
		name:  name,
		cases: cases,
		docs:  docs,
	})
	ex.environment.Declare(name, constantValue)
	return constantValue, nil
}

func (ex *EvaluationContext) EvaluateIdentifier() (*LazyRuntimeValue, LocatableError) {
	content := ex.node.Content(ex.source)
	return NewLazyRuntimeValue(func() (RuntimeValue, LocatableError) {
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
					result, err := value.Call(nil)
					return result, ex.LocatableErrorOrConvert(err)
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

func (ex *EvaluationContext) EvaluateMemberAccess() (*LazyRuntimeValue, LocatableError) {
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
	lazyResult := NewLazyRuntimeValue(func() (RuntimeValue, LocatableError) {
		var err error
		object, err := lazyObject.Evaluate()
		if err != nil {
			return nil, ex.LocatableErrorOrConvert(err)
		}
		for i := 0; i < len(keyPath); i++ {
			lazyObject, err = object.Lookup(keyPath[i])
			if err != nil {
				return nil, ex.LocatableErrorOrConvert(err)
			}
			object, err = lazyObject.Evaluate()
			if err != nil {
				return nil, ex.LocatableErrorOrConvert(err)
			}
		}
		return object, nil
	})
	return lazyResult, nil
}

func (ex *EvaluationContext) EvaluateTypeExpression() (*LazyRuntimeValue, LocatableError) {
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
	lazyCheckedTypeExpression := NewLazyRuntimeValue(func() (RuntimeValue, LocatableError) {
		typeValue, err := lazyTypeValue.Evaluate()
		if err != nil {
			return nil, err
		}
		enumDecl, ok := typeValue.(EnumDeclRuntimeValue)
		if !ok {
			return nil, ex.RuntimeErrorf("expected enum type, got %s", typeValue)
		}
		typeExpression := TypeExpression{typeValue: enumDecl, cases: typeCases}
		if len(enumDecl.cases) != len(typeCases) && typeCases["Any"] == nil {
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

func (ex *EvaluationContext) EvaluateGroup() (*LazyRuntimeValue, LocatableError) {
	expressionNode := ex.node.ChildByFieldName("expression")
	if expressionNode == nil {
		return nil, ex.SyntaxErrorf("expected expression")
	}
	return ex.ChildNodeExecutionContext(expressionNode).EvaluateNode()
}

func (ex *EvaluationContext) EvaluateStringLiteral() (*LazyRuntimeValue, LocatableError) {
	string, err := strconv.Unquote(ex.node.Content(ex.source))
	if err != nil {
		return nil, ex.LocatableErrorOrConvert(err)
	}
	return NewConstantRuntimeValue(PreludeString(string)), nil
}

func (ex *EvaluationContext) EvaluateBinaryExpression() (*LazyRuntimeValue, LocatableError) {
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
	return NewLazyRuntimeValue(func() (RuntimeValue, LocatableError) {
		return impl(lazyLeft, lazyRight)
	}), nil
}

func (ex *EvaluationContext) EvaluateUnaryExpression() (*LazyRuntimeValue, LocatableError) {
	return nil, ex.SyntaxErrorf("unimplemented")
}

func (ex *EvaluationContext) ParseFunctionLiteral(name string, docs DocString) (Function, LocatableError) {
	parametersNode := ex.node.ChildByFieldName("parameters")
	bodyNode := ex.node.ChildByFieldName("body")

	var (
		params []string
		err    error
	)
	if parametersNode != nil {
		params, err = ex.ChildNodeExecutionContext(parametersNode).ParseParamterList()
		if err != nil {
			return Function{}, ex.LocatableErrorOrConvert(err)
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
		docs:      Docs{name: name, docs: docs},
		parent:    ex,
		body: func(i *EvaluationContext) ([]*LazyRuntimeValue, error) {
			if bodyNode == nil {
				return []*LazyRuntimeValue{}, nil
			}
			stmts, err := i.ChildNodeExecutionContext(bodyNode).ParseStatementList()
			if err != nil {
				return nil, err
			}
			return stmts, nil
		},
	}, nil
}

func (ex *EvaluationContext) EvaluateFunctionLiteral() (*LazyRuntimeValue, LocatableError) {
	function, err := ex.ParseFunctionLiteral("", "")
	if err != nil {
		return nil, err
	}
	return NewConstantRuntimeValue(function), nil
}

func (ex *EvaluationContext) EvaluateArrayLiteral() (*LazyRuntimeValue, LocatableError) {
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
	return NewLazyRuntimeValue(func() (RuntimeValue, LocatableError) {
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
				"tail": NewLazyRuntimeValue(func() (RuntimeValue, LocatableError) {
					return SliceToList(consDecl, nilDecl, slice[1:]), nil
				}),
			},
		}
	}
}

func (ex *EvaluationContext) EvaluateNode() (*LazyRuntimeValue, LocatableError) {
	if ex.node.Type() == parser.TYPE_NODE_COMMENT {
		ex.globalComments = append(ex.globalComments, ex.node.Content(ex.source))
		return nil, nil
	}
	value, err := ex.evaluateNodeWithoutComments()
	if len(ex.globalComments) > 0 {
		ex.globalComments = make([]string, 0)
	}
	return value, err
}

func (ex *EvaluationContext) evaluateNodeWithoutComments() (*LazyRuntimeValue, LocatableError) {
	switch ex.node.Type() {
	case parser.TYPE_NODE_SOURCE_FILE:
		return ex.EvaluateSourceFile()
	case parser.TYPE_NODE_MODULE_DECLARATION:
		return ex.EvaluateModule(MakeDocString(ex.globalComments))
	case parser.TYPE_NODE_IMPORT_DECLARATION:
		return ex.EvaluateImport()
	case parser.TYPE_NODE_LET_DECLARATION:
		return ex.EvaluateLetDeclaration()
	case parser.TYPE_NODE_FUNCTION_DECLARATION:
		return ex.EvaluateFunctionDeclaration(MakeDocString(ex.globalComments))
	case parser.TYPE_NODE_DATA_DECLARATION:
		return ex.EvaluateDataDeclaration(MakeDocString(ex.globalComments))
	case parser.TYPE_NODE_EXTERN_DECLARATION:
		return ex.EvaluateExternDeclaration(MakeDocString(ex.globalComments))
	case parser.TYPE_NODE_ENUM_DECLARATION:
		return ex.EvaluateEnumDeclaration(MakeDocString(ex.globalComments))
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
