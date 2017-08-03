package plugin

import "github.com/vito/booklit"

func init() {
	booklit.RegisterPlugin("arbitrary-style", NewPlugin)
}

func NewPlugin(section *booklit.Section) booklit.Plugin {
	return Plugin{
		section: section,
	}
}

type Plugin struct {
	section *booklit.Section
}

func (plugin Plugin) ArbitraryStyle(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Style:   "arbitrary",
		Content: content,
	}
}
