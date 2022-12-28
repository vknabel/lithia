# New Concurrency Model

- Proposal: LE-004
- Status: **Draft**
- Author: [**@vknabel**](https://github.com/vknabel)

## Introduction

The initially planned concurrency model for Lithia was based on Functional Reactive Programming. In practice this turned out to not fit the language very well. FRP is complex, yet Lithia tries to be simple. This proposal describes a new concurrency model, which is inspired by Go's channels and coroutines.

## Motivation

Lithia has a strong and dynamic type system. The current approach of FRP wraps all values and therefore forcing users to unwrap them. This forces developers to double check the docs for the return type. Otherwise this leads to avoidable runtime errors.

When comparing this with the Go concurrency model, where async function calls are blocking, but still return the same value as normal function calls. This should lead to a better developer experience in Lithia programs.

## Proposed Solution

The new concurrency model is based on channels and coroutines. Coroutines are similar to async functions, but they are not blocking. They can be called using the `async` function. The `async` function is followed by an expression, which is executed in background.

```lithia
// create an unbuffered channel
let channel = Channel 0

// start a new coroutine
async { =>
    // send a value to the channel
    // as unbuffered, send is blocking until received
    channel.send "Hello World"
}

// start selecting on the channel
// blocking until fulfilled
select { on, closed =>
    // on new value in channel
    on channel, { value =>
        print value
    }
    // channel is closed
    closed channel, { reason =>
        print reason
    }
    // default case
    on None, { _ =>
        print "No value"
    }
}
```

When spinning up a new async task, a `Task` is returned. This can be used to `await` for the task to finish.
But in contrast to other languages, long running tasks are not async by default and returning Tasks is considered an anti-pattern. Invocations of `async` and `await` are designed to take place on the call site.

```lithia
// start all jobs in background
let tasks = lists.map jobs, { job => async run job }
// wait for all tasks to finish
lists.map tasks, await
```

The goal of a `Task` is to avoid local channels to wait for multiple results, requested in parallel. They should not escape your function scope.

## Detailed Design

The new concurrency model is based around three basic building blocks:

1. Channels - represent a communication channel between different parts of the program
2. Select - allows to wait for certain events to take place
3. Async Tasks - represent a small unit of work, which can be executed in background

### Channels

Channels are used to communicate between different parts of the program. They are created using the `Channel` type. The `Channel` type is a wrapper around a Go-channel. It can be created with a buffer size. If the buffer size is 0, the channel is unbuffered. Otherwise it is buffered.

- Sending a value to a channel is blocking if the channel is unbuffered or if the buffer is full until there is a receiver.
- Receiving a value from a channel is blocking if the channel is unbuffered or if the buffer is empty until there is a sender.
- Closing a channel not blocking.
- Waiting for a channel to be closed is blocking until it has been blocked.
- Closing a channel requires a reason.
- Sending a value to a closed channel is a runtime error.
- Closing an already closed channel is a runtime error. Even if the reason is the same.

### Select

The `select` function allows to handle the fastest case in a set of channels or Async Tasks. It is blocking until one of the cases is fulfilled. The `select` function takes a function as argument. This function is called with two functions `on` and `closed`. The `on` function is used to register a case for a channel. The `closed` function is used to register a case for a closed channel.

If `on` is called with `None` or `Nil`, it is called if no other case is fulfilled immediately. If it has been omitted, the `select` function will block until one of the cases is fulfilled.

### Async Tasks

Async Tasks are used to execute a small unit of work in background. They are created using the `async` function. The `async` function takes an expression as argument. This expression is executed in background. The `async` function returns a `Task` which can be used to `await` for the task to finish.

When awaiting a task, the current coroutine is blocked until the task has finished. The result of the task is returned. Tasks cannot fail.

## Changes to the Standard Library

The `rx` module will be removed including the types `rx.Async`.
`rx.Variable` will be moved to `prelude.Variable`. `rx.Future` will be moved to `prelude.Future`.

In module `prelude`:

- _extern_ `Channel`
  - can be called as `Channel bufferSize`
  - _func_ `send value`
  - _func_ `close reason`
- _func_ `select func`
- _extern_ `async expr`
- _func_ `await future`
- _extern_ `Variable` defined as before
- _extern_ `Task` defined as `rx.Future` before, but with all references to `results.Result` removed. They can't fail anymore.

## Alternatives Considered

Most alternatives are more complex and might impact developer experience.
The `Task` could be dropped, but it might be useful.

## Acknowledgements

`Channel`, `select`, but also `async` are heavily inspired by Go. But in contrast to Go, they require a reason to close a channel.

`Task`, `async` and `await` are inspired by many implementations like in TypeScript or Swift. Especially decoupling errors from async operations like in Swift really fit into Lithia's design.
