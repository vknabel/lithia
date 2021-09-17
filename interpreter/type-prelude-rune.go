package interpreter

import "fmt"

var _ RuntimeValue = PreludeRune('l')

type PreludeRune rune

func (r PreludeRune) String() string {
	return string(r)
}

func (PreludeRune) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Rune",
		moduleName: "prelude",
	}
}

func (r PreludeRune) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("rune %q has no member %s", r, member)
}
