package booklit

type Content interface {
	Visit(Visitor)
}

type Visitor interface {
	VisitString(String)
	VisitSequence(Sequence)
	VisitSection(*Section)
}
