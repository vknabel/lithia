package runtime

import (
	"context"
	"sync"
)

type LazyEvaluationCache struct {
	ctx   context.Context
	once  *sync.Once
	value RuntimeValue
	err   *RuntimeError
}

func NewLazyEvaluationCache(ctx context.Context) *LazyEvaluationCache {
	return &LazyEvaluationCache{ctx: ctx, once: &sync.Once{}}
}

func (c *LazyEvaluationCache) Evaluate(f func() (RuntimeValue, *RuntimeError)) (RuntimeValue, *RuntimeError) {
	c.once.Do(func() {
		done := make(chan struct{})
		go func() {
			c.value, c.err = f()
			close(done)
		}()
		select {
		case <-c.ctx.Done():
			c.err = NewRuntimeError(c.ctx.Err())
		case <-done:
		}
	})
	return c.value, c.err
}
