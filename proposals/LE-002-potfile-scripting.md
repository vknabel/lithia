# Potfile Scripting

- Proposal: LE-002
- Status: **Draft**
- Author: [**@vknabel**](https://github.com/vknabel)

## Introduction

This proposal adds some basic subcommands to the `lithia` CLI to allow user defined scripts to be executed.

## Motivation

Every project has specific jobs to be done. These jobs are often repetitive and can be automated. This proposal aims to provide simple access to these scripts.

## Proposed Solution

Within the `Potfile`, new commands that execute a specific Lithia file as a script can be defined using the `pots.cmds.script` function. Additional help context can be provided on the existing command. The command can be executed using `lithia [script]` or `lithia run [script]`.

```lithia
let testCmd = cmds.script "test", "cmd/test.lithia"
testCmd.summary "runs all tests"
testCmd.flag "verbose", "verbose logging"
testCmd.env "LITHIA_TESTS", "1"
```

Similar, the `pots.cmds.bin` function can be used to define a command that executes a binary. Under the hood `/usr/bin/env` will be used to lookup the binary. In case of `lithia` itself, the currently executed `lithia` binary will be used.

```lithia
cmds.bin "example", ["lithia", "test"]
```

## Detailed Design

The user adds the commands of their choice to their `Potfile`.
The Lithia interpreter will now execute it and will proceed working on the internal store structure.

The execution of the Potfile takes place in a sandboxed environment which whitelists possible external module imports to avoid direct access to the file system and other resources. It will also be time-restricted to at most one second using Go `context.Context`.

To avoid privilege escalation, scripts declared in the `Potfile` with overridden `LITHIA_TIMEOUT` or `LITHIA_EXTERNAL_DEFINITIONS`, can only restrict the environment further. They cannot allow more access to the environment.

## Changes to the Standard Library

Adds a new module `pot.cmds`:

- _func_ `script alias, path` - returns a `MutableCommand`.
- _func_ `bin alias, path` - returns a `MutableCommand`.
- _data_ `MutableCommand` - a mutable command that can be used to configure the command.
  - _func_ `summary text` - sets the summary text for the command.
  - _func_ `flag name, description` - adds a flag to the command.
  - _func_ `env name, value` - adds an environment variable to the command.

Other non-final additions are internal and include:

- _let_ `store` - a stateful application store which stores the commands and in the future more commands.

## Alternatives Considered

- `pot [script]`, `lithia run [script]` or `lithia pot [script]`

## Acknowledgements

`npm run [script]`, `yarn [script]`, `swift run [target]` and
`archery [script]`.
