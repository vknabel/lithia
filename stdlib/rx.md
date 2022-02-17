# rx

_module_
A very early concept of implementing functional reactive programming.
Currently only used to provide mutability.

- _extern_ [Variable](#Variable)

## Variable

_extern_
Holds a value and enables replacing it.
Planned to propagate value changes to observers, but not implemented, yet.

### Properties

- `accept value` - Changes the currently hold value of the variable.
- `current` - Returns the currently hold value.

