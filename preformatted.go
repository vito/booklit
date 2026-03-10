package booklit

import "strings"

// Preformatted is block content representing preformatted text, e.g. a code
// block.
type Preformatted []Content

// String summarizes the content for debugging purposes.
func (con Preformatted) String() string {
	var str strings.Builder

	for _, seq := range con {
		str.WriteString(seq.String())
		str.WriteString("\n")
	}

	return str.String()
}

// IsFlow returns false.
func (con Preformatted) IsFlow() bool {
	return false
}

// Visit calls VisitPreformatted.
func (con Preformatted) Visit(visitor Visitor) error {
	return visitor.VisitPreformatted(con)
}
