package langsrv_test

// import (
// 	"testing"

// 	"github.com/TobiasYin/go-lsp/lsp/defines"
// 	"github.com/vknabel/lithia/langsrv"
// 	"github.com/vknabel/lithia/parser"
// )

// func TestNodeAtPosition(t *testing.T) {
// 	p := parser.NewParser()
// 	fp, err := p.Parse("module", "test.go", `module test`)
// 	if len(err) > 0 {
// 		t.Error(err)
// 		return
// 	}

// 	root := fp.Tree.RootNode()
// 	node := langsrv.NodeAtPosition(root, defines.Position{Line: 0, Character: 9})
// 	if node == nil {
// 		t.Error("node not found")
// 		return
// 	}
// 	if node.Type() != parser.TYPE_NODE_IDENTIFIER {
// 		t.Errorf("expected TYPE_NODE_IDENTIFIER, got %s", node.Type())
// 	}
// }
