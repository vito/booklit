package ast

type Node interface {
	Visit(Visitor)
}

type Visitor interface {
	VisitString(String)
	VisitInvoke(Invoke)
	VisitSequence(Sequence)
	VisitParagraph(Paragraph)
}

type String string

func (node String) Visit(visitor Visitor) {
	visitor.VisitString(node)
}

type Invoke struct {
	Method    string
	Arguments []Node
}

func (node Invoke) Visit(visitor Visitor) {
	visitor.VisitInvoke(node)
}

type Sequence []Node

func (node Sequence) Visit(visitor Visitor) {
	visitor.VisitSequence(node)
}

type Paragraph []Sequence

func (node Paragraph) Visit(visitor Visitor) {
	visitor.VisitParagraph(node)
}
