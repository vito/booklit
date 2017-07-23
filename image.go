package booklit

import "fmt"

type Image struct {
	Path        string
	Description string
}

func (con Image) IsFlow() bool {
	return true
}

func (con Image) String() string {
	return fmt.Sprintf("{image: %s}", con.Path)
}

func (con Image) Visit(visitor Visitor) error {
	return visitor.VisitImage(con)
}
