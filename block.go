package booklit

import "fmt"

type Block struct {
	Content Content

	Class string
}

func (con Block) IsFlow() bool {
	return false
}

func (con Block) String() string {
	return fmt.Sprintf("{block (%s): %s}", con.Class, con.Content)
}

func (con Block) Visit(visitor Visitor) error {
	return visitor.VisitBlock(con)
}
