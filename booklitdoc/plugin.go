package booklitdoc

import (
	"fmt"
	"strings"

	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/baselit"
)

func init() {
	booklit.RegisterPlugin("booklitdoc", NewPlugin)
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

func (plugin Plugin) Godoc(ref string) booklit.Content {
	spl := strings.SplitN(ref, ".", 2)

	pkg := strings.TrimLeft(spl[0], "*")

	return plugin.base.Link(
		plugin.base.Code(booklit.Sequence{
			booklit.String(spl[0] + "."),
			plugin.base.Bold(booklit.String(spl[1])),
		}),
		"https://godoc.org/github.com/vito/"+pkg+"#"+spl[1],
	)
}

func (plugin Plugin) Define(node ast.Node, content booklit.Content) booklit.Content {
	invoke := node.(ast.Sequence)[0].(ast.Invoke)

	return booklit.Sequence{
		booklit.Block{
			Class: "definition-thumb",
			Content: booklit.Sequence{
				booklit.Target{
					TagName: invoke.Function,
					Display: plugin.base.Code(booklit.Sequence{
						booklit.String("\\"),
						plugin.base.Bold(booklit.String(invoke.Function)),
					}),
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
