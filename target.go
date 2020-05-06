package booklit

import "github.com/vito/booklit/ast"

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
