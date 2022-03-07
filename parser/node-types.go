package parser

var (
	TYPE_NODE_SOURCE_FILE                   = "source_file"
	TYPE_NODE_MODULE_DECLARATION            = "module_declaration"
	TYPE_NODE_IMPORT_DECLARATION            = "import_declaration"
	TYPE_NODE_LET_DECLARATION               = "let_declaration"
	TYPE_NODE_FUNCTION_DECLARATION          = "function_declaration"
	TYPE_NODE_DATA_DECLARATION              = "data_declaration"
	TYPE_NODE_EXTERN_DECLARATION            = "extern_declaration"
	TYPE_NODE_DATA_PROPERTY_LIST            = "data_property_list"
	TYPE_NODE_DATA_PROPERTY_VALUE           = "data_property_value"
	TYPE_NODE_DATA_PROPERTY_FUNCTION        = "data_property_function"
	TYPE_NODE_ENUM_DECLARATION              = "enum_declaration"
	TYPE_NODE_ENUM_CASE_LIST                = "enum_case_list"
	TYPE_NODE_COMPLEX_INVOCATION_EXPRESSION = "complex_invocation_expression"
	TYPE_NODE_SIMPLE_INVOCATION_EXPRESSION  = "simple_invocation_expression"
	TYPE_NODE_UNARY_EXPRESSION              = "unary_expression"
	TYPE_NODE_BINARY_EXPRESSION             = "binary_expression"
	TYPE_NODE_MEMBER_ACCESS                 = "member_access"
	TYPE_NODE_TYPE_EXPRESSION               = "type_expression"
	TYPE_NODE_TYPE_BODY                     = "type_body"
	TYPE_NODE_TYPE_CASE                     = "type_case"
	TYPE_NODE_STRING_LITERAL                = "string_literal"
	TYPE_NODE_ESCAPE_SEQUENCE               = "escape_sequence"
	TYPE_NODE_GROUP_LITERAL                 = "group_literal"
	TYPE_NODE_INT_LITERAL                   = "int_literal"
	TYPE_NODE_FLOAT_LITERAL                 = "float_literal"
	TYPE_NODE_ARRAY_LITERAL                 = "array_literal"
	TYPE_NODE_DICT_LITERAL                  = "dict_literal"
	TYPE_NODE_DICT_ENTRY                    = "dict_entry"
	TYPE_NODE_FUNCTION_LITERAL              = "function_literal"
	TYPE_NODE_PARAMETER_LIST                = "parameter_list"
	TYPE_NODE_IDENTIFIER                    = "identifier"
	TYPE_NODE_COMMENT                       = "comment"

	// alias
	TYPE_NODE_ENUM_CASE_REFERENCE = "enum_case_reference"

	// errors
	TYPE_NODE_ERROR      = "ERROR"
	TYPE_NODE_MISSING    = "MISSING"
	TYPE_NODE_UNEXPECTED = "UNEXPECTED"
)
