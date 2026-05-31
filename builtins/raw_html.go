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
// inert HTML scaffolding between JSX components in a `.md` template, but
// it's also available to content authors as an escape hatch when
// markdown can't express a particular bit of HTML.
func rawHTMLFunc(ctx *Context, _ map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	content, err := EvaluateChildren(ctx, children)
	if err != nil {
		return nil, err
	}
	if content == nil {
		content = booklit.Empty
	}
	return booklit.Styled{Style: "raw-html", Content: content}, nil
}
