package rx

import (
	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/runtime"
)

var _ runtime.RuntimeValue = RxVariableType{}
var _ runtime.DeclRuntimeValue = RxVariableType{}
var _ runtime.RuntimeType = RxVariableType{}
var _ runtime.CallableRuntimeValue = RxVariableType{}

var RxVariableTypeRef = runtime.MakeRuntimeTypeRef("Variable", "rx")

type RxVariableType struct {
	ast.DeclExternType
}

func (RxVariableType) RuntimeType() runtime.RuntimeTypeRef {
	return runtime.PreludeAnyTypeRef
}

func (RxVariableType) String() string {
	return RxVariableTypeRef.String()
}

func (t RxVariableType) Declaration() (ast.Decl, *runtime.RuntimeError) {
	return t.DeclExternType, nil
}

func (d RxVariableType) HasInstance(value runtime.RuntimeValue) (bool, *runtime.RuntimeError) {
	if _, ok := value.(RxVariable); ok {
		return true, nil
	} else {
		return false, nil
	}
}

func (RxVariableType) Lookup(member string) (runtime.Evaluatable, *runtime.RuntimeError) {
	return nil, runtime.NewRuntimeErrorf("%s is not a member of %s", member, RxVariableTypeRef.String())
}

func (RxVariableType) Arity() int {
	return 1
}

func (t RxVariableType) Call(arguments []runtime.Evaluatable, fromExpr ast.Expr) (runtime.RuntimeValue, *runtime.RuntimeError) {
	if len(arguments) != 1 {
		return nil, runtime.NewRuntimeErrorf("too many arguments for variable type %s", t)
	}
	value, err := arguments[0].Evaluate()
	if err != nil {
		return nil, err.CascadeDecl(t.DeclExternType)
	}
	return MakeRxVariable(&t, value), nil
}

func (t RxVariableType) Source() *ast.Source {
	return t.Meta().Source
}
