package builtins

import (
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/treehighlight"
)

func init() {
	Register("CodeBlock", codeBlockFunc)
	Register("Syntax", codeBlockFunc)
}

// codeBlockFunc — `<CodeBlock language="go">code</CodeBlock>` or
// `<Syntax language="go">code</Syntax>`. Tree-sitter-driven syntax
// highlighting; identifiers that resolve to existing tags are linkified.
func codeBlockFunc(ctx *Context, props map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	language, err := requireStringProp(ctx, props, "language", "CodeBlock")
	if err != nil {
		return nil, err
	}

	code, err := EvaluateChildren(ctx, children)
	if err != nil {
		return nil, err
	}
	if code == nil {
		code = booklit.Empty
	}

	return syntax(ctx.Section, language, code)
}

// syntax runs the source through treehighlight and wraps the resulting
// chunks in a Styled code-block/code-flow container. Raw HTML chunks pass
// through untouched so the highlighter's <span> markup survives intact;
// link chunks become \reference invocations so they auto-link to matching
// tags defined elsewhere in the document.
//
// References are marked Optional so an unresolved tag falls back to
// rendering the original source text — most identifiers in a code sample
// won't have a matching target, and we don't want unresolved names to
// fail the build.
//
// NOTE: the open/close wrapper here mirrors treehighlight.HTML and
// treehighlight.PlainHTML; keep all three in sync if the markup changes.
func syntax(section *booklit.Section, language string, code booklit.Content) (booklit.Content, error) {
	chunks, err := treehighlight.Chunks(language, code.String(), treehighlight.Options{LinkReferences: true})
	if err != nil {
		return nil, err
	}

	inline := code.IsFlow()

	var style booklit.Style
	if inline {
		style = booklit.StyleCodeFlow
	} else {
		style = booklit.StyleCodeBlock
	}

	open, close := codeWrapper(inline)
	content := booklit.Sequence{rawHTML(open)}
	for _, chunk := range chunks {
		switch {
		case chunk.HTML != "":
			content = append(content, rawHTML(chunk.HTML))
		case chunk.LinkTag != "":
			content = append(content, &booklit.Reference{
				Section:  section,
				TagName:  chunk.LinkTag,
				Content:  booklit.String(chunk.LinkText),
				Optional: true,
				Location: section.InvokeLocation,
			})
		}
	}
	content = append(content, rawHTML(close))

	return booklit.Styled{
		Style:   style,
		Block:   !inline,
		Content: content,
		Partials: booklit.Partials{
			"Language": booklit.String(language),
		},
	}, nil
}

func codeWrapper(inline bool) (open, close string) {
	if inline {
		return `<code style=";-webkit-text-size-adjust:none;">`, `</code>`
	}
	return `<pre style=";-webkit-text-size-adjust:none;"><code>`, `</code></pre>`
}

func rawHTML(s string) booklit.Content {
	return booklit.Styled{Style: "raw-html", Content: booklit.String(s)}
}
