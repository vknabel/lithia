package runtime

func (env *Environment) MakeEmptyDataRuntimeValue(name string) (DataRuntimeValue, *RuntimeError) {
	evaluatable, ok := env.GetPrivte(name)
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
	return MakeDataRuntimeValueMemberwise(&dataDecl, make(map[string]Evaluatable))
}
