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
	// the first heading at titleLevel is the <Title> call, deeper headings
	// create <Section> blocks.
	return c.structureSections(children, titleLevel)
}

// sectionChild holds a goldmark child node and its heading (if any).
type sectionChild struct {
	gmNode  gast.Node
	heading *gast.Heading
}

// structureSections takes a flat list of goldmark children and a "title level"
// and produces a Booklit AST with proper <Title> and <Section> nesting.
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

	// First heading at titleLevel → <Title> element
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
			// in a single document but handle it as another <Title>
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

// buildSection creates a <Section>...</Section> JSX element from a group of
// children starting with a heading. The heading becomes the <Title> inside
// the section; deeper headings create nested sub-sections.
func (c *converter) buildSection(children []sectionChild, headingLevel int) ast.Node {
	body := c.structureSections(children, headingLevel)

	return ast.Paragraph{ast.Sequence{ast.JSXElement{
		Name:      "Section",
		Children:  []ast.Node{body},
		MultiLine: true,
	}}}
}

// headingToTitle converts a goldmark Heading into a <Title>...</Title> JSX
// element. If the heading has an {#id} attribute, it becomes the tag prop.
func (c *converter) headingToTitle(h *gast.Heading) ast.Node {
	titleContent := c.collectInlines(h)

	elem := ast.JSXElement{
		Name:     "Title",
		Children: titleContent,
	}

	if id, ok := h.AttributeString("id"); ok {
		if idStr, ok := id.([]byte); ok {
			elem.Props = map[string]ast.Node{"tag": ast.String(idStr)}
		}
	}

	return ast.Paragraph{ast.Sequence{elem}}
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
		return ast.JSXElement{Name: "hr", MultiLine: true}

	case gast.KindHTMLBlock:
		return c.convertHTMLBlock(n.(*gast.HTMLBlock))

	case gast.KindRawHTML:
		return c.convertRawHTML(n.(*gast.RawHTML))

	case KindJSXExprInline:
		e := n.(*jsxExprInlineNode)
		return ast.JSXExpression{Raw: e.Raw}

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
	// back to a simple <Title> element.
	return c.headingToTitle(h)
}

func (c *converter) convertEmphasis(e *gast.Emphasis) ast.Node {
	name := "em"
	if e.Level >= 2 {
		name = "strong"
	}
	return ast.JSXElement{
		Name:     name,
		Children: c.collectInlines(e),
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
	return ast.JSXElement{
		Name:     "code",
		Children: []ast.Node{ast.String(text.String())},
	}
}

func (c *converter) convertLink(l *gast.Link) ast.Node {
	dest := string(l.Destination)

	// [text](#tag) → <Reference tag="tag">text</Reference>
	if len(dest) > 1 && dest[0] == '#' {
		tag := dest[1:]
		return ast.JSXElement{
			Name:     "Reference",
			Props:    map[string]ast.Node{"tag": ast.String(tag)},
			Children: c.collectInlines(l),
		}
	}

	return ast.JSXElement{
		Name:     "a",
		Props:    map[string]ast.Node{"href": ast.String(dest)},
		Children: c.collectInlines(l),
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

	return ast.JSXElement{
		Name: "img",
		Props: map[string]ast.Node{
			"src": ast.String(img.Destination),
			"alt": ast.String(alt),
		},
	}
}

func (c *converter) convertAutoLink(al *gast.AutoLink) ast.Node {
	url := al.URL(c.source)
	label := al.Label(c.source)
	return ast.JSXElement{
		Name:     "a",
		Props:    map[string]ast.Node{"href": ast.String(url)},
		Children: []ast.Node{ast.String(label)},
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

	// For fenced code blocks with a language, route through <CodeBlock> so
	// the syntax-highlighting builtin runs over the body.
	if fcb, ok := n.(*gast.FencedCodeBlock); ok {
		lang := fcb.Language(c.source)
		if len(lang) > 0 {
			return ast.JSXElement{
				Name:      "CodeBlock",
				Props:     map[string]ast.Node{"language": ast.String(lang)},
				Children:  []ast.Node{pre},
				MultiLine: true,
			}
		}
	}

	return ast.JSXElement{
		Name:      "pre",
		Children:  []ast.Node{pre},
		MultiLine: true,
	}
}

func (c *converter) convertBlockquote(n gast.Node) ast.Node {
	// Routed through <Inset> rather than a plain <blockquote> so the
	// rendered output keeps the `class="inset"` wrapper that the rest of
	// the docs (and existing CSS) expect.
	return ast.JSXElement{
		Name:      "Inset",
		Children:  []ast.Node{c.convertChildren(n)},
		MultiLine: true,
	}
}

func (c *converter) convertList(l *gast.List) ast.Node {
	name := "ul"
	if l.IsOrdered() {
		name = "ol"
	}

	var items []ast.Node
	for child := l.FirstChild(); child != nil; child = child.NextSibling() {
		items = append(items, ast.JSXElement{
			Name:      "li",
			Children:  []ast.Node{c.convertChildren(child)},
			MultiLine: true,
		})
	}

	return ast.JSXElement{
		Name:      name,
		Children:  items,
		MultiLine: true,
	}
}

func (c *converter) convertTable(t *east.Table) ast.Node {
	var rows []ast.Node
	for child := t.FirstChild(); child != nil; child = child.NextSibling() {
		rows = append(rows, c.convertTableRow(child))
	}
	return ast.JSXElement{
		Name:      "table",
		Children:  rows,
		MultiLine: true,
	}
}

func (c *converter) convertTableRow(row gast.Node) ast.Node {
	var cells []ast.Node
	for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
		inlines := c.collectInlines(cell)
		var content ast.Node
		switch len(inlines) {
		case 0:
			content = ast.String("")
		case 1:
			content = inlines[0]
		default:
			content = ast.Sequence(inlines)
		}
		cells = append(cells, ast.JSXElement{
			Name:      "td",
			Children:  []ast.Node{content},
			MultiLine: true,
		})
	}
	return ast.JSXElement{
		Name:      "tr",
		Children:  cells,
		MultiLine: true,
	}
}

// convertHTMLBlock is a fallback for cases where goldmark still produces
// an HTMLBlock node — primarily HTML comments (`<!-- -->`) and similar
// edge cases that our JSX block parser doesn't claim. The body is wrapped
// in a `<RawHTML block="true">…</RawHTML>` element so the bytes survive
// untouched and the surrounding paragraph layout treats it as block
// content.
func (c *converter) convertHTMLBlock(n *gast.HTMLBlock) ast.Node {
	var text []byte
	for i := 0; i < n.Lines().Len(); i++ {
		seg := n.Lines().At(i)
		text = append(text, seg.Value(c.source)...)
	}
	if n.HasClosure() {
		text = append(text, n.ClosureLine.Value(c.source)...)
	}
	return ast.JSXElement{
		Name:      "RawHTML",
		Props:     map[string]ast.Node{"block": ast.String("true")},
		Children:  []ast.Node{ast.String(text)},
		MultiLine: true,
	}
}

// convertRawHTML is the fallback for inline raw HTML that goldmark's
// default parser claims when our JSX inline parser doesn't (e.g. `<br>`
// without an explicit `<br/>`, comments). Bytes pass through verbatim.
func (c *converter) convertRawHTML(n *gast.RawHTML) ast.Node {
	var text []byte
	for i := 0; i < n.Segments.Len(); i++ {
		seg := n.Segments.At(i)
		text = append(text, seg.Value(c.source)...)
	}
	return ast.JSXElement{
		Name:     "RawHTML",
		Children: []ast.Node{ast.String(text)},
	}
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
		Name:      name,
		Props:     make(map[string]ast.Node, len(props)),
		MultiLine: block,
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
