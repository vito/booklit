package builtins

import (
	"fmt"

	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
)

func init() {
	Register("List", listFunc(false))
	Register("OrderedList", listFunc(true))
	Register("Item", itemFunc)
	Register("Definition", definitionFunc)
	Register("Definitions", definitionsFunc)
}

// listFunc — `<List><Item>a</Item><Item>b</Item></List>`. Each <Item>
// child contributes one entry; non-item children (whitespace text
// between items) are dropped.
func listFunc(ordered bool) Func {
	return func(ctx *Context, _ map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
		var items []booklit.Content
		for _, child := range children {
			val, err := ctx.Evaluate(child)
			if err != nil {
				return nil, err
			}
			if val == nil {
				continue
			}
			item, ok := val.(itemContent)
			if !ok {
				continue
			}
			items = append(items, item.Content)
		}
		return booklit.List{Items: items, Ordered: ordered}, nil
	}
}

// itemContent wraps a single list/definition entry so the parent
// container can distinguish authored items from incidental whitespace
// or stray text.
type itemContent struct {
	Content booklit.Content
}

func (i itemContent) IsFlow() bool                          { return i.Content.IsFlow() }
func (i itemContent) String() string                        { return i.Content.String() }
func (i itemContent) Visit(v booklit.Visitor) error         { return i.Content.Visit(v) }

// itemFunc — `<Item>content</Item>`. Carries one entry for <List> or
// <OrderedList>; outside those it just renders as its content.
func itemFunc(ctx *Context, _ map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	content, err := EvaluateChildren(ctx, children)
	if err != nil {
		return nil, err
	}
	if content == nil {
		content = booklit.Empty
	}
	return itemContent{Content: content}, nil
}

// definitionFunc — `<Definition term="x">body</Definition>`. Produces a
// 2-element List that <Definitions> recognizes; mirrors \definition{}{}.
func definitionFunc(ctx *Context, props map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	t, ok := props["term"]
	if !ok {
		return nil, fmt.Errorf("<Definition> requires prop \"term\"")
	}
	term, err := ctx.Evaluate(t)
	if err != nil {
		return nil, err
	}
	body, err := EvaluateChildren(ctx, children)
	if err != nil {
		return nil, err
	}
	if body == nil {
		body = booklit.Empty
	}
	return booklit.List{Items: []booklit.Content{term, body}}, nil
}

// definitionsFunc — `<Definitions><Definition .../></Definitions>`.
// Mirrors \definitions{}.
func definitionsFunc(ctx *Context, _ map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	defs := booklit.Definitions{}
	for _, child := range children {
		val, err := ctx.Evaluate(child)
		if err != nil {
			return nil, err
		}
		if val == nil {
			continue
		}
		list, ok := val.(booklit.List)
		if !ok {
			continue
		}
		if len(list.Items) != 2 {
			return nil, fmt.Errorf("<Definitions> entry has %d items, expected 2", len(list.Items))
		}
		defs = append(defs, booklit.Definition{
			Subject:    list.Items[0],
			Definition: list.Items[1],
		})
	}
	return defs, nil
}
