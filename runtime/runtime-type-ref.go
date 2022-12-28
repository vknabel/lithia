package runtime

import (
	"fmt"

	"github.com/vknabel/lithia/ast"
)

var _ RuntimeType = RuntimeTypeRef{}

type RuntimeTypeRef struct {
	Name   ast.Identifier
	Module ast.ModuleName
}

func MakeRuntimeTypeRef(name ast.Identifier, module ast.ModuleName) RuntimeTypeRef {
	return RuntimeTypeRef{name, module}
}

func (r RuntimeTypeRef) String() string {
	return fmt.Sprintf("%s.%s", r.Module, r.Name)
}

func (r RuntimeTypeRef) Declaration(inter *Interpreter) (ast.Decl, *RuntimeError) {
	valueType, err := r.ResolveType(inter)
	if err != nil {
		return nil, err
	}
	if runtimeType, ok := valueType.(RuntimeType); ok {
		return runtimeType.Declaration(inter)
	}
	panic(fmt.Errorf("TODO: decl runtime value %s has no declaration", valueType))
}

func (r RuntimeTypeRef) ResolveType(inter *Interpreter) (DeclRuntimeValue, *RuntimeError) {
	module, ok := inter.Modules[r.Module]
	if !ok {
		return nil, NewRuntimeErrorf("module not found %s", r.Module)
	}
	// TODO: non-local types!
	value, err := module.Environment.GetExportedEvaluatedRuntimeValue(string(r.Name))
	if err != nil {
		return nil, err
	}

	if typeValue, ok := value.(DeclRuntimeValue); ok {
		return typeValue, nil
	} else {
		return nil, NewRuntimeErrorf("not a valid type %s", r)
	}
}

func (ref RuntimeTypeRef) HasInstance(inter *Interpreter, value RuntimeValue) (bool, *RuntimeError) {
	if ref == PreludeAnyTypeRef {
		return true, nil
	}
	runtimeType, err := ref.ResolveType(inter)
	if err != nil {
		return false, err
	}
	return runtimeType.HasInstance(inter, value)
}
