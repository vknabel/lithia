# Changelog

## v0.0.18-next

- stdlib: `docs.docsToMarkup` now sorts extern properties
- stdlib: new function `prelude.eager` which recursively evaluates a given value
- stdlib: new function `lists.prepend` for prepending lists
- fix: improved type switch error message
- fix: improved stdlib not found error message
- fix: source locations were off
- lsp: complete importable modules #35
- lsp: field completions #35
- lsp: type switch completions #35
- lsp: improved statement completions
- lsp: outline symbols for current file
- lsp: workspace symbols for all files
- chore: bump dependencies

## v0.0.17

- compiler: alias imports `import alias = module.name` #32
- compiler: dictionary literals `[:]` and `["a": "b"]` #9
- stdlib: new type `prelude.Dict` #9
- stdlib: new type `prelude.Pair`
- stdlib: moved `results.Result` to `prelude.Result`
- stdlib: Float-support! #24
- stdlib: `!` operator #22
- stdlib: only the first leading space will be trimmed in docs
- stdlib: new type `rx.Future` and `rx.Async` and associated functions

## v0.0.16

- fix: undeclared enum case error due to file order within a module
- compiler: new concept of packages containing modules, optionally marked by a `Potfile`
- compiler: implicit `src` module for package imports if folder exists
- compiler: `import root` imports the current package (imports `src` if folder exists)
- lsp: improved autocompletion and hover information
- lsp: autocompletion and hover information across module and file boundaries
- lsp: local autocompletion #30

## v0.0.15

- stdlib: removed `docs.moduleMemberDocsToMarkup`, `docs.dataFieldDocsToMarkup` and `docs.enumCaseDocsToMarkup` and marked them as internal
- stdlib: renamed `markup.MarkupNode` to `markup.Markup`
- stdlib: renamed `markup.Serializer` to `markup.Format`
- stdlib: renamed `markup.SerializerWitness` to `markup.FormatWitness`
- stdlib: renamed `markup.serialize` to `markup.convert`
- stdlib: renamed `markdown.serializer` to `markup.format`
- stdlib: renamed `markdown.serialize` to `markup.convert`
- stdlib: removed `markdown.asMarkdown`
- stdlib: removed `optionals.Optional`. Use `prelude.Optional` instead
- stdlib: added `prelude.Maybe` as `prelude.Some`, `prelude.None` or `prelude.Any`
- stdlib: improved `optionals` functions to also support `prelude.Maybe`
- stdlib: removed `tests.testCases` and `tests.runTestCase` and marked them as internal
- fix: `docs` of `docs.ExternFunctionDocs` has always been empty
- docs: improved for all modules

## v0.0.14

- lsp: fix jump to definition
- lsp: fix multiline block comments
- lsp: fix logging too many errors
- lsp: deleting files, deletes its diagnostics
- fix: extern docs generation
- stdlib: new `markup` and `markdown` library
- stdlib: better docs generation

## v0.0.13

- lsp: semantic syntax highlighting #28
- lsp: diagnostics #29

## v0.0.12

- cli: new CLI interface, including, help and version
- fix: imports and declarations in repl were broken
- lsp: basic, proof of concept language server #19

## v0.0.11

- compiler: binaries for macOS, linux, windows

## v0.0.10

- stdlib: added `arity` to `prelude.Function`
- compiler: massive performance improvements
- compiler: large refactoring
- fix: rare equality bugs

## v0.0.9

- fix: `==` didn't work for independent values
- fix: `||` sometimes lead to wrong results
- stdlib: tests will now print all error messages
- stdlib: adds new functions to `lists`: `dropFirst`, `dropN`, `dropWhile`, `filter`, `first`, `isEmpty` and `zipWith`
- stdlib: adds new functions to `prelude`: `identity` and `const`

## v0.0.8

- fix: type expressions with `Any` were not deterministic
- fix: type expressions didn't allow comments
- fix: type expressions didn't always recognize `Module` and `Function`
- stdlib: renamed `booleans.negated` to `booleans.not`
- stdlib: fix `comparables.pullback` returned wrong type
- stdlib: added `comparables.numeric`
- stdlib: renamed `comparables` to `cmp`
- stdlib: renamed `equatables` to `eq`
- stdlib: moved `prelude.Equatable` to `eq.Equatable`
- stdlib: renamed `prelude.sameEquatable` to `eq.strict`
- stdlib: added `prelude.Never`

## v0.0.7

- docs: overall improvements
- stdlib: `strings.join` returns strings when just one element given
- stdlib: removed `prelude.reduceList`, use `lists.reduce` instead
- stdlib: moved `prelude.Result` to `results.Result`
- stdlib: moved `prelude.Functor` and `prelude.Monad` to `controls.Functor` and `controls.Monad`
- stdlib: moved `comparables.negated` and `comparables.pullback` to `equatables`
- stdlib: added `controls.Functor`, `controls.Monad` witnesses and implementations for `results.Result` and `lists.List`
- stdlib: added `comparables.Comparable`, `Order`, `Ascending`, `Equal`, `Descending`, `equatableFrom`
- stdlib: added `controls.Contravariant` and witnesses `contravariant` and `pullback` in `equatables`, `comparables`
- stdlib: added `lists.prependList`, `lists.concat`, `lists.replicate`
- stdlib: `controls` types now define `*From` functions and `*Witness` types declaring compatability regarding `Function` and `Module`
- stdlib: `controls.map`, `controls.flatMap`, `controls.pure` and `controls.pullback` convert their witness types before accessing the implementation

## v0.0.6

- stdlib: added `comparables.pullback`
- stdlib: removed `comparables.map`
- docs: improved docs for prelude
- docker: introduced [vknabel/lithia](https://hub.docker.com/repository/docker/vknabel/lithia/)

## v0.0.5

- docs: improved docs
- runtime: better print of functions
- runtime: improved runtime error messages
- compiler: type-expressions with `Any` to match multiple cases
- compiler: memberwise imports

## v0.0.4

- docs: generated docs for stdlib
- stdlib: added new modules
- stdlib: renamed `Rune` to `Char`
- improved printing of values
- stdlib: `Variable` now moved to `rx.Variable`
- stdlib: `osEnv` and `osExit` now moved to `os.env` and `os.exit`
- compiler: `extern` fails when they can't be resolved

## v0.0.3

- `module` instead of `package` to avoid confusion

## v0.0.2

- Absolute modules and multiple roots #7

## v0.0.1

- Proof of concept
- Most basic stdlib
