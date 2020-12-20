package booklit

// Paragraph is block content containing flow content as sentences.
//
// When rendered, the sentences are joined by a single space in between and a
// blank line following the paragraph.
type Paragraph []Content

// String summarizes the content for debugging purposes.
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

// IsFlow returns true.
func (con Paragraph) IsFlow() bool {
	return false
}

// Visit calls VisitParagraph.
func (con Paragraph) Visit(visitor Visitor) error {
	return visitor.VisitParagraph(con)
}
