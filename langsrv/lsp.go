package langsrv

// func (ls *LangServer) Start() {
// 	ls.server = lsp.NewServer(&lsp.Options{
// 		Network: ls.network,
// 		Address: ls.address,
// 		CompletionProvider: &defines.CompletionOptions{
// 			TriggerCharacters: &[]string{".", "type "},
// 		},
// 	})

// 	ls.server.OnTypeDefinition(func(ctx context.Context, req *defines.TypeDefinitionParams) (result *[]defines.LocationLink, err error) {
// 		rc := NewReqContextAtPosition(&req.TextDocumentPositionParams)
// 		_, err = rc.parseSourceFile()
// 		if err != nil {
// 			return nil, err
// 		}
// 		return nil, nil
// 	})

// 	ls.server.OnHover(func(ctx context.Context, req *defines.HoverParams) (result *defines.Hover, err error) {
// 		rc := NewReqContextAtPosition(&req.TextDocumentPositionParams)
// 		sourceFile, err := rc.parseSourceFile()
// 		if err != nil {
// 			return nil, err
// 		}
// 		name, tokenRange, err := rc.findToken()
// 		if err != nil {
// 			return nil, err
// 		}
// 		for _, decl := range sourceFile.Declarations {
// 			if string(decl.DeclName()) != name {
// 				continue
// 			}
// 			var docs string
// 			if documented, ok := decl.(ast.Documented); ok {
// 				docs = documented.ProvidedDocs().Content + "\n"
// 			}
// 			return &defines.Hover{
// 				Contents: defines.MarkupContent{
// 					Kind: defines.MarkupKindMarkdown,
// 					Value: docs + "```lithia\n" + string(decl.DeclName()) + "\n```\n" +
// 						string(decl.Meta().ModuleName),
// 				},
// 				Range: tokenRange,
// 			}, nil
// 		}
// 		return nil, nil
// 	})

// 	ls.server.OnCompletion(func(ctx context.Context, req *defines.CompletionParams) (result *[]defines.CompletionItem, err error) {
// 		rc := NewReqContextAtPosition(&req.TextDocumentPositionParams)
// 		sourceFile, err := rc.parseSourceFile()
// 		if err != nil {
// 			return nil, err
// 		}
// 		completionItems := []defines.CompletionItem{}
// 		for _, decl := range sourceFile.Declarations {
// 			var docs string
// 			if documented, ok := decl.(ast.Documented); ok {
// 				docs = documented.ProvidedDocs().Content + "\n"
// 			}
// 			var detail string
// 			if decl.Meta().ModuleName != "" {
// 				detail = string(decl.Meta().ModuleName) +
// 					"." +
// 					string(decl.DeclName())
// 			}
// 			kind := defines.CompletionItemKindEnum
// 			insertText := string(decl.DeclName())
// 			completionItems = append(completionItems, defines.CompletionItem{
// 				Label:      string(decl.DeclName()),
// 				Kind:       &kind,
// 				InsertText: &insertText,
// 				Detail:     &detail,
// 				Documentation: &defines.MarkupContent{
// 					Kind: defines.MarkupKindMarkdown,
// 					Value: docs + "```lithia\n" + string(decl.DeclName()) + "\n```\n" +
// 						string(decl.Meta().ModuleName),
// 				},
// 			})
// 		}
// 		return &completionItems, nil
// 	})

// 	ls.server.OnDefinition(func(ctx context.Context, req *defines.DefinitionParams) (result *[]defines.LocationLink, err error) {
// 		rc := NewReqContextAtPosition(&req.TextDocumentPositionParams)
// 		sourceFile, err := rc.parseSourceFile()
// 		if err != nil {
// 			return nil, err
// 		}
// 		token, _, err := rc.findToken()
// 		if err != nil {
// 			return nil, err
// 		}
// 		for _, decl := range sourceFile.Declarations {
// 			if string(decl.DeclName()) != token || decl.Meta().Source == nil {
// 				continue
// 			}
// 			return &[]defines.LocationLink{
// 				{
// 					TargetUri: defines.DocumentUri(decl.Meta().Source.FileName),
// 					TargetRange: defines.Range{
// 						Start: defines.Position{
// 							Line:      uint(decl.Meta().Source.Start.Line),
// 							Character: uint(decl.Meta().Source.Start.Column),
// 						},
// 						End: defines.Position{
// 							Line:      uint(decl.Meta().Source.End.Line),
// 							Character: uint(decl.Meta().Source.End.Line),
// 						},
// 					},
// 				},
// 			}, nil
// 		}
// 		return nil, nil
// 	})

// 	ls.server.OnDeclaration(func(ctx context.Context, req *defines.DeclarationParams) (result *[]defines.LocationLink, err error) {
// 		rc := NewReqContextAtPosition(&req.TextDocumentPositionParams)
// 		sourceFile, err := rc.parseSourceFile()
// 		if err != nil {
// 			return nil, err
// 		}
// 		token, _, err := rc.findToken()
// 		if err != nil {
// 			return nil, err
// 		}
// 		for _, decl := range sourceFile.Declarations {
// 			if string(decl.DeclName()) != token || decl.Meta().Source == nil {
// 				continue
// 			}
// 			return &[]defines.LocationLink{
// 				{
// 					TargetUri: defines.DocumentUri(decl.Meta().Source.FileName),
// 					TargetRange: defines.Range{
// 						Start: defines.Position{
// 							Line:      uint(decl.Meta().Source.Start.Line),
// 							Character: uint(decl.Meta().Source.Start.Column),
// 						},
// 						End: defines.Position{
// 							Line:      uint(decl.Meta().Source.End.Line),
// 							Character: uint(decl.Meta().Source.End.Line),
// 						},
// 					},
// 				},
// 			}, nil
// 		}
// 		return nil, nil
// 	})

// 	// TODO: textDocument/semanticTokens/full
// 	// server.OnDocumentHighlight(func(ctx context.Context, req *defines.DocumentHighlightParams) (result *[]defines.DocumentHighlight, err error) {
// 	// 	panic("implement me")
// 	// })

// 	// server.OnDidSaveTextDocument(func(ctx context.Context, req *defines.DidSaveTextDocumentParams) (err error) {
// 	// 	panic("implement me")
// 	// })
// 	rpc := jsonrpc.NewServer()
// 	rpc.ConnComeIn(lsp.NewStdio())

// 	ls.server.Run()
// }
