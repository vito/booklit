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
	return fmt.Sprintf("{definitions: %s}", []Definition(con))
}

func (con Definitions) Visit(visitor Visitor) error {
	return visitor.VisitDefinitions(con)
}
