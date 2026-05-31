package templates

import "github.com/vito/booklit/ast"

// rawHTML is the ast.Node a template parser emits for the bytes between
// JSX elements and {expressions}. Visiting it dispatches as a JSX element
// named `RawHTML` whose body is the literal text; the matching built-in
// (builtins/raw_html.go) wraps it in `Styled{Style: "raw-html"}` so the
// HTML engine's raw-html template passes the bytes through unchanged.
//
// Templates emit this directly rather than relying on `\raw-html` invokes
// or markdown processing — they are HTML scaffolding, not prose.
type rawHTML struct {
	text string
	loc  ast.Location
}

// Visit dispatches as a synthetic <RawHTML> JSX element.
func (r rawHTML) Visit(v ast.Visitor) error {
	return ast.JSXElement{
		Name:     "RawHTML",
		Children: []ast.Node{ast.String(r.text)},
		Location: r.loc,
	}.Visit(v)
}
