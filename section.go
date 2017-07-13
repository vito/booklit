package booklit

import (
	"fmt"
	"regexp"
	"strings"
)

type Section struct {
	Path string

	Title Content
	Body  Content

	PrimaryTag Tag
	Tags       []Tag

	Parent   *Section
	Children []*Section

	SplitSections                   bool
	OmitChildrenFromTableOfContents bool

	PluginFactories []PluginFactory
	Plugins         []Plugin
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

	con.Tags = []Tag{}
	for _, name := range tags {
		con.Tags = append(con.Tags, Tag{
			Name:    name,
			Display: title,

			Section: con,
		})
	}

	con.Title = title
	con.PrimaryTag = con.Tags[0]
}

func (con *Section) Number() string {
	if con.Parent == nil {
		return ""
	}

	parentNumber := con.Parent.Number()
	selfIndex := 1
	for i := 0; con.Parent.Children[i] != con; i++ {
		selfIndex++
	}

	if parentNumber == "" {
		return fmt.Sprintf("%d", selfIndex)
	}

	return fmt.Sprintf("%s.%d", parentNumber, selfIndex)
}

func (con *Section) HasAnchors() bool {
	for _, tag := range con.Tags {
		if tag.Anchor != "" {
			return true
		}
	}

	if con.SplitSections {
		return false
	}

	for _, child := range con.Children {
		if child.HasAnchors() {
			return true
		}
	}

	return false
}

func (con *Section) AnchorTags() []Tag {
	tags := []Tag{}

	for _, tag := range con.Tags {
		if tag.Anchor == "" {
			continue
		}

		tags = append(tags, tag)
	}

	return tags
}

func (con *Section) Top() *Section {
	if con.Parent != nil {
		return con.Parent.Top()
	}

	return con
}

func (con *Section) Contains(sub *Section) bool {
	if con == sub {
		return true
	}

	for _, child := range con.Children {
		if child.Contains(sub) {
			return true
		}
	}

	return false
}

func (con *Section) FindTag(tagName string) (Tag, bool) {
	return con.findTag(tagName, true, nil)
}

func (con *Section) UsePlugin(pf PluginFactory) {
	con.PluginFactories = append(con.PluginFactories, pf)
	con.Plugins = append(con.Plugins, pf.NewPlugin(con))
}

func (con *Section) findTag(tagName string, up bool, exclude *Section) (Tag, bool) {
	for _, t := range con.Tags {
		if t.Name == tagName {
			return t, true
		}
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
