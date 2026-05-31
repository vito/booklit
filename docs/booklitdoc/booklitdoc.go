// Package booklitdoc provides the JSX built-ins used by Booklit's own
// documentation site. It is imported by cmd/booklit-docs so the helpers
// register themselves; the main cmd/booklit binary does not include them.
package booklitdoc

import (
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/baselit"
	"github.com/vito/booklit/builtins"
)

func init() {
	// LitSyntax is the only JSX built-in still implemented in Go here:
	// it needs chroma highlighting plus the regex-based `\function`
	// linkifier below. Everything else the docs add is expressed as an
	// mdx template in docs/html/ and dispatched via the template tier:
	// Columns / Column / ColumnHeader (plain <div> wrappers around their
	// children), Define / Godoc, OutputFrame, and TemplateLink.
	builtins.Register("LitSyntax", litSyntaxFunc)

	styles.Fallback = chroma.MustNewStyle("booklitdoc", chroma.StyleEntries{
		chroma.Comment:               "#c29d7c italic",
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

// linkTransformer rewrites `\function-name` occurrences inside
// syntax-highlighted Booklit code into references to the function's
// definition tag (set up by <Define>).
func linkTransformer(sec *booklit.Section) chromaTransformer {
	return chromaTransformer{
		Pattern: regexp.MustCompile(`\\([a-z][a-z0-9-]*)`),
		Transform: func(invoke string) booklit.Content {
			fn := strings.TrimPrefix(invoke, `\`)
			return booklit.Sequence{
				booklit.String(`\`),
				&booklit.Reference{
					Section:  sec,
					TagName:  fn,
					Content:  booklit.String(fn),
					Optional: true,
				},
			}
		},
	}
}

// chromaTransformer mirrors the old chroma.Transformer locally since
// the chroma plugin package is gone.
type chromaTransformer struct {
	Pattern   *regexp.Regexp
	Transform func(string) booklit.Content
}

func (t chromaTransformer) transformAll(str string) booklit.Sequence {
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

// syntaxTransform runs chroma syntax highlighting via baselit (which
// uses styles.Fallback — set in init() to the booklitdoc palette) and
// then applies a chain of textual transformers (e.g. linkify function
// names) over the highlighted output. Re-implemented locally so the
// docs aren't blocked on the deleted chroma plugin package.
func syntaxTransform(section *booklit.Section, language string, code booklit.Content, transformers ...chromaTransformer) (booklit.Content, error) {
	base := baselit.NewPlugin(section).(baselit.Plugin)
	highlighted, err := base.Syntax(language, code)
	if err != nil {
		return nil, err
	}

	styled, ok := highlighted.(booklit.Styled)
	if !ok {
		return highlighted, nil
	}
	inner, ok := styled.Content.(booklit.Styled)
	if !ok {
		return highlighted, nil
	}
	raw, ok := inner.Content.(booklit.String)
	if !ok {
		return highlighted, nil
	}

	transformed := booklit.Sequence{booklit.String(string(raw))}
	for _, t := range transformers {
		var next booklit.Sequence
		for _, con := range transformed {
			if s, ok := con.(booklit.String); ok {
				next = append(next, t.transformAll(s.String())...)
			} else {
				next = append(next, con)
			}
		}
		transformed = next
	}
	for i, con := range transformed {
		if _, ok := con.(booklit.String); ok {
			transformed[i] = booklit.Styled{Style: "raw-html", Content: con}
		}
	}
	styled.Content = transformed
	return styled, nil
}

// litSyntaxFunc — `<LitSyntax>code</LitSyntax>`. Highlights Booklit
// source as `lit`, then linkifies `\function` references.
func litSyntaxFunc(ctx *builtins.Context, _ map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	code, err := evalChildren(ctx, children)
	if err != nil {
		return nil, err
	}
	syntax, err := syntaxTransform(ctx.Section, "lit", code, linkTransformer(ctx.Section))
	if err != nil {
		return nil, err
	}
	var style booklit.Style = "lit-block"
	if code.IsFlow() {
		style = "lit-flow"
	}
	return booklit.Styled{Style: style, Content: syntax}, nil
}

// evalChildren wraps the children list in a Sequence and evaluates it.
func evalChildren(ctx *builtins.Context, children []ast.Node) (booklit.Content, error) {
	if len(children) == 0 {
		return booklit.Empty, nil
	}
	return ctx.Evaluate(ast.Sequence(children))
}
