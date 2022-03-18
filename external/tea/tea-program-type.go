package tea

import (
	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/runtime"
)

var _ runtime.RuntimeValue = TeaProgramType{}
var _ runtime.DeclRuntimeValue = TeaProgramType{}
var _ runtime.RuntimeType = TeaProgramType{}
var _ runtime.CallableRuntimeValue = TeaProgramType{}

var TeaProgramTypeRef = runtime.MakeRuntimeTypeRef("Program", "tea")

type TeaProgramType struct {
	ast.DeclExternType
}

func (TeaProgramType) RuntimeType() runtime.RuntimeTypeRef {
	return runtime.PreludeAnyTypeRef
}

func (TeaProgramType) String() string {
	return TeaProgramTypeRef.String()
}

func (t TeaProgramType) Declaration() (ast.Decl, *runtime.RuntimeError) {
	return t.DeclExternType, nil
}

func (d TeaProgramType) HasInstance(value runtime.RuntimeValue) (bool, *runtime.RuntimeError) {
	if _, ok := value.(TeaProgram); ok {
		return true, nil
	} else {
		return false, nil
	}
}

func (TeaProgramType) Lookup(member string) (runtime.Evaluatable, *runtime.RuntimeError) {
	return nil, runtime.NewRuntimeErrorf("%s is not a member of %s", member, TeaProgramTypeRef.String())
}

func (TeaProgramType) Arity() int {
	return 1
}

func (t TeaProgramType) Call(arguments []runtime.Evaluatable, fromExpr ast.Expr) (runtime.RuntimeValue, *runtime.RuntimeError) {
	if len(arguments) != 1 {
		return nil, runtime.NewRuntimeErrorf("too many arguments for variable type %s", t)
	}
	value, err := arguments[0].Evaluate()
	if err != nil {
		return nil, err.CascadeDecl(t.DeclExternType)
	}
	model, ok := value.(TeaModel)
	if !ok {
		return nil, runtime.NewRuntimeErrorf("%s is not a Model", value)
	}
	return MakeTeaProgram(&t, model), nil
}

func (t TeaProgramType) Source() *ast.Source {
	return t.Meta().Source
}
