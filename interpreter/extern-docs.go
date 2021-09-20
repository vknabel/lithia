package interpreter

import (
	"fmt"
	"sort"
)

var _ ExternalDefinition = ExternalDocs{}

type ExternalDocs struct{}

func (e ExternalDocs) Lookup(name string, env *Environment, docs Docs) (DocumentedRuntimeValue, bool) {
	switch name {
	case "inspect":
		return e.docsInspectValueFunction(env, docs), true
	default:
		return nil, false
	}
}

func (e ExternalDocs) docsInspectValueFunction(env *Environment, docs Docs) DocumentedRuntimeValue {
	return NewBuiltinFunction(
		"inspect",
		[]string{"value"},
		docs,
		func(args []*LazyRuntimeValue) (RuntimeValue, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("inspect takes exactly one argument")
			}
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err
			}
			return docsInspectValue(value, env)
		},
	)
}

func docsInspectValue(value RuntimeValue, env *Environment) (RuntimeValue, error) {
	switch value := value.(type) {
	case RuntimeModule:
		sortedKeys := make([]string, 0, len(value.module.environment.Scope))
		for key := range value.module.environment.Scope {
			sortedKeys = append(sortedKeys, key)
		}
		sort.Strings(sortedKeys)

		children := make([]RuntimeValue, 0)
		for _, key := range sortedKeys {
			decl := value.module.environment.Scope[key]
			var err error
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

		childrenList, err := env.eagerSliceToList(children)
		if err != nil {
			return nil, err
		}
		return env.MakeDataRuntimeValue("ModuleDocs", map[string]*LazyRuntimeValue{
			"name":  NewConstantRuntimeValue(PreludeString(value.module.name)),
			"types": NewConstantRuntimeValue(childrenList),
			"docs":  NewConstantRuntimeValue(PreludeString(value.module.docs)),
		})
	case DataDeclRuntimeValue:
		fieldDocs := make([]RuntimeValue, 0)
		for _, field := range value.fields {
			params := make([]RuntimeValue, 0)
			for _, param := range field.params {
				params = append(params, PreludeString(param))
			}
			paramsList, err := env.eagerSliceToList(params)
			if err != nil {
				return nil, err
			}
			doc, err := env.MakeDataRuntimeValue("DataFieldDocs", map[string]*LazyRuntimeValue{
				"name":   NewConstantRuntimeValue(PreludeString(field.name)),
				"docs":   NewConstantRuntimeValue(PreludeString(field.docs)),
				"params": NewConstantRuntimeValue(paramsList),
			})
			if err != nil {
				return nil, err
			}
			fieldDocs = append(fieldDocs, doc)
		}
		fieldDocsList, err := env.eagerSliceToList(fieldDocs)
		if err != nil {
			return nil, err
		}
		return env.MakeDataRuntimeValue("DataDocs", map[string]*LazyRuntimeValue{
			"name":   NewConstantRuntimeValue(PreludeString(value.name)),
			"docs":   NewConstantRuntimeValue(PreludeString(value.docs)),
			"fields": NewConstantRuntimeValue(fieldDocsList),
		})
	case EnumDeclRuntimeValue:
		casesDocs := make([]RuntimeValue, 0)
		for key, lazyEnumCase := range value.cases {
			var err error
			enumCase, err := lazyEnumCase.Evaluate()
			if err != nil {
				return nil, err
			}

			enumCaseValueDocs, err := docsInspectValue(enumCase, env)
			if err != nil {
				return nil, err
			}
			doc, err := env.MakeDataRuntimeValue("EnumCaseDocs", map[string]*LazyRuntimeValue{
				"name": NewConstantRuntimeValue(PreludeString(key)),
				"docs": NewConstantRuntimeValue(PreludeString("")),
				"type": NewConstantRuntimeValue(enumCaseValueDocs),
			})
			if err != nil {
				return nil, err
			}
			casesDocs = append(casesDocs, doc)
		}
		casesList, err := env.eagerSliceToList(casesDocs)
		if err != nil {
			return nil, err
		}
		return env.MakeDataRuntimeValue("EnumDocs", map[string]*LazyRuntimeValue{
			"name":  NewConstantRuntimeValue(PreludeString(value.name)),
			"docs":  NewConstantRuntimeValue(PreludeString(value.docs)),
			"cases": NewConstantRuntimeValue(casesList),
		})
	case Function:
		params := make([]RuntimeValue, 0)
		for _, param := range value.arguments {
			params = append(params, PreludeString(param))
		}
		paramsList, err := env.eagerSliceToList(params)
		if err != nil {
			return nil, err
		}
		return env.MakeDataRuntimeValue("FunctionDocs", map[string]*LazyRuntimeValue{
			"name":   NewConstantRuntimeValue(PreludeString(value.name)),
			"docs":   NewConstantRuntimeValue(PreludeString(value.docs.docs)),
			"params": NewConstantRuntimeValue(paramsList),
		})
	case BuiltinFunction:
		params := make([]RuntimeValue, 0)
		for _, param := range value.args {
			params = append(params, PreludeString(param))
		}
		paramsList, err := env.eagerSliceToList(params)
		if err != nil {
			return nil, err
		}
		return env.MakeDataRuntimeValue("FunctionDocs", map[string]*LazyRuntimeValue{
			"name":   NewConstantRuntimeValue(PreludeString(value.name)),
			"docs":   NewConstantRuntimeValue(PreludeString("")),
			"params": NewConstantRuntimeValue(paramsList),
		})
	case DocumentedRuntimeValue:
		docs := value.GetDocs()
		return env.MakeDataRuntimeValue("ExternDocs", map[string]*LazyRuntimeValue{
			"name": NewConstantRuntimeValue(PreludeString(docs.name)),
			"docs": NewConstantRuntimeValue(PreludeString(docs.docs)),
		})
	default:
		return env.MakeDataRuntimeValue("None", map[string]*LazyRuntimeValue{})
	}
}

func (env *Environment) eagerSliceToList(slice []RuntimeValue) (RuntimeValue, error) {
	if len(slice) == 0 {
		return env.MakeEmptyDataRuntimeValue("Nil")
	} else {
		tail, err := env.eagerSliceToList(slice[1:])
		if err != nil {
			return nil, err
		}
		return env.MakeDataRuntimeValue("Cons", map[string]*LazyRuntimeValue{
			"head": NewConstantRuntimeValue(slice[0]),
			"tail": NewConstantRuntimeValue(tail),
		})
	}
}
