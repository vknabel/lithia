package runtime

import (
	"fmt"
	"strings"

	"github.com/vknabel/lithia/ast"
)

var _ RuntimeValue = PreludeEnumDecl{}
var _ DeclRuntimeValue = PreludeEnumDecl{}

type PreludeEnumDecl struct {
	context *InterpreterContext
	Decl    ast.DeclEnum

	caseLookups map[ast.Identifier]Evaluatable
}

func MakeEnumDecl(context *InterpreterContext, decl ast.DeclEnum) (PreludeEnumDecl, *RuntimeError) {
	enumContext := context.NestedInterpreterContext(string(decl.Name))
	caseLookups := make(map[ast.Identifier]Evaluatable)
	for _, caseDecl := range decl.Cases {
		lookedUp, ok := context.environment.GetPrivate(string(caseDecl.Name))
		if !ok {
			return PreludeEnumDecl{}, NewTypeErrorf(
				"undeclared enum case %s in %s",
				caseDecl.Name,
				strings.Join(enumContext.path, "."),
			).CascadeDecl(decl)
		}
		caseLookups[caseDecl.Name] = lookedUp
	}
	return PreludeEnumDecl{
		context:     context,
		Decl:        decl,
		caseLookups: caseLookups,
	}, nil
}

func (e PreludeEnumDecl) Lookup(member string) (Evaluatable, *RuntimeError) {
	return nil, NewRuntimeErrorf("cannot access member %s of enum type %s, see https://github.com/vknabel/lithia/discussions/25", member, e.Decl.Name)
}

func (e PreludeEnumDecl) LookupCase(member string) (Evaluatable, *RuntimeError) {
	value, ok := e.caseLookups[ast.Identifier(member)]
	if !ok {
		return nil, NewRuntimeErrorf("enum %s has no member %s", e, member).CascadeDecl(e.Decl)
	}
	return value, nil
}

func (PreludeEnumDecl) RuntimeType() RuntimeTypeRef {
	return PreludeAnyTypeRef
}

func (e PreludeEnumDecl) String() string {
	return fmt.Sprint(e.Decl)
}

func (e PreludeEnumDecl) HasInstance(value RuntimeValue) (bool, *RuntimeError) {
	for identifier, evalCase := range e.caseLookups {
		caseValue, err := evalCase.Evaluate()
		if err != nil {
			return false, err.CascadeDecl(e.Decl)
		}
		caseDeclValue, ok := caseValue.(DeclRuntimeValue)
		if !ok {
			return false, NewRuntimeErrorf("enum case not a declaration %s, got: %s", identifier, caseDeclValue).CascadeDecl(e.Decl)
		}
		ok, err = caseDeclValue.HasInstance(value)
		if err != nil {
			return false, err.CascadeDecl(e.Decl)
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}
