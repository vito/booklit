package booklitdoc

import (
	"errors"
	"fmt"
	"regexp"
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
		chroma.Comment:               "#c29d7c italic", // lighten(@background, 50%)
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

var linkTransformer = chromap.Transformer{
	Pattern: regexp.MustCompile(`\\([a-z-]+)`),
	Transform: func(invoke string) booklit.Content {
		function := strings.TrimPrefix(invoke, `\`)

		return booklit.Sequence{
			booklit.String(`\`),
			&booklit.Reference{
				TagName:  function,
				Content:  booklit.String(function),
				Optional: true,
			},
		}
	},
}

var argTransformer = chromap.Transformer{
	Pattern: regexp.MustCompile(`\{[^\}]+\}`),
	Transform: func(farg string) booklit.Content {
		arg := strings.TrimPrefix(strings.TrimSuffix(farg, "}"), "{")
		return booklit.Sequence{
			booklit.String("{"),
			booklit.Styled{
				Style:   booklit.StyleItalic,
				Content: booklit.String(arg),
			},
			booklit.String("}"),
		}
	},
}

func (plugin Plugin) LitSyntax(code booklit.Content) (booklit.Content, error) {
	syntax, err := plugin.chroma.SyntaxTransform("lit", code, styles.Fallback, linkTransformer)
	if err != nil {
		return nil, err
	}

	var style booklit.Style = "lit-block"
	if code.IsFlow() {
		style = "lit-flow"
	}

	return booklit.Styled{
		Style:   style,
		Content: syntax,
	}, nil
}

func (plugin Plugin) Godoc(ref string) (booklit.Content, error) {
	spl := strings.SplitN(ref, ".", 2)

	pkg := strings.TrimLeft(spl[0], "*")

	return plugin.base.Link(
		booklit.Styled{
			Style: booklit.StyleVerbatim,
			Content: booklit.Sequence{
				booklit.String(spl[0] + "."),
				plugin.base.Bold(booklit.String(spl[1])),
			},
		},
		"https://pkg.go.dev/github.com/vito/"+pkg+"#"+spl[1],
	), nil
}

func (plugin Plugin) TemplateLink(tmpl string) (booklit.Content, error) {
	return plugin.base.Link(
		booklit.Styled{
			Style:   booklit.StyleVerbatim,
			Content: booklit.String(tmpl),
		},
		"https://github.com/vito/booklit/blob/master/render/html/"+tmpl,
	), nil
}

func (plugin Plugin) Define(node ast.Node, content booklit.Content) (booklit.Content, error) {
	invoke := node.(ast.Sequence)[0].(ast.Invoke)

	title, err := plugin.chroma.Syntax("lit", booklit.String("\\"+invoke.Function))
	if err != nil {
		return nil, err
	}

	thumb, err := plugin.chroma.SyntaxTransform("lit", plugin.renderInvoke(invoke), styles.Fallback, linkTransformer, argTransformer)
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
	str := fmt.Sprintf(`\%s`, invoke.Function)

	for _, arg := range invoke.Arguments {
		str += "{"

		for _, n := range arg.(ast.Sequence) {
			str += fmt.Sprintf("%s", n)
		}

		str += "}"
	}

	return booklit.String(str)
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
