package booklit

type Content interface {
	String() string

	IsSentence() bool

	Visit(Visitor) error
}

type Visitor interface {
	VisitString(String) error
	VisitSequence(Sequence) error
	VisitSection(*Section) error
	VisitParagraph(Paragraph) error
}
