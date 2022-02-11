# prelude

_module_
Implements the most basic data types.
Espcially those needed for built-in functionality and for the compiler.
Will always be imported implicitly.

- _extern_ [Any](#Any)
- _enum_ [Bool](#Bool)
- _extern_ [Char](#Char)
- _data_ [Cons](#Cons)
- _data_ [False](#False)
- _extern_ [Float](#Float)
- _extern_ [Function](#Function)
- _extern_ [Int](#Int)
- _enum_ [List](#List)
- _extern_ [Module](#Module)
- _enum_ [Never](#Never)
- _data_ [Nil](#Nil)
- _data_ [None](#None)
- _enum_ [Optional](#Optional)
- _data_ [Some](#Some)
- _extern_ [String](#String)
- _data_ [True](#True)
- _data_ [Void](#Void)
- _func_ [compose](#compose) f, g, value
- _func_ [const](#const) value, _
- _extern_ [debug](#debug) message
- _func_ [identity](#identity) value
- _func_ [if](#if) condition, then, else
- _func_ [pipe](#pipe) functions, initial
- _extern_ [print](#print) message
- _func_ [unless](#unless) condition, then
- _func_ [when](#when) condition, then
- _func_ [with](#with) value, body

## Any

_extern_

## Bool

_enum_
Represents boolean values like `True` and `False`.
Typically used for conditionals and flags.

### Cases

- [True](#True)
- [False](#False)

## Char

_extern_

## Cons

_data_ Represents a non-empty List.

### Properties

- `head` - The first element
- `tail` - The remaining list.
@type List

## False

_data_ A constant to represent invalid conditions.

## Float

_extern_

## Function

_extern_

### Properties

- `arity`

## Int

_extern_

## List

_enum_
A list of arbiratry elements.

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

## Never

_enum_
An enum with no valid values.
Allows empty, but valid type expressions.

## Nil

_data_ Marks the end of the list.

## None

_data_

## Optional

_enum_
An optional value. Either some value or none.

### Cases

- [Some](#Some)
- [None](#None)

## Some

_data_

### Properties

- `value`

## String

_extern_

### Properties

- `length`
- `append str`

## True

_data_ A constant to represent valid conditions.

## Void

_data_ Represents a single value.

## compose

_func_ `compose f, g, value`

Composes two given functions.
Calls the second function first and pipes the result into the second one.

## const

_func_ `const value, _`

Always returns the first argument.

## debug


_extern_ `debug message`

## identity

_func_ `identity value`

Always returns the given argument.

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


_extern_ `print message`

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

