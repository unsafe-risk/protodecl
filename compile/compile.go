package compile

import (
	"log"

	"github.com/unsafe-risk/protodecl/ast"
)

func Compile(t *ast.Tree) ([]byte, error) {
	for i := range t.Nodes {
		switch node := t.Nodes[i].(type) {
		case *ast.EnumerationType:
			for j := range node.Values {
				_ = j
			}
		default:
			// skip
			log.Println("Warning: unknown node type:", node)
		}
	}
	return nil, nil
}
