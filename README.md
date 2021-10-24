# Lithia Programming Language

Lithia is an experimental functional programming language with an implicit but strong and dynamic type system.
Lithia is designed around a few core concepts in mind all language features contribute to.

- Composition instead inheritance
- Predictability
- Readability

## Is Lithia for you?

No. Unless you want to play around with new language concepts for local non-production projects with a proof of concept programming language. I’d be happy to hear your feedback!

### Roadmap

Currently Lithia is an early proof of concept. Basic language features exist, but the current tooling and standard libraries are far from being feature complete or stable.

- [x] Module imports
- [x] Testing library
- [x] Easy installation
- [x] Prebuilt docker image
- [ ] Prebuilt linux binaries
- [ ] Docs generator _in progress_
- [ ] Improved performance
- [ ] Stack traces
- [ ] Creating a custom language server
- [ ] ... with diagnostics
- [ ] ... with syntax highlighting
- [ ] ... with auto completion
- [ ] A package manager
- [ ] Move stdlib to a package
- [ ] Custom plugins for external declarations
- [ ] More static type safety

Not all features end up on the list above. Espcially improving the standard libraries and documentation is an ongoing process.

Lithia will hit **0.1.0** once there is basic tooling support, all language features are implemented and the Standard Library is mature. From there on, every source breaking change will result in a minor release. To reach **1.0.0**, a small ecosystem needs to evolve.

## Installation

If you want to give Lithia a try, the easiest way to get started is using Homebrew. By default Lithia is now ready to go.

```bash
$ brew install vknabel/lithia/lithia
```

To get syntax highlighting, download and install the latest version of [Syntax Highlighter with Lithia](https://github.com/vknabel/syntax-highlighter/releases) for VS Code.

## Which features does Lithia provide?

Lithia is built around the belief, that a language is not only defined by its features, but also by the features it lacks, how it instead approaches these cases and by its ecosystem. And that every feature comes with its own tradeoffs.
As you might expect there aren’t a lot language features to cover:

- Data and enum types
- First-class modules and functions
- Currying
- Module imports

On the other hand we explicitly opted out a pretty long list of features: mutability by default, interfaces, classes, inheritance, type extensions, methods, generics, custom operators, null, instance checks, importing all members of a module, exceptions and tuples.

> Curios? Head over to the generated [Standard Library documentation](./stdlib/README.md).

### Functions

Lithia supports currying and lazy evaluation: functions can be called parameter by parameter. If all parameters have been provided and the resulting expression will be used, the functions itself will be called.

To reflect this behavior, functions are called braceless. Every parameter is separated by a comma.

```
func add { l, r => l + r }

add 1, 2 // 3

// with currying
let incr = add 1 // { r => 1 + r }
incr 2 // 3
```

As parameters of a function call are comma separated, you can compose single arguments.
All operators bind stronger than parameters and function calls.

```
when True, print "will be printed"

// here you can see lazy evaluation in action:
// print will never be executed.
when False, print "won't be printed"

// parens required
when (incr 1) == 2, print "will be printed"

// when needed, single parameter calls can be nested
// fun (a (b (c d)
fun a b c d
// fun (a b), (c d)
fun a b, c d
```

### Data Types

are structured data with named properties. In most other languages they are called `struct`.

```
data Person {
  name
  age
}
```

As data types don’t have any methods, you declare global functions that act on your data.

```
func greet { person =>
  print (strings.append "Hello ", person.name)
}
```

### Enum Types

in Lithia are different than you might know them from other languages.
Some languages define enums as a list of constant values. Others allow associated values for each named case.
Though in Lithia, an enum is an enumeration of types.

To make it easier to use for the value enumeration use case, there is a special syntax to directly declare an enum case and the associated type.

```
enum JuristicPerson {
  Person
  data Company {
    name
    corporateForm
  }
}
```

Instead of a classic switch-case statement, there is a `type`-expression instead.
It requires you to list all types of the enum type. It returns a function which takes a valid enum type.

```
import strings

let nameOf = type JuristicPerson {
  Person: { person => person.name },
  Company: { company =>
    strings.concat [
      company.name, " ", company.corporateForm
    ]
  }
}

nameOf you
```

If you are interested in special cases, you can use the `Any` case.

> _**Attention:** If the given value is not valid, your programm will crash. If you might have arbitrary values, you can add an `Any` case. As it matches all values, make sure it is always the last value._

### Modules

are defined by the folder structure. Once there is a folder with Lithia files in it, you can import it. No additional configuration overhead required.

```
import strings

strings.join " ", []
```

Or alternatively, import members directly. But use this sparingly: it might lead to name collisions.

```
import strings {
  join
}

join " ", []
```

### Current module

Sometimes you might want to pass the whole module as parameter, or to avoid naming collisions.

```
module current

let map = functor.map

doSomeStuff current
```

As shown, a common use case is to pass the module itself instead of multiple witnesses, if all members are also defined on the module itself.

### Module resolution

Lithia will search for a folder containing source files at the following locations:

- when executing a file, relative to it
- when in REPL, inside the current working directory
- at `$LITHIA_LOCALS` if set
- at `$LITHIA_PACKAGES` if set
- at `$LITHIA_STDLIB` or `/usr/local/opt/lithia/stdlib`

> _**Nice to know:** the special module `prelude` will always be imported implicitly. It contains types like `Bool` or `List`. Beyond that Lithia treats the `prelude` as any other module. And you can even override and update the standard library._

Modules and their members can be treated like any other value. Just pass them around as parameters.

## Why is this feature missing?

### Why no Methods?

In theory methods are plain old functions which implicitly receive one additional parameter called `self` or `this`.

In practice you aren‘t able to compose methods as you can compose free functions.

Another important aspect of methods is scoping functions with their data. Here the approach is to create more and smaller modules. In practice we‘d create a new file for every class anyway.

```
data Account {
  balance
}

func withdraw { debit, account =>
  Account account.balance - debit
}
```

```
import accounts {
  Account
}

accounts.withdraw 500, Account 1000

with Account 150, pipe [
  accounts.withdraw 100,
  accounts.withdraw 50,
]
```

### Why no Interfaces?

Interfaces only allow one single implementation per type. The only way to make the implementing types composable is to define more types, requiring more ceremony than plain old functions.

Instead of an interface you simply create a new data type, assign your implementation and pass it alongside to your argument. The instance containing the implementation is called a witness.

```
data Greetable {
  greeting ofValue
}

let shortPersonGreetable = Greetable { person =>
  strings.append "Hi " person.name
}


func greet { greetable, object =>
  print greetable.greeting object
}

greet shortPersonGreetable, someone

```

The benefit of using witnesses instead of plain interface lies in flexibility as one can define multiple implementations of the protocol.

```
let longPersonGreetable = Greetable { person =>
  strings.prepend "Hello " person.name
}
```

And when it comes to composition, this approach really shines! That way we can define our own map or similar functions to transform existing witnesses.

```
import strings

func map { transform, witness =>
  Greetable { object =>
    transform witness.greeting object
  }
}

let uppercased = map strings.uppercased

let screamed = map strings.append "!"
```

As seen above, we can easily rely on existing implementations, compose them and always receive the same data types until we have built complete algorithms!

### Why no class inheritance?

Classes and inheritance have their use cases and benefits, but as Lithia separates data from behavior, inheritance doesn’t serve us well anymore.

For data we have two options:

1.  Copying all members to another data. Though enums must also include this new data type.
2.  Nesting the data. Especially useful if the data is only used outside the default context. This is especially great if you need to combine many different witnesses or data types as with multi-inheritance.

```
data Base { value }

// copying
data CopiedBase {
  value
  other
}

// nesting
data NestedBase {
  base
  other
}
```

When regarding witnesses, we have a third option: modules. We create our witnesses and bind each data value directly to the module itself.

```
module strings

let map = functor.map
let flatMap = monad.flatMap

doSomething strings, ""
```

The `controls` module explicitly embraces the use of modules as valid witnesses.
Alongside `Functor` it defines a function `functorFrom` constructing a `Functor` from an enum `FunctorWitness` which allows modules, functors and even monads.

```
enum FunctorWitness {
    Functor
    Module
    Function
    Monad
}

func functorFrom { moduleWitness =>
    with moduleWitness, type FunctorWitness {
        Functor: { witness => witness },
        Module: { module =>
            Functor module.map
        },
        Function: { fmap =>
            Functor fmap
        },
        Monad: { monad =>
            Functor { f, instance => monad.pure (monad.flatMap f, instance) }
        }
    }
}
```

Though the defaults should be used wisely: for example the `Result` type has two different valid implementations of `map`! On the other hand `List` only has one valid implementation.

One additional feature of class inheritance is calling functionality of the super class. In Lithia the approach looks different, but behaves similar:
We create a whole new witness, which calls the initial one under the hood.

### Why no dynamic type tests?

Most languages allow type casts and checks. Lithia does only support the type switch expression for enums.

These checks are unstructured and therefore tempt to be used in the wrong places. Though type checks should be used sparingly. Lithia prefers to move required decisions to the edge of the code. Witnesses should implement decisions for the provided data and desired behavior.

If there is one type to focus on, the tooling and the developer can understand all cases much easier and faster.

## License

Lithia is available under the [MIT](./LICENSE) license.
