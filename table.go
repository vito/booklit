package booklit

type Table struct {
	Rows [][]Content
}

func (con Table) IsFlow() bool {
	return false
}

func (con Table) String() string {
	var text string
	for _, cols := range con.Rows {
		row := "|"
		for _, col := range cols {
			row += " " + col.String() + " |"
		}

		text += row + "\n"
	}

	text += "\n"

	return text
}

func (con Table) Visit(visitor Visitor) error {
	return visitor.VisitTable(con)
}
