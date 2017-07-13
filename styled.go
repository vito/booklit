package booklit

type Styled struct {
	Content

	Style Style
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
	StyleNote        Style = "note"
)

func (con Styled) Visit(visitor Visitor) error {
	return visitor.VisitStyled(con)
}
