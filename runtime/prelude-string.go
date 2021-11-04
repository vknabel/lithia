package runtime

var _ RuntimeValue = PreludeString("")
var PreludeStringTypeRef = MakeRuntimeTypeRef("String", "prelude")

type PreludeString string

func (PreludeString) RuntimeType() RuntimeTypeRef {
	return PreludeStringTypeRef
}

func (s PreludeString) String() string {
	return string(s)
}

func (i PreludeString) Lookup(member string) (Evaluatable, *RuntimeError) {
	panic("TODO: length, append, chars")
}
