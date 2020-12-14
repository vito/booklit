package booklit

import "github.com/vito/booklit/ast"

// Target is flow content which creates a tag within the section and renders an
// anchor element for the tag to target.
type Target struct {
	TagName  string
	Location ast.Location
	Title    Content
	Content  Content
}

func (con Target) IsFlow() bool {
	return true
}

func (con Target) String() string {
	return ""
}

func (con Target) Visit(visitor Visitor) error {
	return visitor.VisitTarget(con)
}
