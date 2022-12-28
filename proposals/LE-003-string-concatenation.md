# String Concatenation

- Proposal: LE-003
- Status: **Draft**
- Author: [**@vknabel**](https://github.com/vknabel)

## Introduction

Adds the ability to concatenate strings using the `..` operator.

## Motivation

Strings need to be concatnated often. The current way to do this is to use the `strings.concat` function, which feels a bit verbose in many cases.

A simple `..` might improve the readability of the code.

## Proposed Solution

We add a new left-associative operator `..` to the language, which is used to concatenate strings with a low precendence.

## Detailed Design

The `..` operator is left-associative, and has a precedence of 5.
This allows comparisions of concatenated strings.

```lithia
"Hello" .. " " .. "World" == "Hello World"
```

## Changes to the Standard Library

`strings.concat` will still be available, and will never be deprecated.

Implementations of string related functions should use `..` internally.

Deprecation of `String.append`.

## Alternatives Considered

We could overload `+`. But due to the dynamic typing it may be confusing in practice.
`++` is also used in different languages, but we do not want to mix up the meaning of `+` and `++`.

## Acknowledgements

Mimics the `..` operator in Lua.
