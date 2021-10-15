# strings

_module_ 

## concat

_func_ Concatenates a list of given strings in order.

```
strings.concat ["Hello ", "World", "!"]
// "Hello World!"
```

### Parameters

- listOfStrings

## join

_func_ Joins a list of strings with a given separator.
The separator will only be inserted between two elements.
If there are none or just one element, there won't be any separator.

```
strings.join " ", ["Hello", "World!"]
// "Hello World!"
```

### Parameters

- separator
- listOfStrings
