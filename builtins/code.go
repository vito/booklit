package builtins

import (
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/baselit"
)

func init() {
	Register("CodeBlock", codeBlockFunc)
	Register("Syntax", codeBlockFunc)
}

// codeBlockFunc — `<CodeBlock language="go">code</CodeBlock>` or
// `<Syntax language="go" style="monokai">code</Syntax>`. Tree-sitter-driven
// syntax highlighting; delegates to baselit's existing implementation.
func codeBlockFunc(ctx *Context, props map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	language, err := requireStringProp(ctx, props, "language", "CodeBlock")
	if err != nil {
		return nil, err
	}

	code, err := EvaluateChildren(ctx, children)
	if err != nil {
		return nil, err
	}
	if code == nil {
		code = booklit.Empty
	}

	base := baselit.NewPlugin(ctx.Section)

	if s, ok := props["style"]; ok {
		styleContent, err := ctx.Evaluate(s)
		if err != nil {
			return nil, err
		}
		return base.Syntax(language, code, styleContent.String())
	}
	return base.Syntax(language, code)
}
