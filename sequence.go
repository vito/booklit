package booklit

type Sequence []Content

func (seq Sequence) Contents() []Content {
	return seq
}

func (seq Sequence) Visit(visitor Visitor) {
	visitor.VisitSequence(seq)
}
