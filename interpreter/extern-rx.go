package interpreter

var _ ExternalDefinition = ExternalRx{}

type ExternalRx struct{}

func (e ExternalRx) Lookup(name string, env *Environment) (RuntimeValue, bool) {
	switch name {
	case "Variable":
		return RxVariableType{}, true
	default:
		return nil, false
	}
}
