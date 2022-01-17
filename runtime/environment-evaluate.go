package runtime

func (env *Environment) GetEvaluatedRuntimeValue(key string) (RuntimeValue, *RuntimeError) {
	if lazyValue, ok := env.GetPrivate(key); ok {
		value, err := lazyValue.Evaluate()
		if err != nil {
			return nil, err
		}
		return value, nil
	} else {
		return nil, NewRuntimeErrorf("undefined %s", key)
	}
}

func (env *Environment) GetExportedEvaluatedRuntimeValue(key string) (RuntimeValue, *RuntimeError) {
	if lazyValue, ok := env.GetPrivate(key); ok {
		value, err := lazyValue.Evaluate()
		if err != nil {
			return nil, err
		}
		return value, nil
	} else {
		return nil, NewRuntimeErrorf("undefined %s", key)
	}
}
