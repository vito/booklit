package baselit

import (
	"fmt"
	"path/filepath"

	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/load"
)

func init() {
	booklit.RegisterPlugin("base", booklit.PluginFactoryFunc(NewPlugin))
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

func (plugin Plugin) Title(title booklit.Content, tags ...string) {
	plugin.section.SetTitle(title, tags...)
}

func (plugin Plugin) Section(node ast.Node) error {
	processor := &load.Processor{
		PluginFactories: plugin.section.PluginFactories,
	}

	section, err := processor.EvaluateSection(plugin.section.Path, node)
	if err != nil {
		return err
	}

	section.Parent = plugin.section

	plugin.section.Children = append(plugin.section.Children, section)

	return nil
}

func (plugin Plugin) IncludeSection(path string) error {
	sectionPath := filepath.Join(filepath.Dir(plugin.section.Path), path)

	result, err := ast.ParseFile(sectionPath)
	if err != nil {
		return err
	}

	processor := &load.Processor{
		PluginFactories: []booklit.PluginFactory{
			booklit.PluginFactoryFunc(NewPlugin),
		},
	}

	section, err := processor.EvaluateSection(sectionPath, result.(ast.Node))
	if err != nil {
		return err
	}

	section.Parent = plugin.section

	plugin.section.Children = append(plugin.section.Children, section)

	return nil
}

func (plugin Plugin) SplitSections() {
	plugin.section.SplitSections = true
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

func (plugin Plugin) Reference(tag string, content ...booklit.Content) booklit.Content {
	ref := &booklit.Reference{
		TagName: tag,
	}

	if len(content) > 0 {
		ref.Content = content[0]
	}

	return ref
}

func (plugin Plugin) Target(tag string, display ...booklit.Content) booklit.Content {
	ref := &booklit.Target{
		TagName: tag,
	}

	if len(display) > 0 {
		ref.Display = display[0]
	} else {
		ref.Display = booklit.String(tag)
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
