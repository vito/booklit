package templates

import (
	"fmt"

	"github.com/vito/booklit/ast"
)

// Parse parses a template source. The template body is a free mix of:
//
//   - raw HTML text (passed through to the renderer unchanged)
//   - JSX components, <Pascal ...>...</Pascal> or <Pascal .../>, which
//     dispatch via the JSX evaluator at evaluation time
//   - {expr} interpolations, evaluated as Dang expressions
//
// Markdown is NOT processed: templates are HTML scaffolding around prop
// holes, not prose. Authors who want Markdown can wrap a `{children}` in
// their own JSX wrapper components.
func Parse(source []byte) (ast.Node, error) {
	p := &parser{src: source}
	return p.parseTopLevel()
}

// parser tokenizes a template source byte-by-byte. It tracks a 1-based
// line/column for error messages.
type parser struct {
	src  []byte
	pos  int
	line int
	col  int
}

func (p *parser) loc() ast.Location {
	line := p.line
	col := p.col
	if line == 0 {
		line = 1
		col = 1
	}
	return ast.Location{Line: line, Col: col, Offset: p.pos}
}

func (p *parser) advance(n int) {
	for i := 0; i < n && p.pos < len(p.src); i++ {
		if p.src[p.pos] == '\n' {
			p.line++
			p.col = 1
		} else {
			p.col++
		}
		p.pos++
	}
}

func (p *parser) initLoc() {
	if p.line == 0 {
		p.line = 1
		p.col = 1
	}
}

func (p *parser) parseTopLevel() (ast.Node, error) {
	p.initLoc()
	nodes, err := p.parseChildren("")
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 0:
		return ast.String(""), nil
	case 1:
		return nodes[0], nil
	default:
		return ast.Sequence(nodes), nil
	}
}

// parseChildren reads child nodes until it hits </parentName> (or EOF if
// parentName is empty). Raw text chunks become rawHTML nodes; <Pascal>
// elements parse recursively; {expr} become ast.JSXExpression.
func (p *parser) parseChildren(parentName string) ([]ast.Node, error) {
	var children []ast.Node
	var textBuf []byte
	textStart := p.loc()

	flushText := func() {
		if len(textBuf) > 0 {
			children = append(children, rawHTML{text: string(textBuf), loc: textStart})
			textBuf = nil
		}
		textStart = p.loc()
	}

	for p.pos < len(p.src) {
		c := p.src[p.pos]
		if c == '<' && p.pos+1 < len(p.src) {
			next := p.src[p.pos+1]
			if next == '/' {
				// Close tag — only valid if it matches parentName.
				name, n, ok := scanCloseTag(p.src[p.pos:])
				if !ok {
					textBuf = append(textBuf, c)
					p.advance(1)
					continue
				}
				if parentName == "" {
					return nil, fmt.Errorf("unexpected close tag </%s> at line %d", name, p.line)
				}
				if name != parentName {
					return nil, fmt.Errorf("mismatched close tag </%s> at line %d (expected </%s>)", name, p.line, parentName)
				}
				flushText()
				p.advance(n)
				return children, nil
			}
			if isUpperAlpha(next) {
				flushText()
				elem, err := p.parseElement()
				if err != nil {
					return nil, err
				}
				children = append(children, elem)
				continue
			}
		}
		if c == '{' {
			flushText()
			loc := p.loc()
			p.advance(1)
			raw, err := p.readBraceExpr()
			if err != nil {
				return nil, err
			}
			children = append(children, ast.JSXExpression{Raw: string(raw), Location: loc})
			continue
		}
		textBuf = append(textBuf, c)
		p.advance(1)
	}

	if parentName != "" {
		return nil, fmt.Errorf("unexpected EOF: missing </%s>", parentName)
	}
	flushText()
	return children, nil
}

// parseElement parses a single `<Foo ...>...</Foo>` or `<Foo .../>`.
// p.pos must be at the opening '<'.
func (p *parser) parseElement() (ast.JSXElement, error) {
	loc := p.loc()
	if p.pos >= len(p.src) || p.src[p.pos] != '<' {
		return ast.JSXElement{}, fmt.Errorf("internal: parseElement called at non-'<' byte")
	}
	p.advance(1)

	name, ok := p.readJSXName(true)
	if !ok {
		return ast.JSXElement{}, fmt.Errorf("invalid JSX name at line %d", loc.Line)
	}

	props := map[string]ast.Node{}
	for {
		p.skipWS()
		if p.pos >= len(p.src) {
			return ast.JSXElement{}, fmt.Errorf("unexpected EOF in <%s> at line %d", name, loc.Line)
		}
		c := p.src[p.pos]
		if c == '/' || c == '>' {
			break
		}
		propName, value, err := p.parseAttribute()
		if err != nil {
			return ast.JSXElement{}, fmt.Errorf("in <%s> at line %d: %w", name, loc.Line, err)
		}
		props[propName] = value
	}

	if p.src[p.pos] == '/' {
		p.advance(1)
		if p.pos >= len(p.src) || p.src[p.pos] != '>' {
			return ast.JSXElement{}, fmt.Errorf("expected '>' after '/' in <%s/> at line %d", name, loc.Line)
		}
		p.advance(1)
		return ast.JSXElement{Name: name, Props: props, Location: loc}, nil
	}
	// p.src[p.pos] == '>'
	p.advance(1)

	children, err := p.parseChildren(name)
	if err != nil {
		return ast.JSXElement{}, err
	}
	return ast.JSXElement{Name: name, Props: props, Children: children, Location: loc}, nil
}

func (p *parser) parseAttribute() (string, ast.Node, error) {
	name, ok := p.readJSXName(false)
	if !ok {
		return "", nil, fmt.Errorf("invalid attribute name at line %d", p.line)
	}
	if p.pos >= len(p.src) || p.src[p.pos] != '=' {
		return "", nil, fmt.Errorf("attribute %q missing '=' at line %d", name, p.line)
	}
	p.advance(1)
	if p.pos >= len(p.src) {
		return "", nil, fmt.Errorf("attribute %q missing value at line %d", name, p.line)
	}
	switch p.src[p.pos] {
	case '"', '\'':
		quote := p.src[p.pos]
		p.advance(1)
		val, err := p.readQuoted(quote)
		if err != nil {
			return "", nil, err
		}
		return name, ast.String(val), nil
	case '{':
		loc := p.loc()
		p.advance(1)
		raw, err := p.readBraceExpr()
		if err != nil {
			return "", nil, err
		}
		return name, ast.JSXExpression{Raw: string(raw), Location: loc}, nil
	}
	return "", nil, fmt.Errorf("attribute %q has unsupported value at line %d", name, p.line)
}

func (p *parser) readJSXName(requireUpper bool) (string, bool) {
	if p.pos >= len(p.src) {
		return "", false
	}
	first := p.src[p.pos]
	if requireUpper {
		if !isUpperAlpha(first) {
			return "", false
		}
	} else if !isAlpha(first) {
		return "", false
	}
	start := p.pos
	for p.pos < len(p.src) {
		c := p.src[p.pos]
		if isAlpha(c) || isDigit(c) || c == '_' {
			p.advance(1)
			continue
		}
		break
	}
	return string(p.src[start:p.pos]), true
}

func (p *parser) readQuoted(term byte) (string, error) {
	var buf []byte
	for p.pos < len(p.src) {
		c := p.src[p.pos]
		if c == '\\' && p.pos+1 < len(p.src) {
			buf = append(buf, c, p.src[p.pos+1])
			p.advance(2)
			continue
		}
		if c == term {
			p.advance(1)
			return string(buf), nil
		}
		if c == '\n' || c == '\r' {
			return "", fmt.Errorf("unterminated string attribute at line %d", p.line)
		}
		buf = append(buf, c)
		p.advance(1)
	}
	return "", fmt.Errorf("unterminated string attribute at line %d", p.line)
}

// readBraceExpr reads bytes between { and the matching }, with brace
// depth tracking. The opening { has already been consumed. Double- and
// single-quoted string literals inside the expression are skipped over
// so braces in string contents don't close the expression.
func (p *parser) readBraceExpr() ([]byte, error) {
	var buf []byte
	depth := 1
	for p.pos < len(p.src) {
		c := p.src[p.pos]
		if c == '"' || c == '\'' {
			term := c
			buf = append(buf, c)
			p.advance(1)
			for p.pos < len(p.src) {
				d := p.src[p.pos]
				buf = append(buf, d)
				if d == '\\' && p.pos+1 < len(p.src) {
					buf = append(buf, p.src[p.pos+1])
					p.advance(2)
					continue
				}
				p.advance(1)
				if d == term {
					break
				}
			}
			continue
		}
		if c == '{' {
			depth++
			buf = append(buf, c)
			p.advance(1)
			continue
		}
		if c == '}' {
			depth--
			p.advance(1)
			if depth == 0 {
				return buf, nil
			}
			buf = append(buf, c)
			continue
		}
		buf = append(buf, c)
		p.advance(1)
	}
	return nil, fmt.Errorf("unterminated {expression}")
}

func (p *parser) skipWS() {
	for p.pos < len(p.src) {
		c := p.src[p.pos]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			p.advance(1)
			continue
		}
		break
	}
}

// scanCloseTag examines a buffer starting with `</` and reports the byte
// length of the close tag (through `>`) and the tag name. Returns ok=false
// if the buffer doesn't contain a valid close tag.
func scanCloseTag(line []byte) (name string, length int, ok bool) {
	if len(line) < 4 || line[0] != '<' || line[1] != '/' {
		return "", 0, false
	}
	if !isUpperAlpha(line[2]) {
		return "", 0, false
	}
	i := 2
	for i < len(line) && (isAlpha(line[i]) || isDigit(line[i]) || line[i] == '_') {
		i++
	}
	if i >= len(line) || line[i] != '>' {
		return "", 0, false
	}
	return string(line[2:i]), i + 1, true
}

func isUpperAlpha(c byte) bool { return c >= 'A' && c <= 'Z' }
func isLowerAlpha(c byte) bool { return c >= 'a' && c <= 'z' }
func isAlpha(c byte) bool      { return isLowerAlpha(c) || isUpperAlpha(c) }
func isDigit(c byte) bool      { return c >= '0' && c <= '9' }
