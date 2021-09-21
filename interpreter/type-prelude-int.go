package interpreter

import "fmt"

var _ RuntimeValue = PreludeInt(0)

type PreludeInt int64

func (i PreludeInt) String() string {
	return fmt.Sprintf("%d", i)
}

func (PreludeInt) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Int",
		moduleName: "prelude",
	}
}

func (i PreludeInt) Lookup(member string) (Evaluatable, error) {
	return nil, fmt.Errorf("int %d has no member %s", i, member)
}
