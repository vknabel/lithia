package interpreter

import "fmt"

var _ RuntimeValue = PreludeFloat(0)

type PreludeFloat float64

func (PreludeFloat) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Float",
		modulePath: []string{"prelude"},
	}
}

func (f PreludeFloat) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("float %f has no member %s", f, member)
}
