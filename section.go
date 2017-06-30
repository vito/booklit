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

func (sec *Section) String() string {
	return fmt.Sprintf("{section: %s}", sec.Title)
}

func (sec *Section) Visit(visitor Visitor) {
	visitor.VisitSection(sec)
}

var whitespaceRegexp = regexp.MustCompile(`\s+`)
var specialCharsRegexp = regexp.MustCompile(`[^[:alnum:]_\-]`)

func (sec *Section) PrimaryTag() string {
	if len(sec.Tags) > 0 {
		return sec.Tags[0]
	}

	return strings.ToLower(
		specialCharsRegexp.ReplaceAllString(
			whitespaceRegexp.ReplaceAllString(
				strings.Replace(
					sec.Title.String(),
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
