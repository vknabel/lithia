package runtime

import "sync"

var _ Evaluatable = &LazyEvaluatable{}

type LazyEvaluatable struct {
	once  *sync.Once
	value RuntimeValue
	err   *RuntimeError
	eval  func() (RuntimeValue, *RuntimeError)
}

func NewLazyRuntimeValue(eval func() (RuntimeValue, *RuntimeError)) *LazyEvaluatable {
	return &LazyEvaluatable{
		once:  &sync.Once{},
		eval:  eval,
		value: nil,
	}
}

func NewConstantRuntimeValue(value RuntimeValue) *LazyEvaluatable {
	return &LazyEvaluatable{
		once:  &sync.Once{},
		eval:  func() (RuntimeValue, *RuntimeError) { return value, nil },
		value: value,
	}
}

func (l *LazyEvaluatable) Evaluate() (RuntimeValue, *RuntimeError) {
	l.once.Do(func() {
		l.value, l.err = l.eval()
	})
	return l.value, l.err
}
