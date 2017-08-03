package plugin

import "github.com/vito/booklit"

func init() {
	booklit.RegisterPlugin("stringer", NewPlugin)
}

func NewPlugin(section *booklit.Section) booklit.Plugin {
	return Plugin{
		section: section,
	}
}

type Plugin struct {
	section *booklit.Section
}

func (plugin Plugin) String(arg string) booklit.Content {
	return booklit.String(arg)
}
