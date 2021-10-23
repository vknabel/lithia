# Changelog

## Unreleased

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
