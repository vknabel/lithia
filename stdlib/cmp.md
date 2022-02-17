# cmp

_module_
Defines comparision operations, ascending and descending of values.

- _data_ [Ascending](#Ascending)
- _data_ [Comparable](#Comparable)
- _data_ [Descending](#Descending)
- _data_ [Equal](#Equal)
- _enum_ [Order](#Order)
- _func_ [equatableFrom](#equatableFrom) comparableWitness
- _func_ [pullback](#pullback) f, witness

## Ascending

_data_ Indicates an ascending order of two values.
For example 1 and 2.

## Comparable

_data_ Instances compare values regarding the order.
Witnesses are typically only defined for specific types.

### Properties

- `compare lhs, rhs` - Compares two values.
@returns Order

## Descending

_data_ Indicates an descending order of two values.
For example 2 and 1.

## Equal

_data_ Both values are ordered equally.
In context of Order, it doesn't necessarily require equality.

## Order

_enum_
Represents the order of two values.

### Cases

- [Ascending](#Ascending)
- [Equal](#Equal)
- [Descending](#Descending)

## equatableFrom

_func_ `equatableFrom comparableWitness`

Creates an `eq.Equatable` from a `cmp.Comparable`.
`cmp.Equal` will result in `True`,
`cmp.Ascending` and `cmp.Descending` will be `False`.

@returns eq.Equatable

## pullback

_func_ `pullback f, witness`

Lifts an existing `cmp.Comparable` witness to a different type.
Can be used to pick a specific property of complex data.

```lithia
let compareUsersById = cmp.pullback { user => user.id }, cmp.numeric
```

