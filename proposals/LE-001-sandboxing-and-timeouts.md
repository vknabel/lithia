# Sandboxing and Timeouts

- Proposal: LE-001
- Status: **Draft**
- Author: [**@vknabel**](https://github.com/vknabel)

## Introduction

In this proposal we want to discuss the possibility of adding sandboxing and timeouts to the Lithia command line interface. This will be accomplished by implementing a whitelist for external declarations on a module level.

To avoid infinitely running scripts, it should also be possible to set a timeout for the execution of a script.

## Motivation

In security critical or embedded environments, it is important to be able to limit the power of a program.
The first practical use case for this will be the evaluation of the `Potfile`. It will be executed whenever opening the project and before actually starting the program itself.

The Potfile execution needs to be fast and secure.

## Proposed Solution

The Lithia CLI will get additional configuration options to allow limiting external definitions and setting timeouts.

```bash
LITHIA_TIMEOUT=1s \
LITHIA_EXTERNAL_DEFINITIONS=docs,fs,os,rx \
lithia run [script]
```

## Detailed Design

`LITHIA_TIMEOUT` will be used to set a timeout for the execution of a script. The default value will be `0s` which means no timeout.

`LITHIA_EXTERNAL_DEFINITIONS` will be used to set a whitelist for external definitions. The default value will be `*` which means all external definitions are allowed. External definitions of `prelude` are always allowed.

The `Potfile` will always be executed with `LITHIA_TIMEOUT=1s` and `LITHIA_EXTERNAL_DEFINITIONS=rx`. No matter what the user has set.

### Merging the Environment Variables

If the user has manually set the environment variables more strict than the default values for the `Potfile` execution, the more strict values will be used.

For example if the user has set `LITHIA_EXTERNAL_DEFINITIONS=docs,os,rx` and the `Potfile` contains `LITHIA_EXTERNAL_DEFINITIONS=docs,fs,os,rx`, the final value will be `docs,os,rx`.

## Changes to the Standard Library

No changes to the standard library are required.

## Alternatives Considered

In the future command line flags or configuration files might be used to set these values.

## Acknowledgements

List people who have contributed to the design of this proposal. Also mention any prior art, such as how other languages have solved this problem.
