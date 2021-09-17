package interpreter

func (ex *EvaluationContext) SyntaxErrorf(format string, args ...interface{}) LocatableError {
	return ex.LocatableErrorf("syntax error", format, args...)
}
