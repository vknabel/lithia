package interpreter

var _ Evaluatable = ExprBinaryOp{}

type ExprBinaryOp struct {
	OpName string
	Left   Evaluatable
	Right  Evaluatable
	OpImpl func(Evaluatable, Evaluatable) (RuntimeValue, LocatableError)

	ex    *EvaluationContext
	cache *LazyEvaluationCache
}

func (ex *EvaluationContext) NewExprBinaryOp(
	OpName string,
	OpImpl func(Evaluatable, Evaluatable) (RuntimeValue, LocatableError),
	Left Evaluatable,
	Right Evaluatable,
) *ExprBinaryOp {
	return &ExprBinaryOp{
		OpName: OpName,
		Left:   Left,
		Right:  Right,
		OpImpl: OpImpl,
		ex:     ex,
		cache:  NewLazyEvaluationCache(),
	}
}

func (e ExprBinaryOp) Evaluate() (RuntimeValue, LocatableError) {
	return e.cache.Evaluate(func() (RuntimeValue, LocatableError) {
		return e.OpImpl(e.Left, e.Right)
	})
}
