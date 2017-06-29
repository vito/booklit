package ast

type Node interface{}

type String string

type Invoke struct {
	Name      string
	Arguments []Node
}

type Sequence []Node
