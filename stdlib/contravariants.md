# contravariants

_module_ 

## Contravariant

_data_ A contravariant wraps behavior to handle inputs depending on a context.
Put simply, a contravariant maps inputs, while a functor maps outputs.

```
import equatables

let personByNameEquatbale = equatables.contravariant.pullback { person => person.name }, sameEquatable

personByNameEquatbale.equal Person "Alice", Person "Bob"
// > False
```

Invariants:
1. Identity: `(pullback { a => a }, value) == value`
2. Associative: `(pipe [pullback f, pullback g], value) == pullback pipe [f, g], value`

### Properties

- `pullback transform, value`

## ContravariantWitness

_enum_ Defines all valid witnesses for a contravariant.

```
import comparables

pullback { person => person.name }, comparables
pullback { person => person.name }, comparables.pullback
pullback { person => person.name }, comparables.contravariant
```

### Cases

- [Contravariant](#Contravariant)
- [Function](#Function)
- [Module](#Module)

## from

_func_ `from moduleWitness`


## pullback

_func_ `pullback f, witness, value`

