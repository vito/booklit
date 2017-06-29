package booklit

type String string

func (str String) String() string {
	return string(str)
}

func (str String) Visit(visitor Visitor) {
	visitor.VisitString(str)
}
