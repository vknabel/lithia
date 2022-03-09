package runtime

import (
	"fmt"
	"sync"
)

var _ RuntimeValue = &RxFuture{}

type internalPromise struct {
	lock    *sync.RWMutex
	result  *promiseResult
	channel chan promiseResult
}

type promiseResult struct {
	value *RuntimeValue
	err   *RuntimeError
}

type RxFuture struct {
	futureType *RxFutureType
	storage    *internalPromise
}

func MakeRxFuture(futureType *RxFutureType, configure CallableRuntimeValue) RxFuture {
	future := RxFuture{
		futureType: futureType,
		storage: &internalPromise{
			lock:    &sync.RWMutex{},
			result:  nil,
			channel: make(chan promiseResult, 1),
		},
	}
	go func() {
		receive := MakeAnonymousFunction("receive", []string{"event"}, func(args []Evaluatable) (RuntimeValue, *RuntimeError) {
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err.CascadeDecl(futureType.DeclExternType)
			}
			resultTypeRef := MakeRuntimeTypeRef("Result", "results")
			isResult, err := resultTypeRef.HasInstance(value)
			if err != nil {
				return nil, err.CascadeDecl(futureType.DeclExternType)
			} else if !isResult {
				return nil, NewRuntimeErrorf("future %s received non-result %s", fmt.Sprint(future), fmt.Sprint(value)).CascadeDecl(futureType.DeclExternType)
			}
			future.accept(value)
			return future, nil
		})
		_, err := configure.Call([]Evaluatable{
			NewConstantRuntimeValue(receive),
		}, nil)
		if err != nil {
			future.fail(*err)
		}
	}()
	return future
}

func (RxFuture) RuntimeType() RuntimeTypeRef {
	return RxVariableTypeRef
}

func (v RxFuture) String() string {
	v.storage.lock.RLock()
	defer v.storage.lock.RUnlock()
	if v.storage.result.err != nil {
		return fmt.Sprintf("(%s %s)", v.RuntimeType().Name, *v.storage.result.err)
	} else if v.storage.result.value != nil {
		return fmt.Sprintf("(%s %s)", v.RuntimeType().Name, *v.storage.result.value)
	} else {
		return "Future"
	}
}

func (v RxFuture) Lookup(member string) (Evaluatable, *RuntimeError) {
	switch member {
	case "await":
		return NewLazyRuntimeValue(func() (RuntimeValue, *RuntimeError) {
			return v.Await()
		}), nil
	default:
		return nil, NewRuntimeErrorf("future %s has no member %s", fmt.Sprint(v), member)
	}
}

func (v RxFuture) Await() (RuntimeValue, *RuntimeError) {
	v.storage.lock.RLock()

	if v.storage.result != nil {
		defer v.storage.lock.RUnlock()
		if v.storage.result.err != nil {
			return nil, v.storage.result.err
		} else {
			return *v.storage.result.value, nil
		}
	} else {
		v.storage.lock.RUnlock()
		result := <-v.storage.channel
		v.storage.lock.Lock()
		defer v.storage.lock.Unlock()

		v.storage.result = &result
		if result.err != nil {
			return nil, result.err
		} else {
			return *result.value, nil
		}
	}
}

func (v RxFuture) accept(value RuntimeValue) (RuntimeValue, *RuntimeError) {
	v.storage.lock.RLock()
	if v.storage.result != nil {
		defer v.storage.lock.RUnlock()
		if v.storage.result.err != nil {
			return nil, v.storage.result.err
		} else {
			return *v.storage.result.value, nil
		}
	} else {
		v.storage.lock.RUnlock()
		v.storage.lock.Lock()

		v.storage.result = &promiseResult{value: &value}
		v.storage.lock.Unlock()
		v.storage.channel <- promiseResult{value: &value}
		close(v.storage.channel)
		return value, nil
	}
}

func (v RxFuture) fail(err RuntimeError) (RuntimeValue, *RuntimeError) {
	v.storage.lock.Lock()
	defer v.storage.lock.Unlock()
	if v.storage.result != nil {
		if v.storage.result.err != nil {
			return nil, v.storage.result.err
		} else {
			return *v.storage.result.value, nil
		}
	} else {
		v.storage.result = &promiseResult{err: &err}
		v.storage.lock.Unlock()

		v.storage.channel <- promiseResult{err: &err}
		close(v.storage.channel)
		return nil, &err
	}
}
