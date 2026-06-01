package builtins

import (
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
)

func init() {
	Register("RawHTML", rawHTMLFunc)
}

// rawHTMLFunc — `<RawHTML>...literal HTML...</RawHTML>`. Wraps its body
// in `Styled{Style: "raw-html"}` so the renderer's raw-html template
// passes the bytes through unescaped.
//
// Primarily used by the mdx template engine (templates/) to emit the
// inert HTML scaffolding between JSX components in a `.md` template,
// and by the markdown converter as the carrier for inline / block raw
// HTML (`<!-- comments -->`, bare `<br>`). It's also available to content
// authors as an escape hatch when markdown can't express a particular
// bit of HTML.
//
// The `block` prop flags the wrapped content as block-level so the
// surrounding paragraph layout skips its `<p>` wrap.
func rawHTMLFunc(ctx *Context, props map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	content, err := EvaluateChildren(ctx, children)
	if err != nil {
		return nil, err
	}
	if content == nil {
		content = booklit.Empty
	}
	block := false
	if b, ok := props["block"]; ok {
		v, err := ctx.Evaluate(b)
		if err != nil {
			return nil, err
		}
		block = v.String() == "true"
	}
	return booklit.Styled{Style: "raw-html", Block: block, Content: content}, nil
}
