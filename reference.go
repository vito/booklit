package booklit

import "fmt"

type Reference struct {
	TagName string
	Content Content

	Tag *Tag
}

func (con *Reference) IsFlow() bool {
	return true
}

func (con *Reference) String() string {
	if con.Content != nil {
		return con.Content.String()
	}

	if con.Tag != nil {
		return con.Tag.Title.String()
	}

	return fmt.Sprintf("{reference: %s}", con.TagName)
}

func (con *Reference) Visit(visitor Visitor) error {
	return visitor.VisitReference(con)
}

func (con *Reference) Display() Content {
	if con.Content != nil {
		return con.Content
	}

	return con.Tag.Title
}
