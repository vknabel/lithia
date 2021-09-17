package interpreter

type Callable interface {
	Call(arguments []*LazyRuntimeValue) (RuntimeValue, error)
	String() string
}
