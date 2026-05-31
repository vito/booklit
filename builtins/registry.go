// Package builtins is the first dispatch tier for JSX elements. Each
// built-in is a Go function bound to a PascalCase name (e.g. "Title",
// "Section"); the evaluator consults this registry before falling back to
// template-default rendering.
package builtins

import (
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
)

// Context carries everything a built-in might need to do its work: the
// current section (some built-ins mutate the section tree) and an Evaluate
// helper that turns an arbitrary ast.Node into Content (so built-ins can
// recursively evaluate props and children when they choose to).
type Context struct {
	Section  *booklit.Section
	Evaluate func(ast.Node) (booklit.Content, error)
}

// Func is the shape of a built-in. Props and children are passed as raw
// AST so built-ins can decide whether to evaluate them; built-ins that
// just want content should call ctx.Evaluate on each.
type Func func(ctx *Context, props map[string]ast.Node, children []ast.Node) (booklit.Content, error)

var registry = map[string]Func{}

// Register binds a name to a built-in implementation. Called from each
// builtin file's init() function.
func Register(name string, fn Func) {
	registry[name] = fn
}

// Lookup returns the registered built-in for name, if any.
func Lookup(name string) (Func, bool) {
	fn, ok := registry[name]
	return fn, ok
}

// EvaluateChildren is a convenience for built-ins that just want to
// evaluate their entire children list into a single Content (typically
// a Sequence). Returns nil if children is empty.
func EvaluateChildren(ctx *Context, children []ast.Node) (booklit.Content, error) {
	if len(children) == 0 {
		return nil, nil
	}
	return ctx.Evaluate(ast.Sequence(children))
}
