package builtins

import (
	"fmt"

	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
)

func init() {
	Register("Definition", definitionFunc)
	Register("Definitions", definitionsFunc)
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
// Mirrors \definitions{}. Definition lists are JSX-only because
// CommonMark doesn't have a definition-list syntax; for plain ordered
// or unordered lists, use the Markdown forms (`- item` and `1. item`).
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
