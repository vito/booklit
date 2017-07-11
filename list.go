package booklit

import "fmt"

type List struct {
	Items []Content

	Ordered bool
}

func (con List) IsSentence() bool {
	return false
}

func (con List) String() string {
	return fmt.Sprintf("{list: %s}", con.Items)
}

func (con List) Visit(visitor Visitor) error {
	return visitor.VisitList(con)
}
