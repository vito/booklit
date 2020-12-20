package booklit

// Sequence is a generic slice of content which will be concatenated together
// upon rendering.
type Sequence []Content

// IsFlow returns true if the sequence contains only flow content or is empty.
func (con Sequence) IsFlow() bool {
	for _, c := range con {
		if !c.IsFlow() {
			return false
		}
	}

	return true
}

// String summarizes the content for debugging purposes.
func (con Sequence) String() string {
	str := ""
	for _, content := range con {
		str += content.String()
	}

	return str
}

// Contents returns the content as a slice.
func (con Sequence) Contents() []Content {
	return con
}

// Visit calls VisitSequence.
func (con Sequence) Visit(visitor Visitor) error {
	return visitor.VisitSequence(con)
}
