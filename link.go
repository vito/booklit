package booklit

type Link struct {
	Content

	Target string
}

func (con Link) Visit(visitor Visitor) error {
	return visitor.VisitLink(con)
}
