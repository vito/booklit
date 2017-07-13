package booklit

import "fmt"

type Element struct {
	Content Content

	Class string
}

func (con Element) IsSentence() bool {
	return true
}

func (con Element) String() string {
	return fmt.Sprintf("{element (%s): %s}", con.Class, con.Content)
}

func (con Element) Visit(visitor Visitor) error {
	return visitor.VisitElement(con)
}
