package interpreter

type PreludeInt int64
type PreludeFloat float64
type PreludeString string
type PreludeRune rune
type PreludeFunctionType struct{}
type PreludeAnyType struct{}

func (PreludeInt) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Int",
		modulePath: []string{"prelude"},
	}
}

func (PreludeFloat) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Float",
		modulePath: []string{"prelude"},
	}
}

func (PreludeString) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "String",
		modulePath: []string{"prelude"},
	}
}

func (PreludeRune) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Rune",
		modulePath: []string{"prelude"},
	}
}

func (PreludeFunctionType) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Function",
		modulePath: []string{"prelude"},
	}
}

func (PreludeAnyType) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Any",
		modulePath: []string{"prelude"},
	}
}

func (t RuntimeType) RuntimeType() RuntimeType {
	typeValue := t.typeValue
	if typeValue == nil {
		return PreludeAnyType{}.RuntimeType()
	} else {
		return (*typeValue).RuntimeType()
	}
}
