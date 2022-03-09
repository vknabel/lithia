package rx

import (
	"fmt"
	"sync"

	"github.com/vknabel/lithia/runtime"
)

var _ runtime.RuntimeValue = &RxFuture{}

type internalPromise struct {
	lock    *sync.RWMutex
	result  *promiseResult
	channel chan promiseResult
}

type promiseResult struct {
	value *runtime.RuntimeValue
	err   *runtime.RuntimeError
}

type RxFuture struct {
	futureType *RxFutureType
	storage    *internalPromise
}

func MakeRxFuture(futureType *RxFutureType, configure runtime.CallableRuntimeValue) RxFuture {
	future := RxFuture{
		futureType: futureType,
		storage: &internalPromise{
			lock:    &sync.RWMutex{},
			result:  nil,
			channel: make(chan promiseResult, 1),
		},
	}
	go func() {
		receive := runtime.MakeAnonymousFunction("receive", []string{"event"}, func(args []runtime.Evaluatable) (runtime.RuntimeValue, *runtime.RuntimeError) {
			value, err := args[0].Evaluate()
			if err != nil {
				return nil, err.CascadeDecl(futureType.DeclExternType)
			}
			resultTypeRef := runtime.MakeRuntimeTypeRef("Result", "results")
			isResult, err := resultTypeRef.HasInstance(value)
			if err != nil {
				return nil, err.CascadeDecl(futureType.DeclExternType)
			} else if !isResult {
				return nil, runtime.NewRuntimeErrorf("future %s received non-result %s", fmt.Sprint(future), fmt.Sprint(value)).CascadeDecl(futureType.DeclExternType)
			}
			future.accept(value)
			return future, nil
		})
		_, err := configure.Call([]runtime.Evaluatable{
			runtime.NewConstantRuntimeValue(receive),
		}, nil)
		if err != nil {
			future.fail(*err)
		}
	}()
	return future
}

func (RxFuture) RuntimeType() runtime.RuntimeTypeRef {
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

func (v RxFuture) Lookup(member string) (runtime.Evaluatable, *runtime.RuntimeError) {
	switch member {
	case "await":
		return runtime.NewLazyRuntimeValue(func() (runtime.RuntimeValue, *runtime.RuntimeError) {
			return v.Await()
		}), nil
	default:
		return nil, runtime.NewRuntimeErrorf("future %s has no member %s", fmt.Sprint(v), member)
	}
}

func (v RxFuture) Await() (runtime.RuntimeValue, *runtime.RuntimeError) {
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

func (v RxFuture) accept(value runtime.RuntimeValue) (runtime.RuntimeValue, *runtime.RuntimeError) {
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

func (v RxFuture) fail(err runtime.RuntimeError) (runtime.RuntimeValue, *runtime.RuntimeError) {
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
