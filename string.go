package booklit

type String string

var Empty String

func (str String) IsFlow() bool {
	return true
}

func (str String) String() string {
	return string(str)
}

func (str String) Visit(visitor Visitor) error {
	return visitor.VisitString(str)
}
