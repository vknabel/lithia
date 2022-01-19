package runtime

import (
	"fmt"

	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeType = RuntimeTypeRef{}

type RuntimeTypeRef struct {
	Name ast.Identifier
	// TODO: rename to ModuleName
	Module ast.ModuleName
}

func MakeRuntimeTypeRef(name ast.Identifier, module ast.ModuleName) RuntimeTypeRef {
	return RuntimeTypeRef{name, module}
}

func (r RuntimeTypeRef) String() string {
	return fmt.Sprintf("%s.%s", r.Module, r.Name)
}

func (r RuntimeTypeRef) Declaration(interpreter *Interpreter) (ast.Decl, *RuntimeError) {
	module, ok := interpreter.Modules[r.Module]
	if !ok {
		return nil, NewRuntimeErrorf("module not found %s", r.Module)
	}
	value, err := module.Environment.GetExportedEvaluatedRuntimeValue(string(r.Name))
	if err != nil {
		return nil, err
	}

	switch value := value.(type) {
	case PreludeEnumDecl:
		return value.Decl, nil
	case PreludeDataDecl:
		return value.Decl, nil
	default:
		// TODO: External Data?
		return nil, NewRuntimeErrorf("not a valid type %s", r.Name)
	}
}

func (ref RuntimeTypeRef) IncludesValue(interpreter *Interpreter, value RuntimeValue) (bool, *RuntimeError) {
	// TODO: Enums
	if ref == value.RuntimeType() {
		return true, nil
	}
	refDecl, err := ref.Declaration(interpreter)
	if err != nil {
		return false, err
	}
	valueDecl, err := value.RuntimeType().Declaration(interpreter)
	if err != nil {
		return false, err
	}
	return refDecl == valueDecl, nil
}
