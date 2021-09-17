package interpreter

type ExternalDefinition interface {
	Lookup(name string, env *Environment) (RuntimeValue, bool)
}
