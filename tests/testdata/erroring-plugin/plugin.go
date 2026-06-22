package plugin

import (
	"errors"

	"github.com/vito/booklit"
)

func init() {
	booklit.RegisterPlugin("errer", NewPlugin)
}

func NewPlugin(section *booklit.Section) booklit.Plugin {
	return Plugin{
		section: section,
	}
}

type Plugin struct {
	section *booklit.Section
}

func (plugin Plugin) SingleFail(arg string) error {
	return errors.New("oh no")
}

func (plugin Plugin) MultiFail(arg string) (booklit.Content, error) {
	return nil, errors.New("oh no")
}
