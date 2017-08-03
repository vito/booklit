package plugin

import "github.com/vito/booklit"

func init() {
	booklit.RegisterPlugin("set-partials", NewPlugin)
}

func NewPlugin(section *booklit.Section) booklit.Plugin {
	return Plugin{
		section: section,
	}
}

type Plugin struct {
	section *booklit.Section
}

func (plugin Plugin) SetThePartial() {
	plugin.section.SetPartial("FooBar", booklit.Paragraph{
		booklit.String("I'm a partial!"),
	})
}
