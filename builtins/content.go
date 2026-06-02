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

// codeFunc — `<Code>x</Code>`. Inline-or-block verbatim content; the
// shape of the body decides. Flow content (a short snippet, a phrase)
// wraps in `<code>`; block content (a multi-paragraph or other block
// container) wraps in `<pre><code>` so the result is valid HTML and
// renders monospaced as a block.
func codeFunc(ctx *Context, _ map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	content, err := EvaluateChildren(ctx, children)
	if err != nil {
		return nil, err
	}
	if content == nil {
		content = booklit.Empty
	}
	if content.IsFlow() {
		return booklit.RawElement{Tag: "code", Content: content}, nil
	}
	return booklit.RawElement{
		Tag:     "pre",
		Content: booklit.RawElement{Tag: "code", Content: content},
	}, nil
}

// linkFunc — `<Link target="url">content</Link>`. Lowers to the lowercase
// `<a href="url">…</a>` JSX shape so the standard raw-HTML dispatcher
// handles attribute escaping; the PascalCase form exists for parity with
// the rest of the builtin API.
func linkFunc(ctx *Context, props map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	target, err := requireStringProp(ctx, props, "target", "Link")
	if err != nil {
		return nil, err
	}
	if len(children) == 0 {
		children = []ast.Node{ast.String(target)}
	}
	return ctx.Evaluate(ast.JSXElement{
		Name:     "a",
		Props:    map[string]ast.Node{"href": ast.String(target)},
		Children: children,
	})
}

// imageFunc — `<Image path="..." description="..."/>`. Lowers to
// `<img src="..." alt="..."/>`.
func imageFunc(ctx *Context, props map[string]ast.Node, _ []ast.Node) (booklit.Content, error) {
	path, err := requireStringProp(ctx, props, "path", "Image")
	if err != nil {
		return nil, err
	}
	imgProps := map[string]ast.Node{
		"src": ast.String(path),
		"alt": ast.String(""),
	}
	if d, ok := props["description"]; ok {
		desc, err := ctx.Evaluate(d)
		if err != nil {
			return nil, err
		}
		imgProps["alt"] = ast.String(desc.String())
	}
	return ctx.Evaluate(ast.JSXElement{
		Name:  "img",
		Props: imgProps,
	})
}
