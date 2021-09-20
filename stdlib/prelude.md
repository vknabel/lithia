# prelude

_module_ Implements the most basic data types.
Espcially those needed for built-in functionality and for the compiler.
Will always be imported implicitly.

## Any

_extern_ 

## Bool

_enum_ 

### Cases

- [True](#True)
- [False](#False)

## Char

_extern_ 

## Cons

_data_ Represents a non-empty List.

### Properties

- head - The first element
- tail - The remaining list.
@type List

## Equatable

_data_ 

### Properties

- equal lhs, rhs

## Failure

_data_ 

### Properties

- error

## False

_data_ 
## Float

_extern_ 

## Function

_extern_ 

## Functor

_data_ 

### Properties

- map f, value

## Int

_extern_ 

## List

_enum_ A list of arbiratry elements.

```
import lists

let myList = [1, 2, 3, 4]
lists.reduce { l, r => l + r }, 0, myList
```

### Cases

- [Cons](#Cons)
- [Nil](#Nil)

## Module

_extern_ 

## Monad

_data_ 

### Properties

- pure value
- flatMap f, instance

## Nil

_data_ Marks the end of the list.
## None

_data_ 
## Optional

_enum_ 

### Cases

- [Some](#Some)
- [None](#None)

## Result

_enum_ 

### Cases

- [Success](#Success)
- [Failure](#Failure)

## Some

_data_ 

### Properties

- value

## String

_extern_ 

## Success

_data_ 

### Properties

- value

## True

_data_ 
## Void

_data_ 
## compose

_func_ 

### Parameters

- f
- g
- value

## debug

_func_ 

### Parameters

- message

## flatMap

_func_ 

### Parameters

- f
- witness
- instance

## if

_func_ 

### Parameters

- condition
- then
- else

## map

_func_ 

### Parameters

- f
- witness
- value

## pipe

_func_ Pipes a given value through a list of functions.
The first function is applied to the value, the second to the result of the first, etc.

@param functions
@param initial

### Parameters

- functions
- initial

## print

_func_ 

### Parameters

- message

## pure

_func_ 

### Parameters

- value
- witness

## reduceList

_func_ Recursively walk a tree of nodes, calling a function on each node.
The given accumulator function merges each element into a new one for the next call.

@param accumulator into, next
@param initial
@param list

### Parameters

- accumulator
- initial


## unless

_func_ 

### Parameters

- condition
- then

## when

_func_ 

### Parameters

- condition
- then

## with

_func_ Applies the given body to the given value.
Mostly useful for readability, e.g. in destructings.

```
with True, Bool(True: { _ => }, False: { _ => })
```

@param value
@param body

### Parameters

- value
- body
