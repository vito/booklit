package booklit

import (
	"fmt"
	"strings"
)

type List struct {
	Items []Content

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
