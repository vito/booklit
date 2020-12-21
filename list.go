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

// IsFlow returns false.
func (con List) IsFlow() bool {
	return false
}

// String summarizes the content for debugging purposes.
func (con List) String() string {
	var str string
	for i, c := range con.Items {
		var text string
		for _, line := range strings.Split(strings.TrimRight(c.String(), "\n"), "\n") {
			if text == "" {
				text = line
			} else if line == "" {
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

// Visit calls VisitList.
func (con List) Visit(visitor Visitor) error {
	return visitor.VisitList(con)
}
