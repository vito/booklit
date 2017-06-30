package booklit

import (
	"fmt"
	"regexp"
	"strings"
)

type Section struct {
	Title Content
	Body  Content

	Tags []string

	Parent   *Section
	Children []*Section
}

func (con *Section) String() string {
	return fmt.Sprintf("{section: %s}", con.Title)
}

func (con *Section) IsSentence() bool { return false }

func (con *Section) Visit(visitor Visitor) {
	visitor.VisitSection(con)
}

var whitespaceRegexp = regexp.MustCompile(`\s+`)
var specialCharsRegexp = regexp.MustCompile(`[^[:alnum:]_\-]`)

func (con *Section) PrimaryTag() string {
	if len(con.Tags) > 0 {
		return con.Tags[0]
	}

	return strings.ToLower(
		specialCharsRegexp.ReplaceAllString(
			whitespaceRegexp.ReplaceAllString(
				strings.Replace(
					con.Title.String(),
					" & ",
					" and ",
					-1,
				),
				"-",
			),
			"",
		),
	)
}
