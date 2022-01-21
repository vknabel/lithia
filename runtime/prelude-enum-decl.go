package runtime

import (
	"fmt"
	"strings"

	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeValue = PreludeEnumDecl{}
var _ DeclRuntimeValue = PreludeEnumDecl{}

type PreludeEnumDecl struct {
	context *InterpreterContext
	Decl    ast.DeclEnum

	caseLookups map[ast.Identifier]Evaluatable
}

func MakeEnumDecl(context *InterpreterContext, decl ast.DeclEnum) PreludeEnumDecl {
	enumContext := context.NestedInterpreterContext(string(decl.Name))
	caseLookups := make(map[ast.Identifier]Evaluatable)
	for _, caseDecl := range decl.Cases {
		lookedUp, ok := context.environment.GetPrivate(string(caseDecl.Name))
		if !ok {
			panic(fmt.Sprintf(
				"undeclared enum case %s in %s",
				caseDecl.Name,
				strings.Join(enumContext.path, "."),
			))
		}
		caseLookups[caseDecl.Name] = lookedUp
	}
	return PreludeEnumDecl{
		context:     context,
		Decl:        decl,
		caseLookups: caseLookups,
	}
}

func (e PreludeEnumDecl) Lookup(member string) (Evaluatable, *RuntimeError) {
	value, ok := e.caseLookups[ast.Identifier(member)]
	if !ok {
		return nil, NewRuntimeErrorf("enum %s has no member %s", e, member).Cascade(*e.Decl.MetaInfo.Source)
	}
	return value, nil
}

func (PreludeEnumDecl) RuntimeType() RuntimeTypeRef {
	return PreludeAnyTypeRef
}

func (e PreludeEnumDecl) String() string {
	return fmt.Sprint(e.Decl)
}

func (e PreludeEnumDecl) HasInstance(interpreter *Interpreter, value RuntimeValue) (bool, *RuntimeError) {
	for identifier, evalCase := range e.caseLookups {
		caseValue, err := evalCase.Evaluate()
		if err != nil {
			return false, err.Cascade(*e.Decl.MetaInfo.Source)
		}
		caseDeclValue, ok := caseValue.(DeclRuntimeValue)
		if !ok {
			return false, NewRuntimeErrorf("enum case not a declaration %s, got: %s", identifier, caseDeclValue).Cascade(*e.Decl.MetaInfo.Source)
		}
		ok, err = caseDeclValue.HasInstance(interpreter, value)
		if err != nil {
			return false, err.Cascade(*e.Decl.MetaInfo.Source)
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}
