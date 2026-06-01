package marklit

import (
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// jsxBlockParser claims block-level JSX elements, including lowercase
// raw-HTML tags like `<div>`. With goldmark's HTMLBlockParser disabled
// (see newParser), this is the only thing that recognizes `<x...>` at
// the start of a line.
//
// The block parser consumes the full element (possibly across multiple
// lines) via the same scanner used by the inline parser. Goldmark's rule
// that "Open must not parse beyond the current line" is in tension with
// this; in practice consuming the whole element in Open and reporting
// Close from Continue works because the block has no further children to
// gather.
type jsxBlockParser struct{}

var defaultJSXBlockParser = &jsxBlockParser{}

// NewJSXBlockParser returns a new BlockParser for block-level JSX.
func NewJSXBlockParser() parser.BlockParser {
	return defaultJSXBlockParser
}

func (p *jsxBlockParser) Trigger() []byte {
	return []byte{'<'}
}

func (p *jsxBlockParser) Open(parent gast.Node, reader text.Reader, pc parser.Context) (gast.Node, parser.State) {
	line, _ := reader.PeekLine()
	pos := pc.BlockOffset()
	if pos < 0 || pos+1 >= len(line) || line[pos] != '<' || !isAlpha(line[pos+1]) {
		return nil, parser.NoChildren
	}

	saveLine, savePos := reader.Position()
	reader.Advance(pos)

	s := &jsxScanner{block: reader}
	inline, ok := parseJSXElement(s)
	if !ok {
		reader.SetPosition(saveLine, savePos)
		return nil, parser.NoChildren
	}
	// If the element ended on the line it started on and there's still
	// non-whitespace text after it, the element is part of a flow of
	// inline content (e.g. `<Italic>x</Italic> and <Bold>y</Bold>`).
	// Bail back to inline parsing — claiming a block here would let the
	// remainder of the line silently fall off.
	if !inline.MultiLine && hasTrailingContent(reader) {
		reader.SetPosition(saveLine, savePos)
		return nil, parser.NoChildren
	}
	block := NewJSXBlockElementNode(inline.Name)
	block.Props = inline.Props
	block.SelfClosing = inline.SelfClosing
	block.MultiLine = inline.MultiLine
	block.Children = inline.Children
	block.Line = inline.Line
	block.Col = inline.Col
	return block, parser.NoChildren
}

// hasTrailingContent reports whether the rest of the current line in
// reader contains anything besides whitespace.
func hasTrailingContent(reader text.Reader) bool {
	line, _ := reader.PeekLine()
	for _, b := range line {
		if b == '\n' {
			return false
		}
		if b != ' ' && b != '\t' && b != '\r' {
			return true
		}
	}
	return false
}

func (p *jsxBlockParser) Continue(node gast.Node, reader text.Reader, pc parser.Context) parser.State {
	return parser.Close
}

func (p *jsxBlockParser) Close(node gast.Node, reader text.Reader, pc parser.Context) {}

func (p *jsxBlockParser) CanInterruptParagraph() bool { return false }
func (p *jsxBlockParser) CanAcceptIndentedLine() bool { return false }
