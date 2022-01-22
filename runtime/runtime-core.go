package runtime

import "github.com/vknabel/go-lithia/ast"

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
	HasInstance(interpreter *Interpreter, value RuntimeValue) (bool, *RuntimeError)
}

type CallableRuntimeValue interface {
	RuntimeValue
	Arity() int
	Call(args []Evaluatable) (RuntimeValue, *RuntimeError)
}

type DeclRuntimeValue interface {
	RuntimeValue
	HasInstance(interpreter *Interpreter, value RuntimeValue) (bool, *RuntimeError)
}
