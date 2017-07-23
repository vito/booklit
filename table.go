package booklit

import "fmt"

type Table struct {
	Rows [][]Content
}

func (con Table) IsFlow() bool {
	return false
}

func (con Table) String() string {
	return fmt.Sprintf("{table: %s}", con.Rows)
}

func (con Table) Visit(visitor Visitor) error {
	return visitor.VisitTable(con)
}
