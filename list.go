package booklit

import (
	"fmt"
	"strings"
)

// List is a block content containing a list of content, either ordered or
// unordered.
type List struct {
	// The items in the list.
	Items []Content

	// Whether the list is ordered.
	Ordered bool
}

func (con List) IsFlow() bool {
	return false
}

func (con List) String() string {
	var str string
	for i, c := range con.Items {
		var text string
		for _, line := range strings.Split(strings.TrimRight(c.String(), "\n"), "\n") {
			if len(text) == 0 {
				text = line
			} else if len(line) == 0 {
				text += "\n"
			} else {
				text += "\n  " + line
			}
		}

		if con.Ordered {
			str += fmt.Sprintf("%d. %s\n\n", i+1, text)
		} else {
			str += fmt.Sprintf("* %s\n\n", text)
		}
	}

	return str
}

func (con List) Visit(visitor Visitor) error {
	return visitor.VisitList(con)
}
