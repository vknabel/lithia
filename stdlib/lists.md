# lists

_module_ 

- _func_ [append](#append) element
- _func_ [concat](#concat) nestedLists
- _func_ [count](#count) list
- _func_ [dropFirst](#dropFirst) list
- _func_ [dropN](#dropN) n, list
- _func_ [dropWhile](#dropWhile) predicate, list
- _func_ [filter](#filter) predicate, list
- _func_ [first](#first) list
- _func_ [flatMap](#flatMap) transform, list
- _func_ [foldr](#foldr) accumulator, initial
- _func_ [forEach](#forEach) list, action

- _func_ [isEmpty](#isEmpty) list
- _func_ [map](#map) transform, list

- _func_ [prependList](#prependList) prefix, postfix
- _func_ [pure](#pure) value
- _func_ [reduce](#reduce) accumulator, initial
- _func_ [replicate](#replicate) n, element
- _func_ [zipWith](#zipWith) combine, left, right

## append

_func_ `append element`

Appends an element to the end of a list.
## concat

_func_ `concat nestedLists`


## count

_func_ `count list`

Counts all elements of a list.
## dropFirst

_func_ `dropFirst list`

Returns the list without the first element.
## dropN

_func_ `dropN n, list`

Returns the list without the first `n` elements.
Negative numbers will be treated as `0`.
## dropWhile

_func_ `dropWhile predicate, list`

Returns the list without the first elements where the given predicate is True.
## filter

_func_ `filter predicate, list`

Returns the list only with elements where the given predicate is True.
## first

_func_ `first list`

Gets the first element of the list as `Optional`.
When empty `Nil`, otherwise `Some`.
## flatMap

_func_ `flatMap transform, list`


## foldr

_func_ `foldr accumulator, initial`

Starts folding at the end of the list, from the right.

`accumulator nextValue, transformedTail`
## forEach

_func_ `forEach list, action`

Iterates over all values of a list to generate side effects.
Returns last result or `prelude.Void`.

## isEmpty

_func_ `isEmpty list`


## map

_func_ `map transform, list`



## prependList

_func_ `prependList prefix, postfix`


## pure

_func_ `pure value`


## reduce

_func_ `reduce accumulator, initial`

Recursively walk a tree of nodes, calling a function on each node.
The given accumulator function merges each element into a new one for the next call.

```
lists.reduce { into, next => into + next.length }, 0, ["count", "chars"]
```
## replicate

_func_ `replicate n, element`

Creates a list with a count of `n`. Every item will be the given element.
Negative numbers will be treated as `0`.
## zipWith

_func_ `zipWith combine, left, right`

Combines two lists pairwise. Has the length of the shortest list.
`combine left, right`.