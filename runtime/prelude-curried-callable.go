package runtime

import (
	"fmt"

	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeValue = PreludeCurriedCallable{}
var _ CallableRuntimeValue = PreludeCurriedCallable{}

type PreludeCurriedCallable struct {
	actual         CallableRuntimeValue
	arguments      []Evaluatable
	remainingArity int
}

func MakeCurriedCallable(actual CallableRuntimeValue, arguments []Evaluatable) PreludeCurriedCallable {
	return PreludeCurriedCallable{
		actual:         actual,
		arguments:      arguments,
		remainingArity: actual.Arity() - len(arguments),
	}
}

func (f PreludeCurriedCallable) Lookup(member string) (Evaluatable, *RuntimeError) {
	switch member {
	case "arity":
		return NewConstantRuntimeValue(PreludeInt(f.Arity())), nil
	default:
		return nil, NewRuntimeErrorf("no such member: %s", member)
	}
}

func (PreludeCurriedCallable) RuntimeType() RuntimeTypeRef {
	return PreludeFunctionTypeRef
}

func (f PreludeCurriedCallable) String() string {
	return fmt.Sprintf("%s curried by %d", f.actual.String(), f.remainingArity)
}

func (f PreludeCurriedCallable) Arity() int {
	return f.remainingArity
}

func (f PreludeCurriedCallable) Call(args []Evaluatable, fromExpr ast.Expr) (RuntimeValue, *RuntimeError) {
	if len(args) != f.Arity() {
		panic("use Call to call functions!")
	}
	allArgs := append(f.arguments, args...)
	return Call(f.actual, allArgs, fromExpr)
}

func (f PreludeCurriedCallable) Source() *ast.Source {
	return nil
}
