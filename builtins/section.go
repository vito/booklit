package builtins

import (
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
)

func init() {
	Register("Section", sectionFunc)
}

// sectionFunc creates a sub-section from the children. Mirrors \section{}:
// children are evaluated in a fresh section context (so an inner <Title>
// sets that sub-section's title, not the parent's). The sub-section is
// appended to the parent's Children.
func sectionFunc(ctx *Context, props map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	body := ast.Sequence(children)
	section, err := ctx.Section.Processor.EvaluateNode(ctx.Section, body)
	if err != nil {
		return nil, err
	}

	section.Location = ctx.Section.InvokeLocation
	ctx.Section.Children = append(ctx.Section.Children, section)
	return nil, nil
}
