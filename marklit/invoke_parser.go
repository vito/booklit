package marklit

import (
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// invokeInlineParser parses \function-name{arg1}{arg2} syntax inline.
// It triggers on '\' and parses when followed by a lowercase letter.
// When '\' is followed by punctuation, the parser returns nil and
// goldmark handles it as a standard Markdown backslash escape.
type invokeInlineParser struct{}

var defaultInvokeInlineParser = &invokeInlineParser{}

// NewInvokeInlineParser returns a new InlineParser for \invoke syntax.
func NewInvokeInlineParser() parser.InlineParser {
	return defaultInvokeInlineParser
}

// Trigger returns '\' as the trigger character.
func (p *invokeInlineParser) Trigger() []byte {
	return []byte{'\\'}
}

// Parse parses a \function{arg1}{arg2} invocation.
func (p *invokeInlineParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	line, _ := block.PeekLine()
	if len(line) == 0 || line[0] != '\\' {
		return nil
	}

	// If next char is not a lowercase letter, bail out and let goldmark
	// handle it (e.g. \\ escape, \* escape, etc.)
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

	// Parse zero or more {arg} / {{arg}} / {{{arg}}} sequences
	for {
		line, _ = block.PeekLine()
		if len(line) == 0 || line[0] != '{' {
			break
		}

		raw, argType, consumed := parseBracedArg(block)
		if consumed == 0 {
			break
		}

		node.RawArgs = append(node.RawArgs, raw)
		node.ArgTypes = append(node.ArgTypes, argType)
	}

	return node
}

// parseBracedArg reads a single braced argument from the reader. Detects
// {{{…}}} (verbatim), {{…}} (preformatted, requires newline), and {…}
// (normal). Returns the inner content, the argument type, and a consumed
// flag (0 = failure, >0 = success).
func parseBracedArg(block text.Reader) (content []byte, argType ArgType, consumed int) {
	saveLine, savePos := block.Position()

	line, _ := block.PeekLine()

	// Check for {{{…}}} verbatim
	if len(line) >= 3 && line[0] == '{' && line[1] == '{' && line[2] == '{' {
		block.Advance(3)
		var buf []byte
		for {
			line, _ = block.PeekLine()
			if line == nil {
				block.SetPosition(saveLine, savePos)
				return nil, ArgVerbatim, 0
			}
			for i := 0; i+2 < len(line); i++ {
				if line[i] == '}' && line[i+1] == '}' && line[i+2] == '}' {
					buf = append(buf, line[:i]...)
					block.Advance(i + 3)
					return buf, ArgVerbatim, 1
				}
			}
			buf = append(buf, line...)
			block.AdvanceLine()
		}
	}

	// Normal {…} with brace depth ({{…}} without newline falls through here
	// and is treated as { followed by content starting with {, matching old
	// parser behavior)
	var buf []byte
	depth := 0
	started := false

	for {
		line, _ = block.PeekLine()
		if line == nil {
			block.SetPosition(saveLine, savePos)
			return nil, ArgNormal, 0
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
					block.Advance(i + 1)
					return buf, ArgNormal, 1
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
