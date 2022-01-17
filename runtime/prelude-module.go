package runtime

var _ RuntimeValue = PreludeModule{}
var PreludeModuleTypeRef = MakeRuntimeTypeRef("Module", "prelude")

type PreludeModule struct {
	*Module
}

func (m PreludeModule) Lookup(member string) (Evaluatable, *RuntimeError) {
	value, ok := m.Module.Environment.GetExported(member)
	if !ok {
		return nil, NewRuntimeErrorf("module %s has no member %s", m.Name, member)
	}
	return value, nil
}

func (PreludeModule) RuntimeType() RuntimeTypeRef {
	return PreludeModuleTypeRef
}

func (m PreludeModule) String() string {
	return string(m.Module.Name)
}
