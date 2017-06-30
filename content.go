package booklit

type Content interface {
	String() string

	IsSentence() bool

	Visit(Visitor)
}

type Visitor interface {
	VisitString(String)
	VisitSequence(Sequence)
	VisitSection(*Section)
	VisitParagraph(Paragraph)
}
