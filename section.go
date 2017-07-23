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

	Partials Partials

	SplitSections        bool
	PreventSplitSections bool

	ResetDepth bool

	OmitChildrenFromTableOfContents bool

	PluginFactories []PluginFactory
	Plugins         []Plugin
}

type Partials map[string]Content

type Tag struct {
	Name    string
	Display Content

	Section *Section
	Anchor  string
}

func (con *Section) String() string {
	return fmt.Sprintf("{section (%s): %s}", con.Path, con.Title)
}

func (con *Section) IsFlow() bool {
	return false
}

func (con *Section) Visit(visitor Visitor) error {
	return visitor.VisitSection(con)
}

func (con *Section) SetTitle(title Content, tags ...string) {
	if len(tags) == 0 {
		tags = []string{con.defaultTag(title)}
	}

	con.Tags = []Tag{}
	for _, name := range tags {
		con.SetTag(name, title)
	}

	con.Title = title
	con.PrimaryTag = con.Tags[0]
}

func (con *Section) SetTag(name string, display Content, optionalAnchor ...string) {
	anchor := ""
	if len(optionalAnchor) > 0 {
		anchor = optionalAnchor[0]
	}

	con.Tags = append(con.Tags, Tag{
		Section: con,

		Name:    name,
		Display: display,
		Anchor:  anchor,
	})
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

func (con *Section) SetPartial(name string, value Content) {
	if con.Partials == nil {
		con.Partials = Partials{}
	}

	con.Partials[name] = value
}

func (con *Section) Partial(name string) Content {
	return con.Partials[name]
}

func (con *Section) UsePlugin(pf PluginFactory) {
	con.PluginFactories = append(con.PluginFactories, pf)
	con.Plugins = append(con.Plugins, pf(con))
}

func (con *Section) PageDepth() int {
	if con.Parent == nil || con.Parent.ResetDepth {
		return 0
	}

	return con.Parent.PageDepth() + 1
}

func (con *Section) SplitSectionsPrevented() bool {
	if con.PreventSplitSections {
		return true
	}

	if con.Parent != nil && con.Parent.SplitSectionsPrevented() {
		return true
	}

	return false
}

func (con *Section) findTag(tagName string, up bool, exclude *Section) (Tag, bool) {
	if tagName == con.Title.String() {
		return con.PrimaryTag, true
	}

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
					StripAux(title).String(),
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
