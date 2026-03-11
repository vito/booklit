package marklit

import (
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// referenceInlineParser parses [#tag] shorthand into a \reference{tag}
// invocation. It triggers on '[' and matches [#tag] where tag is composed
// of alphanumeric characters, hyphens, and underscores.
type referenceInlineParser struct{}

var defaultReferenceInlineParser = &referenceInlineParser{}

// NewReferenceInlineParser returns a new InlineParser for [#tag] syntax.
func NewReferenceInlineParser() parser.InlineParser {
	return defaultReferenceInlineParser
}

// Trigger returns '[' as the trigger character.
func (p *referenceInlineParser) Trigger() []byte {
	return []byte{'['}
}

// Parse parses a [#tag] reference shorthand.
func (p *referenceInlineParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	line, _ := block.PeekLine()
	if len(line) < 4 || line[0] != '[' || line[1] != '#' {
		return nil
	}

	// Find the closing ]
	i := 2
	for i < len(line) && line[i] != ']' {
		if line[i] == '\n' || line[i] == '[' {
			return nil
		}
		i++
	}
	if i >= len(line) || line[i] != ']' {
		return nil
	}

	tag := string(line[2:i])
	if len(tag) == 0 {
		return nil
	}

	// Make sure this is not followed by '(' which would make it a regular
	// link like [#foo](url) — let goldmark handle that case.
	if i+1 < len(line) && line[i+1] == '(' {
		return nil
	}

	block.Advance(i + 1) // consume [#tag]

	node := NewInvokeNode("reference")
	node.RawArgs = [][]byte{[]byte(tag)}
	node.ArgTypes = []ArgType{ArgNormal}
	return node
}
