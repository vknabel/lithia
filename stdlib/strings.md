# strings

_module_
Provides convenience functions around the basic String type.
Currently pretty limited.

- _func_ [concat](#concat) listOfStrings
- _func_ [join](#join) separator, listOfStrings

## concat

_func_ `concat listOfStrings`

Concatenates a list of given strings in order.

```
strings.concat ["Hello ", "World", "!"]
// "Hello World!"
```

## join

_func_ `join separator, listOfStrings`

Joins a list of strings with a given separator.
The separator will only be inserted between two elements.
If there are none or just one element, there won't be any separator.

```
strings.join " ", ["Hello", "World!"]
// "Hello World!"
```

