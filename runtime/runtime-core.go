package runtime

import "github.com/vknabel/go-lithia/ast"

type Evaluatable interface {
	Evaluate() (RuntimeValue, RuntimeError)
}

type RuntimeValue interface {
	// for printing
	String() string
	// Member access
	Lookup(name string) (Evaluatable, error)
	RuntimeType() RuntimeTypeRef
}

type RuntimeTypeRef struct {
	Name   ast.Identifier
	Module ast.ModuleName
}

func MakeRuntimeTypeRef(name ast.Identifier, module ast.ModuleName) RuntimeTypeRef {
	return RuntimeTypeRef{name, module}
}

type RuntimeType interface {
	RuntimeValue

	Declaration() *ast.Decl
	IncludesValue(value RuntimeValue) bool
}

type CallableRuntimeValue interface {
	RuntimeValue
	Call(args []Evaluatable) (RuntimeValue, RuntimeError)
}
