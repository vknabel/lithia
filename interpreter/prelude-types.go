package interpreter

import (
	"fmt"
	"sync"
)

type PreludeInt int64
type PreludeFloat float64
type PreludeString string
type PreludeRune rune
type PreludeFunctionType struct{}
type PreludeVariableType struct{}
type PreludeModuleType struct{}
type PreludeAnyType struct{}

func (PreludeInt) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Int",
		modulePath: []string{"prelude"},
	}
}

func (i PreludeInt) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("int %d has no member %s", i, member)
}

func (PreludeFloat) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Float",
		modulePath: []string{"prelude"},
	}
}

func (f PreludeFloat) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("float %f has no member %s", f, member)
}

func (PreludeString) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "String",
		modulePath: []string{"prelude"},
	}
}

func (s PreludeString) Lookup(member string) (*LazyRuntimeValue, error) {
	switch member {
	case "length":
		return NewConstantRuntimeValue(PreludeInt(len(s))), nil
	case "append":
		return NewConstantRuntimeValue(BuiltinFunction{
			name: "append",
			args: []string{""},
			impl: func(arguments []*LazyRuntimeValue) (RuntimeValue, error) {
				arg, err := arguments[0].Evaluate()
				if err != nil {
					return nil, err
				}
				switch arg := arg.(type) {
				case PreludeString:
					return PreludeString(s + arg), nil
				case PreludeInt:
					return PreludeString(s + PreludeString(fmt.Sprint(arg))), nil
				default:
					return nil, fmt.Errorf("append expects string argument, got %s", arg)
				}
			},
		}), nil
	}
	return nil, fmt.Errorf("string %q has no member %s", s, member)
}

func (PreludeRune) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Rune",
		modulePath: []string{"prelude"},
	}
}

func (r PreludeRune) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("rune %q has no member %s", r, member)
}

func (PreludeFunctionType) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Function",
		modulePath: []string{"prelude"},
	}
}

func (f PreludeFunctionType) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("function type %s has no member %s", f, member)
}

func (PreludeVariableType) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Variable",
		modulePath: []string{"prelude"},
	}
}

func (v PreludeVariableType) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("variable type %s has no member %s", v, member)
}

func (v PreludeVariableType) String() string {
	return "(extern Variable)"
}

func (v PreludeVariableType) Call(arguments []*LazyRuntimeValue) (RuntimeValue, error) {
	var _ Callable = v
	if len(arguments) != 1 {
		return nil, fmt.Errorf("too many arguments for variable type %s", v)
	}
	value, err := arguments[0].Evaluate()
	if err != nil {
		return nil, err
	}
	return &RuntimeVariable{current: value, lock: &sync.RWMutex{}}, nil
}

func (PreludeModuleType) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Module",
		modulePath: []string{"prelude"},
	}
}

func (a PreludeModuleType) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("module type %s has no member %s", a, member)
}

func (PreludeAnyType) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Any",
		modulePath: []string{"prelude"},
	}
}

func (a PreludeAnyType) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("any type %s has no member %s", a, member)
}

func (t RuntimeType) RuntimeType() RuntimeType {
	typeValue := t.typeValue
	if typeValue == nil {
		return PreludeAnyType{}.RuntimeType()
	} else {
		return (*typeValue).RuntimeType()
	}
}

func (t RuntimeType) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("type %s has no member %s", t, member)
}
