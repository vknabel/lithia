# lists

_module_ 

- _func_ [append](#append) element
- _func_ [concat](#concat) nestedLists
- _func_ [count](#count) list
- _func_ [flatMap](#flatMap) transform, list
- _func_ [foldr](#foldr) accumulator, initial
- _func_ [forEach](#forEach) list, action

- _func_ [map](#map) transform, list

- _func_ [prependList](#prependList) prefix, postfix
- _func_ [pure](#pure) value
- _func_ [reduce](#reduce) accumulator, initial
- _func_ [replicate](#replicate) n, element

## append

_func_ `append element`

Appends an element to the end of a list.
## concat

_func_ `concat nestedLists`


## count

_func_ `count list`

Counts all elements of a list.
## flatMap

_func_ `flatMap transform, list`


## foldr

_func_ `foldr accumulator, initial`

Starts folding at the end of the list, from the right.

`accumulator nextValue, transformedTail`
## forEach

_func_ `forEach list, action`

Iterates over all values of a list to generate side effects.

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

