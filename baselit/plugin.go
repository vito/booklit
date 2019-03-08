package baselit

import (
	"fmt"
	"path/filepath"

	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
)

func init() {
	booklit.RegisterPlugin("base", NewPlugin)
}

func NewPlugin(section *booklit.Section) booklit.Plugin {
	return Plugin{
		section: section,
	}
}

type Plugin struct {
	section *booklit.Section
}

func (plugin Plugin) UsePlugin(name string) error {
	pluginFactory, found := booklit.LookupPlugin(name)
	if !found {
		return fmt.Errorf("unknown plugin '%s'", name)
	}

	plugin.section.UsePlugin(pluginFactory)

	return nil
}

func (plugin Plugin) Styled(name string) {
	plugin.section.Style = name
}

func (plugin Plugin) Title(title booklit.Content, tags ...string) {
	plugin.section.SetTitle(title, tags...)
}

func (plugin Plugin) Aux(content booklit.Content) booklit.Content {
	return booklit.Aux{content}
}

func (plugin Plugin) Section(node ast.Node) error {
	section, err := plugin.section.Processor.EvaluateNode(plugin.section, node, plugin.section.PluginFactories)
	if err != nil {
		return err
	}

	plugin.section.Children = append(plugin.section.Children, section)

	return nil
}

func (plugin Plugin) IncludeSection(path string) error {
	sectionPath := filepath.Join(filepath.Dir(plugin.section.Path), path)

	section, err := plugin.section.Processor.EvaluateFile(plugin.section, sectionPath, []booklit.PluginFactory{NewPlugin})
	if err != nil {
		return err
	}

	plugin.section.Children = append(plugin.section.Children, section)

	return nil
}

func (plugin Plugin) SinglePage() {
	plugin.section.PreventSplitSections = true
}

func (plugin Plugin) SplitSections() {
	plugin.section.ResetDepth = true

	if !plugin.section.SplitSectionsPrevented() {
		plugin.section.SplitSections = true
	}
}

func (plugin Plugin) OmitChildrenFromTableOfContents() {
	plugin.section.OmitChildrenFromTableOfContents = true
}

func (plugin Plugin) TableOfContents() booklit.Content {
	return booklit.TableOfContents{
		Section: plugin.section,
	}
}

func (plugin Plugin) Code(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Content: content,
		Style:   booklit.StyleVerbatim,
	}
}

func (plugin Plugin) Italic(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Content: content,
		Style:   booklit.StyleItalic,
	}
}

func (plugin Plugin) Bold(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Content: content,
		Style:   booklit.StyleBold,
	}
}

func (plugin Plugin) Larger(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Content: content,
		Style:   booklit.StyleLarger,
	}
}

func (plugin Plugin) Smaller(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Content: content,
		Style:   booklit.StyleSmaller,
	}
}

func (plugin Plugin) Strike(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Content: content,
		Style:   booklit.StyleStrike,
	}
}

func (plugin Plugin) Superscript(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Content: content,
		Style:   booklit.StyleSuperscript,
	}
}

func (plugin Plugin) Subscript(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Content: content,
		Style:   booklit.StyleSubscript,
	}
}

func (plugin Plugin) Inset(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Content: content,
		Style:   booklit.StyleInset,
	}
}

func (plugin Plugin) Aside(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Content: content,
		Style:   booklit.StyleAside,
	}
}

func (plugin Plugin) Link(content booklit.Content, target string) booklit.Content {
	return booklit.Link{
		Content: content,
		Target:  target,
	}
}

func (plugin Plugin) Reference(tag string, content ...booklit.Content) booklit.Content {
	ref := &booklit.Reference{
		TagName: tag,
	}

	if len(content) > 0 {
		ref.Content = content[0]
	}

	return ref
}

func (plugin Plugin) Target(tag string, titleAndContent ...booklit.Content) booklit.Content {
	ref := &booklit.Target{
		TagName: tag,
	}

	switch len(titleAndContent) {
	case 2:
		ref.Title = titleAndContent[0]
		ref.Content = titleAndContent[1]
	case 1:
		ref.Title = titleAndContent[0]
	default:
		ref.Title = booklit.String(tag)
	}

	return ref
}

func (plugin Plugin) List(items ...booklit.Content) booklit.Content {
	return booklit.List{
		Items: items,
	}
}

func (plugin Plugin) OrderedList(items ...booklit.Content) booklit.Content {
	return booklit.List{
		Items:   items,
		Ordered: true,
	}
}

func (plugin Plugin) Image(path string, description ...string) booklit.Content {
	img := booklit.Image{
		Path: path,
	}

	if len(description) > 0 {
		img.Description = description[0]
	}

	return img
}

func (plugin Plugin) SetPartial(name string, content booklit.Content) {
	plugin.section.SetPartial(name, content)
}

func (plugin Plugin) Table(rows ...booklit.Content) (booklit.Content, error) {
	table := booklit.Table{}

	for _, row := range rows {
		list, ok := row.(booklit.List)
		if !ok {
			return nil, fmt.Errorf("table row is not a list: %s", row)
		}

		table.Rows = append(table.Rows, list.Items)
	}

	return table, nil
}

func (plugin Plugin) TableRow(cols ...booklit.Content) booklit.Content {
	return plugin.List(cols...)
}

func (plugin Plugin) Definitions(items ...booklit.Content) (booklit.Content, error) {
	defs := booklit.Definitions{}
	for _, item := range items {
		list, ok := item.(booklit.List)
		if !ok {
			return nil, fmt.Errorf("definition item is not a list: %s", item)
		}

		if len(list.Items) != 2 {
			return nil, fmt.Errorf("definition item must have two entries: %s", item)
		}

		defs = append(defs, booklit.Definition{
			Subject:    list.Items[0],
			Definition: list.Items[1],
		})
	}

	return defs, nil
}

func (plugin Plugin) Definition(subject booklit.Content, definition booklit.Content) booklit.Content {
	return plugin.List(subject, definition)
}
