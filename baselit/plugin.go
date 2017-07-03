package baselit

import (
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/load"
)

type PluginFactory struct {
	Processor *load.Processor
}

func (pf PluginFactory) NewPlugin(section *booklit.Section) booklit.Plugin {
	return Plugin{
		processor: pf.Processor,
		section:   section,
	}
}

type Plugin struct {
	processor *load.Processor
	section   *booklit.Section
}

func (plugin Plugin) Title(title booklit.Content, tags ...string) {
	plugin.section.SetTitle(title, tags...)
}

func (plugin Plugin) Section(node ast.Node) error {
	section, err := plugin.processor.EvaluateSection(node)
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

func (plugin Plugin) Reference(tag string, content ...booklit.Content) booklit.Content {
	ref := &booklit.Reference{
		TagName: tag,
	}

	if len(content) > 0 {
		ref.Content = content[0]
	}

	return ref
}
