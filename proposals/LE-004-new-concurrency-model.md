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

The new concurrency model is based on channels and coroutines. Coroutines are similar to async functions, but they are not blocking. They can be called using the `async` function. The `async` function is followed by an expression, which is executed in a new coroutine. The coroutine can be suspended using the `await` function. The `await` keyword is followed by an expression, which is evaluated in the current coroutine. The result of the expression is returned to the coroutine, which called `await`.

```lithia
let task = async { =>
    let channel = Channel 0
    select { on, closed =>
        on channel, { value =>
            print value
        }
        on None, { _ =>
            print "No value"
        }
        closed channel, { reason =>
            print reason
        }
    }
}
await task
```

## Detailed Design

`Channel` is just a wrapper around a Go-channel. But in contrast to `close` in Go, `close` in Lithia takes a reason.

The `async` function directly creates a new Goroutine with a `Task` under the hood. It is implemented similar to `rx.Future`.

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

This is heavily inspired by Go. `Task` is inspired by Swift.
