package runtime

import "fmt"

var _ RuntimeValue = PreludeFloat(0.0)
var PreludeFloatTypeRef = MakeRuntimeTypeRef("Float", "prelude")

type PreludeFloat float64

func (i PreludeFloat) Lookup(member string) (Evaluatable, error) {
	return nil, fmt.Errorf("float %f has no member %s", i, member)
}

func (PreludeFloat) RuntimeType() RuntimeTypeRef {
	return PreludeFloatTypeRef
}

func (i PreludeFloat) String() string {
	return fmt.Sprintf("%f", i)
}
