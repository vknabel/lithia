package interpreter

var _ ExternalDefinition = ExternalRx{}

type ExternalRx struct{}

func (e ExternalRx) Lookup(name string, env *Environment, docs Docs) (DocumentedRuntimeValue, bool) {
	switch name {
	case "Variable":
		return RxVariableType{docs}, true
	default:
		return nil, false
	}
}
