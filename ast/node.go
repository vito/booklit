package ast

import "strings"

type Node interface {
	Visit(Visitor) error
}

type Visitor interface {
	VisitString(String) error
	VisitInvoke(Invoke) error
	VisitSequence(Sequence) error
	VisitParagraph(Paragraph) error
	VisitPreformatted(Preformatted) error
}

type String string

func (node String) Visit(visitor Visitor) error {
	return visitor.VisitString(node)
}

type Location struct {
	Line   int
	Col    int
	Offset int
}

type Invoke struct {
	Function  string
	Arguments []Node

	Location Location
}

func (node Invoke) Visit(visitor Visitor) error {
	return visitor.VisitInvoke(node)
}

func (node Invoke) Method() string {
	camel := ""
	for _, word := range strings.Split(node.Function, "-") {
		camel += strings.Title(word)
	}

	return camel
}

type Sequence []Node

func (node Sequence) Visit(visitor Visitor) error {
	return visitor.VisitSequence(node)
}

type Paragraph []Sequence

func (node Paragraph) Visit(visitor Visitor) error {
	return visitor.VisitParagraph(node)
}

type Preformatted []Sequence

func (node Preformatted) Visit(visitor Visitor) error {
	return visitor.VisitPreformatted(node)
}
