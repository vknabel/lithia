package runtime

import (
	"github.com/vknabel/go-lithia/ast"
)

type Evaluatable interface {
	Evaluate() (RuntimeValue, *RuntimeError)
}

type RuntimeValue interface {
	// for printing
	String() string
	// Member access
	Lookup(name string) (Evaluatable, *RuntimeError)
	RuntimeType() RuntimeTypeRef
}

type RuntimeType interface {
	Declaration(*Interpreter) (ast.Decl, *RuntimeError)
	IncludesValue(interpreter *Interpreter, value RuntimeValue) (bool, *RuntimeError)
}

type CallableRuntimeValue interface {
	RuntimeValue
	Call(args []Evaluatable) (RuntimeValue, RuntimeError)
}
