package booklitdoc

import (
	"fmt"

	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/baselit"
)

func init() {
	booklit.RegisterPlugin("booklitdoc", booklit.PluginFactoryFunc(NewPlugin))
}

func NewPlugin(section *booklit.Section) booklit.Plugin {
	return Plugin{
		section: section,
		base:    baselit.NewPlugin(section).(baselit.Plugin),
	}
}

type Plugin struct {
	section *booklit.Section
	base    baselit.Plugin
}

func (plugin Plugin) Define(node ast.Node, content booklit.Content) booklit.Content {
	invoke := node.(ast.Sequence)[0].(ast.Invoke)

	return booklit.Sequence{
		booklit.Block{
			Class: "definition-thumb",
			Content: booklit.Sequence{
				booklit.Target{
					TagName: invoke.Function,
					Display: plugin.base.Code(booklit.String("\\" + invoke.Function)),
				},
				plugin.base.Code(plugin.renderInvoke(invoke)),
			},
		},
		booklit.Block{
			Class:   "definition-content",
			Content: content,
		},
	}
}

func (plugin Plugin) renderInvoke(invoke ast.Invoke) booklit.Content {
	str := booklit.Sequence{booklit.String("\\")}

	str = append(str, &booklit.Reference{
		TagName: invoke.Function,
		Content: plugin.base.Bold(booklit.String(invoke.Function)),
	})

	for _, arg := range invoke.Arguments {
		str = append(str, booklit.String("{"))

		for _, n := range arg.(ast.Sequence) {
			str = append(str, plugin.base.Italic(booklit.String(fmt.Sprintf("%s", n))))
		}

		str = append(str, booklit.String("}"))
	}

	return str
}
