# eq

_module_ 

- _data_ [Equatable](#Equatable)

- _func_ [negated](#negated) witness
- _func_ [pullback](#pullback) transform, witness


## Equatable

_data_ Allows comparision of values for equality.
Declare and pass a witness for custom equality.

In contrast to the default equality operator ==, you can define custom equality.
If you explicitly want the strict behavior, pick the `strict` witness.

### Properties

- `equal lhs, rhs`


## negated

_func_ `negated witness`

Negates the result of the given `Equatable`.
## pullback

_func_ `pullback transform, witness`

Transforms the inputs of an `Equatable`-witness.

```
cmp.pullback { person => person.name }, insensitiveEquatable, Person "Somebody"
```
