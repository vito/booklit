package booklit

import "fmt"

type Target struct {
	TagName string
	Display Content
}

func (con Target) IsFlow() bool {
	return true
}

func (con Target) String() string {
	return fmt.Sprintf("{target: %s}", con.TagName)
}

func (con Target) Visit(visitor Visitor) error {
	return visitor.VisitTarget(con)
}
