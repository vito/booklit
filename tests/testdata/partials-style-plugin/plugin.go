package plugin

import "github.com/vito/booklit"

func init() {
	booklit.RegisterPlugin("partial-style", NewPlugin)
}

func NewPlugin(section *booklit.Section) booklit.Plugin {
	return Plugin{
		section: section,
	}
}

type Plugin struct {
	section *booklit.Section
}

type MyStruct struct {
	MyField booklit.Content
}

func (plugin Plugin) BlockStyle(title, content booklit.Content) booklit.Content {
	return booklit.Styled{
		Style: "custom-style",
		Block: true,

		Content: content,

		Partials: booklit.Partials{
			"Title": title,
		},
	}
}

func (plugin Plugin) InlineStyle(title, content booklit.Content) booklit.Content {
	return booklit.Styled{
		Style: "inline-custom-style",

		Content: content,

		Partials: booklit.Partials{
			"Title": title,
		},
	}
}
