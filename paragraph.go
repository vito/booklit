package booklit

type Paragraph struct {
	Content
}

func (con Paragraph) Visit(visitor Visitor) {
	visitor.VisitParagraph(con)
}
