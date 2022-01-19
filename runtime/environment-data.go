package runtime

func (env *Environment) MakeEmptyDataRuntimeValue(name string) (DataRuntimeValue, *RuntimeError) {
	return env.MakeDataRuntimeValue(name, make(map[string]Evaluatable))
}

func (env *Environment) MakeDataRuntimeValue(name string, members map[string]Evaluatable) (DataRuntimeValue, *RuntimeError) {
	evaluatable, ok := env.GetPrivate(name)
	if !ok {
		return DataRuntimeValue{}, NewRuntimeErrorf("not declared: %s", name)
	}
	typeValue, err := evaluatable.Evaluate()
	if err != nil {
		return DataRuntimeValue{}, err
	}
	dataDecl, ok := typeValue.(PreludeDataDecl)
	if !ok {
		return DataRuntimeValue{}, NewRuntimeErrorf("not declared: %s", name)
	}
	return MakeDataRuntimeValueMemberwise(&dataDecl, members)
}

func (env *Environment) MakeList(slice []Evaluatable) (DataRuntimeValue, *RuntimeError) {
	if len(slice) == 0 {
		return env.MakeEmptyDataRuntimeValue("Nil")
	} else {
		return env.MakeDataRuntimeValue("Cons", map[string]Evaluatable{
			"head": slice[0],
			"tail": NewLazyRuntimeValue(func() (RuntimeValue, *RuntimeError) {
				return env.MakeList(slice[1:])
			}),
		})
	}
}
