package booklit

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/agext/levenshtein"
	"github.com/vito/booklit/ast"
)

// Section is a Booklit document, typically loaded from a .lit file.
type Section struct {
	// The file path that the section was loaded from.
	Path string

	// The section title and body.
	Title Content
	Body  Content

	// The the file source location where the title was set.
	TitleLocation ast.Location

	// The primary tag and additional tags for the section.
	PrimaryTag Tag
	Tags       []Tag

	// The section's parent and child sections, if any.
	Parent   *Section
	Children []*Section

	// The section's style. Used at the rendering stage to e.g. use different
	// templates.
	//
	// Set with \style{foo}.
	Style string

	// Arbitrary named content.
	//
	// Partials are typically rendered in templates via {{.Partial "Foo" | render}}.
	//
	// Set with \set-partial{name}{content}.
	Partials Partials

	// Instructs the renderer to render child sections on their own pages.
	//
	// Set with \split-sections.
	SplitSections bool

	// Instructs the renderer to never render child sections on their own pages,
	// even if they set SplitSections.
	//
	// Typically used to have a single-page view or printout of the content.
	//
	// Set with \single-page.
	PreventSplitSections bool

	// Instructs .PageDepth to pretend the section is the lowest point.
	//
	// Set with \split-sections, even if PreventSplitSections is true, to ensure
	// section numbers remain consistent.
	ResetDepth bool

	// Omit child sections from table-of-contents lists.
	//
	// Set with \omit-children-from-table-of-contents.
	OmitChildrenFromTableOfContents bool

	// Plugins used by the section.
	PluginFactories []PluginFactory
	Plugins         []Plugin

	// Location is the source location where the section was defined, if it was
	// defined as an inline section.
	Location ast.Location

	// InvokeLocation is set prior to each function call so that the plugin's
	// method can assign it on content that it produces, e.g. Reference and
	// Target, so that error messages can be annotated with the source of the
	// error.
	//
	// XXX: make this an optional argument instead?
	InvokeLocation ast.Location

	// Processor is used for evaluating additional child sections.
	Processor SectionProcessor
}

// Partials is a map of named snippets of content. By convention, partial names
// are camel-cased like FooBar.
type Partials map[string]Content

// SectionProcessor evaluates a file or an inline syntax node to produce a
// child section.
type SectionProcessor interface {
	EvaluateFile(*Section, string, []PluginFactory) (*Section, error)
	EvaluateNode(*Section, ast.Node, []PluginFactory) (*Section, error)
}

// Tag is something which can be referenced (by its name) from elsewhere in the
// section tree to produce a link.
type Tag struct {
	// The name of the tag, to be referenced with \reference.
	Name string

	// The title of the tag to display when no display is given by \reference.
	Title Content

	// The section the tag resides in.
	Section *Section

	// An optional anchor for the tag.
	//
	// Anchored tags correspond to page anchors and are typically displayed in
	// the section body, as opposed to tags corresponding to a section.
	Anchor string

	// Content that the tag corresponds to. For example, the section body or the
	// text that an anchored tag is for.
	//
	// Typically used for showing content previews in search results.
	Content Content

	// The source location of the tag's definition.
	Location ast.Location
}

// String summarizes the content for debugging purposes.
func (con *Section) String() string {
	return fmt.Sprintf("{section (%s): %s}", con.Path, con.Title)
}

// FilePath is the file that the section is contained in.
func (con *Section) FilePath() string {
	if con.Path != "" {
		return con.Path
	}

	if con.Parent != nil {
		return con.Parent.FilePath()
	}

	return ""
}

// IsFlow returns false.
func (con *Section) IsFlow() bool {
	return false
}

// Visit calls VisitSection on visitor.
func (con *Section) Visit(visitor Visitor) error {
	return visitor.VisitSection(con)
}

// SetTitle sets the section title and tags, setting a default tag based on the
// title if no tags are specified.
func (con *Section) SetTitle(title Content, loc ast.Location, tags ...string) {
	if len(tags) == 0 {
		tags = []string{con.defaultTag(title)}
	}

	con.Tags = []Tag{}
	for _, name := range tags {
		con.SetTag(name, title, loc)
	}

	con.Title = title
	con.TitleLocation = loc
	con.PrimaryTag = con.Tags[0]
}

// SetTag adds a tag to the section.
func (con *Section) SetTag(name string, title Content, loc ast.Location, optionalAnchor ...string) {
	anchor := ""
	if len(optionalAnchor) > 0 {
		anchor = optionalAnchor[0]
	}

	con.Tags = append(con.Tags, Tag{
		Section:  con,
		Location: loc,

		Name:   name,
		Title:  title,
		Anchor: anchor,
	})
}

// SetTagAnchored adds an anchored tag to the section, along with content
// associated to the anchor.
func (con *Section) SetTagAnchored(name string, title Content, loc ast.Location, content Content, anchor string) {
	con.Tags = append(con.Tags, Tag{
		Section:  con,
		Location: loc,

		Name:   name,
		Title:  title,
		Anchor: anchor,

		Content: content,
	})
}

// Number returns a string denoting the section's unique numbering for use in
// titles and tables of contents, e.g. "3.2".
func (con *Section) Number() string {
	if con.Parent == nil {
		return ""
	}

	parentNumber := con.Parent.Number()
	selfIndex := 1
	for _, child := range con.Parent.Children {
		if child == con {
			break
		}

		selfIndex++
	}

	if parentNumber == "" {
		return fmt.Sprintf("%d", selfIndex)
	}

	return fmt.Sprintf("%s.%d", parentNumber, selfIndex)
}

// HasAnchors returns true if the section has any anchored tags or if any
// inline child sections have anchors.
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

// AnchorTags returns the section's tags which have anchors.
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

// Top returns the top-level section.
func (con *Section) Top() *Section {
	if con.Parent != nil {
		return con.Parent.Top()
	}

	return con
}

// Contains returns true if the section is sub or if any of the children
// Contains sub.
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

// IsOrHasChild returns true if the section is sub or has sub as a direct
// child.
func (con *Section) IsOrHasChild(sub *Section) bool {
	if con == sub {
		return true
	}

	for _, child := range con.Children {
		if child == sub {
			return true
		}
	}

	return false
}

// Prev returns the previous section, i.e. the previous sibling section or the
// parent section if it is the first child.
//
// If the section has no parent, Prev returns nil.
func (con *Section) Prev() *Section {
	if con.Parent == nil {
		return nil
	}

	var lastChild *Section
	for _, child := range con.Parent.Children {
		if lastChild != nil && child == con {
			return lastChild
		}

		lastChild = child
	}

	return con.Parent
}

// Next returns the next section, i.e. if SplitSections return the first child,
// otherwise return the NextSibling.
func (con *Section) Next() *Section {
	if con.SplitSections {
		if len(con.Children) > 0 {
			return con.Children[0]
		}
	}

	return con.NextSibling()
}

// NextSibling returns the section after the current section in the parent's
// children.
//
// If there is no section after this section, the parent's NextSibling is
// returned.
func (con *Section) NextSibling() *Section {
	if con.Parent == nil {
		return nil
	}

	var sawSelf bool
	for _, child := range con.Parent.Children {
		if sawSelf {
			return child
		}

		if child == con {
			sawSelf = true
		}
	}

	return con.Parent.NextSibling()
}

// FindTag searches for a given tag name and returns all matching tags.
func (con *Section) FindTag(tagName string) []Tag {
	return con.filterTags(true, nil, func(other string) bool {
		return other == tagName
	})
}

// SimilarTags searches for a given tag name and returns all similar tags, i.e.
// tags with a levenshtein distance > 0.5.
func (con *Section) SimilarTags(tagName string) []Tag {
	return con.filterTags(true, nil, func(other string) bool {
		return levenshtein.Match(tagName, other, nil) > 0.5
	})
}

// SetPartial assigns a partial within the section.
func (con *Section) SetPartial(name string, value Content) {
	if con.Partials == nil {
		con.Partials = Partials{}
	}

	con.Partials[name] = value
}

// Partial returns the given partial, or nil if it does not exist.
func (con *Section) Partial(name string) Content {
	return con.Partials[name]
}

// UsePlugin registers the plugin within the section.
func (con *Section) UsePlugin(pf PluginFactory) {
	con.PluginFactories = append(con.PluginFactories, pf)
	con.Plugins = append(con.Plugins, pf(con))
}

// Depth returns the absolute depth of the section.
func (con *Section) Depth() int {
	if con.Parent == nil {
		return 0
	}

	return con.Parent.Depth() + 1
}

// PageDepth returns the depth within the page that the section will be
// rendered on, i.e. accounting for ResetDepth being set on the parent section.
func (con *Section) PageDepth() int {
	if con.Parent == nil || con.Parent.ResetDepth {
		return 0
	}

	return con.Parent.PageDepth() + 1
}

// SplitSectionsPrevented returns true if PreventSplitSections is true or if
// the parent has SplitSectionsPrevented.
func (con *Section) SplitSectionsPrevented() bool {
	if con.PreventSplitSections {
		return true
	}

	if con.Parent != nil && con.Parent.SplitSectionsPrevented() {
		return true
	}

	return false
}

func (con *Section) filterTags(up bool, exclude *Section, match func(string) bool) []Tag {
	tags := []Tag{}

	for _, t := range con.Tags {
		if match(t.Name) {
			tags = append(tags, t)
		}
	}

	for _, sub := range con.Children {
		if sub != exclude {
			tags = append(tags, sub.filterTags(false, nil, match)...)
		}
	}

	if up && con.Parent != nil {
		tags = append(tags, con.Parent.filterTags(true, con, match)...)
	}

	return tags
}

var whitespaceRegexp = regexp.MustCompile(`\s+`)
var specialCharsRegexp = regexp.MustCompile(`[^[:alnum:]_\-]`)

func (con *Section) defaultTag(title Content) string {
	return strings.ToLower(
		specialCharsRegexp.ReplaceAllString(
			whitespaceRegexp.ReplaceAllString(
				strings.ReplaceAll(
					StripAux(title).String(),
					" & ",
					" and ",
				),
				"-",
			),
			"",
		),
	)
}
