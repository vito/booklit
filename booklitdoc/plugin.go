package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/styles"
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/baselit"
)

func init() {
	booklit.RegisterPlugin("booklitdoc", NewPlugin)

	styles.Fallback = chroma.MustNewStyle("booklitdoc", chroma.StyleEntries{
		chroma.Comment:               "italic",
		chroma.CommentPreproc:        "noitalic",
		chroma.Keyword:               "bold",
		chroma.KeywordPseudo:         "nobold",
		chroma.KeywordType:           "nobold",
		chroma.OperatorWord:          "bold",
		chroma.NameClass:             "bold",
		chroma.NameNamespace:         "bold",
		chroma.NameException:         "bold",
		chroma.NameEntity:            "bold",
		chroma.NameTag:               "bold",
		chroma.LiteralString:         "italic",
		chroma.LiteralStringInterpol: "bold",
		chroma.LiteralStringEscape:   "bold",
		chroma.GenericHeading:        "bold",
		chroma.GenericSubheading:     "bold",
		chroma.GenericEmph:           "italic",
		chroma.GenericStrong:         "bold",
		chroma.GenericPrompt:         "bold",
		chroma.Error:                 "border:#FF0000",
	})
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

	return booklit.Styled{
		Style: "definition",

		Content: content,

		Partials: booklit.Partials{
			"Thumb": booklit.Sequence{
				booklit.Target{
					TagName: invoke.Function,
					Title: plugin.base.Code(booklit.Sequence{
						booklit.String("\\"),
						plugin.base.Bold(booklit.String(invoke.Function)),
					}),
					Content: content,
				},
				plugin.base.Code(plugin.renderInvoke(invoke)),
			},
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

func (plugin Plugin) DescribeFruit(
	name string,
	definition booklit.Content,
	tags ...string,
) (booklit.Content, error) {
	if name == "" {
		return nil, errors.New("name cannot be blank")
	}

	content := booklit.Sequence{}
	if len(tags) == 0 {
		tags = []string{name}
	}

	for _, tag := range tags {
		content = append(content, booklit.Target{
			TagName: tag,
			Title:   booklit.String(name),
			Content: definition,
		})
	}

	content = append(content, booklit.Paragraph{
		booklit.Styled{
			Style:   booklit.StyleBold,
			Content: booklit.String(name),
		},
	})

	content = append(content, definition)

	return content, nil
}
