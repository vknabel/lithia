# <img src="./assets/lithia.png" width="50"> Lithia Programming Language

Lithia is an experimental functional programming language with an implicit but strong and dynamic type system.
Lithia is designed around a few core concepts in mind all language features contribute to.

- Composition instead inheritance
- Predictability
- Readability

## Is Lithia for you?

No, unless you want to play around with new language concepts for local non-production projects with a proof of concept programming language. I’d be happy to hear your feedback!

### Roadmap

Currently Lithia is an early proof of concept. Basic language features exist, but the current tooling and standard libraries are far from being feature complete or stable.

- [x] Module imports
- [x] Testing library
- [x] Easy installation
- [x] Prebuilt docker image
- [x] Prebuilt linux binaries
- [x] Generated docs for stdlib
- [x] Improved performance [#20](https://github.com/vknabel/lithia/pull/20)
- [x] Stack traces [#20](https://github.com/vknabel/lithia/pull/20)
- [x] Creating a custom language server
- [x] ... with diagnostics
- [x] ... with syntax highlighting
- [x] ... with auto completion _basic_
- [x] ... with symbols
- [ ] ... with refactorings
- [ ] ... with formatter
- [ ] A package manager
- [ ] A debugger
- [ ] Custom plugins for external declarations
- [ ] More static type safety

Not all features end up on the list above. Improving the standard libraries and documentation is an ongoing process.

To hit version **0.1.0**, Lithia needs [all planned language features](https://github.com/vknabel/lithia/milestone/2), a [rich standard library](https://github.com/vknabel/lithia/milestone/1) and at least a [basic language server](https://github.com/vknabel/lithia/milestone/3) implementation.
To reach **1.0.0**, we need a stable standard library, documentation, solid and broad tooling.

### Breaking Changes

Until we reach **0.1.0** every update is considered breaking.
Upcoming **0.x.Patch**-updates may fix bugs and add features. Existing Lithia source code will not break, but extensions may.
**0.Minor.0** releases are breaking updates.

#### What is considered a breaking change?

- renaming, moving or removing declarations
- adding cases to *enum*s, that do not contain `Any`
- renaming, adding or removing fields to _data_
- changing the order of fields to _data_

#### What is not considered a breaking change?

- renaming function parameters
- moving parameters through definitions using currying
- importing new modules

## Installation

If you want to give Lithia a try, the easiest way to get started is using Homebrew. By default Lithia is now ready to go.

```bash
$ brew install vknabel/lithia/lithia
```

To get syntax highlighting, use the [Lithia for VS Code extension](https://marketplace.visualstudio.com/items?itemName=vknabel.vscode-lithia).

> Not using Visual Studio Code? Get started with `lithia lsp --help`.

### asdf

If you are using the [asdf version manager](https://asdf-vm.com), there is a [lithia plugin](https://github.com/vknabel/asdf-lithia)!

```bash
# Install the plugin
$ asdf plugin add lithia https://github.com/vknabel/asdf-lithia.git

# Show all installable versions
asdf list-all lithia

# Install specific version
asdf install lithia latest

# Set a version globally (on your ~/.tool-versions file)
asdf global lithia latest

# Now lithia commands are available
lithia --version
```

### Docker

To give Lithia a try, you can use our docker container to start the REPL:

```bash
$ docker run --rm -it vknabel/lithia
> print "Hello World"
Hello World
- Hello World
>
```

To deploy your own application built with Lithia, create your own Dockerfile.

```docker
FROM vknabel/lithia:latest

WORKDIR /app
ENV LITHIA_PACKAGES /app/packages
COPY ./packages /app/packages
ENV LITHIA_LOCALS /app/src
COPY ./src /app/src
COPY ./main.lithia /app/main.lithia

RUN lithia main.lithia
```

## Which features does Lithia provide?

Lithia is built around the belief, that a language is not only defined by its features, but by the features it lacks, how it instead approaches these cases and by its ecosystem. Every feature comes with its own tradeoffs.
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
Other languages define enums as a list of constant values. A few allow associated values for each named case.
A Lithia _enum_ is an enumeration of types.

There is syntactic sugar for value enumerations, to directly declare a case and the associated _enum_ or _data_ type.

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

> _**Attention:** If the given value is not valid, your program will crash. If you might have arbitrary values, you can add an `Any` case. As it matches all values, make sure it is always the last value._

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

To create your own package of modules, you can create a `Potfile`. Every module defined by the package is a directory next to the Potfile, with `.lithia` files in it.

Though there are a few special cases:

- the `cmd`-folder is typically used for individual files rather than a module. You execute them with `$ lithia cmd/<file>`.
- the `src`-folder represents the `root` of the package.
- if the `src`-folder is missing, the package `root` is next to the Potfile.

```
.
├── Potfile
├── cmd
│   ├── main.lithia
│   └── test.lithia
└── src
    └── greet.lithia
```

If there aren't any matching local modules, Lithia will search for a package containing source files at the following locations:

- when in REPL, inside the current working directory
- at `$LITHIA_LOCALS` if set
- at `$LITHIA_PACKAGES` if set
- at `$LITHIA_STDLIB` or `/usr/local/opt/lithia/stdlib`

> _**Nice to know:** the special module `prelude` will always be imported implicitly. It contains types like `Bool` or `List`. Beyond that Lithia treats the `prelude` as any other module. And you can even override and update the standard library._

Modules and their members can be treated like any other value. Pass them around as parameters.

## Why is this feature missing?

### Why no Methods?

In theory methods are plain old functions which implicitly receive one additional parameter called `self` or `this`.

In practice you aren‘t able to compose methods as you can compose free functions.

Another important aspect of methods is scoping functions with their data. Here the approach is to create more and smaller modules. In practice we create separate files for every class.

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

Interfaces allow one implementation per type. The only way to make the implementing types composable is to define more types, requiring more ceremony than plain old functions.

Instead of an interface you create a new _data_ type, assign your implementation and pass it alongside to your argument. The instance containing the implementation is called a witness.

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

As seen above, we can rely on existing implementations, compose them and always receive the same data types until we have built complete algorithms!

### Why no class inheritance?

Classes and inheritance have their use cases and benefits, but as Lithia separates data from behavior, inheritance doesn’t serve us anymore.

For data we have two options:

1.  Copying all members to another _data_. *enum*s must include this new data type.
2.  Nesting the data. Useful if the data is used outside the default context and is great if you need to combine many different witnesses or data types as you would with multi-inheritance.

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

Though the defaults should be used wisely: for example the `Result` type has two different valid implementations of `map`! On the other hand `List` has one valid implementation.

One additional feature of class inheritance is calling functionality of the super class. In Lithia the approach looks different, but behaves similar:
We create a whole new witness, which calls the initial one under the hood.

### Why no dynamic type tests?

Most languages allow type casts and checks. Lithia does only support the type switch expression for enums.

These checks are unstructured and therefore tempt to be used in the wrong places. Though type checks should be used sparingly. Lithia prefers to move required decisions to the edge of the code. Witnesses should implement decisions for the provided data and desired behavior.

If there is one type to focus on, the developer and the tooling can understand all cases much easier and faster.

## License

Lithia is available under the [MIT](./LICENSE) license.
