package booklit

// Preformatted is block content representing preformatted text, e.g. a code
// block.
type Preformatted []Content

func (con Preformatted) String() string {
	str := ""

	for _, seq := range con {
		str += seq.String()
		str += "\n"
	}

	return str
}

func (con Preformatted) IsFlow() bool {
	return false
}

func (con Preformatted) Visit(visitor Visitor) error {
	return visitor.VisitPreformatted(con)
}
