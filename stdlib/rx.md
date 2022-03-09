# rx

_module_
A very early concept of implementing functional reactive programming.
Currently only used to provide mutability.

- _enum_ [Async](#Async)
- _extern_ [Future](#Future)
- _extern_ [Variable](#Variable)
- _func_ [await](#await) async
- _func_ [catch](#catch) transform
- _func_ [flatMap](#flatMap) transform
- _func_ [flatten](#flatten) async
- _func_ [map](#map) transform
- _func_ [onSuccess](#onSuccess) sideEffect

## Async

_enum_

### Cases

- [Result](#Result)
- [Future](#Future)
- [Variable](#Variable)

## Future

_extern_
Represents a value calculated in background.
It will arrive some time in the future.

```lithia
import results
import rx

let future = rx.Future { receive =>
    // will be performed in background
    receive results.Success 42
}

// the .await will block and wait for the result
with future.await, type results.Result {
    Success: { value => value },
    Failure: { err =>
        print err
        0 // as default
    },
}
```

### Properties

- `await` - Waits for the future to complete.
This will lock the current function until the result has arrived.
At the end, returns the `results.Result`.

## Variable

_extern_
Holds a value and enables replacing it.
Planned to propagate value changes to observers, but not implemented, yet.

### Properties

- `accept value` - Changes the currently hold value of the variable.
- `current` - Returns the currently hold value.

## await

_func_ `await async`

## catch

_func_ `catch transform`

## flatMap

_func_ `flatMap transform`

## flatten

_func_ `flatten async`

## map

_func_ `map transform`

## onSuccess

_func_ `onSuccess sideEffect`

