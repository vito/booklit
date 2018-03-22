package booklit

import "fmt"

type Definitions []Definition

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
