// Package ast contains functions and types for parsing a Booklit document into
// a Node.
package ast

// Node is a Booklit syntax tree.
type Node interface {
	Visit(Visitor) error
}

// Visitor is implemented in order to traverse Node.
type Visitor interface {
	VisitString(String) error
	VisitSequence(Sequence) error
	VisitParagraph(Paragraph) error
	VisitJSXElement(JSXElement) error
	VisitJSXExpression(JSXExpression) error
}

// String is literal text, not including linebreaks.
type String string

// Visit calls VisitString.
func (node String) Visit(visitor Visitor) error {
	return visitor.VisitString(node)
}

// Location represents the location of a syntax node within the
// document.
type Location struct {
	Line   int
	Col    int
	Offset int
}

// Sequence represents adjacent nodes typically within a single
// line.
type Sequence []Node

// Visit calls VisitSequence.
func (node Sequence) Visit(visitor Visitor) error {
	return visitor.VisitSequence(node)
}

// Paragraph is a grouping of lines separated by two linebreaks.
type Paragraph []Sequence

// Visit calls VisitParagraph.
func (node Paragraph) Visit(visitor Visitor) error {
	return visitor.VisitParagraph(node)
}

// JSXElement is a JSX-style invocation, e.g. <Foo bar="x">body</Foo>.
//
// Name is the tag as authored. PascalCase tags dispatch through the
// JSX tier (built-in, Dang scope, mdx template); lowercase tags wrap
// their evaluated children in literal HTML opening/closing tags at
// evaluation time.
//
// Props maps prop names to their values. A value is either an
// ast.String (literal "..." attribute) or an ast.JSXExpression (attribute
// of the form name={expr}). The map is unordered — raw-HTML attribute
// emission sorts the names alphabetically for determinism.
//
// Children is the flat list of nodes between the opening and closing tags,
// in source order. Empty for self-closing elements.
//
// MultiLine reports whether the element spanned more than one source
// line. Block-level lowercase tags use this to decide whether to flag
// their output as Block content (so the surrounding paragraph wrapper
// skips its `<p>` wrap).
type JSXElement struct {
	Name      string
	Props     map[string]Node
	Children  []Node
	MultiLine bool

	Location Location
}

// Visit calls VisitJSXElement.
func (node JSXElement) Visit(visitor Visitor) error {
	return visitor.VisitJSXElement(node)
}

// JSXExpression is an unparsed {expr} occurrence inside JSX. The Raw field
// holds the source text between the braces; it will eventually be parsed and
// evaluated by Dang (Phase 3). For now consumers treat it as opaque.
type JSXExpression struct {
	Raw string

	Location Location
}

// Visit calls VisitJSXExpression.
func (node JSXExpression) Visit(visitor Visitor) error {
	return visitor.VisitJSXExpression(node)
}
