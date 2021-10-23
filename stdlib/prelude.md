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

- [False](#False)
- [True](#True)

## Char

_extern_ 

## Cons

_data_ Represents a non-empty List.

### Properties

- `head` - The first element
- `tail` - The remaining list.
@type List

## Equatable

_data_ Allows comparision of values for equality.
Declare and pass a witness for custom equality.

In contrast to the default equality operator ==, you can define custom equality.
If you explicitly want the strict behavior, pick the `sameEquatable` witness.

### Properties

- `equal lhs, rhs`

## False

_data_ A constant to represent invalid conditions.
## Float

_extern_ 

## Function

_extern_ 

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

## Nil

_data_ Marks the end of the list.
## None

_data_ 
## Optional

_enum_ 

### Cases

- [None](#None)
- [Some](#Some)

## Some

_data_ 

### Properties

- `value`

## String

_extern_ 

## True

_data_ A constant to represent valid conditions.
## Void

_data_ 
## compose

_func_ `compose f, g, value`


## debug

_func_ `debug message`


## if

_func_ `if condition, then, else`

When the given condition evaluates to `True`, returns `then`. Otherwise `false`.
Both, `then` and `else` are evaluted lazily.

```
if True, print "Succeeded", exit 1
```
## pipe

_func_ `pipe functions, initial`

Pipes a given value through a list of functions.
The first function is applied to the value, the second to the result of the first, etc.
## print

_func_ `print message`



## unless

_func_ `unless condition, then`


## when

_func_ `when condition, then`


## with

_func_ `with value, body`

Applies the given body to the given value.
Mostly useful for readability, e.g. in destructings.

```
with True, Bool(True: { _ => }, False: { _ => })
```