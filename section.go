package booklit

import (
	"fmt"
	"regexp"
	"strings"
)

type Section struct {
	Title Content
	Body  Content

	PrimaryTag Tag
	Tags       map[string]Tag

	Parent   *Section
	Children []*Section

	SplitSections bool
}

type Tag struct {
	Name    string
	Display Content

	Section *Section
	Anchor  string
}

func (con *Section) String() string {
	return fmt.Sprintf("{section: %s}", con.Title)
}

func (con *Section) IsSentence() bool { return false }

func (con *Section) Visit(visitor Visitor) error {
	return visitor.VisitSection(con)
}

func (con *Section) SetTitle(title Content, tags ...string) {
	if len(tags) == 0 {
		tags = []string{con.defaultTag(title)}
	}

	con.Tags = map[string]Tag{}
	for _, name := range tags {
		con.Tags[name] = Tag{
			Name:    name,
			Display: title,

			Section: con,
		}
	}

	con.Title = title
	con.PrimaryTag = con.Tags[tags[0]]
}

func (con *Section) FindTag(tagName string) (Tag, bool) {
	return con.findTag(tagName, true, nil)
}

func (con *Section) findTag(tagName string, up bool, exclude *Section) (Tag, bool) {
	tag, found := con.Tags[tagName]
	if found {
		return tag, true
	}

	for _, sub := range con.Children {
		if sub != exclude {
			tag, found := sub.findTag(tagName, false, nil)
			if found {
				return tag, true
			}
		}
	}

	if up && con.Parent != nil {
		return con.Parent.findTag(tagName, true, con)
	}

	return Tag{}, false
}

var whitespaceRegexp = regexp.MustCompile(`\s+`)
var specialCharsRegexp = regexp.MustCompile(`[^[:alnum:]_\-]`)

func (con *Section) defaultTag(title Content) string {
	return strings.ToLower(
		specialCharsRegexp.ReplaceAllString(
			whitespaceRegexp.ReplaceAllString(
				strings.Replace(
					title.String(),
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
