package baselit

import (
	"github.com/vito/booklit"
)

type BaselitPluginFactory struct{}

func (BaselitPluginFactory) NewPlugin(section *booklit.Section) booklit.Plugin {
	return baselitPlugin{
		section: section,
	}
}

type baselitPlugin struct {
	section *booklit.Section
}

func (plugin baselitPlugin) UseTemplate(name string) {
}

func (plugin baselitPlugin) Title(title booklit.Content, tags ...string) {
	plugin.section.Title = title
	plugin.section.Tags = tags
}

func (plugin baselitPlugin) Section(title booklit.Content, content booklit.Content) {
	section := &booklit.Section{
		Title: title,
		Body:  content,

		Parent: plugin.section,
	}

	plugin.section.Children = append(plugin.section.Children, section)
}

func (plugin baselitPlugin) Reference(tag string, content booklit.Content) {
}

func (plugin baselitPlugin) IncludeSection(path string) booklit.Content {
	return nil
}

func (plugin baselitPlugin) SplitSections() {
	plugin.section.SplitSections = true
}

func (plugin baselitPlugin) Something() booklit.Content {
	return booklit.String("hello from plugin")
}
