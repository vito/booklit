package booklit

type Style string

const (
	StyleVerbatim Style = "verbatim"
)

type Styled struct {
	Content

	Style Style
}

func (con Styled) Visit(visitor Visitor) error {
	return visitor.VisitStyled(con)
}
