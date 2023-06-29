# apps

_module_
A simple, uni-directional data flow architecture.
All global functions inside this module can be used on all statefuls.

## Example

We start by defining a model. A model is a record with an initial state and an update function.
The update function takes a state and a message and returns an `Update` record.

In this example we have a simple counter model.

```lithia
import apps

enum Message {
    data Inc
    data Dec
}

let model = apps.defineModel 0, { state, msg =>
    with msg, type apps.Message {
        Inc: { _ => apps.Update state + 1, Nil },
        Dec: { _ => apps.Update state - 1, Nil },
    }
}

let view = apps.render model, { dispatch, state =>
    print strings.join ["Current state: ", state]
}

import tests

tests.test "apps", { fail =>
    apps.send view, Inc

    unless (apps.currentState view) == 1, fail "should be 1"
}
```

- _data_ [Batch](#Batch)
- _enum_ [Command](#Command)
- _enum_ [Message](#Message)
- _data_ [Model](#Model)
- _data_ [Quit](#Quit)
- _enum_ [Stateful](#Stateful)
- _data_ [Store](#Store)
- _data_ [Update](#Update)
- _data_ [View](#View)
- _func_ [currentState](#currentState) stateful
- _func_ [defineModel](#defineModel) initial, update
- _func_ [render](#render) model, view
- _func_ [send](#send) stateful, cmd
- _func_ [storeFrom](#storeFrom) model
- _func_ [subscribe](#subscribe) observer, stateful
- _func_ [update](#update) stateful, msg

## Batch

_data_ A batch of Commands

### Properties

- `cmds`

## Command

_enum_
A command can be either a message, a batch of messages or an async message.

### Cases

- [Async](#Async)
- [Batch](#Batch)
- [Message](#Message)

## Message

_enum_
The messages any stateful can handle.
`prelude.Nil` and `prelude.None` are used to indicate that no message should be send.
`prelude.Any` is used to indicate that any message should be send.

### Cases

- [Quit](#Quit)
- [Nil](#Nil)
- [None](#None)
- [Any](#Any)

## Model

_data_ A model defines the abstract and underlaying behaviour of a stateful.
It has an initial state and an update function.

### Properties

- `initial` - The initial state
- `update state, msg` - The update function.
Takes a state and a message and returns an `Update` record.

## Quit

_data_ Quits the whole program

## Stateful

_enum_
Represents all statefuls like Stores and Views.

### Cases

- [Store](#Store)
- [View](#View)

## Store

_data_ A stateful which holds a state and can be observed.

### Properties

- `model` - The Model
- `states` - A rx.Variable of a State
- `observers` - A rx.Variable of functions to be notified.

## Update

_data_ The update record includes the new state after applying a message and a command to be send to the store.

### Properties

- `state` - The new state
- `cmd` - A command to be send to the store.

## View

_data_ A stateful which renders a view.
The result of the view differs from use case to use case.

Whenever you have one representation of your whole state, the `View` is the right choice.

### Properties

- `store` - The underlying and observed store.
- `view dispatch, state` - A function which renders the current state.
Might `dispatch cmd`.

## currentState

_func_ `currentState stateful`

Returns the current state of a stateful.

## defineModel

_func_ `defineModel initial, update`

Creates a new model.

## render

_func_ `render model, view`

Creates a new view by rendering a store on every change.

## send

_func_ `send stateful, cmd`

Sends an eventually async Command to a stateful.
Always returns a Result.

## storeFrom

_func_ `storeFrom model`

Creates a new store from a model.
This essentially adds dynamic behaviour to a model.

## subscribe

_func_ `subscribe observer, stateful`

Observes changes of a stateful.

## update

_func_ `update stateful, msg`

Directly applies a message to a stateful.

