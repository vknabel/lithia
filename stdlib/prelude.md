# prelude

_module_ Implements the most basic data types.
Espcially those needed for built-in functionality and for the compiler.
Will always be imported implicitly.

## Any

_extern_ 

## Bool

_enum_ Represents boolean values like `True` and `False`.
Typically used for conditionals and flags.

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

_data_ Allows comparision of values for equality.
Declare and pass a witness for custom equality.

In contrast to the default equality operator ==, you can define custom equality.
If you explicitly want the strict behavior, pick the `sameEquatable` witness.

### Properties

- equal lhs, rhs

## Failure

_data_ 

### Properties

- error

## False

_data_ A constant to represent invalid conditions.
## Float

_extern_ 

## Function

_extern_ 

## Functor

_data_ A functor wraps values in a context and allows different decisions depending on the context.
For example, the types `Optional` and `List` have functors.

```
import lists
import optionals

let incr = { i => i + 1 }
lists.functor.map incr, [1, 2, 3]
// > [2, 3, 4]
optionals.functor.map incr, Some 41
// > Some 42
optionals.functor.map incr, None
// > None
```

Invariants:
1. Identity: `(map { a => a}, value) == value`
2. Associative: `(pipe [map f, map g], value) == map pipe [f, g], value`

### Properties

- map f, value - Transforms a wrapped value using a function depending context of the functor

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

- [Nil](#Nil)
- [Cons](#Cons)

## Module

_extern_ 

## Monad

_data_ Left-Identity: `(pipe [pure, flatMap f], value) == f value`
Right-Identity: `(pipe [pure, flatMap { x => x }], value) == pure value`
Associative: `(pipe [pure, flatMap f, flatMap g], value) == pipe [pure, flatMap g, flatMap f], value`

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

_data_ A constant to represent valid conditions.
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

_func_ When the given condition evaluates to `True`, returns `then`. Otherwise `false`.
Both, `then` and `else` are evaluted lazily.

```
if True, print "Succeeded", exit 1
```

### Parameters

- condition
- then
- else

## map

_func_ Transforms a wrapped value using a functor witness.
Essentially just uses the map of the given witness,
but allows to defer the decision regarding the witness itself.

```
import lists

let incr = { i => i + 1 }
map incr, lists.functor, [1, 2, 3]
```

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
