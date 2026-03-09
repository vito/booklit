package marklit

import (
	"fmt"
	"strings"

	gast "github.com/yuin/goldmark/ast"
)

// KindInvoke is a NodeKind for Booklit function invocations.
var KindInvoke = gast.NewNodeKind("BooklitInvoke")

// InvokeNode represents an @function{arg1}{arg2} invocation in the goldmark
// AST. Each argument is stored as a child node.
type InvokeNode struct {
	gast.BaseInline

	// Function is the name of the function being invoked, e.g. "title".
	Function string

	// Each argument's raw source text. Arguments are parsed recursively
	// and stored as children, but we also keep the raw bytes for arguments
	// that need to be passed as ast.Node (unevaluated).
	RawArgs [][]byte

	// Line/Col for error reporting.
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

// KindInvokeBlock is a NodeKind for block-level Booklit invocations.
var KindInvokeBlock = gast.NewNodeKind("BooklitInvokeBlock")

// InvokeBlockNode represents a block-level @function{...} invocation.
// This is used when an invocation starts at the beginning of a line and
// should be treated as block content rather than inline.
type InvokeBlockNode struct {
	gast.BaseBlock

	// Function is the name of the function being invoked.
	Function string

	// RawArgs holds the raw source text of each argument.
	RawArgs [][]byte

	// Line/Col for error reporting.
	Line int
	Col  int
}

// Kind implements ast.Node.Kind.
func (n *InvokeBlockNode) Kind() gast.NodeKind {
	return KindInvokeBlock
}

// Dump implements ast.Node.Dump.
func (n *InvokeBlockNode) Dump(source []byte, level int) {
	indent := strings.Repeat("    ", level)
	fmt.Printf("%sBooklitInvokeBlock {\n", indent)
	fmt.Printf("%s    Function: %q\n", indent, n.Function)
	fmt.Printf("%s    Args: %d\n", indent, len(n.RawArgs))
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		c.Dump(source, level+1)
	}
	fmt.Printf("%s}\n", indent)
}

// NewInvokeBlockNode returns a new InvokeBlockNode.
func NewInvokeBlockNode(function string) *InvokeBlockNode {
	return &InvokeBlockNode{
		Function: function,
	}
}
