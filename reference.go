package booklit

import (
	"fmt"

	"github.com/vito/booklit/ast"
)

// Reference is flow content linking to a tag defined elsewhere.
type Reference struct {
	// The tag to link to.
	TagName string

	// Optional content to display for the reference. If not present, the tag's
	// own display content will be used.
	Content Content

	// The tag that the name resolved to in the "resolving" phase.
	Tag *Tag

	// The original source location of the reference. Used when generating error
	// messages.
	Location ast.Location
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
