package booklit

import "fmt"

type TableOfContents struct {
	Section *Section
}

func (con TableOfContents) IsFlow() bool {
	return false
}

func (con TableOfContents) String() string {
	return fmt.Sprintf("{toc: %s}", con.Section)
}

func (con TableOfContents) Visit(visitor Visitor) error {
	return visitor.VisitTableOfContents(con)
}
