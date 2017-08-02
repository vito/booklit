package baselit

import "github.com/vito/booklit"

func init() {
	booklit.RegisterPlugin("data-style", NewPlugin)
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

func (plugin Plugin) StructStyle(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Style: "data-style",

		Data: MyStruct{
			MyField: content,
		},
	}
}

func (plugin Plugin) MapStyle(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Style: "data-style",

		Data: map[string]booklit.Content{
			"MyField": content,
		},
	}
}

func (plugin Plugin) InlineStyle(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Style: "inline-data-style",

		Flow: true,
		Data: map[string]booklit.Content{
			"MyField": content,
		},
	}
}
