package marklit

import (
	"bytes"
	"strings"

	"github.com/vito/booklit/ast"
	gast "github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
)

type converter struct {
	source []byte
}

// convertChildren collects the Booklit AST for all children of a goldmark
// node. Block-level children become paragraphs in a sequence; inline children
// become a flat sequence.
//
// For Document nodes, headings are used to structure sections: the first
// heading sets the section title, and subsequent headings at the same or
// deeper level create sub-sections with their content grouped until the
// next heading of the same or shallower level.
func (c *converter) convertChildren(n gast.Node) ast.Node {
	if n.Kind() == gast.KindDocument {
		return c.convertDocument(n)
	}
	return c.convertChildrenFlat(n)
}

// convertChildrenFlat collects children without section structuring.
func (c *converter) convertChildrenFlat(n gast.Node) ast.Node {
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

// convertDocument handles section structuring based on headings.
//
// Headings create a hierarchy: `# Title` sets the top-level section title,
// `## Sub` creates a sub-section, `### SubSub` creates a sub-sub-section,
// etc. Content between headings becomes the body of the most recent section.
// The `{#tag}` attribute on a heading becomes the tag for that section.
func (c *converter) convertDocument(n gast.Node) ast.Node {
	// Collect all children into a flat list
	var children []sectionChild
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		var h *gast.Heading
		if child.Kind() == gast.KindHeading {
			h = child.(*gast.Heading)
		}
		children = append(children, sectionChild{gmNode: child, heading: h})
	}

	// If there are no headings, just convert everything flat
	hasHeading := false
	for _, ch := range children {
		if ch.heading != nil {
			hasHeading = true
			break
		}
	}
	if !hasHeading {
		return c.convertChildrenFlat(n)
	}

	// Find the first heading — it determines the "title level"
	titleLevel := 0
	for _, ch := range children {
		if ch.heading != nil {
			titleLevel = ch.heading.Level
			break
		}
	}

	// Build the result: content before the first heading is top-level body,
	// the first heading at titleLevel is the \title call, deeper headings
	// create \section blocks.
	return c.structureSections(children, titleLevel)
}

// sectionChild holds a goldmark child node and its heading (if any).
type sectionChild struct {
	gmNode  gast.Node
	heading *gast.Heading
}

// structureSections takes a flat list of goldmark children and a "title level"
// and produces a Booklit AST with proper \title and \section nesting.
func (c *converter) structureSections(children []sectionChild, titleLevel int) ast.Node {
	var result []ast.Node

	i := 0

	// Content before the first heading at titleLevel is top-level body
	for i < len(children) {
		ch := children[i]
		if ch.heading != nil && ch.heading.Level == titleLevel {
			break
		}
		converted := c.convertNode(ch.gmNode)
		if converted != nil {
			result = append(result, converted)
		}
		i++
	}

	// First heading at titleLevel → \title invocation
	if i < len(children) && children[i].heading != nil && children[i].heading.Level == titleLevel {
		result = append(result, c.headingToTitle(children[i].heading))
		i++
	}

	// Remaining content: nodes at body level, or sub-headings creating sections
	for i < len(children) {
		ch := children[i]
		if ch.heading != nil && ch.heading.Level > titleLevel {
			// Sub-heading: collect everything until the next heading at the
			// same or shallower level
			sectionLevel := ch.heading.Level
			sectionStart := i
			i++
			for i < len(children) {
				if children[i].heading != nil && children[i].heading.Level <= sectionLevel {
					break
				}
				i++
			}
			sectionChildren := children[sectionStart:i]
			sectionNode := c.buildSection(sectionChildren, sectionLevel)
			result = append(result, sectionNode)
		} else if ch.heading != nil && ch.heading.Level == titleLevel {
			// Another heading at the same level — this would be unusual
			// in a single document but handle it as another \title
			result = append(result, c.headingToTitle(ch.heading))
			i++
		} else {
			// Regular content
			converted := c.convertNode(ch.gmNode)
			if converted != nil {
				result = append(result, converted)
			}
			i++
		}
	}

	switch len(result) {
	case 0:
		return ast.Sequence{}
	case 1:
		return result[0]
	default:
		return ast.Sequence(result)
	}
}

// buildSection creates a \section{...} invoke from a group of children
// starting with a heading. The heading becomes the \title inside the section;
// deeper headings create nested sub-sections.
func (c *converter) buildSection(children []sectionChild, headingLevel int) ast.Node {
	// The section body is built by recursively structuring the children
	body := c.structureSections(children, headingLevel)

	return ast.Paragraph{ast.Sequence{ast.Invoke{
		Function:  "section",
		Arguments: []ast.Node{body},
	}}}
}

// headingToTitle converts a goldmark Heading into a \title invocation.
// If the heading has an {#id} attribute, it becomes the tag argument.
func (c *converter) headingToTitle(h *gast.Heading) ast.Node {
	titleContent := c.collectInlines(h)
	args := []ast.Node{ast.Sequence(titleContent)}

	// Check for {#id} attribute → becomes the tag
	if id, ok := h.AttributeString("id"); ok {
		if idStr, ok := id.([]byte); ok {
			args = append(args, ast.String(idStr))
		}
	}

	return ast.Paragraph{
		ast.Sequence{ast.Invoke{
			Function:  "title",
			Arguments: args,
		}},
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

	case KindJSXElement:
		j := n.(*JSXElementNode)
		return c.buildJSXElement(j.Name, j.Props, j.Children, j.Line, j.Col, j.MultiLine)

	case KindJSXBlockElement:
		j := n.(*JSXBlockElementNode)
		return c.buildJSXElement(j.Name, j.Props, j.Children, j.Line, j.Col, j.MultiLine)

	default:
		if n.Kind() == east.KindTable {
			return c.convertTable(n.(*east.Table))
		}
		// Unknown node type — try converting children
		return c.convertChildren(n)
	}
}

func (c *converter) convertParagraph(n gast.Node) ast.Node {
	// A paragraph containing only block-shaped JSX elements (with whitespace
	// between them) is split back into separate paragraphs. Goldmark glues
	// consecutive single-element lines into one paragraph, but
	// `<TableOfContents/>\n<Section>…</Section>` is two top-level blocks,
	// not one paragraph with a soft-break between them.
	if split := c.splitElementOnlyParagraph(n); split != nil {
		return split
	}

	inlines := c.collectInlines(n)
	if len(inlines) == 0 {
		return nil
	}
	return ast.Paragraph{ast.Sequence(inlines)}
}

// splitElementOnlyParagraph splits a paragraph that contains only JSX
// elements (and whitespace between them) back into one paragraph per
// element. Returns nil if the paragraph has any non-element content.
func (c *converter) splitElementOnlyParagraph(n gast.Node) ast.Node {
	elemCount := 0
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		switch child.Kind() {
		case KindJSXElement:
			elemCount++
		case gast.KindText:
			t := child.(*gast.Text)
			if strings.TrimSpace(string(t.Value(c.source))) != "" {
				return nil
			}
		default:
			return nil
		}
	}
	if elemCount < 2 {
		return nil
	}

	var nodes []ast.Node
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		if child.Kind() != KindJSXElement {
			continue
		}
		converted := c.convertNode(child)
		if converted != nil {
			nodes = append(nodes, ast.Paragraph{ast.Sequence{converted}})
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

	// Strip Markdown backslash escapes. Goldmark's text segments contain
	// raw escape sequences (e.g. \* \[ \\) — its own renderer strips them,
	// but since we produce Booklit AST we need to strip them here.
	s := ast.String(stripBackslashEscapes(value))
	if t.SoftLineBreak() {
		// Soft line breaks become spaces in flow content
		return ast.Sequence{s, ast.String(" ")}
	}
	if t.HardLineBreak() {
		return ast.Sequence{s, ast.String("\n")}
	}
	return s
}

// stripBackslashEscapes removes Markdown backslash escapes from text.
// In CommonMark, \ followed by an ASCII punctuation character produces the
// literal punctuation character. The escape backslash is stripped.
func stripBackslashEscapes(b []byte) []byte {
	if !bytes.ContainsRune(b, '\\') {
		return b
	}
	var out []byte
	for i := 0; i < len(b); i++ {
		if b[i] == '\\' && i+1 < len(b) && isASCIIPunct(b[i+1]) {
			// Skip the escape backslash; emit the escaped char directly
			i++
			out = append(out, b[i])
			continue
		}
		out = append(out, b[i])
	}
	return out
}

func isASCIIPunct(c byte) bool {
	return (c >= '!' && c <= '/') || (c >= ':' && c <= '@') ||
		(c >= '[' && c <= '`') || (c >= '{' && c <= '~')
}

func (c *converter) convertHeading(h *gast.Heading) ast.Node {
	// When a heading is encountered outside of document-level section
	// structuring (e.g. inside a blockquote or a JSX child), it falls
	// back to a simple \title invocation.
	return c.headingToTitle(h)
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
	var text strings.Builder
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		if t, ok := child.(*gast.Text); ok {
			text.WriteString(string(t.Value(c.source)))
		}
	}
	return ast.Invoke{
		Function:  "code",
		Arguments: []ast.Node{ast.String(text.String())},
	}
}

func (c *converter) convertLink(l *gast.Link) ast.Node {
	dest := string(l.Destination)

	// [text](#tag) → \reference{tag}{text}
	if len(dest) > 1 && dest[0] == '#' {
		tag := dest[1:]
		inner := c.collectInlines(l)
		return ast.Invoke{
			Function: "reference",
			Arguments: []ast.Node{
				ast.String(tag),
				ast.Sequence(inner),
			},
		}
	}

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
		line := seg.Value(c.source)
		// Strip trailing newline — Preformatted renders its own line separators
		line = bytes.TrimRight(line, "\n")
		lines = append(lines, ast.Sequence{ast.String(line)})
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

func (c *converter) convertTable(t *east.Table) ast.Node {
	var rows []ast.Node
	for child := t.FirstChild(); child != nil; child = child.NextSibling() {
		// Each child is a TableHeader or TableRow; both contain TableCells
		row := c.convertTableRow(child)
		rows = append(rows, row)
	}
	return ast.Invoke{
		Function:  "table",
		Arguments: rows,
	}
}

func (c *converter) convertTableRow(row gast.Node) ast.Node {
	var cells []ast.Node
	for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
		inlines := c.collectInlines(cell)
		if len(inlines) == 0 {
			cells = append(cells, ast.String(""))
		} else {
			cells = append(cells, ast.Sequence(inlines))
		}
	}
	return ast.Invoke{
		Function:  "table-row",
		Arguments: cells,
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
	return ast.Invoke{
		Function:  "raw-html-block",
		Arguments: []ast.Node{ast.String(text)},
	}
}

func (c *converter) convertRawHTML(n *gast.RawHTML) ast.Node {
	var text []byte
	for i := 0; i < n.Segments.Len(); i++ {
		seg := n.Segments.At(i)
		text = append(text, seg.Value(c.source)...)
	}
	return ast.Invoke{
		Function:  "raw-html",
		Arguments: []ast.Node{ast.String(text)},
	}
}

func (c *converter) convertInvoke(n *InvokeNode) ast.Node {
	invoke := ast.Invoke{
		Function: n.Function,
		Location: ast.Location{
			Line: n.Line,
			Col:  n.Col,
		},
	}

	for _, raw := range n.RawArgs {
		invoke.Arguments = append(invoke.Arguments, ParseInlineArg(raw))
	}

	return invoke
}

// buildJSXElement turns parsed JSX node data into ast.JSXElement. Shared
// between inline (JSXElementNode) and block (JSXBlockElementNode) goldmark
// nodes — both carry the same fields. The block flag controls how text
// chunks between children are parsed: block context uses the full block
// parser (so blank lines yield paragraphs), inline context uses the inline
// parser (newlines collapse to spaces).
//
// String-attribute values become ast.String (no markdown parsing — attributes
// are data, not content). Expression-attribute values become ast.JSXExpression.
// Nested elements inherit the parent's block context.
func (c *converter) buildJSXElement(name string, props []JSXProp, children []JSXChild, line, col int, block bool) ast.Node {
	elem := ast.JSXElement{
		Name:  name,
		Props: make(map[string]ast.Node, len(props)),
		Location: ast.Location{
			Line: line,
			Col:  col,
		},
	}

	for _, p := range props {
		switch p.Kind {
		case JSXPropExpression:
			elem.Props[p.Name] = ast.JSXExpression{Raw: string(p.Value)}
		default:
			elem.Props[p.Name] = ast.String(p.Value)
		}
	}

	for _, child := range children {
		switch child.Kind {
		case JSXChildElement:
			j := child.Elem
			// Nested elements use their own multi-line status, not the
			// parent's: <Section> can hold an inline <Title> whose
			// children should still be inline-parsed.
			elem.Children = append(elem.Children, c.buildJSXElement(j.Name, j.Props, j.Children, j.Line, j.Col, j.MultiLine))
		case JSXChildExpression:
			elem.Children = append(elem.Children, ast.JSXExpression{Raw: string(child.Text)})
		default:
			var node ast.Node
			if block {
				node = ParseArg(child.Text)
			} else {
				node = ParseInlineArg(child.Text)
			}
			// Skip empty results (whitespace-only chunks between block
			// children would otherwise add stray empty sequences).
			if s, ok := node.(ast.String); ok && len(s) == 0 {
				continue
			}
			if seq, ok := node.(ast.Sequence); ok {
				if len(seq) == 0 {
					continue
				}
				elem.Children = append(elem.Children, seq...)
			} else {
				elem.Children = append(elem.Children, node)
			}
		}
	}

	return elem
}
