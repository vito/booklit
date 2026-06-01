package marklit

import (
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// jsxExprInlineParser parses top-level `{expr}` interpolations in paragraph
// content. The brace-balanced scanner skips over double- and single-quoted
// string literals so that `{ "a}b" }` parses as one expression. Returns nil
// (allowing goldmark to keep the byte as text) on any failure.
type jsxExprInlineParser struct{}

var defaultJSXExprInlineParser = &jsxExprInlineParser{}

// NewJSXExpressionInlineParser returns a new InlineParser for `{expr}`
// syntax in MarkDangJSX paragraph content.
func NewJSXExpressionInlineParser() parser.InlineParser {
	return defaultJSXExprInlineParser
}

func (p *jsxExprInlineParser) Trigger() []byte {
	return []byte{'{'}
}

func (p *jsxExprInlineParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	saveLine, savePos := block.Position()

	line, _ := block.PeekLine()
	if len(line) < 2 || line[0] != '{' {
		return nil
	}

	s := &jsxScanner{block: block}
	if b, ok := s.next(); !ok || b != '{' {
		block.SetPosition(saveLine, savePos)
		return nil
	}

	var buf []byte
	depth := 1
	for {
		b, ok := s.next()
		if !ok {
			block.SetPosition(saveLine, savePos)
			return nil
		}
		if b == '"' || b == '\'' {
			term := b
			buf = append(buf, b)
			for {
				c, ok := s.next()
				if !ok {
					block.SetPosition(saveLine, savePos)
					return nil
				}
				buf = append(buf, c)
				if c == '\\' {
					esc, ok := s.next()
					if !ok {
						block.SetPosition(saveLine, savePos)
						return nil
					}
					buf = append(buf, esc)
					continue
				}
				if c == term {
					break
				}
			}
			continue
		}
		if b == '{' {
			depth++
			buf = append(buf, b)
			continue
		}
		if b == '}' {
			depth--
			if depth == 0 {
				return NewJSXExpressionInlineNode(string(buf))
			}
			buf = append(buf, b)
			continue
		}
		buf = append(buf, b)
	}
}

// jsxExprInlineNode is a goldmark AST inline node carrying a Dang
// expression text. It is converted to `ast.JSXExpression` by the
// converter.
type jsxExprInlineNode struct {
	gast.BaseInline
	Raw string
}

// KindJSXExprInline is the NodeKind for inline `{expr}` interpolations.
var KindJSXExprInline = gast.NewNodeKind("BooklitJSXExprInline")

func (n *jsxExprInlineNode) Kind() gast.NodeKind { return KindJSXExprInline }

func (n *jsxExprInlineNode) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, map[string]string{"Raw": n.Raw}, nil)
}

func NewJSXExpressionInlineNode(raw string) *jsxExprInlineNode {
	return &jsxExprInlineNode{Raw: raw}
}
