package marklit

import (
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// invokeInlineParser parses @function-name{arg1}{arg2} syntax inline.
type invokeInlineParser struct{}

var defaultInvokeInlineParser = &invokeInlineParser{}

// NewInvokeInlineParser returns a new InlineParser for @invoke syntax.
func NewInvokeInlineParser() parser.InlineParser {
	return defaultInvokeInlineParser
}

// Trigger returns '@' as the trigger character.
func (p *invokeInlineParser) Trigger() []byte {
	return []byte{'@'}
}

// Parse parses an @function{arg1}{arg2} invocation.
func (p *invokeInlineParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	line, segment := block.PeekLine()
	if len(line) == 0 || line[0] != '@' {
		return nil
	}

	// Check for @@ escape sequence
	if len(line) > 1 && line[1] == '@' {
		block.Advance(2)
		return gast.NewTextSegment(text.NewSegment(segment.Start+1, segment.Start+2))
	}

	// Parse function name: @[a-z][a-z0-9-]*
	i := 1
	if i >= len(line) || !isLowerAlpha(line[i]) {
		return nil
	}
	for i < len(line) && isNameChar(line[i]) {
		i++
	}

	funcName := string(line[1:i])
	block.Advance(i)

	node := NewInvokeNode(funcName)

	// Parse zero or more {arg} sequences
	for {
		line, _ = block.PeekLine()
		if len(line) == 0 || line[0] != '{' {
			break
		}

		raw, consumed := parseBracedArg(block)
		if consumed == 0 {
			break
		}

		node.RawArgs = append(node.RawArgs, raw)
	}

	return node
}

// parseBracedArg reads a single {...} argument from the reader, handling
// nested braces. Returns the inner content (without outer braces) and the
// total bytes consumed. Returns 0 consumed if the braces are unbalanced.
func parseBracedArg(block text.Reader) (content []byte, consumed int) {
	// We need to read potentially across multiple lines to find the
	// matching close brace.
	saveLine, savePos := block.Position()

	var buf []byte
	depth := 0
	started := false

	for {
		line, _ := block.PeekLine()
		if line == nil {
			// EOF without matching close brace — rollback
			block.SetPosition(saveLine, savePos)
			return nil, 0
		}

		for i := 0; i < len(line); i++ {
			ch := line[i]
			if ch == '{' {
				depth++
				if !started {
					started = true
					continue // skip the opening brace
				}
			} else if ch == '}' {
				depth--
				if depth == 0 {
					// Found matching close brace
					block.Advance(i + 1)
					return buf, 1 // consumed > 0 signals success
				}
			}
			if started {
				buf = append(buf, ch)
			}
		}

		block.AdvanceLine()
	}
}

func isLowerAlpha(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func isNameChar(c byte) bool {
	return isLowerAlpha(c) || (c >= '0' && c <= '9') || c == '-'
}
