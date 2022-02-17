# fs

_module_

- _extern_ [delete](#delete) atPath
- _extern_ [deleteAll](#deleteAll) atPath
- _extern_ [exists](#exists) atPath
- _extern_ [readString](#readString) fromPath
- _extern_ [writeString](#writeString) toPath, newContents

## delete

_func_ `delete atPath`

Deletes a file or empty directory at a given path.
Returns a `result.Result`.

## deleteAll

_func_ `deleteAll atPath`

Deletes a file or any directory at a given path.
Returns a `result.Result`.

## exists

_func_ `exists atPath`

Returns a Bool if a file or directory exists at a given path.

## readString

_func_ `readString fromPath`

Reads a file at a given path.
Returns `result.Result`.

## writeString

_func_ `writeString toPath, newContents`

Write a given string to a path.
Returns `results.Result`.

