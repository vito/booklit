package booklit

import "strings"

import "fmt"

// Definitions is a list of definitions, e.g. a glossary.
type Definitions []Definition

// Definition is a subject and its definition.
type Definition struct {
	Subject    Content
	Definition Content
}

// IsFlow returns false.
func (con Definitions) IsFlow() bool {
	return false
}

// String summarizes the content for debugging purposes.
func (con Definitions) String() string {
	var text strings.Builder
	for _, def := range con {
		text.WriteString(fmt.Sprintf("%s: %s\n", def.Subject, def.Definition))
	}

	text.WriteString("\n")

	return text.String()
}

// Visit calls VisitDefinitions.
func (con Definitions) Visit(visitor Visitor) error {
	return visitor.VisitDefinitions(con)
}
