package litmd

import (
	"github.com/vito/booklit"
)

func init() {
	booklit.RegisterPlugin("litmd", NewPlugin)
}

func NewPlugin(section *booklit.Section) booklit.Plugin {
	return Plugin{
		section: section,
	}
}

type Plugin struct {
	section *booklit.Section
}

func (plugin Plugin) Quote(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Style:   "quote",
		Block:   true,
		Content: content,
	}
}

func (plugin Plugin) Header(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Style:   "header",
		Block:   true,
		Content: content,
		Partials: booklit.Partials{
			"Level": booklit.String("1"),
		},
	}
}

func (plugin Plugin) Subheader(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Style:   "header",
		Block:   true,
		Content: content,
		Partials: booklit.Partials{
			"Level": booklit.String("2"),
		},
	}
}

func (plugin Plugin) Subsubheader(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Style:   "header",
		Block:   true,
		Content: content,
		Partials: booklit.Partials{
			"Level": booklit.String("3"),
		},
	}
}

func (plugin Plugin) Subsubsubheader(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Style:   "header",
		Block:   true,
		Content: content,
		Partials: booklit.Partials{
			"Level": booklit.String("4"),
		},
	}
}

func (plugin Plugin) Subsubsubsubheader(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Style:   "header",
		Block:   true,
		Content: content,
		Partials: booklit.Partials{
			"Level": booklit.String("5"),
		},
	}
}

func (plugin Plugin) Subsubsubsubsubheader(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Style:   "header",
		Block:   true,
		Content: content,
		Partials: booklit.Partials{
			"Level": booklit.String("6"),
		},
	}
}
