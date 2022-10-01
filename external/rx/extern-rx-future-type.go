package rx

import (
	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/runtime"
)

var _ runtime.RuntimeValue = RxFutureType{}
var _ runtime.DeclRuntimeValue = RxFutureType{}
var _ runtime.RuntimeType = RxFutureType{}
var _ runtime.CallableRuntimeValue = RxFutureType{}

var RxFutureTypeRef = runtime.MakeRuntimeTypeRef("Future", "rx")

type RxFutureType struct {
	ast.DeclExternType
}

func (RxFutureType) RuntimeType() runtime.RuntimeTypeRef {
	return runtime.PreludeAnyTypeRef
}

func (RxFutureType) String() string {
	return RxVariableTypeRef.String()
}

func (t RxFutureType) Declaration() (ast.Decl, *runtime.RuntimeError) {
	return t.DeclExternType, nil
}

func (d RxFutureType) HasInstance(value runtime.RuntimeValue) (bool, *runtime.RuntimeError) {
	if _, ok := value.(RxFuture); ok {
		return true, nil
	} else {
		return false, nil
	}
}

func (RxFutureType) Lookup(member string) (runtime.Evaluatable, *runtime.RuntimeError) {
	return nil, runtime.NewRuntimeErrorf("%s is not a member of %s", member, RxVariableTypeRef.String())
}

func (RxFutureType) Arity() int {
	return 1
}

func (t RxFutureType) Call(arguments []runtime.Evaluatable, fromExpr ast.Expr) (runtime.RuntimeValue, *runtime.RuntimeError) {
	if len(arguments) != 1 {
		return nil, runtime.NewRuntimeErrorf("too many arguments for variable type %s", t)
	}
	receive, err := arguments[0].Evaluate()
	if err != nil {
		return nil, err.CascadeDecl(t.DeclExternType)
	}
	if receive, ok := receive.(runtime.CallableRuntimeValue); ok {
		return MakeRxFuture(&t, receive), nil
	} else {
		return nil, runtime.NewRuntimeErrorf("%s is not callable", receive)
	}
}

func (t RxFutureType) Source() *ast.Source {
	return t.Meta().Source
}
