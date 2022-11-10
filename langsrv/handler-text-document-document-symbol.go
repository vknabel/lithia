package langsrv

import (
	"fmt"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/ast"
)

func textDocumentDocumentSymbol(context *glsp.Context, params *protocol.DocumentSymbolParams) (any, error) {
	rc := NewReqContext(protocol.TextDocumentIdentifier{URI: params.TextDocument.URI})
	if rc.sourceFile == nil {
		return nil, fmt.Errorf("no source file for %s", params.TextDocument.URI)
	}
	symbols := make([]protocol.DocumentSymbol, 0)
	for _, decl := range rc.sourceFileDeclarations() {
		var children []protocol.DocumentSymbol
		switch decl := decl.(type) {
		case ast.DeclData:
			children = make([]protocol.DocumentSymbol, len(decl.Fields))
			for i, field := range decl.Fields {
				children[i] = protocol.DocumentSymbol{
					Name:           string(field.Name),
					Kind:           symbolKindClassForDecl(field),
					Range:          rangeFromAstSourceLocation(field.Meta().Source),
					SelectionRange: rangeFromAstSourceLocation(field.Meta().Source),
				}
			}
		case ast.DeclExternType:
			children = make([]protocol.DocumentSymbol, 0, len(decl.Fields))
			for _, field := range decl.Fields {
				children = append(children, protocol.DocumentSymbol{
					Name:           string(field.Name),
					Kind:           symbolKindClassForDecl(field),
					Range:          rangeFromAstSourceLocation(field.Meta().Source),
					SelectionRange: rangeFromAstSourceLocation(field.Meta().Source),
				})
			}
		}
		symbols = append(symbols, protocol.DocumentSymbol{
			Name:           string(decl.DeclName()),
			Kind:           symbolKindClassForDecl(decl),
			Range:          rangeFromAstSourceLocation(decl.Meta().Source),
			SelectionRange: rangeFromAstSourceLocation(decl.Meta().Source),
			Children:       children,
		})
	}
	return symbols, nil
}

func symbolKindClassForDecl(decl ast.Decl) protocol.SymbolKind {
	switch decl.(type) {
	case ast.DeclConstant:
		return protocol.SymbolKindConstant
	case ast.DeclData:
		return protocol.SymbolKindStruct
	case ast.DeclEnum:
		return protocol.SymbolKindEnum
	case ast.DeclEnumCase:
		return protocol.SymbolKindEnumMember
	case ast.DeclExternFunc:
		return protocol.SymbolKindFunction
	case ast.DeclExternType:
		return protocol.SymbolKindClass
	case ast.DeclField:
		return protocol.SymbolKindField
	case ast.DeclFunc:
		return protocol.SymbolKindFunction
	case ast.DeclImport:
		return protocol.SymbolKindNamespace
	case ast.DeclImportMember:
		return protocol.SymbolKindKey
	case ast.DeclModule:
		return protocol.SymbolKindModule
	case ast.DeclParameter:
		return protocol.SymbolKindVariable
	default:
		return protocol.SymbolKindVariable
	}
}
