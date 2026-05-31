package marklit

import (
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// jsxInlineParser parses JSX-style invocations: <Name attr="x" foo={expr}>
// children</Name> or <Name.../>. It triggers on '<' and only claims the
// input when the following byte is an uppercase ASCII letter (lowercase
// tags fall through to goldmark's raw-HTML handling).
type jsxInlineParser struct{}

var defaultJSXInlineParser = &jsxInlineParser{}

// NewJSXInlineParser returns a new InlineParser for JSX syntax.
func NewJSXInlineParser() parser.InlineParser {
	return defaultJSXInlineParser
}

// Trigger returns '<' as the trigger character.
func (p *jsxInlineParser) Trigger() []byte {
	return []byte{'<'}
}

// Parse parses one JSX element. On failure the reader position is restored.
func (p *jsxInlineParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	line, _ := block.PeekLine()
	if len(line) < 2 || line[0] != '<' || !isUpperAlpha(line[1]) {
		return nil
	}

	saveLine, savePos := block.Position()

	s := &jsxScanner{block: block}
	node, ok := parseJSXElement(s)
	if !ok {
		block.SetPosition(saveLine, savePos)
		return nil
	}
	return node
}

// jsxScanner is a byte-at-a-time view over a goldmark text.Reader. Each
// next/peek operation crosses line boundaries automatically.
type jsxScanner struct {
	block text.Reader
}

func (s *jsxScanner) peek() (byte, bool) {
	for {
		line, _ := s.block.PeekLine()
		if line == nil {
			return 0, false
		}
		if len(line) == 0 {
			s.block.AdvanceLine()
			continue
		}
		return line[0], true
	}
}

func (s *jsxScanner) next() (byte, bool) {
	for {
		line, _ := s.block.PeekLine()
		if line == nil {
			return 0, false
		}
		if len(line) == 0 {
			s.block.AdvanceLine()
			continue
		}
		b := line[0]
		if len(line) == 1 {
			s.block.AdvanceLine()
		} else {
			s.block.Advance(1)
		}
		return b, true
	}
}

func (s *jsxScanner) skipWS() {
	for {
		b, ok := s.peek()
		if !ok {
			return
		}
		if b == ' ' || b == '\t' || b == '\n' || b == '\r' {
			s.next()
			continue
		}
		return
	}
}

// parseJSXElement parses a single <Name...> element starting at the current
// reader position. On success the reader is left just past the closing '>'.
// On failure the reader position is undefined; the outer Parse restores it.
func parseJSXElement(s *jsxScanner) (*JSXElementNode, bool) {
	b, ok := s.next()
	if !ok || b != '<' {
		return nil, false
	}

	name, ok := readJSXName(s, true)
	if !ok {
		return nil, false
	}

	node := NewJSXElementNode(name)

	for {
		// Allow newlines between attributes for multi-line tags.
		s.skipWS()
		b, ok := s.peek()
		if !ok {
			return nil, false
		}
		if b == '/' || b == '>' {
			break
		}
		attr, ok := parseAttribute(s)
		if !ok {
			return nil, false
		}
		node.Props = append(node.Props, attr)
	}

	b, ok = s.next()
	if !ok {
		return nil, false
	}
	if b == '/' {
		next, ok := s.next()
		if !ok || next != '>' {
			return nil, false
		}
		node.SelfClosing = true
		return node, true
	}
	if b != '>' {
		return nil, false
	}

	if !parseChildren(s, name, &node.Children) {
		return nil, false
	}
	return node, true
}

// readJSXName reads a tag or attribute name. If requireUpper is true the
// first byte must be uppercase ASCII (component names).
func readJSXName(s *jsxScanner, requireUpper bool) (string, bool) {
	b, ok := s.peek()
	if !ok {
		return "", false
	}
	if requireUpper {
		if !isUpperAlpha(b) {
			return "", false
		}
	} else if !isAlpha(b) {
		return "", false
	}
	var buf []byte
	for {
		b, ok := s.peek()
		if !ok {
			break
		}
		if isAlpha(b) || isDigit(b) || b == '_' {
			buf = append(buf, b)
			s.next()
			continue
		}
		break
	}
	return string(buf), true
}

// parseAttribute parses a `name="value"` or `name={expr}` attribute.
func parseAttribute(s *jsxScanner) (JSXProp, bool) {
	name, ok := readJSXName(s, false)
	if !ok {
		return JSXProp{}, false
	}
	b, ok := s.next()
	if !ok || b != '=' {
		return JSXProp{}, false
	}
	b, ok = s.peek()
	if !ok {
		return JSXProp{}, false
	}
	switch b {
	case '"':
		s.next()
		value, ok := readQuoted(s, '"')
		if !ok {
			return JSXProp{}, false
		}
		return JSXProp{Name: name, Kind: JSXPropString, Value: value}, true
	case '{':
		s.next()
		value, ok := readBraceExpr(s)
		if !ok {
			return JSXProp{}, false
		}
		return JSXProp{Name: name, Kind: JSXPropExpression, Value: value}, true
	}
	return JSXProp{}, false
}

// readQuoted reads bytes up to terminator (a single byte). Backslash escapes
// are preserved verbatim in the output (unescaping happens at convert time).
// Newlines inside attribute strings are not permitted.
func readQuoted(s *jsxScanner, term byte) ([]byte, bool) {
	var buf []byte
	for {
		b, ok := s.next()
		if !ok {
			return nil, false
		}
		if b == '\\' {
			n, ok := s.next()
			if !ok {
				return nil, false
			}
			buf = append(buf, b, n)
			continue
		}
		if b == '\n' || b == '\r' {
			return nil, false
		}
		if b == term {
			return buf, true
		}
		buf = append(buf, b)
	}
}

// readBraceExpr reads bytes between { and the matching }, with brace depth.
// The opening { has already been consumed. Double- and single-quoted string
// literals are tracked so braces inside string contents don't close the
// expression. This is a best-effort scan for MVP; full Dang parsing happens
// later.
func readBraceExpr(s *jsxScanner) ([]byte, bool) {
	var buf []byte
	depth := 1
	for {
		b, ok := s.next()
		if !ok {
			return nil, false
		}
		if b == '"' || b == '\'' {
			buf = append(buf, b)
			term := b
			for {
				c, ok := s.next()
				if !ok {
					return nil, false
				}
				buf = append(buf, c)
				if c == '\\' {
					n, ok := s.next()
					if !ok {
						return nil, false
					}
					buf = append(buf, n)
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
				return buf, true
			}
			buf = append(buf, b)
			continue
		}
		buf = append(buf, b)
	}
}

// parseChildren reads element body until a matching </parentName> close tag.
// Children are partitioned into text chunks, nested JSX elements, and
// {expression} captures.
func parseChildren(s *jsxScanner, parentName string, out *[]JSXChild) bool {
	var textBuf []byte
	flushText := func() {
		if len(textBuf) > 0 {
			*out = append(*out, JSXChild{Kind: JSXChildText, Text: textBuf})
			textBuf = nil
		}
	}
	for {
		line, _ := s.block.PeekLine()
		if line == nil {
			return false
		}
		if len(line) == 0 {
			s.block.AdvanceLine()
			continue
		}
		c := line[0]
		if c == '<' {
			if len(line) >= 2 && line[1] == '/' {
				tagLen, name, ok := scanCloseTag(line)
				if !ok {
					return false
				}
				if name != parentName {
					return false
				}
				flushText()
				s.block.Advance(tagLen)
				return true
			}
			if len(line) >= 2 && isUpperAlpha(line[1]) {
				flushText()
				child, ok := parseJSXElement(s)
				if !ok {
					return false
				}
				*out = append(*out, JSXChild{Kind: JSXChildElement, Elem: child})
				continue
			}
			// '<' followed by anything else — literal text (raw HTML, etc.)
			textBuf = append(textBuf, '<')
			s.next()
			continue
		}
		if c == '{' {
			flushText()
			s.next()
			expr, ok := readBraceExpr(s)
			if !ok {
				return false
			}
			*out = append(*out, JSXChild{Kind: JSXChildExpression, Text: expr})
			continue
		}
		// Run of plain text until the next '<' or '{' or end of line.
		i := 0
		for i < len(line) && line[i] != '<' && line[i] != '{' {
			i++
		}
		textBuf = append(textBuf, line[:i]...)
		if i < len(line) {
			s.block.Advance(i)
		} else {
			s.block.AdvanceLine()
		}
	}
}

// scanCloseTag examines a buffer starting with `</` and reports the byte
// length of the close tag (through the `>`) and the tag name. Returns
// ok=false if the buffer doesn't contain a valid close tag.
func scanCloseTag(line []byte) (length int, name string, ok bool) {
	if len(line) < 4 || line[0] != '<' || line[1] != '/' {
		return 0, "", false
	}
	if !isUpperAlpha(line[2]) {
		return 0, "", false
	}
	i := 2
	for i < len(line) && (isAlpha(line[i]) || isDigit(line[i])) {
		i++
	}
	if i >= len(line) || line[i] != '>' {
		return 0, "", false
	}
	return i + 1, string(line[2:i]), true
}

func isUpperAlpha(c byte) bool { return c >= 'A' && c <= 'Z' }
func isAlpha(c byte) bool      { return isLowerAlpha(c) || isUpperAlpha(c) }
func isDigit(c byte) bool      { return c >= '0' && c <= '9' }
