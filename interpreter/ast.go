package interpreter

import "sync"

type Evaluatable interface {
	Evaluate() (RuntimeValue, LocatableError)
}

type LazyEvaluationCache struct {
	once  *sync.Once
	value RuntimeValue
	err   LocatableError
}

func NewLazyEvaluationCache() *LazyEvaluationCache {
	return &LazyEvaluationCache{once: &sync.Once{}}
}

func (c *LazyEvaluationCache) Evaluate(f func() (RuntimeValue, LocatableError)) (RuntimeValue, LocatableError) {
	c.once.Do(func() {
		c.value, c.err = f()
	})
	return c.value, c.err
}

type ExprMemberAccess struct {
	Expr Evaluatable
}

type ExprInvocation struct {
	Function  Evaluatable
	Arguments []Evaluatable
}

type Statement interface{}
