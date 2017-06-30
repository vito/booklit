package booklit

type Paragraph []Content

func (con Paragraph) String() string {
	str := ""

	for _, seq := range con {
		str += seq.String()
	}

	return str
}

func (con Paragraph) IsSentence() bool {
	return false
}

func (con Paragraph) Visit(visitor Visitor) {
	visitor.VisitParagraph(con)
}
