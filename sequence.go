package booklit

type Sequence []Content

func (seq Sequence) String() string {
	str := ""
	for _, content := range seq {
		str += content.String()
	}

	return str
}

func (seq Sequence) Contents() []Content {
	return seq
}

func (seq Sequence) Visit(visitor Visitor) {
	visitor.VisitSequence(seq)
}
