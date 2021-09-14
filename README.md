# Lithia Language Definition

Lithia is an experimental functional programming language with an implicit but strong and dynamic type system.
Lithia is designed around a few core concepts in mind all language features contribute to.

- Composition instead inheritance
- Predictability
- Readability

## Is Lithia for you?

No. Unless you want to play around with new language concepts for some local non-production projects with a proof of concept programming language. If so, I’d be very happy to hear your feedback!

## Which features does Lithia provide?

Lithia is built around the belief, that a language is not only defined by its features, but also by the features it lacks, how it instead approaches these cases and by its ecosystem. And that every feature comes with its own tradeoffs.
As you might expect there aren’t a lot language features to cover:

- Data and enum types
- First-class modules and functions
- Currying
- Imports

On the other hand we explicitly opted out a pretty long list of features: mutability by default, interfaces, classes, inheritance, type extensions, methods, generics, custom operators, null, instance checks, importing all members of a module, exceptions and tuples.

### Data Types

are structured data with named properties. In most other languages they are called `struct`.

```
data Person {
  name
  age
}
```

As data types don’t have any methods, you simply declare global functions that act on your data.

```
func greet { person =>
  print (strings.append "Hello " person.name)
}
```

### Enum Types

in Lithia are a little bit different than you might know them from other languages. 
Some languages define enums just as a list of constant values. Others allow associated values for each named case.
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

## Why is this feature missing?

### Why no Methods?

In theory methods are just plain old functions which implicitly receive one additional parameter often called `self` or `this`.

In practice you often aren‘t able to compose methods as you can compose free functions.

Another important aspect of methods is scoping functions with their data. Here the approach is to simply create more and smaller modules. In practice we‘d create a new file for every class anyway.

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

::Overriding methods = patching witness::
::Calling super = delegate to other witness::
::Adding properties = copy::

### Why no dynamic type tests?

::(hint to modules for witness implementations, only caveat: witnesses in enums)::
