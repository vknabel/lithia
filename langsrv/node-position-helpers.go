package langsrv

import (
	sitter "github.com/smacker/go-tree-sitter"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func NodeAtPosition(node *sitter.Node, position protocol.Position) *sitter.Node {
	if !includesNodePosition(node, position) {
		return nil
	}
	childCount := int(node.ChildCount())
	commonFields := []string{
		"name",
		"members",
		"value",
		"function",
		"properties",
		"parameters",
		"cases",
		"operator",
		"object",
		"type",
		"body",
		"label",
		"expression",
	}
	for _, field := range commonFields {
		if node.ChildByFieldName(field) != nil {
			child := node.ChildByFieldName(field)
			if includesNodePosition(child, position) {
				return NodeAtPosition(child, position)
			}
		}
	}
	for i := 0; i < childCount; i++ {
		child := node.Child(i)
		if includesNodePosition(child, position) {
			return NodeAtPosition(child, position)
		}
	}
	return node
}

func includesNodePosition(node *sitter.Node, position protocol.Position) bool {
	if node == nil {
		return false
	}
	startRow := uint32(node.StartPoint().Row)
	startCol := uint32(node.StartPoint().Column)
	endRow := uint32(node.EndPoint().Row)
	endCol := uint32(node.EndPoint().Column)

	if startRow <= position.Line && endRow >= position.Line {
		if startRow == position.Line && startCol > position.Character {
			return false
		}
		if endRow == position.Line && endCol < position.Character {
			return false
		}
		return true
	}
	return false
}
