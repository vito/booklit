// Package chroma provides a basic plugin for implementing syntax highlighting
// using Chroma (https://github.com/alecthomas/chroma).
//
// To use this plugin, pass `--plugin github.com/vito/booklit/chroma/plugin`
// and use it like so:
//
//   \use-plugin{chroma}
//
//   \syntax{go}{{{
//   package chroma
//
//   // ...
//   }}}
//
// An optional style name may be specified as the third argument. To use a
// custom style you may write your own plugin that embeds this plugin, or
// re-assign github.com/alecthomas/chroma/styles.Fallback to change the
// default.
package chroma

import (
	"bytes"
	"regexp"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/vito/booklit"
)

func NewPlugin(section *booklit.Section) booklit.Plugin {
	return Plugin{
		section: section,
	}
}

type Plugin struct {
	section *booklit.Section
}

func (plugin Plugin) Syntax(language string, code booklit.Content, styleName ...string) (booklit.Content, error) {
	chromaStyle := styles.Fallback
	if len(styleName) > 0 {
		chromaStyle = styles.Get(styleName[0])
	}

	return plugin.SyntaxTransform(language, code, chromaStyle)
}

type Transformer struct {
	Pattern   *regexp.Regexp
	Transform func(string) booklit.Content
}

func (t Transformer) TransformAll(str string) booklit.Sequence {
	matches := t.Pattern.FindAllStringIndex(str, -1)

	out := booklit.Sequence{}
	last := 0
	for _, match := range matches {
		if match[0] > last {
			out = append(out, booklit.String(str[last:match[0]]))
		}

		out = append(out, t.Transform(str[match[0]:match[1]]))

		last = match[1]
	}

	if len(str) > last {
		out = append(out, booklit.String(str[last:]))
	}

	return out
}

func (plugin Plugin) SyntaxTransform(language string, code booklit.Content, chromaStyle *chroma.Style, transformers ...Transformer) (booklit.Content, error) {
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	iterator, err := lexer.Tokenise(nil, code.String())
	if err != nil {
		return nil, err
	}

	formatter := html.New(html.PreventSurroundingPre(code.IsFlow()))

	buf := new(bytes.Buffer)
	err = formatter.Format(buf, chromaStyle, iterator)
	if err != nil {
		return nil, err
	}

	var style booklit.Style
	if code.IsFlow() {
		style = "code-flow"
	} else {
		style = "code-block"
	}

	highlighted := booklit.Sequence{booklit.String(buf.String())}

	for _, t := range transformers {
		var newHighlighted booklit.Sequence
		for _, con := range highlighted {
			switch val := con.(type) {
			case booklit.String:
				newHighlighted = append(newHighlighted, t.TransformAll(val.String())...)
			default:
				newHighlighted = append(newHighlighted, con)
			}
		}

		highlighted = newHighlighted
	}

	for i, con := range highlighted {
		if _, ok := con.(booklit.String); ok {
			highlighted[i] = booklit.Styled{
				Style:   "raw-html",
				Content: con,
			}
		}
	}

	return booklit.Styled{
		Style:   style,
		Block:   !code.IsFlow(),
		Content: highlighted,
		Partials: booklit.Partials{
			"Language": booklit.String(language),
		},
	}, nil
}
