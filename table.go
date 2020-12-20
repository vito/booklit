package booklit

// Table is block content containing tabular data, i.e. rows and columns.
type Table struct {
	Rows [][]Content
}

// IsFlow returns false.
func (con Table) IsFlow() bool {
	return false
}

// String summarizes the content for debugging purposes.
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

// Visit calls VisitTable.
func (con Table) Visit(visitor Visitor) error {
	return visitor.VisitTable(con)
}
