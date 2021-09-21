package interpreter

var _ Evaluatable = ExprFunction{}

type ExprFunction struct {
	Name      string
	Arguments []string
	Docs      Docs
	ex        *EvaluationContext
	cache     *LazyEvaluationCache
}

func (ex *EvaluationContext) NewExprFunction(Name string) *ExprFunction {
	return &ExprFunction{Name: Name, ex: ex, cache: NewLazyEvaluationCache()}
}

func (e ExprFunction) Evaluate() (RuntimeValue, LocatableError) {
	return e.cache.Evaluate(func() (RuntimeValue, LocatableError) {
		return nil, nil
	})
}
