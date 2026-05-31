package builtins

import (
	"fmt"

	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
)

func init() {
	Register("Reference", referenceFunc)
	Register("Target", targetFunc)
}

// referenceFunc is `<Reference tag="foo">display text</Reference>`. The tag
// prop is required. Children, if any, become the display content.
func referenceFunc(ctx *Context, props map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	tag, err := requireStringProp(ctx, props, "tag", "Reference")
	if err != nil {
		return nil, err
	}

	ref := &booklit.Reference{
		Section:  ctx.Section,
		TagName:  tag,
		Location: ctx.Section.InvokeLocation,
	}

	if len(children) > 0 {
		content, err := EvaluateChildren(ctx, children)
		if err != nil {
			return nil, err
		}
		ref.Content = content
	}

	return ref, nil
}

// targetFunc is `<Target tag="foo" title="optional">rich title</Target>`.
// Mirrors baselit's `\target{tag}{title}{content}` semantics: children
// become the Title content (the display text references fall back to
// when they don't provide their own), allowing structured titles like
// `<Target tag={tag}><Syntax language="html">&lt;{tag}&gt;</Syntax></Target>`.
// The `title` prop is a shorthand for a plain-string title and loses
// to children when both are given.
func targetFunc(ctx *Context, props map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	tag, err := requireStringProp(ctx, props, "tag", "Target")
	if err != nil {
		return nil, err
	}

	ref := &booklit.Target{
		TagName:  tag,
		Location: ctx.Section.InvokeLocation,
	}

	if len(children) > 0 {
		title, err := EvaluateChildren(ctx, children)
		if err != nil {
			return nil, err
		}
		ref.Title = title
	} else if t, ok := props["title"]; ok {
		title, err := ctx.Evaluate(t)
		if err != nil {
			return nil, err
		}
		ref.Title = title
	} else {
		ref.Title = booklit.String(tag)
	}

	return ref, nil
}

// requireStringProp evaluates a named prop and stringifies it, erroring if
// the prop is missing.
func requireStringProp(ctx *Context, props map[string]ast.Node, name, component string) (string, error) {
	node, ok := props[name]
	if !ok {
		return "", fmt.Errorf("<%s> requires prop %q", component, name)
	}
	val, err := ctx.Evaluate(node)
	if err != nil {
		return "", err
	}
	return val.String(), nil
}
