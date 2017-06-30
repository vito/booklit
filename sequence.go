package booklit

type Sequence []Content

func (con Sequence) IsSentence() bool {
	for _, c := range con {
		if !c.IsSentence() {
			return false
		}
	}

	return true
}

func (con Sequence) String() string {
	str := ""
	for _, content := range con {
		str += content.String()
	}

	return str
}

func (con Sequence) Contents() []Content {
	return con
}

func (con Sequence) Visit(visitor Visitor) error {
	return visitor.VisitSequence(con)
}
