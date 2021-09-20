# docs

_module_ 

## ConstantDocs

_data_ 

### Properties

- name
- docs

## DataDocs

_data_ 

### Properties

- name - The name of the data declaration
- docs
- fields - A List of fields of DataFieldDocs

## DataFieldDocs

_data_ Describes a field of a data type

### Properties

- name - name of the data property
- docs
- params

## EnumCaseDocs

_data_ The documentation for a specific enum case.

### Properties

- name
- docs
- type

## EnumDocs

_data_ Contains all documentation information for a given enum.

### Properties

- name - Name of the enum
- docs
- cases

## ExternDocs

_data_ 

### Properties

- name
- docs

## FunctionDocs

_data_ 

### Properties

- name
- docs
- params

## ModuleDocs

_data_ 

### Properties

- name - name of the module
- types - all types declared in the module
- docs

## TypeDocs

_enum_ All possible docs inspection values.

### Cases

- [EnumDocs](#EnumDocs)
- [DataDocs](#DataDocs)
- [ModuleDocs](#ModuleDocs)
- [FunctionDocs](#FunctionDocs)
- [ConstantDocs](#ConstantDocs)
- [ExternDocs](#ExternDocs)
- [None](#None)

## dataFieldDocsToMarkdown

_func_ 

### Parameters

- field

## docsToMarkdown

_func_ 

### Parameters

- docs

## enumCaseDocsToMarkdown

_func_ 

### Parameters

- case

## inspect

_func_ 

### Parameters

- value

## markdownDocsSuffix

_func_ 

### Parameters

- docsString
