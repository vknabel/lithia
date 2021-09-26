package parser

import sitter "github.com/smacker/go-tree-sitter"

func (fp *FileParser) ParseChildCommentIfNeeded(child *sitter.Node) bool {
	if child.Type() == TYPE_NODE_COMMENT {
		fp.Comments = append(fp.Comments, child.Content(fp.Source))
		return true
	} else {
		return false
	}
}
