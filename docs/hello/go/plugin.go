package plugin

import "github.com/vito/booklit"

func init() {
	booklit.RegisterPlugin("hello", NewPlugin)
}

type MyPlugin struct {
	section *booklit.Section
}

func NewPlugin(section *booklit.Section) booklit.Plugin {
	return MyPlugin{section: section}
}

func (p MyPlugin) Testimonial(
	source, quote booklit.Content,
) booklit.Content {
	return booklit.Styled{
		Style:   "testimonial",
		Content: quote,
		Partials: booklit.Partials{
			"Source": source,
		},
	}
}
