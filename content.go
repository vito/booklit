package booklit

type Content interface {
	String() string

	IsFlow() bool

	Visit(Visitor) error
}

type Visitor interface {
	VisitString(String) error
	VisitSequence(Sequence) error
	VisitReference(*Reference) error
	VisitLink(Link) error
	VisitSection(*Section) error
	VisitParagraph(Paragraph) error
	VisitTableOfContents(TableOfContents) error
	VisitPreformatted(Preformatted) error
	VisitStyled(Styled) error
	VisitTarget(Target) error
	VisitImage(Image) error
	VisitList(List) error
	VisitTable(Table) error
	VisitDefinitions(Definitions) error
}
