// Package ast contains functions and types for parsing a Booklit document into
// a Node.
package ast

import "strings"

// Node is a Booklit syntax tree.
type Node interface {
	Visit(Visitor) error
}

// Visitor is implemented in order to traverse Node.
type Visitor interface {
	VisitString(String) error
	VisitInvoke(Invoke) error
	VisitSequence(Sequence) error
	VisitParagraph(Paragraph) error
	VisitPreformatted(Preformatted) error
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

// Invoke is a function call, e.g. \foo{bar}.
type Invoke struct {
	Function  string
	Arguments []Node

	Location Location
}

// Visit calls VisitInvoke.
func (node Invoke) Visit(visitor Visitor) error {
	return visitor.VisitInvoke(node)
}

// Method returns a method name by splitting Function on '-' and title-casing
// each word.
func (node Invoke) Method() string {
	camel := ""
	for _, word := range strings.Split(node.Function, "-") {
		camel += strings.Title(word)
	}

	return camel
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

// Preformatted is a grouping of lines, typically parsed
type Preformatted []Sequence

// Visit calls VisitPreformatted.
func (node Preformatted) Visit(visitor Visitor) error {
	return visitor.VisitPreformatted(node)
}
