package booklit

type Styled struct {
	Content

	Style Style
}

type Style string

const (
	StyleVerbatim Style = "verbatim"
	StyleItalic   Style = "italic"
	StyleBold     Style = "bold"
)

func (con Styled) Visit(visitor Visitor) error {
	return visitor.VisitStyled(con)
}
