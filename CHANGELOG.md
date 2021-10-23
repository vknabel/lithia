# Changelog

## Unreleased

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
