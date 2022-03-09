# docs

_module_

- _data_ [ConstantDocs](#ConstantDocs)
- _data_ [DataDocs](#DataDocs)
- _data_ [DataFieldDocs](#DataFieldDocs)
- _data_ [EnumCaseDocs](#EnumCaseDocs)
- _data_ [EnumDocs](#EnumDocs)
- _data_ [ExternFieldDocs](#ExternFieldDocs)
- _data_ [ExternFunctionDocs](#ExternFunctionDocs)
- _data_ [ExternTypeDocs](#ExternTypeDocs)
- _data_ [FunctionDocs](#FunctionDocs)
- _data_ [ModuleDocs](#ModuleDocs)
- _enum_ [TypeDocs](#TypeDocs)
- _func_ [docsToMarkdown](#docsToMarkdown) docs
- _func_ [docsToMarkup](#docsToMarkup)
- _extern_ [inspect](#inspect) value

## ConstantDocs

_data_ The documentation of a constant declaration.
If a constant resolves to a function, it can still be called.

```lithia
/// The docs string
let constantName = 42
```

### Properties

- `name` - The name of the constant declaration.
- `docs` - The docs as text.

## DataDocs

_data_ The documentation of a data declaration.

```lithia
/// The docs string
data DataName {
   someFields with, params
   orWithout
}
```

### Properties

- `name` - The name of the data declaration
- `docs` - The docs as text.
- `fields` - A List of fields of DataFieldDocs

## DataFieldDocs

_data_ Describes a field of a data type

### Properties

- `name` - name of the data property
- `docs` - The docs as text.
- `params` - A List of the data field's parameter names.
Can be treated as function if not Nil.

## EnumCaseDocs

_data_ The documentation for a specific enum case.

### Properties

- `name` - The name of the enum case.
- `docs` - The docs as text.
- `type` - The underlaying type of the enum case.

## EnumDocs

_data_ Contains all documentation information for a given enum.

### Properties

- `name` - Name of the enum.
- `docs` - The docs as text.
- `cases` - Each enum case as EnumCaseDocs.

## ExternFieldDocs

_data_ Describes a field of a data type

### Properties

- `name` - name of the data property
- `docs` - The docs as text.
- `params` - A List of the extern field's parameter names.
Can be treated as function if not Nil.

## ExternFunctionDocs

_data_ The documentation of an extern function declaration.

```lithia
/// The docs string
extern externFunctionName params, asList
```

### Properties

- `name` - The name of the external function.
- `docs` - The docs as text.
- `params` - A list of its paramter names. Can be treated as constant if Nil.

## ExternTypeDocs

_data_ The documentation of an extern type declaration.

```lithia
/// The docs string
extern ExternTypeName {
   someFields with, params
   orWithout
}
```

### Properties

- `name` - The name of the external type.
- `docs` - The docs as text.
- `fields` - a list of ExternFieldDocs.

## FunctionDocs

_data_ The documentation of a function declaration.

```lithia
/// The docs string
func functionName { params => }
```

### Properties

- `name` - name of the function.
- `docs` - The docs as text.
- `params` - A list of the function's parameter names.

## ModuleDocs

_data_ The documentation of a module declaration.

```lithia
/// The docs string
module internalModuleName
```

### Properties

- `name` - name of the module
- `types` - all types declared in the module
- `docs` - The docs as text.

## TypeDocs

_enum_
All possible docs inspection values.

### Cases

- [EnumDocs](#EnumDocs)
- [DataDocs](#DataDocs)
- [ModuleDocs](#ModuleDocs)
- [FunctionDocs](#FunctionDocs)
- [ConstantDocs](#ConstantDocs)
- [ExternTypeDocs](#ExternTypeDocs)
- [ExternFunctionDocs](#ExternFunctionDocs)
- [None](#None)

## docsToMarkdown

_func_ `docsToMarkdown docs`

## docsToMarkup

_func_ `docsToMarkup`

## inspect

_func_ `inspect value`

Inspects a given value for documentation.
@returns TypeDocs

