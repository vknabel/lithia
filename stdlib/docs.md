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
- _func_ [dataFieldDocsToMarkup](#dataFieldDocsToMarkup) field
- _func_ [docsToMarkdown](#docsToMarkdown) docs
- _func_ [docsToMarkup](#docsToMarkup)
- _func_ [enumCaseDocsToMarkup](#enumCaseDocsToMarkup) case
- _extern_ [inspect](#inspect) value
- _func_ [moduleMemberDocsToMarkup](#moduleMemberDocsToMarkup)

## ConstantDocs

_data_

### Properties

- `name`
- `docs`

## DataDocs

_data_

### Properties

- `name` - The name of the data declaration
- `docs`
- `fields` - A List of fields of DataFieldDocs

## DataFieldDocs

_data_ Describes a field of a data type

### Properties

- `name` - name of the data property
- `docs`
- `params`

## EnumCaseDocs

_data_ The documentation for a specific enum case.

### Properties

- `name`
- `docs`
- `type`

## EnumDocs

_data_ Contains all documentation information for a given enum.

### Properties

- `name` - Name of the enum
- `docs`
- `cases`

## ExternFieldDocs

_data_ Describes a field of a data type

### Properties

- `name` - name of the data property
- `docs`
- `params`

## ExternFunctionDocs

_data_

### Properties

- `name`
- `docs`
- `params`

## ExternTypeDocs

_data_

### Properties

- `name`
- `docs`
- `fields` - a list of ExternFieldDocs

## FunctionDocs

_data_

### Properties

- `name`
- `docs`
- `params`

## ModuleDocs

_data_

### Properties

- `name` - name of the module
- `types` - all types declared in the module
- `docs`

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

## dataFieldDocsToMarkup

_func_ `dataFieldDocsToMarkup field`

## docsToMarkdown

_func_ `docsToMarkdown docs`

## docsToMarkup

_func_ `docsToMarkup`

## enumCaseDocsToMarkup

_func_ `enumCaseDocsToMarkup case`

## inspect


_extern_ `inspect value`

## moduleMemberDocsToMarkup

_func_ `moduleMemberDocsToMarkup`

