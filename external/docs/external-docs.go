package docs

import (
	"sort"
	"strings"

	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/runtime"
)

var _ runtime.ExternalDefinition = ExternalDocs{}

type ExternalDocs struct{}

func (e ExternalDocs) Lookup(name string, env *runtime.Environment, decl ast.Decl) (runtime.RuntimeValue, bool) {
	switch name {
	case "inspect":
		return e.docsInspectValueFunction(env, decl), true
	default:
		return nil, false
	}
}

func (e ExternalDocs) docsInspectValueFunction(env *runtime.Environment, decl ast.Decl) runtime.PreludeExternFunction {
	return runtime.MakeExternFunction(
		decl,
		func(args []runtime.Evaluatable) (runtime.RuntimeValue, *runtime.RuntimeError) {
			if len(args) != 1 {
				return nil, runtime.NewRuntimeErrorf("inspect takes exactly one argument")
			}
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			return docsInspectValue(value, env)
		},
	)
}

func docsInspectValue(value runtime.RuntimeValue, env *runtime.Environment) (runtime.RuntimeValue, *runtime.RuntimeError) {
	switch value := value.(type) {
	case runtime.PreludeModule:
		sortedKeys := make([]string, 0, len(value.Module.Environment.Scope))
		for key := range value.Module.Environment.Scope {
			sortedKeys = append(sortedKeys, key)
		}
		sort.Strings(sortedKeys)

		children := make([]runtime.RuntimeValue, 0)
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
		return env.MakeDataRuntimeValue("ModuleDocs", map[string]runtime.Evaluatable{
			"name":  runtime.NewConstantRuntimeValue(runtime.PreludeString(value.Module.Name)),
			"types": runtime.NewConstantRuntimeValue(childrenList),
			"docs":  runtime.NewConstantRuntimeValue(runtime.PreludeString(strings.Join(moduleDocs, "\n"))),
		})
	case runtime.PreludeDataDecl:
		fieldDocs := make([]runtime.RuntimeValue, 0)
		for _, field := range value.Decl.Fields {
			params := make([]runtime.RuntimeValue, 0)
			for _, param := range field.Parameters {
				params = append(params, runtime.PreludeString(param.Name))
			}
			paramsList, err := env.MakeEagerList(params)
			if err != nil {
				return nil, err
			}
			doc, err := env.MakeDataRuntimeValue("DataFieldDocs", map[string]runtime.Evaluatable{
				"name":   runtime.NewConstantRuntimeValue(runtime.PreludeString(field.Name)),
				"docs":   runtime.NewConstantRuntimeValue(runtime.PreludeString(field.Docs.Content)),
				"params": runtime.NewConstantRuntimeValue(paramsList),
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
		return env.MakeDataRuntimeValue("DataDocs", map[string]runtime.Evaluatable{
			"name":   runtime.NewConstantRuntimeValue(runtime.PreludeString(value.Decl.Name)),
			"docs":   runtime.NewConstantRuntimeValue(runtime.PreludeString(value.Decl.Docs.Content)),
			"fields": runtime.NewConstantRuntimeValue(fieldDocsList),
		})
	case runtime.PreludePrimitiveExternType:
		fieldDocs := make([]runtime.RuntimeValue, 0)
		for _, field := range value.Fields {
			params := make([]runtime.RuntimeValue, 0)
			for _, param := range field.Parameters {
				params = append(params, runtime.PreludeString(param.Name))
			}
			paramsList, err := env.MakeEagerList(params)
			if err != nil {
				return nil, err
			}
			doc, err := env.MakeDataRuntimeValue("ExternFieldDocs", map[string]runtime.Evaluatable{
				"name":   runtime.NewConstantRuntimeValue(runtime.PreludeString(field.Name)),
				"docs":   runtime.NewConstantRuntimeValue(runtime.PreludeString(field.Docs.Content)),
				"params": runtime.NewConstantRuntimeValue(paramsList),
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
		return env.MakeDataRuntimeValue("ExternTypeDocs", map[string]runtime.Evaluatable{
			"name":   runtime.NewConstantRuntimeValue(runtime.PreludeString(value.Name)),
			"docs":   runtime.NewConstantRuntimeValue(runtime.PreludeString(value.Docs.Content)),
			"fields": runtime.NewConstantRuntimeValue(fieldDocsList),
		})
	case runtime.PreludeEnumDecl:
		casesDocs := make([]runtime.RuntimeValue, 0)
		for _, enumCaseDecl := range value.Decl.Cases {
			lazyEnumCase, err := value.LookupCase(string(enumCaseDecl.Name))
			if err != nil {
				return nil, err
			}

			enumCase, err := lazyEnumCase.Evaluate()
			if err != nil {
				return nil, err
			}

			enumCaseValueDocs, err := docsInspectValue(enumCase, env)
			if err != nil {
				return nil, err
			}
			doc, err := env.MakeDataRuntimeValue("EnumCaseDocs", map[string]runtime.Evaluatable{
				"name": runtime.NewConstantRuntimeValue(runtime.PreludeString(enumCaseDecl.Name)),
				"docs": runtime.NewConstantRuntimeValue(runtime.PreludeString("")),
				"type": runtime.NewConstantRuntimeValue(enumCaseValueDocs),
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
		return env.MakeDataRuntimeValue("EnumDocs", map[string]runtime.Evaluatable{
			"name":  runtime.NewConstantRuntimeValue(runtime.PreludeString(value.Decl.Name)),
			"docs":  runtime.NewConstantRuntimeValue(runtime.PreludeString(value.Decl.Docs.Content)),
			"cases": runtime.NewConstantRuntimeValue(casesList),
		})
	case runtime.PreludeFuncDecl:
		params := make([]runtime.RuntimeValue, 0)
		for _, param := range value.Decl.Impl.Parameters {
			params = append(params, runtime.PreludeString(param.Name))
		}
		paramsList, err := env.MakeEagerList(params)
		if err != nil {
			return nil, err
		}
		return env.MakeDataRuntimeValue("FunctionDocs", map[string]runtime.Evaluatable{
			"name":   runtime.NewConstantRuntimeValue(runtime.PreludeString(value.Decl.Name)),
			"docs":   runtime.NewConstantRuntimeValue(runtime.PreludeString(value.Decl.Docs.Content)),
			"params": runtime.NewConstantRuntimeValue(paramsList),
		})
	case runtime.PreludeExternFunction:
		params := make([]runtime.RuntimeValue, 0)
		for _, param := range value.Decl.Parameters {
			params = append(params, runtime.PreludeString(param.Name))
		}
		paramsList, err := env.MakeEagerList(params)
		if err != nil {
			return nil, err
		}
		return env.MakeDataRuntimeValue("ExternFunctionDocs", map[string]runtime.Evaluatable{
			"name":   runtime.NewConstantRuntimeValue(runtime.PreludeString(value.Decl.Name)),
			"docs":   runtime.NewConstantRuntimeValue(runtime.PreludeString(value.Decl.Docs.Content)),
			"params": runtime.NewConstantRuntimeValue(paramsList),
		})

	case runtime.RuntimeType:
		decl, err := value.Declaration()
		if err != nil {
			return nil, err
		}
		switch decl := decl.(type) {
		case ast.DeclExternType:
			fieldDocs := make([]runtime.RuntimeValue, 0)
			for _, field := range decl.Fields {
				params := make([]runtime.RuntimeValue, 0)
				for _, param := range field.Parameters {
					params = append(params, runtime.PreludeString(param.Name))
				}
				paramsList, err := env.MakeEagerList(params)
				if err != nil {
					return nil, err
				}
				doc, err := env.MakeDataRuntimeValue("ExternFieldDocs", map[string]runtime.Evaluatable{
					"name":   runtime.NewConstantRuntimeValue(runtime.PreludeString(field.Name)),
					"docs":   runtime.NewConstantRuntimeValue(runtime.PreludeString(field.Docs.Content)),
					"params": runtime.NewConstantRuntimeValue(paramsList),
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
			return env.MakeDataRuntimeValue("ExternTypeDocs", map[string]runtime.Evaluatable{
				"name":   runtime.NewConstantRuntimeValue(runtime.PreludeString(decl.Name)),
				"docs":   runtime.NewConstantRuntimeValue(runtime.PreludeString(decl.Docs.Content)),
				"fields": runtime.NewConstantRuntimeValue(fieldDocsList),
			})
		default:
			return env.MakeDataRuntimeValue("None", map[string]runtime.Evaluatable{})
		}
	// case DocumentedRuntimeValue:
	// 	docs := value.GetDocs()
	// 	return env.MakeDataRuntimeValue("ExternDocs", map[string]Evaluatable{
	// 		"name": runtime.NewConstantRuntimeValue(PreludeString(docs.name)),
	// 		"docs": runtime.NewConstantRuntimeValue(PreludeString(docs.docs)),
	// 	})
	default:
		return env.MakeDataRuntimeValue("None", map[string]runtime.Evaluatable{})
	}
}
