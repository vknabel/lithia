package interpreter

import "fmt"

var _ RuntimeValue = PreludeRune('l')

type PreludeRune rune

func (PreludeRune) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Rune",
		modulePath: []string{"prelude"},
	}
}

func (r PreludeRune) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("rune %q has no member %s", r, member)
}
