# Potfile Scripting

- Proposal: LE-002
- Status: **Implemented** in v0.0.19
- Author: [**@vknabel**](https://github.com/vknabel)

## Introduction

This proposal adds some basic subcommands to the `lithia` CLI to allow user defined scripts to be executed.

## Motivation

Every project has specific jobs to be done. These jobs are often repetitive and can be automated. This proposal aims to provide simple access to these scripts.

## Proposed Solution

Within the `Potfile`, new commands that execute a specific Lithia file as a script can be defined using the `pots.cmds.add` function. Additional help context can be provided inside a builder function. The command can be executed using `lithia [script]` or `lithia run [script]`.

```lithia
cmds.add "test", { c =>
  c.script "cmd/test.lithia"
  c.summary "runs all tests"
  c.env "LITHIA_TESTS", "1"
  c.flag "verbose", { f =>
    f.short "v"
    f.summary "verbose logging"
  }
}

cmds.add "build", { c =>
  c.bin "docker"
  c.args ["build", "-t", "my-image", "."]
  c.summary "builds the project using docker"
}
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

Under the hood, the default `bin` target will be lithia itself.

## Changes to the Standard Library

Adds a new module `pot.cmds`:

- _func_ `add name, conf` - returns a `Command`.
- _data_ `CommandBuilder`
  - _func_ `summary text` - sets the summary text for the command.
  - _func_ `flag name, conf` - adds a flag to the command.
  - _func_ `env name, value` - adds an environment variable to the command.
  - _func_ `bin name` - sets the binary to be executed.
  - _func_ `script path` - sets the script to be executed.
  - _func_ `args args` - sets the arguments to be passed to the binary.
- _data_ `FlagBuilder`
  - _func_ `short name` - sets the short name of the flag.
  - _func_ `summary text` - sets the summary text for the flag.
  - _func_ `default value` - sets the default value for the flag.
  - _func_ `required bool` - sets the flag to be required.
- _data_ `Command`
  - _let_ `name` - the name of the command.
  - _let_ `summary` - the summary text of the command.
  - _let_ `flags` - a dict of flags.
  - _let_ `envs` - a dict of environment variables.
  - _let_ `bin` - the name of the binary to be executed.
  - _let_ `args` - the arguments to be passed to the binary.
- _data_ `Flag`
  - _let_ `name` - the name of the flag.
  - _let_ `short` - the short name of the flag.
  - _let_ `summary` - the summary text of the flag.
  - _let_ `default` - the default value of the flag.
  - _let_ `required` - whether the flag is required.

Other non-final additions are internal and include:

- _let_ `store` - a stateful application store which stores the commands and in the future more commands.

## Alternatives Considered

- `pot [script]`, `lithia run [script]` or `lithia pot [script]`

But also regarding the api used in the `Potfile` itself, there was an alternative considered:

```lithia
let testCmd = cmds.script "test", "cmd/test.lithia"
testCmd.summary "runs all tests"
testCmd.flag "verbose", "verbose logging"
testCmd.env "LITHIA_TESTS", "1"
```

## Acknowledgements

`npm run [script]`, `yarn [script]`, `swift run [target]` and
`archery [script]`.
