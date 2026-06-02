package booklit

// Content is arbitrary content (e.g. text, links, paragraphs) created by
// evaluating a Booklit document or created by a plugin.
type Content interface {
	// String summarizes the content. It is only used for troubleshooting.
	String() string

	// IsFlow must return true if the content is 'flow' content (e.g. anything
	// that fits within a sentence) or false if the content is 'block' content
	// (e.g. a paragraph or table).
	IsFlow() bool

	// Visit calls the VisitX method on Visitor corresponding to the Content's
	// type.
	Visit(Visitor) error
}

// Visitor is implemented in order to traverse Content.
type Visitor interface {
	VisitString(String) error
	VisitSequence(Sequence) error
	VisitReference(*Reference) error
	VisitSection(*Section) error
	VisitParagraph(Paragraph) error
	VisitTableOfContents(TableOfContents) error
	VisitTarget(Target) error
	VisitLazy(*Lazy) error
	VisitRawElement(RawElement) error
	VisitRawFragment(RawFragment) error
}
