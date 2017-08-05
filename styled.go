package booklit

type Styled struct {
	Style Style
	Block bool

	Content  Content
	Partials Partials
}

type Style string

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

func (con Styled) String() string {
	return con.Content.String()
}

func (con Styled) IsFlow() bool {
	if con.Block {
		return false
	}

	return con.Content.IsFlow()
}

func (con Styled) Visit(visitor Visitor) error {
	return visitor.VisitStyled(con)
}

func (con Styled) Partial(name string) Content {
	return con.Partials[name]
}
