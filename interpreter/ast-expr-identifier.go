package interpreter

var _ Evaluatable = ExprIdentifier{}

type ExprIdentifier struct {
	Name  string
	ex    *EvaluationContext
	cache *LazyEvaluationCache
}

func (ex *EvaluationContext) NewExprIdentifier(Name string) *ExprIdentifier {
	return &ExprIdentifier{Name: Name, ex: ex, cache: NewLazyEvaluationCache()}
}

func (e ExprIdentifier) Evaluate() (RuntimeValue, LocatableError) {
	return e.cache.Evaluate(func() (RuntimeValue, LocatableError) {
		if lazyValue, ok := e.ex.environment.Get(e.Name); ok {
			value, err := lazyValue.Evaluate()
			if err != nil {
				return nil, err
			}
			switch value := value.(type) {
			case DataDeclRuntimeValue:
				if len(value.fields) == 0 {
					return DataRuntimeValue{
						typeValue: &value,
						members:   make(map[string]Evaluatable),
					}, nil
				} else {
					return value, nil
				}
			case Function:
				if len(value.arguments) == 0 {
					result, err := value.Call(nil)
					return result, e.ex.LocatableErrorOrConvert(err)
				} else {
					return value, nil
				}
			default:
				return value, nil
			}
		} else {
			return nil, e.ex.RuntimeErrorf("undefined identifier %s", e.Name)
		}
	})
}
