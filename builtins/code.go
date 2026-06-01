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
// chunks in a `<pre><code class="language-X">…</code></pre>` shape (or
// just `<code class="language-X">…</code>` when the source is inline
// flow content). Highlighter chunks become RawFragment so the `<span>`
// markup survives intact; link chunks become \reference invocations so
// they auto-link to matching tags defined elsewhere.
//
// References are marked Optional so an unresolved tag falls back to
// rendering the original source text — most identifiers in a code
// sample won't have a matching target, and we don't want unresolved
// names to fail the build.
//
// `class="language-X"` follows the Prism/highlight.js convention and is
// the natural home for the language metadata that the old
// `Partials{Language}` bag was supposed to carry (set in one place,
// read by zero readers).
//
// NOTE: the open/close wrappers here mirror treehighlight.HTML and
// treehighlight.PlainHTML; keep all three in sync if the markup
// changes.
func syntax(section *booklit.Section, language string, code booklit.Content) (booklit.Content, error) {
	chunks, err := treehighlight.Chunks(language, code.String(), treehighlight.Options{LinkReferences: true})
	if err != nil {
		return nil, err
	}

	inline := code.IsFlow()

	var body booklit.Sequence
	for _, chunk := range chunks {
		switch {
		case chunk.HTML != "":
			body = append(body, booklit.RawFragment{HTML: chunk.HTML})
		case chunk.LinkTag != "":
			body = append(body, &booklit.Reference{
				Section:  section,
				TagName:  chunk.LinkTag,
				Content:  booklit.String(chunk.LinkText),
				Optional: true,
				Location: section.InvokeLocation,
			})
		}
	}

	codeAttrs := ` class="language-` + language + `"`
	codeEl := booklit.RawElement{
		Tag:     "code",
		Attrs:   codeAttrs,
		Content: body,
	}
	if inline {
		codeEl.Attrs = ` style=";-webkit-text-size-adjust:none;" class="language-` + language + `"`
		return codeEl, nil
	}
	return booklit.RawElement{
		Tag:     "pre",
		Attrs:   ` style=";-webkit-text-size-adjust:none;"`,
		Content: codeEl,
	}, nil
}
