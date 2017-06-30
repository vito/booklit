package ast

type Node interface {
	Visit(Visitor) error
}

type Visitor interface {
	VisitString(String) error
	VisitInvoke(Invoke) error
	VisitSequence(Sequence) error
	VisitParagraph(Paragraph) error
}

type String string

func (node String) Visit(visitor Visitor) error {
	return visitor.VisitString(node)
}

type Invoke struct {
	Method    string
	Arguments []Node
}

func (node Invoke) Visit(visitor Visitor) error {
	return visitor.VisitInvoke(node)
}

type Sequence []Node

func (node Sequence) Visit(visitor Visitor) error {
	return visitor.VisitSequence(node)
}

type Paragraph []Sequence

func (node Paragraph) Visit(visitor Visitor) error {
	return visitor.VisitParagraph(node)
}
