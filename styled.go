package booklit

// Styled allows Content to be rendered with custom templates.
type Styled struct {
	// A string identifying the template name.
	Style Style

	// Block may be set to true to force otherwise flow content to be block
	// instead.
	Block bool

	// The content to render with the template.
	Content Content

	// Additional partials to pass to the template.
	Partials Partials
}

// Style identifies a template name.
type Style string

// Common styled templated names.
const (
	StyleVerbatim    Style = "verbatim"
	StyleItalic      Style = "italic"
	StyleBold        Style = "bold"
	StyleLarger      Style = "larger"
	StyleSmaller     Style = "smaller"
	StyleStrike      Style = "strike"
	StyleSuperscript Style = "superscript"
	StyleSubscript   Style = "subscript"
	StyleInset       Style = "inset"
	StyleAside       Style = "aside"
)

// String summarizes the content for debugging purposes.
func (con Styled) String() string {
	return con.Content.String()
}

// IsFlow returns false if Block is true and otherwise delegates to
// content.IsFlow.
func (con Styled) IsFlow() bool {
	if con.Block {
		return false
	}

	return con.Content.IsFlow()
}

// Visit calls VisitStyled.
func (con Styled) Visit(visitor Visitor) error {
	return visitor.VisitStyled(con)
}

// Partial returns the given partial by name, or nil if it does not exist.
func (con Styled) Partial(name string) Content {
	return con.Partials[name]
}
