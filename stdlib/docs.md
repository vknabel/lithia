# docs

_module_ 

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

## ExternDocs

_data_ 

### Properties

- `name`
- `docs`

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

- [ConstantDocs](#ConstantDocs)
- [DataDocs](#DataDocs)
- [EnumDocs](#EnumDocs)
- [ExternDocs](#ExternDocs)
- [FunctionDocs](#FunctionDocs)
- [ModuleDocs](#ModuleDocs)
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

