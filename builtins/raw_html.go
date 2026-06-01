package builtins

import (
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
)

func init() {
	Register("RawHTML", rawHTMLFunc)
}

// rawHTMLFunc — `<RawHTML>...literal HTML...</RawHTML>`. Stringifies
// the body and wraps it in a booklit.RawFragment so the renderer writes
// the bytes through unescaped.
//
// Primarily used by the markdown converter as the carrier for inline
// raw HTML (`<!-- comments -->`, bare `<br>`) and by the mdx template
// engine to emit inert HTML scaffolding between JSX components. Also
// available to content authors as an escape hatch when markdown can't
// express a particular bit of HTML.
//
// The legacy `block` prop is gone: RawFragment is always flow. Block
// raw-HTML use cases (e.g. block-level HTML comments) survive — the
// paragraph layout wraps them in `<p>`, which is harmless for invisible
// markup like comments. Any case that truly needs block emission should
// be written as a lowercase JSX element (`<div>...</div>`) and routed
// through dispatchRawHTML → RawElement, whose IsFlow comes from the
// tag.
func rawHTMLFunc(ctx *Context, props map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	content, err := EvaluateChildren(ctx, children)
	if err != nil {
		return nil, err
	}
	if content == nil {
		return booklit.RawFragment{}, nil
	}
	return booklit.RawFragment{HTML: content.String()}, nil
}
