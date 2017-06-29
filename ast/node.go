package ast

type Node interface {
	Visit(Visitor)
}

type Visitor interface {
	VisitString(String)
	VisitInvoke(Invoke)
	VisitSequence(Sequence)
}

type String string

func (str String) Visit(visitor Visitor) {
	visitor.VisitString(str)
}

type Invoke struct {
	Method    string
	Arguments []Node
}

func (inv Invoke) Visit(visitor Visitor) {
	visitor.VisitInvoke(inv)
}

type Sequence []Node

func (seq Sequence) Visit(visitor Visitor) {
	visitor.VisitSequence(seq)
}
