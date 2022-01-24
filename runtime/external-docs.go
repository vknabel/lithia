package runtime

import (
	"sort"
	"strings"

	"github.com/vknabel/go-lithia/ast"
)

var _ ExternalDefinition = ExternalDocs{}

type ExternalDocs struct{}

func (e ExternalDocs) Lookup(name string, env *Environment, decl ast.Decl) (RuntimeValue, bool) {
	switch name {
	case "inspect":
		return e.docsInspectValueFunction(env, decl), true
	default:
		return nil, false
	}
}

func (e ExternalDocs) docsInspectValueFunction(env *Environment, decl ast.Decl) PreludeExternFunction {
	return MakeExternFunction(
		decl,
		func(args []Evaluatable) (RuntimeValue, *RuntimeError) {
			if len(args) != 1 {
				return nil, NewRuntimeErrorf("inspect takes exactly one argument")
			}
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			return docsInspectValue(value, env)
		},
	)
}

func docsInspectValue(value RuntimeValue, env *Environment) (RuntimeValue, *RuntimeError) {
	switch value := value.(type) {
	case PreludeModule:
		sortedKeys := make([]string, 0, len(value.Module.Environment.Scope))
		for key := range value.Module.Environment.Scope {
			sortedKeys = append(sortedKeys, key)
		}
		sort.Strings(sortedKeys)

		children := make([]RuntimeValue, 0)
		for _, key := range sortedKeys {
			decl := value.Module.Environment.Scope[key]
			lazyChild, err := decl.Evaluate()
			if err != nil {
				return nil, err
			}
			child, err := docsInspectValue(lazyChild, env)
			if err != nil {
				return nil, err
			}
			children = append(children, child)
		}

		childrenList, err := env.MakeEagerList(children)
		if err != nil {
			return nil, err
		}
		moduleDocs := []string{}
		for _, file := range value.Module.Decl.Files {
			for _, decl := range file.Declarations {
				if moduleDecl, ok := decl.(ast.DeclModule); ok {
					if len(moduleDecl.Docs.Content) == 0 {
						continue
					}
					moduleDocs = append(moduleDocs, moduleDecl.Docs.Content)
				}
			}
		}
		return env.MakeDataRuntimeValue("ModuleDocs", map[string]Evaluatable{
			"name":  NewConstantRuntimeValue(PreludeString(value.Module.Name)),
			"types": NewConstantRuntimeValue(childrenList),
			"docs":  NewConstantRuntimeValue(PreludeString(strings.Join(moduleDocs, "\n"))),
		})
	case PreludeDataDecl:
		fieldDocs := make([]RuntimeValue, 0)
		for _, field := range value.Decl.Fields {
			params := make([]RuntimeValue, 0)
			for _, param := range field.Parameters {
				params = append(params, PreludeString(param.Name))
			}
			paramsList, err := env.MakeEagerList(params)
			if err != nil {
				return nil, err
			}
			doc, err := env.MakeDataRuntimeValue("DataFieldDocs", map[string]Evaluatable{
				"name":   NewConstantRuntimeValue(PreludeString(field.Name)),
				"docs":   NewConstantRuntimeValue(PreludeString(field.Docs.Content)),
				"params": NewConstantRuntimeValue(paramsList),
			})
			if err != nil {
				return nil, err
			}
			fieldDocs = append(fieldDocs, doc)
		}
		fieldDocsList, err := env.MakeEagerList(fieldDocs)
		if err != nil {
			return nil, err
		}
		return env.MakeDataRuntimeValue("DataDocs", map[string]Evaluatable{
			"name":   NewConstantRuntimeValue(PreludeString(value.Decl.Name)),
			"docs":   NewConstantRuntimeValue(PreludeString(value.Decl.Docs.Content)),
			"fields": NewConstantRuntimeValue(fieldDocsList),
		})
	case PreludePrimitiveExternType:
		fieldDocs := make([]RuntimeValue, 0)
		for _, field := range value.Fields {
			params := make([]RuntimeValue, 0)
			for _, param := range field.Parameters {
				params = append(params, PreludeString(param.Name))
			}
			paramsList, err := env.MakeEagerList(params)
			if err != nil {
				return nil, err
			}
			doc, err := env.MakeDataRuntimeValue("ExternFieldDocs", map[string]Evaluatable{
				"name":   NewConstantRuntimeValue(PreludeString(field.Name)),
				"docs":   NewConstantRuntimeValue(PreludeString(field.Docs.Content)),
				"params": NewConstantRuntimeValue(paramsList),
			})
			if err != nil {
				return nil, err
			}
			fieldDocs = append(fieldDocs, doc)
		}
		fieldDocsList, err := env.MakeEagerList(fieldDocs)
		if err != nil {
			return nil, err
		}
		return env.MakeDataRuntimeValue("ExternTypeDocs", map[string]Evaluatable{
			"name":   NewConstantRuntimeValue(PreludeString(value.Name)),
			"docs":   NewConstantRuntimeValue(PreludeString(value.Docs.Content)),
			"fields": NewConstantRuntimeValue(fieldDocsList),
		})
	case PreludeEnumDecl:
		casesDocs := make([]RuntimeValue, 0)
		for _, enumCaseDecl := range value.Decl.Cases {
			lazyEnumCase := value.caseLookups[enumCaseDecl.Name]

			enumCase, err := lazyEnumCase.Evaluate()
			if err != nil {
				return nil, err
			}

			enumCaseValueDocs, err := docsInspectValue(enumCase, env)
			if err != nil {
				return nil, err
			}
			doc, err := env.MakeDataRuntimeValue("EnumCaseDocs", map[string]Evaluatable{
				"name": NewConstantRuntimeValue(PreludeString(enumCaseDecl.Name)),
				"docs": NewConstantRuntimeValue(PreludeString("")),
				"type": NewConstantRuntimeValue(enumCaseValueDocs),
			})
			if err != nil {
				return nil, err
			}
			casesDocs = append(casesDocs, doc)
		}
		casesList, err := env.MakeEagerList(casesDocs)
		if err != nil {
			return nil, err
		}
		return env.MakeDataRuntimeValue("EnumDocs", map[string]Evaluatable{
			"name":  NewConstantRuntimeValue(PreludeString(value.Decl.Name)),
			"docs":  NewConstantRuntimeValue(PreludeString(value.Decl.Docs.Content)),
			"cases": NewConstantRuntimeValue(casesList),
		})
	case PreludeFuncDecl:
		params := make([]RuntimeValue, 0)
		for _, param := range value.Decl.Impl.Parameters {
			params = append(params, PreludeString(param.Name))
		}
		paramsList, err := env.MakeEagerList(params)
		if err != nil {
			return nil, err
		}
		return env.MakeDataRuntimeValue("FunctionDocs", map[string]Evaluatable{
			"name":   NewConstantRuntimeValue(PreludeString(value.Decl.Name)),
			"docs":   NewConstantRuntimeValue(PreludeString(value.Decl.Docs.Content)),
			"params": NewConstantRuntimeValue(paramsList),
		})
	case PreludeExternFunction:
		params := make([]RuntimeValue, 0)
		for _, param := range value.Decl.Parameters {
			params = append(params, PreludeString(param.Name))
		}
		paramsList, err := env.MakeEagerList(params)
		if err != nil {
			return nil, err
		}
		return env.MakeDataRuntimeValue("FunctionDocs", map[string]Evaluatable{
			"name":   NewConstantRuntimeValue(PreludeString(value.Decl.Name)),
			"docs":   NewConstantRuntimeValue(PreludeString("")),
			"params": NewConstantRuntimeValue(paramsList),
		})
	case RxVariableType:
		fieldDocs := make([]RuntimeValue, 0)
		for _, field := range value.Fields {
			params := make([]RuntimeValue, 0)
			for _, param := range field.Parameters {
				params = append(params, PreludeString(param.Name))
			}
			paramsList, err := env.MakeEagerList(params)
			if err != nil {
				return nil, err
			}
			doc, err := env.MakeDataRuntimeValue("ExternFieldDocs", map[string]Evaluatable{
				"name":   NewConstantRuntimeValue(PreludeString(field.Name)),
				"docs":   NewConstantRuntimeValue(PreludeString(field.Docs.Content)),
				"params": NewConstantRuntimeValue(paramsList),
			})
			if err != nil {
				return nil, err
			}
			fieldDocs = append(fieldDocs, doc)
		}
		fieldDocsList, err := env.MakeEagerList(fieldDocs)
		if err != nil {
			return nil, err
		}
		return env.MakeDataRuntimeValue("ExternTypeDocs", map[string]Evaluatable{
			"name":   NewConstantRuntimeValue(PreludeString(value.Name)),
			"docs":   NewConstantRuntimeValue(PreludeString(value.Docs.Content)),
			"fields": NewConstantRuntimeValue(fieldDocsList),
		})
	// case DocumentedRuntimeValue:
	// 	docs := value.GetDocs()
	// 	return env.MakeDataRuntimeValue("ExternDocs", map[string]Evaluatable{
	// 		"name": NewConstantRuntimeValue(PreludeString(docs.name)),
	// 		"docs": NewConstantRuntimeValue(PreludeString(docs.docs)),
	// 	})
	default:
		return env.MakeDataRuntimeValue("None", map[string]Evaluatable{})
	}
}
