package booklit

type Section struct {
	Title Content
	Body  Content

	Tags []string

	Parent   *Section
	Children []*Section
}

func (sec *Section) Visit(visitor Visitor) {
	visitor.VisitSection(sec)
}
