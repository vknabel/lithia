# results

_module_ 

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

_data_ 

### Properties

- `error`

## Result

_enum_ 

### Cases

- [Failure](#Failure)
- [Success](#Success)

## Success

_data_ 

### Properties

- `value`



## flatMapSuccess

_func_ `flatMapSuccess transform, result`


## flatMapFailure

_func_ `flatMapFailure transform, result`


## flatMapSuccess

_func_ `flatMapSuccess transform, result`



## mapSuccess

_func_ `mapSuccess transform, result`


## mapFailure

_func_ `mapFailure transform, result`


## mapSuccess

_func_ `mapSuccess transform, result`



## pureFailure

_func_ `pureFailure error`


## pureSuccess

_func_ `pureSuccess value`



