# controls

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
1. Identity: `(map { a => a }, value) == value`
2. Associative: `(pipe [map f, map g], value) == map pipe [f, g], value`

### Properties

- `map f, value` - Transforms a wrapped value using a function depending context of the functor

## FunctorWitness

_enum_ Defines all valid witnesses for a functor.

```
import lists

map incr, lists
map incr, lists.map
map incr, lists.functor
```

### Cases

- [Function](#Function)
- [Functor](#Functor)
- [Module](#Module)
- [Monad](#Monad)

## Monad

_data_ Monads apply a function returning wrapped values to a wrapped value.

Invariants:
1. Left-Identity: `(pipe [pure, flatMap f], value) == f value`
2. Right-Identity: `(pipe [pure, flatMap { x => x }], value) == pure value`
3. Associative: `(pipe [pure, flatMap f, flatMap g], value) == pipe [pure, flatMap g, flatMap f], value`

### Properties

- `pure value` - Wraps a value in a neutral context.
- `flatMap f, instance` - Transforms a wrapped value and merges potential partial results.

## MonadWitness

_enum_ Valid witnesses for a monad.

```
import lists

flatMap repeat 2, lists
flatMap repeat 2, lists.monad
```

### Cases

- [Module](#Module)
- [Monad](#Monad)

## contravariantFrom

_func_ `contravariantFrom moduleWitness`


## flatMap

_func_ `flatMap f, witness, instance`


## functorFrom

_func_ `functorFrom moduleWitness`


## map

_func_ `map f, witness, value`

Transforms a wrapped value using a functor witness.
Essentially just uses the map of the given witness,
but allows to defer the decision regarding the witness itself.

```
import lists

let incr = { i => i + 1 }
map incr, lists, [1, 2, 3]
```
## monadFrom

_func_ `monadFrom monadWitness`


## pullback

_func_ `pullback f, witness, value`


## pure

_func_ `pure value, witness`

