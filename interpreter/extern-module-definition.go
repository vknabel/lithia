package interpreter

type ExternalDefinition interface {
	Lookup(name string, env *Environment, docs Docs) (DocumentedRuntimeValue, bool)
}

type DocumentedRuntimeValue interface {
	RuntimeValue
	GetDocs() Docs
}

type Docs struct {
	name string
	docs DocString
}
