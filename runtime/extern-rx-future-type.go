package runtime

import (
	"github.com/vknabel/lithia/ast"
)

var _ RuntimeValue = RxFutureType{}
var _ DeclRuntimeValue = RxFutureType{}
var _ CallableRuntimeValue = RxFutureType{}

var RxFutureTypeRef = MakeRuntimeTypeRef("Future", "rx")

type RxFutureType struct {
	ast.DeclExternType
}

func (RxFutureType) RuntimeType() RuntimeTypeRef {
	return PreludeAnyTypeRef
}

func (RxFutureType) String() string {
	return RxVariableTypeRef.String()
}

func (t RxFutureType) Declaration(*Interpreter) (ast.Decl, *RuntimeError) {
	return t.DeclExternType, nil
}

func (d RxFutureType) HasInstance(value RuntimeValue) (bool, *RuntimeError) {
	if _, ok := value.(RxVariable); ok {
		return true, nil
	} else {
		return false, nil
	}
}

func (RxFutureType) Lookup(member string) (Evaluatable, *RuntimeError) {
	return nil, NewRuntimeErrorf("%s is not a member of %s", member, RxVariableTypeRef.String())
}

func (RxFutureType) Arity() int {
	return 1
}

func (t RxFutureType) Call(arguments []Evaluatable, fromExpr ast.Expr) (RuntimeValue, *RuntimeError) {
	if len(arguments) != 1 {
		return nil, NewRuntimeErrorf("too many arguments for variable type %s", t)
	}
	receive, err := arguments[0].Evaluate()
	if err != nil {
		return nil, err.CascadeDecl(t.DeclExternType)
	}
	if receive, ok := receive.(CallableRuntimeValue); ok {
		return MakeRxFuture(&t, receive), nil
	} else {
		return nil, NewRuntimeErrorf("%s is not callable", receive)
	}
}

func (t RxFutureType) Source() *ast.Source {
	return t.Meta().Source
}
