package langsrv

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	syntax "github.com/vknabel/tree-sitter-lithia"
)

func textDocumentSemanticTokensFull(context *glsp.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	entry := langserver.documentCache.documents[params.TextDocument.URI]
	if entry == nil {
		return nil, nil
	}
	rootNode := entry.fileParser.Tree.RootNode()
	highlightsQuery, err := sitter.NewQuery([]byte(`
	[
		"func"
		"let"
		"enum"
		"data"
		"module"
		"import"
		"extern"
		"type"
	] @keyword
	
	[
		"=>"
	] @operator
	
	[
		","
		"."
	] @punctuation
	
	[
		"{"
		"}"
		"("
		")"
		"["
		"]"
	] @punctuation.bracket
	
	(binary_expression operator: (["*" "/" "+" "-" "==" "!=" ">=" ">" "<" "<=" "&&" "||"]) @operator) ; extract
	(unary_expression operator: (["!"]) @operator)
	
	(parameter_list (identifier) @variable.parameter)
	(number_literal) @constant.numeric
	(comment) @comment
	(function_declaration name: (identifier) @function)
	(let_declaration name: (identifier) @variable)
	(enum_declaration name: (identifier) @type.enum)
	(enum_case_reference) @type.case
	(data_declaration name: (identifier) @type.data)
	(data_property_function name: (identifier) @function)
	(data_property_value name: (identifier) @property)
	(extern_declaration
		name: (identifier) @variable.builtin
		!properties
		!parameters)
	(extern_declaration
		name: (identifier) @function.builtin
		!properties)
	(extern_declaration
		name: (identifier) @type.builtin
		!parameters)
	(import_declaration name: (import_module) @variable.import)
	(import_members (identifier) @variable.import)
	(module_declaration name: (identifier) @variable.import)
	(complex_invocation_expression function: (identifier) @function)
	(simple_invocation_expression function: (identifier) @function.simple)
	(string_literal) @string
	(escape_sequence) @string.special
	(type_expression type: (identifier) @type.enum)
	(type_case label: (identifier) @type.case)
	(simple_invocation_expression function: (member_access (member_identifier) @function @method))
	(complex_invocation_expression function: (member_access (member_identifier) @function @method))
	(member_identifier) @property
	
	(ERROR) @error
	; (identifier) @variable ; would override every other token
	`), syntax.GetLanguage())
	if err != nil {
		return nil, err
	}
	cursor := sitter.NewQueryCursor()
	cursor.Exec(highlightsQuery, rootNode)
	defer cursor.Close()

	tokens := make([]highlightedToken, 0)
	for match, ok := cursor.NextMatch(); ok; match, ok = cursor.NextMatch() {
		for _, capture := range match.Captures {
			captureName := highlightsQuery.CaptureNameForId(capture.Index)
			capturedNode := capture.Node
			tokenType := tokenTypeForCaptureName(captureName)
			if tokenType == nil {
				continue
			}
			tokenModifiers := tokenModifiersForCaptureName(captureName)
			tokens = append(tokens, highlightedToken{
				line:           uint32(capturedNode.StartPoint().Row),
				column:         uint32(capturedNode.StartPoint().Column),
				length:         capturedNode.EndByte() - capturedNode.StartByte(),
				tokenType:      *tokenType,
				tokenModifiers: tokenModifiers,
			})
		}
	}
	return &protocol.SemanticTokens{
		Data: serializeHighlightedTokens(tokens),
	}, nil
}

func tokenTypeForCaptureName(captureName string) *tokenType {
	switch captureName {
	case "keyword":
		return &token_keyword
	case "operator":
		return &token_operator
	case "punctuation":
		return &token_operator
	case "punctuation.bracket":
		return &token_operator
	case "variable":
		return &token_variable
	case "variable.parameter":
		return &token_parameter
	case "variable.builtin":
		return &token_variable
	case "variable.import":
		return &token_namespace
	case "constant.numeric":
		return &token_number
	case "comment":
		return &token_comment
	case "function":
		return &token_function
	case "function.builtin":
		return &token_function
	case "function.simple":
		return &token_decorator
	case "method":
		return &token_method
	case "type":
		return &token_type
	case "type.enum":
		return &token_enum
	case "type.case":
		return &token_enumMember
	case "type.data":
		return &token_class
	case "type.builtin":
		return &token_type
	case "property":
		return &token_property
	case "string":
		return &token_string
	case "string.special":
		return &token_string
	case "error":
		return nil
	default:
		return nil
	}
}

func tokenModifiersForCaptureName(captureName string) []tokenModifier {
	switch captureName {
	case "variable":
		return []tokenModifier{modifier_readonly}
	case "type.enum", "type.data":
		return []tokenModifier{modifier_declaration}
	case "variable.builtin", "function.builtin", "type.builtin":
		return []tokenModifier{modifier_declaration, modifier_defaultLibrary, modifier_static, modifier_readonly}
	case "string.special":
		return []tokenModifier{modifier_modification}
	default:
		return nil
	}
}
