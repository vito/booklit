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

// IsFlow returns true.
func (con Target) IsFlow() bool {
	return true
}

// String returns an empty string.
//
// XXX: maybe this should summarize it, and the search index should use
// render.TextEngine isntead of String
func (con Target) String() string {
	return ""
}

// Visit calls VisitTarget.
func (con Target) Visit(visitor Visitor) error {
	return visitor.VisitTarget(con)
}
