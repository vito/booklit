package booklit

// TableOfContents is block content which renders a table of contents for the
// section and its children.
type TableOfContents struct {
	Section *Section
}

// IsFlow returns false.
func (con TableOfContents) IsFlow() bool {
	return false
}

// String returns an empty string.
func (con TableOfContents) String() string {
	return ""
}

// Visit calls VisitTableOfContents.
func (con TableOfContents) Visit(visitor Visitor) error {
	return visitor.VisitTableOfContents(con)
}
