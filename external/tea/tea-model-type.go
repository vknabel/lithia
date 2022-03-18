package tea

import (
	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/runtime"
)

var _ runtime.RuntimeValue = TeaModelType{}
var _ runtime.DeclRuntimeValue = TeaModelType{}
var _ runtime.RuntimeType = TeaModelType{}
var _ runtime.CallableRuntimeValue = TeaModelType{}

var TeaModelTypeRef = runtime.MakeRuntimeTypeRef("Model", "tea")

type TeaModelType struct {
	ast.DeclExternType
	env *runtime.Environment
}

func (TeaModelType) RuntimeType() runtime.RuntimeTypeRef {
	return runtime.PreludeAnyTypeRef
}

func (TeaModelType) String() string {
	return TeaModelTypeRef.String()
}

func (t TeaModelType) Declaration() (ast.Decl, *runtime.RuntimeError) {
	return t.DeclExternType, nil
}

func (d TeaModelType) HasInstance(value runtime.RuntimeValue) (bool, *runtime.RuntimeError) {
	if _, ok := value.(TeaModel); ok {
		return true, nil
	} else {
		return false, nil
	}
}

func (TeaModelType) Lookup(member string) (runtime.Evaluatable, *runtime.RuntimeError) {
	return nil, runtime.NewRuntimeErrorf("%s is not a member of %s", member, TeaModelTypeRef.String())
}

func (TeaModelType) Arity() int {
	return 3
}

func (t TeaModelType) Call(arguments []runtime.Evaluatable, fromExpr ast.Expr) (runtime.RuntimeValue, *runtime.RuntimeError) {
	if len(arguments) != 3 {
		return nil, runtime.NewRuntimeErrorf("too many arguments for variable type %s", t)
	}
	init := arguments[0]
	update, err := arguments[1].Evaluate()
	if err != nil {
		return nil, err.CascadeDecl(t.DeclExternType)
	}
	view, err := arguments[2].Evaluate()
	if err != nil {
		return nil, err.CascadeDecl(t.DeclExternType)
	}
	return MakeTeaModel(t.env, init, update, view), nil
}

func (t TeaModelType) Source() *ast.Source {
	return t.Meta().Source
}
