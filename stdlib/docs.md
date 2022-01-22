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
- _func_ [dataFieldDocsToMarkdown](#dataFieldDocsToMarkdown) field
- _func_ [docsToMarkdown](#docsToMarkdown) docs
- _func_ [enumCaseDocsToMarkdown](#enumCaseDocsToMarkdown) case
- _func_ [inspect](#inspect) value
- _func_ [markdownDocsSuffix](#markdownDocsSuffix) docsString
- _func_ [moduleMemberDocsToMarkdown](#moduleMemberDocsToMarkdown) member

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
- `fields` - a list of ExternFieldDocs

## ExternTypeDocs

_data_ 

### Properties

- `name`
- `docs`
- `fields`

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

_enum_ All possible docs inspection values.

### Cases

- [EnumDocs](#EnumDocs)
- [DataDocs](#DataDocs)
- [ModuleDocs](#ModuleDocs)
- [FunctionDocs](#FunctionDocs)
- [ConstantDocs](#ConstantDocs)
- [ExternTypeDocs](#ExternTypeDocs)
- [ExternFunctionDocs](#ExternFunctionDocs)
- [None](#None)

## dataFieldDocsToMarkdown

_func_ `dataFieldDocsToMarkdown field`


## docsToMarkdown

_func_ `docsToMarkdown docs`


## enumCaseDocsToMarkdown

_func_ `enumCaseDocsToMarkdown case`


## inspect

_func_ `inspect value`


## markdownDocsSuffix

_func_ `markdownDocsSuffix docsString`


## moduleMemberDocsToMarkdown

_func_ `moduleMemberDocsToMarkdown member`

