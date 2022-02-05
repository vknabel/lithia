package langsrv

import protocol "github.com/tliron/glsp/protocol_3_16"

type tokenType string

var (
	// For identifiers that declare or reference a namespace, module, or package.
	token_namespace tokenType = "namespace"
	// For identifiers that declare or reference a class type.
	token_class tokenType = "class"
	// For identifiers that declare or reference an enumeration type.
	token_enum tokenType = "enum"
	// For identifiers that declare or reference an interface type.
	token_interface tokenType = "interface"
	// For identifiers that declare or reference a struct type.
	token_struct tokenType = "struct"
	// For identifiers that declare or reference a type parameter.
	token_typeParameter tokenType = "typeParameter"
	// For identifiers that declare or reference a type that is not covered above.
	token_type tokenType = "type"
	// For identifiers that declare or reference a function or method parameters.
	token_parameter tokenType = "parameter"
	// For identifiers that declare or reference a local or global variable.
	token_variable tokenType = "variable"
	// For identifiers that declare or reference a member property, member field, or member variable.
	token_property tokenType = "property"
	// For identifiers that declare or reference an enumeration property, constant, or member.
	token_enumMember tokenType = "enumMember"
	// For identifiers that declare or reference decorators and annotations.
	token_decorator tokenType = "decorator"
	// For identifiers that declare an event property.
	token_event tokenType = "event"
	// For identifiers that declare a function.
	token_function tokenType = "function"
	// For identifiers that declare a member function or method.
	token_method tokenType = "method"
	// For identifiers that declare a macro.
	token_macro tokenType = "macro"
	// For identifiers that declare a label.
	token_label tokenType = "label"
	// For tokens that represent a comment.
	token_comment tokenType = "comment"
	// For tokens that represent a string literal.
	token_string tokenType = "string"
	// For tokens that represent a language keyword.
	token_keyword tokenType = "keyword"
	// For tokens that represent a number literal.
	token_number tokenType = "number"
	// For tokens that represent a regular expression literal.
	token_regexp tokenType = "regexp"
	// For tokens that represent an operator.
	token_operator tokenType = "operator"
)

var allTokenTypes = []tokenType{
	token_namespace,
	token_class,
	token_enum,
	token_interface,
	token_struct,
	token_typeParameter,
	token_type,
	token_parameter,
	token_variable,
	token_property,
	token_enumMember,
	token_decorator,
	token_event,
	token_function,
	token_method,
	token_macro,
	token_label,
	token_comment,
	token_string,
	token_keyword,
	token_number,
	token_regexp,
	token_operator,
}

func (tt tokenType) bitflag() protocol.UInteger {
	switch tt {
	case token_namespace:
		return 1
	case token_class:
		return 2
	case token_enum:
		return 3
	case token_interface:
		return 3
	case token_struct:
		return 4
	case token_typeParameter:
		return 5
	case token_type:
		return 6
	case token_parameter:
		return 7
	case token_variable:
		return 8
	case token_property:
		return 9
	case token_enumMember:
		return 10
	case token_decorator:
		return 11
	case token_event:
		return 12
	case token_function:
		return 13
	case token_method:
		return 14
	case token_macro:
		return 15
	case token_label:
		return 16
	case token_comment:
		return 17
	case token_string:
		return 18
	case token_keyword:
		return 19
	case token_number:
		return 20
	case token_regexp:
		return 21
	case token_operator:
		return 22
	default:
		return 0
	}
}

type tokenModifier string

var (
	// For declarations of symbols.
	modifier_declaration tokenModifier = "declaration"
	// For definitions of symbols, for example, in header files.
	modifier_definition tokenModifier = "definition"
	// For readonly variables and member fields (constants).
	modifier_readonly tokenModifier = "readonly"
	// For class members (static members).
	modifier_static tokenModifier = "static"
	// For symbols that should no longer be used.
	modifier_deprecated tokenModifier = "deprecated"
	// For types and member functions that are abstract.
	modifier_abstract tokenModifier = "abstract"
	// For functions that are marked async.
	modifier_async tokenModifier = "async"
	// For variable references where the variable is assigned to.
	modifier_modification tokenModifier = "modification"
	// For occurrences of symbols in documentation.
	modifier_documentation tokenModifier = "documentation"
	// For symbols that are part of the standard library.
	modifier_defaultLibrary tokenModifier = "defaultLibrary"
)

var allTokenModifiers = []tokenModifier{
	modifier_declaration,
	modifier_definition,
	modifier_readonly,
	modifier_static,
	modifier_deprecated,
	modifier_abstract,
	modifier_async,
	modifier_modification,
	modifier_documentation,
	modifier_defaultLibrary,
}

func (tm tokenModifier) bitflag() protocol.UInteger {
	switch tm {
	case modifier_declaration:
		return 1
	case modifier_definition:
		return 2
	case modifier_readonly:
		return 3
	case modifier_static:
		return 4
	case modifier_deprecated:
		return 5
	case modifier_abstract:
		return 6
	case modifier_async:
		return 7
	case modifier_modification:
		return 8
	case modifier_documentation:
		return 9
	case modifier_defaultLibrary:
		return 10
	default:
		return 0
	}
}

type highlightedToken struct {
	line           uint32
	column         uint32
	length         uint32
	tokenType      tokenType
	tokenModifiers []tokenModifier
}

func (tok *highlightedToken) serialize(previous highlightedToken) []protocol.UInteger {
	modifiers := protocol.UInteger(0)
	for _, modifier := range tok.tokenModifiers {
		modifiers |= modifier.bitflag()
	}
	deltaLine := protocol.UInteger(tok.line - previous.line)
	var deltaStartChar protocol.UInteger
	if deltaLine == 0 {
		deltaStartChar = protocol.UInteger(tok.column - previous.column)
	} else {
		deltaStartChar = protocol.UInteger(tok.column)
	}

	return []protocol.UInteger{
		deltaLine,
		deltaStartChar,
		tok.length,
		tok.tokenType.bitflag(),
		modifiers,
	}
}

func serializeHighlightedTokens(tokens []highlightedToken) []protocol.UInteger {
	var serializedTokens []protocol.UInteger
	for i, tok := range tokens {
		if i == 0 {
			serializedTokens = append(serializedTokens, tok.serialize(highlightedToken{})...)
		} else {
			serializedTokens = append(serializedTokens, tok.serialize(tokens[i-1])...)
		}
	}
	return serializedTokens
}
