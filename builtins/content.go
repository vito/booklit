package builtins

import (
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
)

func init() {
	Register("Aux", auxFunc)
	Register("Code", codeFunc)
	Register("Link", linkFunc)
	Register("Image", imageFunc)
	Register("ThematicBreak", thematicBreakFunc)
}

// auxFunc — `<Aux>content</Aux>`. Auxiliary content; mirrors \aux{}.
func auxFunc(ctx *Context, _ map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	content, err := EvaluateChildren(ctx, children)
	if err != nil {
		return nil, err
	}
	if content == nil {
		content = booklit.Empty
	}
	return booklit.Aux{Content: content}, nil
}

// codeFunc — `<Code>x</Code>`. Inline verbatim content; mirrors \code{}.
func codeFunc(ctx *Context, _ map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	content, err := EvaluateChildren(ctx, children)
	if err != nil {
		return nil, err
	}
	if content == nil {
		content = booklit.Empty
	}
	return booklit.Styled{Content: content, Style: booklit.StyleVerbatim}, nil
}

// linkFunc — `<Link target="url">content</Link>`. Mirrors \link{content}{target}.
func linkFunc(ctx *Context, props map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	target, err := requireStringProp(ctx, props, "target", "Link")
	if err != nil {
		return nil, err
	}
	content, err := EvaluateChildren(ctx, children)
	if err != nil {
		return nil, err
	}
	if content == nil {
		content = booklit.String(target)
	}
	return booklit.Link{Content: content, Target: target}, nil
}

// imageFunc — `<Image path="..." description="..."/>`. Mirrors
// \image{path}{description}.
func imageFunc(ctx *Context, props map[string]ast.Node, _ []ast.Node) (booklit.Content, error) {
	path, err := requireStringProp(ctx, props, "path", "Image")
	if err != nil {
		return nil, err
	}
	img := booklit.Image{Path: path}
	if d, ok := props["description"]; ok {
		desc, err := ctx.Evaluate(d)
		if err != nil {
			return nil, err
		}
		img.Description = desc.String()
	}
	return img, nil
}

// thematicBreakFunc — `<ThematicBreak/>`. Emitted by the markdown parser
// for `---` horizontal rules; renders via a `thematic-break` template if
// one exists.
func thematicBreakFunc(_ *Context, _ map[string]ast.Node, _ []ast.Node) (booklit.Content, error) {
	return booklit.Styled{
		Style: "thematic-break",
		Block: true,
	}, nil
}
