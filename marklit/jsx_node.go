package marklit

import (
	"fmt"
	"strings"

	gast "github.com/yuin/goldmark/ast"
)

// KindJSXElement is a NodeKind for JSX-style component invocations.
var KindJSXElement = gast.NewNodeKind("BooklitJSXElement")

// JSXPropKind distinguishes between string-literal and expression
// attribute values.
type JSXPropKind int

const (
	// JSXPropString indicates a `name="..."` attribute. Value holds
	// the bytes between the quotes.
	JSXPropString JSXPropKind = iota
	// JSXPropExpression indicates a `name={...}` attribute. Value holds
	// the bytes between the braces.
	JSXPropExpression
)

// JSXProp is a single attribute on a JSX element.
type JSXProp struct {
	Name  string
	Kind  JSXPropKind
	Value []byte
}

// JSXChildKind distinguishes the kinds of child nodes inside a JSX element.
type JSXChildKind int

const (
	// JSXChildText is a stretch of raw markdown text between tags or
	// expressions. The bytes are re-parsed as markdown at convert time
	// so emphasis, links, etc. inside JSX children continue to work.
	JSXChildText JSXChildKind = iota
	// JSXChildElement is a nested JSX element.
	JSXChildElement
	// JSXChildExpression is a {...} expression as a child, captured
	// opaquely. Dang will eventually parse the contents.
	JSXChildExpression
)

// JSXChild is one piece of content inside a JSX element's body.
type JSXChild struct {
	Kind JSXChildKind
	Text []byte
	Elem *JSXElementNode
}

// JSXElementNode represents a <Name ...>...</Name> invocation in the goldmark
// AST. Props and children are stored as parsed sub-structures; convert.go
// turns this into ast.JSXElement.
type JSXElementNode struct {
	gast.BaseInline

	Name        string
	Props       []JSXProp
	SelfClosing bool // true for <Foo/>; false for <Foo>...</Foo>, even when Children is empty
	Children    []JSXChild

	Line int
	Col  int
}

// Kind implements ast.Node.Kind.
func (n *JSXElementNode) Kind() gast.NodeKind { return KindJSXElement }

// Dump implements ast.Node.Dump.
func (n *JSXElementNode) Dump(source []byte, level int) {
	indent := strings.Repeat("    ", level)
	fmt.Printf("%sBooklitJSXElement {\n", indent)
	fmt.Printf("%s    Name: %q\n", indent, n.Name)
	fmt.Printf("%s    Props: %d\n", indent, len(n.Props))
	fmt.Printf("%s    Children: %d\n", indent, len(n.Children))
	fmt.Printf("%s}\n", indent)
}

// NewJSXElementNode returns a new JSXElementNode.
func NewJSXElementNode(name string) *JSXElementNode {
	return &JSXElementNode{Name: name}
}

// KindJSXBlockElement is a NodeKind for block-level JSX invocations.
var KindJSXBlockElement = gast.NewNodeKind("BooklitJSXBlockElement")

// JSXBlockElementNode is the block-context twin of JSXElementNode. Same
// data, different goldmark base type. The block parser produces this so
// that goldmark's HTML block parser doesn't claim `<UpperCase` lines as
// raw HTML.
type JSXBlockElementNode struct {
	gast.BaseBlock

	Name        string
	Props       []JSXProp
	SelfClosing bool
	Children    []JSXChild

	Line int
	Col  int
}

// Kind implements ast.Node.Kind.
func (n *JSXBlockElementNode) Kind() gast.NodeKind { return KindJSXBlockElement }

// Dump implements ast.Node.Dump.
func (n *JSXBlockElementNode) Dump(source []byte, level int) {
	indent := strings.Repeat("    ", level)
	fmt.Printf("%sBooklitJSXBlockElement {\n", indent)
	fmt.Printf("%s    Name: %q\n", indent, n.Name)
	fmt.Printf("%s    Props: %d\n", indent, len(n.Props))
	fmt.Printf("%s    Children: %d\n", indent, len(n.Children))
	fmt.Printf("%s}\n", indent)
}

// NewJSXBlockElementNode returns a new JSXBlockElementNode.
func NewJSXBlockElementNode(name string) *JSXBlockElementNode {
	return &JSXBlockElementNode{Name: name}
}
