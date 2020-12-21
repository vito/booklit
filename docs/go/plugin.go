package booklitdoc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/styles"
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/baselit"
	chromap "github.com/vito/booklit/chroma"
)

func init() {
	booklit.RegisterPlugin("booklitdoc", NewPlugin)

	styles.Fallback = chroma.MustNewStyle("booklitdoc", chroma.StyleEntries{
		chroma.Comment:               "italic",
		chroma.CommentPreproc:        "noitalic",
		chroma.Keyword:               "#ed6c30 bold",
		chroma.KeywordPseudo:         "nobold",
		chroma.KeywordType:           "nobold",
		chroma.OperatorWord:          "#fcc21b bold",
		chroma.NameClass:             "#fcc21b bold",
		chroma.NameNamespace:         "#fcc21b bold",
		chroma.NameException:         "#fcc21b bold",
		chroma.NameEntity:            "#fcc21b bold",
		chroma.NameTag:               "#fcc21b bold",
		chroma.LiteralString:         "#fcc21b",
		chroma.LiteralStringInterpol: "bold",
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
		chroma:  chromap.NewPlugin(section).(chromap.Plugin),
	}
}

type Plugin struct {
	section *booklit.Section
	base    baselit.Plugin
	chroma  chromap.Plugin
}

func (plugin Plugin) OutputFrame(url string) booklit.Content {
	return booklit.Styled{
		Style: "output-frame",
		Content: booklit.Link{
			Content: booklit.String(url),
			Target:  url,
		},
		Partials: booklit.Partials{
			"URL": booklit.String(url),
		},
	}
}

func (plugin Plugin) SyntaxHl(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Style:   "syntax-hl",
		Content: content,
	}
}

func (plugin Plugin) ColumnHeader(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Style:   "column-header",
		Block:   true,
		Content: content,
	}
}

func (plugin Plugin) Columns(title booklit.Content, rest ...booklit.Content) booklit.Content {
	return booklit.Styled{
		Style:   "columns",
		Block:   true,
		Content: title,
		Partials: booklit.Partials{
			"Columns": booklit.Sequence(rest),
		},
	}
}

func (plugin Plugin) BigCode(content booklit.Content) booklit.Content {
	return booklit.Styled{
		Style:   "big-code",
		Content: content,
	}
}

func (plugin Plugin) Codefile(path string) booklit.Content {
	return booklit.Styled{
		Style:   "codefile",
		Content: booklit.String(path),
		Block:   true,
	}
}

func (plugin Plugin) Godoc(ref string) (booklit.Content, error) {
	spl := strings.SplitN(ref, ".", 2)

	pkg := strings.TrimLeft(spl[0], "*")

	syntax, err := plugin.chroma.Syntax("go", booklit.Sequence{
		booklit.String(spl[0] + "."),
		plugin.base.Bold(booklit.String(spl[1])),
	})
	if err != nil {
		return nil, err
	}

	return plugin.base.Link(
		syntax,
		"https://pkg.go.dev/github.com/vito/"+pkg+"#"+spl[1],
	), nil
}

func (plugin Plugin) Define(node ast.Node, content booklit.Content) (booklit.Content, error) {
	invoke := node.(ast.Sequence)[0].(ast.Invoke)

	title, err := plugin.chroma.Syntax("lit", booklit.String("\\"+invoke.Function))
	if err != nil {
		return nil, err
	}

	thumb, err := plugin.chroma.Syntax("lit", plugin.renderInvoke(invoke))
	if err != nil {
		return nil, err
	}

	return booklit.Styled{
		Style: "definition",

		Content: content,

		Partials: booklit.Partials{
			"Thumb": booklit.Sequence{
				booklit.Target{
					TagName:  invoke.Function,
					Location: plugin.section.InvokeLocation,
					Title:    title,
					Content:  content,
				},
				thumb,
			},
		},
	}, nil
}

func (plugin Plugin) renderInvoke(invoke ast.Invoke) booklit.Content {
	str := booklit.Sequence{booklit.String("\\")}

	str = append(str, &booklit.Reference{
		TagName:  invoke.Function,
		Content:  plugin.base.Bold(booklit.String(invoke.Function)),
		Location: invoke.Location,
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
			TagName:  tag,
			Location: plugin.section.InvokeLocation,
			Title:    booklit.String(name),
			Content:  definition,
		})
	}

	content = append(
		content,
		booklit.Paragraph{
			booklit.Styled{
				Style:   booklit.StyleBold,
				Content: booklit.String(name),
			},
		},
		definition,
	)

	return content, nil
}
