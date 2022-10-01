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

func (s PreludeString) EagerEvaluate() *RuntimeError {
	return nil
}

func (i PreludeString) Lookup(member string) (Evaluatable, *RuntimeError) {
	switch member {
	case "length":
		return NewConstantRuntimeValue(PreludeInt(len(i))), nil
	case "append":
		return NewConstantRuntimeValue(MakeAnonymousFunction(
			"append",
			[]string{"str"},
			func(args []Evaluatable) (RuntimeValue, *RuntimeError) {
				value, err := args[0].Evaluate()
				if err != nil {
					return nil, err
				}
				return PreludeString(i) + PreludeString(value.String()), nil
			})), nil
	default:
		return nil, NewRuntimeErrorf("no such member: %s", member)
	}
}
