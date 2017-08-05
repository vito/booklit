package pygments

import "github.com/vito/booklit"

func init() {
	booklit.RegisterPlugin("pygments", NewPlugin)
}

func NewPlugin(section *booklit.Section) booklit.Plugin {
	return Plugin{
		section: section,
	}
}

type Plugin struct {
	section *booklit.Section
}

func (plugin Plugin) Syntax(language string, code booklit.Content) (booklit.Content, error) {
	hl, err := pygmentize(language, code.String())
	if err != nil {
		return nil, err
	}

	var style booklit.Style
	if code.IsFlow() {
		style = "inline-code"
	} else {
		style = "code-block"
	}

	return booklit.Styled{
		Style:   style,
		Content: booklit.String(hl),
	}, nil
}
