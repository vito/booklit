package booklit

// Link is flow content that references something typically external to the
// Booklit content, such as another website.
type Link struct {
	// Content to display as the link.
	Content

	// Target (e.g. a URL) that the link points to.
	Target string
}

// Visit calls VisitLink.
func (con Link) Visit(visitor Visitor) error {
	return visitor.VisitLink(con)
}
