package booklit

type Paragraph []Content

func (con Paragraph) String() string {
	str := ""

	for _, seq := range con {
		if len(str) > 0 {
			str += " "
		}

		str += seq.String()
	}

	if len(str) > 0 {
		str += "\n\n"
	}

	return str
}

func (con Paragraph) IsFlow() bool {
	return false
}

func (con Paragraph) Visit(visitor Visitor) error {
	return visitor.VisitParagraph(con)
}
