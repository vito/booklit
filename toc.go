package booklit

// TableOfContents is block content which renders a table of contents for the
// section and its children.
type TableOfContents struct {
	Section *Section
}

func (con TableOfContents) IsFlow() bool {
	return false
}

func (con TableOfContents) String() string {
	return ""
}

func (con TableOfContents) Visit(visitor Visitor) error {
	return visitor.VisitTableOfContents(con)
}
