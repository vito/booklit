package builtins

import (
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
)

func init() {
	Register("Italic", styled(booklit.StyleItalic))
	Register("Bold", styled(booklit.StyleBold))
	Register("Larger", styled(booklit.StyleLarger))
	Register("Smaller", styled(booklit.StyleSmaller))
	Register("Strike", styled(booklit.StyleStrike))
	Register("Superscript", styled(booklit.StyleSuperscript))
	Register("Subscript", styled(booklit.StyleSubscript))
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
