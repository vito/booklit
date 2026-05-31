package builtins

import (
	"fmt"

	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/dangeval"
)

func init() {
	Register("Children", childrenFunc)
}

// childrenFunc — `<Children/>`. Looks up the `children` binding that the
// mdx template engine attached to the current Dang scope and emits its
// content at the call site.
//
// Equivalent in shape to `{children}` as a JSX expression interpolation,
// but usable in JSX child position where an expression sometimes reads
// awkwardly. Errors when called outside an mdx template (no `children`
// binding in scope).
func childrenFunc(ctx *Context, _ map[string]ast.Node, _ []ast.Node) (booklit.Content, error) {
	if ctx.Dang == nil {
		return nil, fmt.Errorf("<Children/> requires a Dang evaluator: not configured")
	}
	val, ok, err := ctx.Dang.LookupValue("children")
	if err != nil {
		return nil, fmt.Errorf("<Children/> lookup: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("<Children/> used outside of an mdx template (no children binding)")
	}
	return dangeval.ToContent(val)
}
