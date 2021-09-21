package interpreter

import "fmt"

var _ RuntimeValue = PreludeChar('l')

type PreludeChar rune

func (r PreludeChar) String() string {
	return string(r)
}

func (PreludeChar) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Char",
		moduleName: "prelude",
	}
}

func (r PreludeChar) Lookup(member string) (Evaluatable, error) {
	return nil, fmt.Errorf("char %q has no member %s", r, member)
}
