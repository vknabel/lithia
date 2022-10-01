package runtime

import "fmt"

var _ RuntimeValue = PreludeInt(0)
var PreludeIntTypeRef = MakeRuntimeTypeRef("Int", "prelude")

type PreludeInt int64

func (PreludeInt) EagerEvaluate() *RuntimeError {
	return nil
}

func (i PreludeInt) Lookup(member string) (Evaluatable, *RuntimeError) {
	return nil, NewRuntimeError(fmt.Errorf("int %d has no member %s", i, member))
}

func (PreludeInt) RuntimeType() RuntimeTypeRef {
	return PreludeIntTypeRef
}

func (i PreludeInt) String() string {
	return fmt.Sprintf("%d", i)
}
