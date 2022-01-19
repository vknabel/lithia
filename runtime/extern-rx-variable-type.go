package runtime

import (
	"github.com/vknabel/go-lithia/ast"
)

var _ RuntimeValue = RxVariableType{}
var _ RuntimeType = RxVariableType{}
var _ CallableRuntimeValue = RxVariableType{}

var RxVariableTypeRef = MakeRuntimeTypeRef("Variable", "rx")

type RxVariableType struct {
	ast.DeclExternType
}

func (RxVariableType) RuntimeType() RuntimeTypeRef {
	return PreludeAnyTypeRef
}

func (RxVariableType) String() string {
	return RxVariableTypeRef.String()
}

func (t RxVariableType) Declaration(*Interpreter) (ast.Decl, *RuntimeError) {
	return t.DeclExternType, nil
}

func (d RxVariableType) IncludesValue(interpreter *Interpreter, value RuntimeValue) (bool, *RuntimeError) {
	if _, ok := value.(RxVariable); ok {
		return true, nil
	} else {
		return false, nil
	}
}

func (RxVariableType) Lookup(member string) (Evaluatable, *RuntimeError) {
	return nil, NewRuntimeErrorf("%s is not a member of %s", member, RxVariableTypeRef.String())
}

func (RxVariableType) Arity() int {
	return 1
}

func (t RxVariableType) Call(arguments []Evaluatable) (RuntimeValue, *RuntimeError) {
	if len(arguments) != 1 {
		return nil, NewRuntimeErrorf("too many arguments for variable type %s", t)
	}
	value, err := arguments[0].Evaluate()
	if err != nil {
		return nil, err.Cascade(*t.MetaInfo.Source)
	}
	return MakeRxVariable(&t, value), nil
}
