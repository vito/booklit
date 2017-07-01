package baselit

import "github.com/vito/booklit"

type PluginFactory struct{}

func (PluginFactory) NewPlugin(section *booklit.Section) booklit.Plugin {
	return Plugin{
		section: section,
	}
}

type Plugin struct {
	section *booklit.Section
}

func (plugin Plugin) Title(title booklit.Content, tags ...string) {
	plugin.section.SetTitle(title, tags...)
}

func (plugin Plugin) Section(title booklit.Content, content booklit.Content) {
	section := &booklit.Section{
		Body: content,

		Parent: plugin.section,
	}

	section.SetTitle(title)

	plugin.section.Children = append(plugin.section.Children, section)
}

func (plugin Plugin) SplitSections() {
	plugin.section.SplitSections = true
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
