package marklit

import (
	"strings"

	"github.com/vito/booklit/ast"
	gast "github.com/yuin/goldmark/ast"
)

type converter struct {
	source      []byte
	extractions []extractedInvoke
}

// convertChildren collects the Booklit AST for all children of a goldmark
// node. Block-level children become paragraphs in a sequence; inline children
// become a flat sequence.
func (c *converter) convertChildren(n gast.Node) ast.Node {
	var paragraphs []ast.Node

	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		converted := c.convertNode(child)
		if converted != nil {
			paragraphs = append(paragraphs, converted)
		}
	}

	switch len(paragraphs) {
	case 0:
		return ast.Sequence{}
	case 1:
		return paragraphs[0]
	default:
		return ast.Sequence(paragraphs)
	}
}

// convertNode dispatches on goldmark node kind to produce a Booklit AST node.
func (c *converter) convertNode(n gast.Node) ast.Node {
	switch n.Kind() {
	case gast.KindDocument:
		return c.convertChildren(n)

	case gast.KindParagraph:
		return c.convertParagraph(n)

	case gast.KindHeading:
		return c.convertHeading(n.(*gast.Heading))

	case gast.KindText:
		return c.convertText(n.(*gast.Text))

	case gast.KindString:
		return ast.String(n.(*gast.String).Value)

	case gast.KindEmphasis:
		return c.convertEmphasis(n.(*gast.Emphasis))

	case gast.KindCodeSpan:
		return c.convertCodeSpan(n)

	case gast.KindLink:
		return c.convertLink(n.(*gast.Link))

	case gast.KindImage:
		return c.convertImage(n.(*gast.Image))

	case gast.KindAutoLink:
		return c.convertAutoLink(n.(*gast.AutoLink))

	case gast.KindCodeBlock, gast.KindFencedCodeBlock:
		return c.convertCodeBlock(n)

	case gast.KindBlockquote:
		return c.convertBlockquote(n)

	case gast.KindList:
		return c.convertList(n.(*gast.List))

	case gast.KindListItem:
		return c.convertChildren(n)

	case gast.KindThematicBreak:
		// Render as a horizontal rule styled element
		return ast.Invoke{
			Function: "thematic-break",
		}

	case gast.KindHTMLBlock:
		return c.convertHTMLBlock(n.(*gast.HTMLBlock))

	case gast.KindRawHTML:
		return c.convertRawHTML(n.(*gast.RawHTML))

	case KindInvoke:
		return c.convertInvoke(n.(*InvokeNode))

	case KindInvokeBlock:
		return c.convertInvokeBlock(n.(*InvokeBlockNode))

	default:
		// Unknown node type — try converting children
		return c.convertChildren(n)
	}
}

func (c *converter) convertParagraph(n gast.Node) ast.Node {
	// Check if this paragraph contains only invocations (and whitespace
	// between them). In the old parser, each function call on its own line
	// was a separate statement. Goldmark groups consecutive lines into one
	// paragraph, so we split them back into individual paragraphs. This
	// prevents spurious <p> </p> when side-effect-only calls like
	// @use-plugin return empty content (the soft-break space between them
	// would otherwise persist).
	if invokes := c.splitInvokeOnlyParagraph(n); invokes != nil {
		return invokes
	}

	inlines := c.collectInlines(n)
	if len(inlines) == 0 {
		return nil
	}
	return ast.Paragraph{ast.Sequence(inlines)}
}

// splitInvokeOnlyParagraph checks if a goldmark paragraph contains only
// invoke nodes (and whitespace/soft-break text between them). If so, it
// returns each invocation wrapped in its own ast.Paragraph (matching the
// old parser behavior). Returns nil if the paragraph has non-invoke content.
func (c *converter) splitInvokeOnlyParagraph(n gast.Node) ast.Node {
	// First pass: verify the paragraph is invoke-only
	invokeCount := 0
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		switch child.Kind() {
		case KindInvoke, KindInvokeBlock:
			invokeCount++
		case gast.KindText:
			t := child.(*gast.Text)
			value := t.Value(c.source)
			text := strings.TrimSpace(string(value))
			if text != "" && !strings.HasPrefix(text, invokePlaceholderPrefix) {
				return nil
			}
		default:
			return nil
		}
	}
	if invokeCount < 2 {
		return nil // single invoke: keep as normal paragraph
	}

	// Second pass: extract each invocation into its own paragraph
	var nodes []ast.Node
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		switch child.Kind() {
		case KindInvoke, KindInvokeBlock:
			converted := c.convertNode(child)
			if converted != nil {
				nodes = append(nodes, ast.Paragraph{ast.Sequence{converted}})
			}
		case gast.KindText:
			t := child.(*gast.Text)
			value := t.Value(c.source)
			if node := c.tryResolvePlaceholder(string(value)); node != nil {
				nodes = append(nodes, ast.Paragraph{ast.Sequence{node}})
			}
		}
	}
	if len(nodes) == 1 {
		return nodes[0]
	}
	return ast.Sequence(nodes)
}

// collectInlines gathers all inline children of a block node into a flat
// slice of Booklit AST nodes.
func (c *converter) collectInlines(n gast.Node) []ast.Node {
	var nodes []ast.Node
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		converted := c.convertNode(child)
		if converted != nil {
			nodes = append(nodes, converted)
		}
	}
	return nodes
}

func (c *converter) convertText(t *gast.Text) ast.Node {
	value := t.Value(c.source)
	text := string(value)

	// Check if the entire text is a placeholder
	if node := c.tryResolvePlaceholder(text); node != nil {
		return node
	}

	// Check if text contains embedded placeholders (block invocations that
	// appeared inline within a sentence, e.g. "easier to \x00BLI1\x00 later")
	if strings.Contains(text, invokePlaceholderPrefix) {
		return c.resolveEmbeddedPlaceholders(text)
	}

	s := ast.String(value)
	if t.SoftLineBreak() {
		// Soft line breaks become spaces in flow content
		return ast.Sequence{s, ast.String(" ")}
	}
	if t.HardLineBreak() {
		return ast.Sequence{s, ast.String("\n")}
	}
	return s
}

// resolveEmbeddedPlaceholders splits text at placeholder markers and returns
// a Sequence with the text segments interspersed with resolved invoke nodes.
func (c *converter) resolveEmbeddedPlaceholders(text string) ast.Node {
	var nodes []ast.Node
	for {
		idx := strings.Index(text, invokePlaceholderPrefix)
		if idx < 0 {
			break
		}
		// Add text before the placeholder
		if idx > 0 {
			nodes = append(nodes, ast.String(text[:idx]))
		}
		// Find the null terminator after the index number
		rest := text[idx:]
		nullIdx := strings.IndexByte(rest[len(invokePlaceholderPrefix):], 0)
		if nullIdx < 0 {
			// Malformed placeholder — emit as text
			nodes = append(nodes, ast.String(rest))
			text = ""
			break
		}
		placeholder := rest[:len(invokePlaceholderPrefix)+nullIdx+1]
		if node := c.tryResolvePlaceholder(placeholder); node != nil {
			nodes = append(nodes, node)
		}
		text = rest[len(placeholder):]
	}
	// Add remaining text
	if len(text) > 0 {
		nodes = append(nodes, ast.String(text))
	}
	if len(nodes) == 1 {
		return nodes[0]
	}
	return ast.Sequence(nodes)
}

// tryResolvePlaceholder checks if a text string is a placeholder marker and
// returns the corresponding extracted invoke node, or nil.
func (c *converter) tryResolvePlaceholder(text string) ast.Node {
	idx, ok := isPlaceholder(text)
	if !ok || idx >= len(c.extractions) {
		return nil
	}

	ext := c.extractions[idx]
	invoke := ast.Invoke{
		Function: ext.Function,
	}
	for i, raw := range ext.RawArgs {
		argType := ArgNormal
		if i < len(ext.ArgTypes) {
			argType = ext.ArgTypes[i]
		}

		switch argType {
		case ArgVerbatim:
			invoke.Arguments = append(invoke.Arguments, verbatimToNode(raw))
		case ArgPreformatted:
			invoke.Arguments = append(invoke.Arguments, ParsePreformattedArg(raw))
		default:
			// Use block parsing only for args that start with a newline
			// (i.e. the content begins on the line after the opening {).
			// Args with newlines in the middle are just soft-wrapped
			// inline text (e.g. @ref{tag}{display text\nwrapped}).
			if len(raw) > 0 && (raw[0] == '\n' || raw[0] == '\r') {
				invoke.Arguments = append(invoke.Arguments, ParseArg(raw))
			} else {
				invoke.Arguments = append(invoke.Arguments, ParseInlineArg(raw))
			}
		}
	}
	return invoke
}

func (c *converter) convertHeading(h *gast.Heading) ast.Node {
	titleContent := c.collectInlines(h)
	args := []ast.Node{ast.Sequence(titleContent)}
	return ast.Paragraph{
		ast.Sequence{ast.Invoke{
			Function:  "title",
			Arguments: args,
		}},
	}
}

func (c *converter) convertEmphasis(e *gast.Emphasis) ast.Node {
	inner := c.collectInlines(e)
	funcName := "italic"
	if e.Level >= 2 {
		funcName = "bold"
	}
	return ast.Invoke{
		Function:  funcName,
		Arguments: []ast.Node{ast.Sequence(inner)},
	}
}

func (c *converter) convertCodeSpan(n gast.Node) ast.Node {
	// Code spans contain raw text children
	var text string
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		if t, ok := child.(*gast.Text); ok {
			text += string(t.Value(c.source))
		}
	}
	return ast.Invoke{
		Function:  "code",
		Arguments: []ast.Node{ast.String(text)},
	}
}

func (c *converter) convertLink(l *gast.Link) ast.Node {
	inner := c.collectInlines(l)
	return ast.Invoke{
		Function: "link",
		Arguments: []ast.Node{
			ast.Sequence(inner),
			ast.String(l.Destination),
		},
	}
}

func (c *converter) convertImage(img *gast.Image) ast.Node {
	// alt text from children
	var altParts []string
	for child := img.FirstChild(); child != nil; child = child.NextSibling() {
		if t, ok := child.(*gast.Text); ok {
			altParts = append(altParts, string(t.Value(c.source)))
		}
	}
	alt := strings.Join(altParts, "")

	args := []ast.Node{ast.String(img.Destination)}
	if alt != "" {
		args = append(args, ast.String(alt))
	}
	return ast.Invoke{
		Function:  "image",
		Arguments: args,
	}
}

func (c *converter) convertAutoLink(al *gast.AutoLink) ast.Node {
	url := al.URL(c.source)
	label := al.Label(c.source)
	return ast.Invoke{
		Function: "link",
		Arguments: []ast.Node{
			ast.String(label),
			ast.String(url),
		},
	}
}

func (c *converter) convertCodeBlock(n gast.Node) ast.Node {
	var lines []ast.Sequence
	for i := 0; i < n.Lines().Len(); i++ {
		seg := n.Lines().At(i)
		lines = append(lines, ast.Sequence{ast.String(seg.Value(c.source))})
	}

	pre := ast.Preformatted(lines)

	// For fenced code blocks, wrap in a code-block invoke with language info
	if fcb, ok := n.(*gast.FencedCodeBlock); ok {
		lang := fcb.Language(c.source)
		if len(lang) > 0 {
			return ast.Invoke{
				Function:  "code-block",
				Arguments: []ast.Node{ast.String(lang), pre},
			}
		}
	}

	return ast.Invoke{
		Function:  "code",
		Arguments: []ast.Node{pre},
	}
}

func (c *converter) convertBlockquote(n gast.Node) ast.Node {
	inner := c.convertChildren(n)
	return ast.Invoke{
		Function:  "inset",
		Arguments: []ast.Node{inner},
	}
}

func (c *converter) convertList(l *gast.List) ast.Node {
	funcName := "list"
	if l.IsOrdered() {
		funcName = "ordered-list"
	}

	var items []ast.Node
	for child := l.FirstChild(); child != nil; child = child.NextSibling() {
		item := c.convertChildren(child)
		items = append(items, item)
	}

	return ast.Invoke{
		Function:  funcName,
		Arguments: items,
	}
}

func (c *converter) convertHTMLBlock(n *gast.HTMLBlock) ast.Node {
	var text []byte
	for i := 0; i < n.Lines().Len(); i++ {
		seg := n.Lines().At(i)
		text = append(text, seg.Value(c.source)...)
	}
	if n.HasClosure() {
		text = append(text, n.ClosureLine.Value(c.source)...)
	}
	return ast.String(text)
}

func (c *converter) convertRawHTML(n *gast.RawHTML) ast.Node {
	var text []byte
	for i := 0; i < n.Segments.Len(); i++ {
		seg := n.Segments.At(i)
		text = append(text, seg.Value(c.source)...)
	}
	return ast.String(text)
}

func (c *converter) convertInvoke(n *InvokeNode) ast.Node {
	invoke := ast.Invoke{
		Function: n.Function,
		Location: ast.Location{
			Line: n.Line,
			Col:  n.Col,
		},
	}

	for i, raw := range n.RawArgs {
		argType := ArgNormal
		if i < len(n.ArgTypes) {
			argType = n.ArgTypes[i]
		}

		var argNode ast.Node
		switch argType {
		case ArgVerbatim:
			argNode = verbatimToNode(raw)
		case ArgPreformatted:
			argNode = ParsePreformattedArg(raw)
		default:
			// Parse the argument content as Markdown to get nested Booklit
			// nodes. This allows things like @title{Hello *world*} to work.
			argNode = ParseInlineArg(raw)
		}

		invoke.Arguments = append(invoke.Arguments, argNode)
	}

	return invoke
}

func (c *converter) convertInvokeBlock(n *InvokeBlockNode) ast.Node {
	invoke := ast.Invoke{
		Function: n.Function,
		Location: ast.Location{
			Line: n.Line,
			Col:  n.Col,
		},
	}

	for _, raw := range n.RawArgs {
		argNode := ParseArg(raw)
		invoke.Arguments = append(invoke.Arguments, argNode)
	}

	return invoke
}
