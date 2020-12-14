package booklit

import "fmt"

// Definitions is a list of definitions, e.g. a glossary.
type Definitions []Definition

// Definition is a subject and its definition.
type Definition struct {
	Subject    Content
	Definition Content
}

func (con Definitions) IsFlow() bool {
	return false
}

func (con Definitions) String() string {
	var text string
	for _, def := range con {
		text += fmt.Sprintf("%s: %s\n", def.Subject, def.Definition)
	}

	text += "\n"

	return text
}

func (con Definitions) Visit(visitor Visitor) error {
	return visitor.VisitDefinitions(con)
}
