# results

_module_
The results module is all about failable operations.

- _data_ [Failure](#Failure)
- _enum_ [Result](#Result)
- _data_ [Success](#Success)
- _func_ [flatMapSuccess](#flatMapSuccess) transform, result
- _func_ [flatMapFailure](#flatMapFailure) transform, result
- _func_ [flatMapSuccess](#flatMapSuccess) transform, result
- _func_ [mapSuccess](#mapSuccess) transform, result
- _func_ [mapFailure](#mapFailure) transform, result
- _func_ [mapSuccess](#mapSuccess) transform, result
- _func_ [pureFailure](#pureFailure) error
- _func_ [pureSuccess](#pureSuccess) value

## Failure

_data_ Represents a failed result due to an error.

### Properties

- `error`

## Result

_enum_
A result of a failable operation.

```lithia
func positive { n =>
if n < 0,
Failure "negative values not supported!",
Success n
}

with positive, type Result {
Success: { success => print success.value },
Failure: { failure => print strings.concat ["failed: ", failure.error] },
}
```

### Cases

- [Success](#Success)
- [Failure](#Failure)

## Success

_data_ Represents a successful result with a value.

### Properties

- `value`

## flatMapSuccess

_func_ `flatMapSuccess transform, result`

When successful, attempts another operation by transforming the result.

## flatMapFailure

_func_ `flatMapFailure transform, result`

When failed, attempts another operation by transforming the error.

## flatMapSuccess

_func_ `flatMapSuccess transform, result`

When successful, attempts another operation by transforming the result.

## mapSuccess

_func_ `mapSuccess transform, result`

Transorms only successful results.

## mapFailure

_func_ `mapFailure transform, result`

Transorms only failed results.

## mapSuccess

_func_ `mapSuccess transform, result`

Transorms only successful results.

## pureFailure

_func_ `pureFailure error`

Creates a pure failure value.

## pureSuccess

_func_ `pureSuccess value`

Creates a pure succcess value.

