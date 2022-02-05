package langsrv

import (
	"io/ioutil"
	"net/url"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/vknabel/lithia/ast"
	"github.com/vknabel/lithia/parser"
)

type ReqContext struct {
	textDocument protocol.TextDocumentIdentifier
	position     protocol.Position
	parser       *parser.Parser

	contents   string
	token      string
	fileParser *parser.FileParser
	sourceFile *ast.SourceFile
}

func NewReqContext(textDocument protocol.TextDocumentIdentifier) *ReqContext {
	return &ReqContext{
		textDocument: textDocument,
		parser:       parser.NewParser(),
	}
}

func NewReqContextAtPosition(position *protocol.TextDocumentPositionParams) *ReqContext {
	return &ReqContext{
		textDocument: position.TextDocument,
		position:     position.Position,
		parser:       parser.NewParser(),
	}
}

func (rc *ReqContext) readFile() (string, error) {
	if rc.contents != "" {
		return rc.contents, nil
	}
	enEscapeUrl, _ := url.QueryUnescape(string(rc.textDocument.URI))
	var data []byte
	var err error
	if strings.HasPrefix(enEscapeUrl, "file://") {
		data, err = ioutil.ReadFile(enEscapeUrl[7:])
	} else {
		data, err = ioutil.ReadFile(enEscapeUrl)
	}
	if err != nil {
		return "", err
	}
	rc.contents = string(data)
	return rc.contents, nil
}

func (rc *ReqContext) findToken() (string, *protocol.Range, error) {
	node, err := rc.findNode()
	if err != nil || node == nil {
		return "", nil, err
	}
	name := node.Content([]byte(rc.contents))
	return name, &protocol.Range{
		Start: protocol.Position{
			Line:      uint32(node.StartPoint().Row),
			Character: uint32(node.StartPoint().Column),
		},
		End: protocol.Position{
			Line:      uint32(node.EndPoint().Row),
			Character: uint32(node.EndPoint().Column),
		},
	}, nil
}

func (rc *ReqContext) createFileParser() (*parser.FileParser, error) {
	if rc.fileParser != nil {
		return rc.fileParser, nil
	}
	contents, err := rc.readFile()
	if err != nil {
		return nil, err
	}
	fileParser, errs := rc.parser.Parse("default-module", string(rc.textDocument.URI), contents)
	if len(errs) > 0 {
		return nil, parser.NewGroupedSyntaxError(errs)
	}
	rc.fileParser = fileParser
	return rc.fileParser, nil
}

func (rc *ReqContext) parseSourceFile() (*ast.SourceFile, error) {
	if rc.sourceFile != nil {
		return rc.sourceFile, nil
	}
	fileParser, err := rc.createFileParser()
	if err != nil {
		return nil, err
	}

	sourceFile, errs := fileParser.ParseSourceFile()
	if len(errs) > 0 {
		return nil, parser.NewGroupedSyntaxError(errs)
	}
	rc.sourceFile = sourceFile
	return rc.sourceFile, nil
}

func (rc *ReqContext) findNode() (*sitter.Node, error) {
	_, err := rc.parseSourceFile()
	if err != nil {
		return nil, err
	}
	node := NodeAtPosition(rc.fileParser.Tree.RootNode(), rc.position)
	return node, nil
}
