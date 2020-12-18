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
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	iterator, err := lexer.Tokenise(nil, code.String())
	if err != nil {
		return nil, err
	}

	formatter := html.New(html.PreventSurroundingPre(code.IsFlow()))

	chromaStyle := styles.Fallback
	if len(styleName) > 0 {
		chromaStyle = styles.Get(styleName[0])
	}

	buf := new(bytes.Buffer)
	err = formatter.Format(buf, chromaStyle, iterator)
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
		Block:   !code.IsFlow(),
		Content: code,
		Partials: booklit.Partials{
			"Language": booklit.String(language),
			"HTML":     booklit.String(buf.String()),
		},
	}, nil
}
