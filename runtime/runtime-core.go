package runtime

import "github.com/vknabel/lithia/ast"

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
	Call(args []Evaluatable, fromExpr ast.Expr) (RuntimeValue, *RuntimeError)
	// An optional source for stack traces
	// If nil, no stack trace will be printed
	Source() *ast.Source
}

type DeclRuntimeValue interface {
	RuntimeValue
	HasInstance(interpreter *Interpreter, value RuntimeValue) (bool, *RuntimeError)
}
