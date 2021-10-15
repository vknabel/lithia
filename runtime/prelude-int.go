package runtime

import "fmt"

var _ RuntimeValue = PreludeInt(0)
var PreludeIntTypeRef = MakeRuntimeTypeRef("Int", "prelude")

type PreludeInt int64

func (i PreludeInt) Lookup(member string) (Evaluatable, error) {
	return nil, fmt.Errorf("int %d has no member %s", i, member)
}

func (PreludeInt) RuntimeType() RuntimeTypeRef {
	return PreludeIntTypeRef
}

func (i PreludeInt) String() string {
	return fmt.Sprintf("%d", i)
}
