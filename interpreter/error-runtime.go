package interpreter

func (ex *EvaluationContext) RuntimeErrorf(format string, args ...interface{}) LocatableError {
	return ex.LocatableErrorf("runtime error", format, args...)
}
