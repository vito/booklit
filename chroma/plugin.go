// Package chroma provides an advanced plugin for syntax highlighting with
// support for content transformers (e.g. linkifying function names).
//
// For basic syntax highlighting, use baselit's built-in \code-block or
// \syntax functions instead — no plugin needed.
//
// This plugin is useful when you need SyntaxTransform with custom Transformers
// to post-process highlighted output.
package chroma

import (
	"bytes"
	"regexp"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/vito/booklit"
	"github.com/vito/booklit/baselit"
)

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

func (plugin Plugin) CodeBlock(language string, code booklit.Content, styleName ...string) (booklit.Content, error) {
	return plugin.base.CodeBlock(language, code, styleName...)
}

func (plugin Plugin) Syntax(language string, code booklit.Content, styleName ...string) (booklit.Content, error) {
	return plugin.base.Syntax(language, code, styleName...)
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

	var formatter *html.Formatter
	if code.IsFlow() {
		formatter = html.New(html.InlineCode(true))
	} else {
		formatter = html.New()
	}

	buf := new(bytes.Buffer)
	err = formatter.Format(buf, chromaStyle, iterator)
	if err != nil {
		return nil, err
	}

	var style booklit.Style
	if code.IsFlow() {
		style = booklit.StyleCodeFlow
	} else {
		style = booklit.StyleCodeBlock
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
