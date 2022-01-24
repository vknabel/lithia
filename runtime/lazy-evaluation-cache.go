package runtime

import "sync"

type LazyEvaluationCache struct {
	once  *sync.Once
	value RuntimeValue
	err   *RuntimeError
}

func NewLazyEvaluationCache() *LazyEvaluationCache {
	return &LazyEvaluationCache{once: &sync.Once{}}
}

func (c *LazyEvaluationCache) Evaluate(f func() (RuntimeValue, *RuntimeError)) (RuntimeValue, *RuntimeError) {
	c.once.Do(func() {
		c.value, c.err = f()
	})
	return c.value, c.err
}
