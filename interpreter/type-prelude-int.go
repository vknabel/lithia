package interpreter

import "fmt"

var _ RuntimeValue = PreludeInt(0)

type PreludeInt int64

func (PreludeInt) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Int",
		modulePath: []string{"prelude"},
	}
}

func (i PreludeInt) Lookup(member string) (*LazyRuntimeValue, error) {
	return nil, fmt.Errorf("int %d has no member %s", i, member)
}
