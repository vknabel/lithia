package interpreter

import "fmt"

var _ RuntimeValue = PreludeString("")

type PreludeString string

func (s PreludeString) String() string {
	return string(s)
}

func (PreludeString) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "String",
		moduleName: "prelude",
	}
}

func (s PreludeString) Lookup(member string) (Evaluatable, error) {
	switch member {
	case "length":
		return NewConstantRuntimeValue(PreludeInt(len(s))), nil
	case "append":
		return NewConstantRuntimeValue(BuiltinFunction{
			name: "append",
			args: []string{""},
			impl: func(arguments []Evaluatable) (RuntimeValue, error) {
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
