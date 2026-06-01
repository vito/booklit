package marklit

import (
	"fmt"
	"strings"

	gast "github.com/yuin/goldmark/ast"
)

// KindInvoke is a NodeKind for Booklit function invocations.
var KindInvoke = gast.NewNodeKind("BooklitInvoke")

// InvokeNode is the goldmark-side carrier for an `ast.Invoke` produced by
// the `[#tag]` reference shorthand. It is the only remaining producer of
// invokes inside the goldmark AST; everything else flows through JSX nodes.
type InvokeNode struct {
	gast.BaseInline

	Function string
	RawArgs  [][]byte

	Line int
	Col  int
}

// Kind implements ast.Node.Kind.
func (n *InvokeNode) Kind() gast.NodeKind {
	return KindInvoke
}

// Dump implements ast.Node.Dump.
func (n *InvokeNode) Dump(source []byte, level int) {
	indent := strings.Repeat("    ", level)
	fmt.Printf("%sBooklitInvoke {\n", indent)
	fmt.Printf("%s    Function: %q\n", indent, n.Function)
	fmt.Printf("%s    Args: %d\n", indent, len(n.RawArgs))
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		c.Dump(source, level+1)
	}
	fmt.Printf("%s}\n", indent)
}

// NewInvokeNode returns a new InvokeNode.
func NewInvokeNode(function string) *InvokeNode {
	return &InvokeNode{
		Function: function,
	}
}
