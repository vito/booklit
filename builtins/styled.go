package builtins

import (
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
)

func init() {
	// Italic / Bold / Subscript / Superscript were redundant aliases for
	// the lowercase <em> / <strong> / <sub> / <sup> elements the markdown
	// converter now emits directly. Larger / Smaller / Strike stay because
	// their templates add inline styling (font-size, text-decoration)
	// that lowercase HTML doesn't carry; Inset / Aside stay because their
	// templates add the `class="inset"` / `class="aside"` wrappers that
	// the docs CSS hooks onto.
	Register("Larger", styled(booklit.StyleLarger))
	Register("Smaller", styled(booklit.StyleSmaller))
	Register("Strike", styled(booklit.StyleStrike))
	Register("Inset", styled(booklit.StyleInset))
	Register("Aside", styled(booklit.StyleAside))
}

// styled is the shape of every "wrap children in a Styled" built-in. The
// child content becomes the Styled's Content; no props are read.
func styled(style booklit.Style) Func {
	return func(ctx *Context, props map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
		content, err := EvaluateChildren(ctx, children)
		if err != nil {
			return nil, err
		}
		if content == nil {
			content = booklit.Empty
		}
		return booklit.Styled{
			Style:   style,
			Content: content,
		}, nil
	}
}
