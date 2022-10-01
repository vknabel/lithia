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
	Declaration() (ast.Decl, *RuntimeError)
	HasInstance(value RuntimeValue) (bool, *RuntimeError)
}

type EagerEvaluatableRuntimeValue interface {
	EagerEvaluate() *RuntimeError
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
	HasInstance(value RuntimeValue) (bool, *RuntimeError)
}
