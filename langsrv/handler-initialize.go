package langsrv

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/info"
)

func initialize(context *glsp.Context, params *protocol.InitializeParams) (interface{}, error) {
	capabilities := handler.CreateServerCapabilities()
	capabilities.CompletionProvider = &protocol.CompletionOptions{
		TriggerCharacters: []string{"."},
	}
	capabilities.SemanticTokensProvider = protocol.SemanticTokensRegistrationOptions{
		SemanticTokensOptions: protocol.SemanticTokensOptions{
			Legend: protocol.SemanticTokensLegend{
				TokenTypes:     tokenTypeLegend(),
				TokenModifiers: tokenModifierLegend(),
			},
			Full: true,
		},
	}

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &info.Version,
		},
	}, nil
}

func tokenTypeLegend() []string {
	legend := make([]string, len(allTokenTypes))
	for i, tokenType := range allTokenTypes {
		legend[i] = string(tokenType)
	}
	return legend
}

func tokenModifierLegend() []string {
	legend := make([]string, len(allTokenModifiers))
	for i, tokenModifier := range allTokenModifiers {
		legend[i] = string(tokenModifier)
	}
	return legend
}
