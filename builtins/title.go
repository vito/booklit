package builtins

import (
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
)

func init() {
	Register("Title", titleFunc)
}

// titleFunc is `<Title tag="optional">content</Title>`. Children become the
// title content. The tag prop becomes an explicit tag string. Multiple
// titles per section are still rejected, mirroring \title{}.
func titleFunc(ctx *Context, props map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	if ctx.Section.Title != booklit.Empty {
		return nil, booklit.TitleTwiceError{
			TitleLocation: booklit.ErrorLocation{
				FilePath:     ctx.Section.FilePath(),
				NodeLocation: ctx.Section.TitleLocation,
				Length:       len("<Title>"),
			},
			ErrorLocation: booklit.ErrorLocation{
				FilePath:     ctx.Section.FilePath(),
				NodeLocation: ctx.Section.InvokeLocation,
				Length:       len("<Title>"),
			},
		}
	}

	title, err := EvaluateChildren(ctx, children)
	if err != nil {
		return nil, err
	}
	if title == nil {
		title = booklit.Empty
	}

	var tags []string
	if t, ok := props["tag"]; ok {
		tagContent, err := ctx.Evaluate(t)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tagContent.String())
	}

	ctx.Section.SetTitle(title, ctx.Section.InvokeLocation, tags...)
	return nil, nil
}
