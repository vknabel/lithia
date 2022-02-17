# optionals

_module_

- _func_ [equalFor](#equalFor) someWitness, lhs, rhs
- _func_ [equatableFor](#equatableFor) someWitness
- _func_ [from](#from) maybe
- _func_ [isNone](#isNone)
- _func_ [map](#map) transform
- _func_ [orDefault](#orDefault) default

## equalFor

_func_ `equalFor someWitness, lhs, rhs`

Creates an equal function, that understands optionals and maybes for a given witness.

## equatableFor

_func_ `equatableFor someWitness`

Creates an Equatable witness for Optionals on top of an existing witness.

## from

_func_ `from maybe`

Creates an optional from a Maybe-value.

## isNone

_func_ `isNone`

True if None. Otherwise False.

## map

_func_ `map transform`

Transforms Some value to a new one.
Keeps None as-is.
Any other values will still be mapped, but not wrapped.

## orDefault

_func_ `orDefault default`

Returns a default, if None given.
Otherwise unwraps Some value or keeps Any as-is.

