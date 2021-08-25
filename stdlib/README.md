# Lithia Standard Library

There is one implicitly imported module for all files, called `prelude`.

## prelude

Implements the most basic data types. Espcially those needed for built-in functionality and for the compiler.

**Defined aliases:**
_None_

**Defined types:**

- Bool
  - True
  - False
- Equatable
- List
  - Cons
  - Nil
- Optional
  - Some
  - None

**Defined functions:**

- pipe
- reduceList - Really needed in prelude?
- with
- compose

**Defined constants:**
_None_

## booleans

_Depends on: prelude_

**Defined aliases:**

- Bool
- True
- False

**Defined types:**
_None_

**Defined functions:**

- negated

**Defined constants:**
_None_

## comparables

_Depends on: prelude_

**Defined aliases:**

- Equatable

**Defined types:**
_None_

**Defined functions:**

- map
- negated

**Defined constants:**
_None_

## monads

_Depends on: prelude_

**Defined aliases:**
_None_

**Defined types:**

- Functor
- Monad

**Defined functions:**

- map
- pure
- flatMap

**Defined constants:**
_None_

## optionals

_Depends on: prelude, booleans, comparables, monads_

**Defined aliases:**

- Optional

**Defined types:**
_None_

**Defined functions:**

- equalFor
- equatableFor
- map
- isNone

**Defined constants:**

- functor

## lists

**Defined aliases:**
_None_

**Defined types:**
_None_

**Defined functions:**
_None_

**Defined constants:**
_None_
