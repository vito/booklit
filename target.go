package booklit

type Target struct {
	TagName string
	Title   Content
	Content Content
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
