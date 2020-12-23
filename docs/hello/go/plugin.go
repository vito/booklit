package example

import "github.com/vito/booklit"

func init() {
	booklit.RegisterPlugin("example", NewPlugin)
}

type Example struct {
	section *booklit.Section
}

func NewPlugin(section *booklit.Section) booklit.Plugin {
	return Example{section: section}
}

func (Example) Quote(
	quote, source booklit.Content,
) booklit.Content {
	return booklit.Styled{
		Style:   "quote",
		Content: quote,
		Partials: booklit.Partials{
			"Source": source,
		},
	}
}
