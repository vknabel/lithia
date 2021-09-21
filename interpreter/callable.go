package interpreter

type Callable interface {
	Call(arguments []Evaluatable) (RuntimeValue, error)
	String() string
}
