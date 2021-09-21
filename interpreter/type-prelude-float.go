package interpreter

import "fmt"

var _ RuntimeValue = PreludeFloat(0)

type PreludeFloat float64

func (f PreludeFloat) String() string {
	return fmt.Sprintf("%f", f)
}

func (PreludeFloat) RuntimeType() RuntimeType {
	return RuntimeType{
		name:       "Float",
		moduleName: "prelude",
	}
}

func (f PreludeFloat) Lookup(member string) (Evaluatable, error) {
	return nil, fmt.Errorf("float %f has no member %s", f, member)
}
