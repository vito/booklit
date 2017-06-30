package booklit

type Content interface {
	String() string

	Visit(Visitor)
}

type Visitor interface {
	VisitString(String)
	VisitSequence(Sequence)
	VisitSection(*Section)
}
