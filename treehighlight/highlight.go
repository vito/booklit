//go:build cgo

// Package treehighlight renders code HTML, using tree-sitter when cgo is available.
package treehighlight

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"iter"
	"sort"
	"strings"
	"unsafe"

	tshighlight "go.gopad.dev/go-tree-sitter-highlight"

	"github.com/tree-sitter/go-tree-sitter"
	tsbash "github.com/tree-sitter/tree-sitter-bash/bindings/go"
	tsgo "github.com/tree-sitter/tree-sitter-go/bindings/go"
	tshtml "github.com/tree-sitter/tree-sitter-html/bindings/go"
	tsjavascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	booklitgrammar "github.com/vito/booklit/treehighlight/internal/tree_sitter_booklit"
)

// Chunk is a rendered fragment of highlighted source. HTML chunks should be
// emitted as raw HTML. Link chunks should be rendered as Booklit references;
// they intentionally contain only source text so the surrounding raw HTML spans
// can color the eventual <a> element.
type Chunk struct {
	HTML     string
	LinkTag  string
	LinkText string
}

// Options controls highlighting behavior.
type Options struct {
	// LinkReferences turns language-specific link captures into Link chunks.
	LinkReferences bool
}

var captureNames = []string{
	"comment",
	"keyword",
	"operator",
	"punctuation",
	"punctuation.bracket",
	"punctuation.delimiter",
	"punctuation.special",
	"string",
	"string.special",
	"number",
	"constant",
	"function",
	"tag",
	"type",
	"constructor",
	"attribute",
	"property",
	"variable",
	"markup.heading",
	"markup.raw",
}

var captureStyles = map[string]string{
	"comment":               "color:#c29d7c;font-style:italic",
	"keyword":               "color:#ed6c30;font-weight:bold",
	"operator":              "color:#fcc21b;font-weight:bold",
	"punctuation.bracket":   "color:#fcc21b",
	"punctuation.delimiter": "color:#fcc21b",
	"punctuation.special":   "color:#fcc21b;font-weight:bold",
	"string":                "color:#fcc21b",
	"string.special":        "color:#fcc21b;font-weight:bold",
	"number":                "color:#fcc21b",
	"constant":              "color:#fcc21b",
	"function":              "color:#ed6c30;font-weight:bold",
	"tag":                   "color:#fcc21b;font-weight:bold",
	"type":                  "color:#fcc21b;font-weight:bold",
	"constructor":           "color:#fcc21b;font-weight:bold",
	"attribute":             "color:#fcc21b",
	"property":              "color:#fcc21b",
	"variable":              "color:#f0f0f0",
	"markup.heading":        "font-weight:bold",
	"markup.raw":            "color:#fcc21b",
}

type languageSpec struct {
	name       string
	language   func() unsafe.Pointer
	highlights string
	links      string
	linkTag    func(string) string
}

var languages = map[string]*languageSpec{
	"booklit": {
		name:       "booklit",
		language:   booklitgrammar.Language,
		highlights: booklitHighlightsQuery,
		links:      booklitLinksQuery,
		linkTag:    func(s string) string { return s },
	},
	"go": {
		name:       "go",
		language:   tsgo.Language,
		highlights: goHighlightsQuery,
	},
	"javascript": {
		name:       "javascript",
		language:   tsjavascript.Language,
		highlights: javascriptHighlightsQuery,
	},
	"html": {
		name:       "html",
		language:   tshtml.Language,
		highlights: htmlHighlightsQuery,
	},
	"bash": {
		name:       "bash",
		language:   tsbash.Language,
		highlights: bashHighlightsQuery,
	},
}

var languageAliases = map[string]string{
	"":                 "",
	"text":             "",
	"txt":              "",
	"plain":            "",
	"lit":              "booklit",
	"booklit":          "booklit",
	"markdown":         "booklit",
	"md":               "booklit",
	"go":               "go",
	"golang":           "go",
	"js":               "javascript",
	"javascript":       "javascript",
	"jsx":              "javascript",
	"html":             "html",
	"go-html-template": "html",
	"gotemplate":       "html",
	"tmpl":             "html",
	"template":         "html",
	"bash":             "bash",
	"sh":               "bash",
	"shell":            "bash",
}

// HTML returns a highlighted HTML fragment wrapped in the same <pre><code> or
// <code> shape that Chroma used to produce for Booklit.
func HTML(language, source string, inline bool) (string, error) {
	chunks, err := Chunks(language, source, Options{})
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	if inline {
		buf.WriteString(`<code style=";-webkit-text-size-adjust:none;">`)
	} else {
		buf.WriteString(`<pre style=";-webkit-text-size-adjust:none;"><code>`)
	}
	for _, chunk := range chunks {
		// HTML() is used where references are not expected; if a caller ever
		// passes LinkReferences through here, render the link text plainly.
		if chunk.HTML != "" {
			buf.WriteString(chunk.HTML)
		} else {
			buf.WriteString(html.EscapeString(chunk.LinkText))
		}
	}
	if inline {
		buf.WriteString(`</code>`)
	} else {
		buf.WriteString(`</code></pre>`)
	}
	return buf.String(), nil
}

// Chunks returns highlighted source split into raw HTML and optional reference
// chunks. Unknown languages fall back to escaped, unhighlighted source.
func Chunks(language, source string, opts Options) ([]Chunk, error) {
	spec := lookupLanguage(language)
	if spec == nil {
		return []Chunk{{HTML: html.EscapeString(source)}}, nil
	}

	lang := tree_sitter.NewLanguage(spec.language())
	cfg, err := tshighlight.NewConfiguration(lang, spec.name, []byte(spec.highlights), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("configure %s highlighter: %w", spec.name, err)
	}
	cfg.Configure(captureNames)
	defer cfg.Query.Close()
	if cfg.CombinedInjectionsQuery != nil {
		defer cfg.CombinedInjectionsQuery.Close()
	}

	var links []linkRange
	if opts.LinkReferences && spec.links != "" {
		links, err = collectLinks(spec, lang, []byte(source))
		if err != nil {
			return nil, err
		}
	}

	highlighter := tshighlight.New()
	defer highlighter.Parser.Close()
	events := highlighter.Highlight(context.Background(), *cfg, []byte(source), func(string) *tshighlight.Configuration { return nil })
	return renderChunks(events, []byte(source), links)
}

func lookupLanguage(language string) *languageSpec {
	name := strings.ToLower(strings.TrimSpace(language))
	if i := strings.IndexAny(name, " \t"); i >= 0 {
		name = name[:i]
	}
	name = languageAliases[name]
	if name == "" {
		return nil
	}
	return languages[name]
}

type linkRange struct {
	Start uint
	End   uint
	Tag   string
}

func collectLinks(spec *languageSpec, lang *tree_sitter.Language, source []byte) ([]linkRange, error) {
	parser := tree_sitter.NewParser()
	defer parser.Close()
	if err := parser.SetLanguage(lang); err != nil {
		return nil, fmt.Errorf("set %s parser language: %w", spec.name, err)
	}
	tree := parser.Parse(source, nil)
	if tree == nil {
		return nil, nil
	}
	defer tree.Close()

	query, err := tree_sitter.NewQuery(lang, spec.links)
	if err != nil {
		return nil, fmt.Errorf("compile %s link query: %w", spec.name, err)
	}
	defer query.Close()

	captureNames := query.CaptureNames()
	cursor := tree_sitter.NewQueryCursor()
	defer cursor.Close()
	captures := cursor.Captures(query, tree.RootNode(), source)

	var links []linkRange
	for {
		match, captureIndex := captures.Next()
		if match == nil {
			break
		}
		capture := match.Captures[captureIndex]
		if captureNames[capture.Index] != "reference" {
			continue
		}
		start, end := capture.Node.ByteRange()
		text := capture.Node.Utf8Text(source)
		tag := text
		if spec.linkTag != nil {
			tag = spec.linkTag(text)
		}
		if tag == "" || start >= end {
			continue
		}
		links = append(links, linkRange{Start: start, End: end, Tag: tag})
	}

	sort.Slice(links, func(i, j int) bool {
		if links[i].Start == links[j].Start {
			return links[i].End < links[j].End
		}
		return links[i].Start < links[j].Start
	})
	return compactLinks(links), nil
}

func compactLinks(links []linkRange) []linkRange {
	if len(links) < 2 {
		return links
	}
	out := links[:0]
	for _, link := range links {
		if len(out) == 0 {
			out = append(out, link)
			continue
		}
		prev := out[len(out)-1]
		if link.Start == prev.Start && link.End == prev.End && link.Tag == prev.Tag {
			continue
		}
		if link.Start < prev.End {
			// Prefer the first, outermost capture rather than producing nested
			// references, which Booklit cannot represent in raw HTML.
			continue
		}
		out = append(out, link)
	}
	return out
}

func renderChunks(events iter.Seq2[tshighlight.Event, error], source []byte, links []linkRange) ([]Chunk, error) {
	var (
		chunks    []Chunk
		raw       strings.Builder
		spanStack []bool
		linkIndex int
	)

	flushRaw := func() {
		if raw.Len() == 0 {
			return
		}
		chunks = append(chunks, Chunk{HTML: raw.String()})
		raw.Reset()
	}

	writeEscaped := func(start, end uint) {
		if start >= end {
			return
		}
		raw.WriteString(html.EscapeString(string(source[start:end])))
	}

	writeSource := func(start, end uint) {
		pos := start
		for pos < end {
			for linkIndex < len(links) && links[linkIndex].End <= pos {
				linkIndex++
			}
			if linkIndex >= len(links) || links[linkIndex].Start >= end {
				writeEscaped(pos, end)
				return
			}

			link := links[linkIndex]
			if link.Start > pos {
				until := minUint(end, link.Start)
				writeEscaped(pos, until)
				pos = until
				continue
			}

			if link.Start < pos {
				// Overlap from a previous capture; render this piece as plain
				// highlighted source and advance past it.
				until := minUint(end, link.End)
				writeEscaped(pos, until)
				pos = until
				continue
			}

			// The link starts at the current source position. It is expected to be
			// contained by this EventSource range because it comes from a single
			// tree-sitter capture. If not, split conservatively rather than
			// disturbing the event stream.
			until := minUint(end, link.End)
			flushRaw()
			chunks = append(chunks, Chunk{
				LinkTag:  link.Tag,
				LinkText: string(source[pos:until]),
			})
			pos = until
			if pos >= link.End {
				linkIndex++
			}
		}
	}

	for event, err := range events {
		if err != nil {
			return nil, fmt.Errorf("render highlighted source: %w", err)
		}

		switch e := event.(type) {
		case tshighlight.EventLayerStart, tshighlight.EventLayerEnd:
			// Language layer events matter for injected languages. Booklit's
			// current queries do not use injections, so there is no raw HTML to
			// emit here.
		case tshighlight.EventCaptureStart:
			opened := false
			if int(e.Highlight) < len(captureNames) {
				if style := captureStyles[captureNames[e.Highlight]]; style != "" {
					raw.WriteString(`<span style="`)
					raw.WriteString(style)
					raw.WriteString(`">`)
					opened = true
				}
			}
			spanStack = append(spanStack, opened)
		case tshighlight.EventCaptureEnd:
			if len(spanStack) == 0 {
				continue
			}
			opened := spanStack[len(spanStack)-1]
			spanStack = spanStack[:len(spanStack)-1]
			if opened {
				raw.WriteString(`</span>`)
			}
		case tshighlight.EventSource:
			writeSource(e.StartByte, e.EndByte)
		}
	}
	flushRaw()

	return chunks, nil
}

func minUint(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}

// PlainHTML returns source escaped into the renderer's code wrapper without
// syntax spans. It is primarily useful for callers that need an explicit
// fallback path.
func PlainHTML(source string, inline bool) string {
	var buf bytes.Buffer
	if inline {
		buf.WriteString(`<code style=";-webkit-text-size-adjust:none;">`)
	} else {
		buf.WriteString(`<pre style=";-webkit-text-size-adjust:none;"><code>`)
	}
	buf.WriteString(html.EscapeString(source))
	if inline {
		buf.WriteString(`</code>`)
	} else {
		buf.WriteString(`</code></pre>`)
	}
	return buf.String()
}

const booklitHighlightsQuery = `
(comment) @comment
(command marker: (backslash) @function name: (identifier) @function)
(tag name: (identifier) @tag)
(heading) @markup.heading
(code_span) @markup.raw
(tag_string) @string
(delimiter) @punctuation.bracket
`

const booklitLinksQuery = `
(command name: (identifier) @reference)
`

const goHighlightsQuery = `
[
  "break"
  "case"
  "chan"
  "const"
  "continue"
  "default"
  "defer"
  "else"
  "fallthrough"
  "for"
  "func"
  "go"
  "goto"
  "if"
  "import"
  "interface"
  "map"
  "package"
  "range"
  "return"
  "select"
  "struct"
  "switch"
  "type"
  "var"
] @keyword

(comment) @comment

[
  (interpreted_string_literal)
  (raw_string_literal)
  (rune_literal)
] @string

[
  (int_literal)
  (float_literal)
  (imaginary_literal)
] @number

(function_declaration name: (identifier) @function)
(method_declaration name: (field_identifier) @function)
(call_expression function: (identifier) @function)
(call_expression function: (selector_expression field: (field_identifier) @function))
(type_identifier) @type
(package_identifier) @module
(field_identifier) @property
`

const javascriptHighlightsQuery = `
[
  "as"
  "async"
  "await"
  "break"
  "case"
  "catch"
  "class"
  "const"
  "continue"
  "debugger"
  "default"
  "delete"
  "do"
  "else"
  "export"
  "extends"
  "finally"
  "for"
  "from"
  "function"
  "get"
  "if"
  "import"
  "in"
  "let"
  "new"
  "of"
  "return"
  "set"
  "static"
  "switch"
  "target"
  "throw"
  "try"
  "typeof"
  "var"
  "void"
  "while"
  "with"
  "yield"
] @keyword

(comment) @comment
[
  (string)
  (template_string)
] @string
(number) @number
(true) @constant
(false) @constant
(null) @constant
(undefined) @constant

(function_declaration name: (identifier) @function)
(generator_function_declaration name: (identifier) @function)
(call_expression function: (identifier) @function)
(call_expression function: (member_expression property: (property_identifier) @function))

(jsx_opening_element name: (identifier) @tag)
(jsx_closing_element name: (identifier) @tag)
(jsx_self_closing_element name: (identifier) @tag)
(jsx_attribute (property_identifier) @attribute)
`

const htmlHighlightsQuery = `
(comment) @comment
(doctype) @constant
(tag_name) @tag
(attribute_name) @attribute
(quoted_attribute_value) @string
`

const bashHighlightsQuery = `
[
  "case"
  "do"
  "done"
  "elif"
  "else"
  "esac"
  "fi"
  "for"
  "function"
  "if"
  "in"
  "select"
  "then"
  "until"
  "while"
] @keyword

(comment) @comment
(command_name) @function
[
  (string)
  (raw_string)
  (ansi_c_string)
  (heredoc_body)
] @string
(number) @number
(variable_name) @variable
(simple_expansion) @variable
`
