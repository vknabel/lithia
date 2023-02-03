# Contributing

First of all: **thank you** for your interest in contributing to this project and open source in general!

At first, contributing to Lithia might sound intimidating. _After all, it's a programming language, right?_
But the bigger a topic is, the more aspects there are to consider and the more possibilities to contribute there are.

I'd really love to hear your feedback, ideas and suggestions. If you have any questions, feel free to ask them in the [Discussions](https://github.com/vknabel/lithia/discussions) or [Issues](https://github.com/vknabel/lithia/issues) section. You are able to shape this project.

**Table of Contents**

- [Contributing](#contributing)
  - [Code of Conduct](#code-of-conduct)
  - [Ways to Contribute](#ways-to-contribute)
    - [Using Lithia](#using-lithia)
    - [Reporting Bugs](#reporting-bugs)
    - [Driving Evolution](#driving-evolution)
    - [Writing Documentation](#writing-documentation)
    - [Writing Code or Tests](#writing-code-or-tests)
    - [Reviewing Pull Requests](#reviewing-pull-requests)
    - [Writing Libraries](#writing-libraries)
  - [Architecture](#architecture)
    - [Development](#development)
    - [Board](#board)

## Ways to Contribute

You can contribute to Lithia in many ways:

- [Using Lithia](#using-lithia)
- [Reporting Bugs](#reporting-bugs)
- [Driving Evolution](#driving-evolution)
- [Writing Documentation](#writing-documentation)
- [Writing Code or Tests](#writing-code-or-tests)
- [Reviewing Pull Requests](#reviewing-pull-requests)
- [Writing Libraries](#writing-libraries)

### Using Lithia

The best way to get comfortable with new tech is actually using it. And Lithia is new for everybody. Help to find bugs and patterns and idioms. What should be best practice? What should be avoided? What is confusing? What is missing? All this is the driver for innovation. If you are interested in this topic, the [Discussions](https://github.com/vknabel/lithia/discussions) section is your friend.

### Reporting Bugs

As every project and especially young projects, Lithia has bugs. But as Lithia is a programming language, it should be tested very thoroughly. If you find a bug, please report it in the [Issues](https://github.com/vknabel/lithia/issues/new) section. Please make sure to include a minimal example that reproduces the bug.

## Driving Evolution

Lithia is still in its early stages and it's still evolving and tries to find its way of getting things done.
You don't have to start with a whole evolution proposal and a PR or even the implementation. Opening an [issue](https://github.com/vknabel/lithia/issues/new) or [discussion](https://github.com/vknabel/lithia/discussions) gets things going. You can also help to improve [existing proposals](https://github.com/vknabel/lithia/tree/main/proposals) or [drafts](https://github.com/vknabel/lithia/issues?q=label%3Aproposal).

### Writing Documentation

Currently the documentation is rough and currently only consists of the [README](README.md), the [proposals](./proposals/README.md), the [examples](./examples/) and generated [stdlib docs](./stdlib/README.md).
What's missing are tutorials, guides and a reference. But also a website would be nice. If you are interested in this topic, open an [issue](https://github.com/vknabel/lithia/issues/new), [discussion](https://github.com/vknabel/lithia/discussions) or directly a pull request with your changes, e.g. if you found ~tpyos~ typos.

### Writing Code or Tests

Lithia is written in Go, uses Tree-Sitter for parsing, has an VSCode extension in Typescript and an embedded language server in Go. You could contribute to any of these topics and more. Details about the architecture can be found in the [Architecture](#architecture) section.

Lithia is tested with [Go's testing framework](https://golang.org/pkg/testing/) and [Tree-Sitter's testing framework](https://tree-sitter.github.io/tree-sitter/creating-parsers#testing-parsers). We use GitHub Actions to run the tests on every push and pull request. The [stdlib itself is tested](./stdlib/stdlib-tests.lithia) in Lithia.

### Reviewing Pull Requests

Nobody is perfect and code isn't either. Feel free to review pull requests and give feedback. If you don't understand specific changes, just ask. The issue might be the code or the documentation.

### Writing Libraries

If you wrote some code with Lithia, we'd like to link to it! Open a pull request with your changes to the [README](README.md) and we'll merge it. If you want to have some feedback and code review, head to the [Discussions](https://github.com/vknabel/lithia/discussions) and we will have a look at it.

## Architecture

Lithia consists of multiple projects:

- The grammar is written in [Tree-Sitter](https://tree-sitter.github.io/tree-sitter/) and lives in its own repository: [vknabel/tree-sitter-lithia](https://github.com/vknabel/tree-sitter-lithia).
- The [VSCode extension](https://github.com/vknabel/vscode-lithia) is written in Typescript and lives in its own repository.
- The core of Lithia is written in Go and lives in this repository.
- The language server protocol implementation, which also lives in this repository.

### Development

We try to keep as much as possible within this repo and easy to install.
This repository is structured as follows:

- The [parser](./parser/) is written in Go produces an Abstract Syntax Tree ([ast](./ast/)) from the Tree-Sitter output.
- The interpreter is implemented in [runtime](./runtime/) and executes the AST.
- Extensions to Lithia that are not part of the runtime core can be found in [external](./external/). They resemble the names of the standard library modules.
- Most types and code is implemented in the [stdlib](./stdlib/). It is written in Lithia itself and is tested with the [stdlib-tests](./stdlib/stdlib-tests.lithia).
- The language server is implemented in [langsrv](./langsrv/). It is the intelligent backend for the [VSCode extension](https://github.com/vknabel/vscode-lithia).

The easiest way to get started is using the devcontainer configured in this repository. It will install all dependencies and start a VSCode instance with the extension installed. You can then run the tests and debug the extension. The devcontainer is configured to run Lithia using `go run ./app/lithia`.
So if you restart the vscode extension, window, devcontainer or the extension host, it will automatically recompile the lithia binary for you and you will get the latest language server implementation.

### Board

If you plan to contribute more frequently, we have a [GitHub project board](https://github.com/users/vknabel/projects/3/views/1) that contains all the issues and pull requests that are currently open.

The development of Lithia follows a brief concept and some rules. Don't take them too serious. We are pragmatic and do this as a hobby. They are mostly for myself.

**Sprints**

- provide a rough timespan
- usually take 4 weeks
- are not set in stone and can be changed
- should result in at least one release

**Issues**

- are a `bug`, `docs`, `feature` or a `proposal`
- features are assigned to an evolution proposal
- are assigned to a sprint
- have a status
  - **Backlog:** a new issue or an issue that could be tackled in theory.
  - **Upcoming:** this task is planned in the near future. Usually within this Sprint.
  - **In Progress:** the work for this task has already been started, but it has not been finished yet.
  - **Done:** every closed issue.

**Lithia Evolution proposals**

- drive the feature development of Lithia
- are referenced by LE-xxx if they left the draft phase
- serve as implementation guidance and documentation
- might start with an `proposal` issue or a discussion in the forum
- are open for everyone
