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
	"unicode"
	"unsafe"

	tshighlight "go.gopad.dev/go-tree-sitter-highlight"

	"github.com/tree-sitter/go-tree-sitter"
	tsbash "github.com/tree-sitter/tree-sitter-bash/bindings/go"
	tsgo "github.com/tree-sitter/tree-sitter-go/bindings/go"
	tshtml "github.com/tree-sitter/tree-sitter-html/bindings/go"
	tsjavascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	booklitgrammar "github.com/vito/booklit/treehighlight/internal/tree_sitter_booklit"
	"github.com/vito/dang/pkg/dang/danglang"
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
	name string
	// language returns the wrapped tree-sitter Language. The closure form
	// (instead of a stored *Language) defers cgo initialization until a
	// language is actually used; some grammars allocate non-trivial state.
	language   func() *tree_sitter.Language
	highlights string
	// injections is an optional tree-sitter injection query whose
	// (`@injection.content`, `@injection.language`) captures (typically set via
	// `#set! injection.language "name"`) hand off byte ranges to another
	// language's highlighter. The resolver in Chunks looks the named language
	// up in this map.
	injections string
	links      string
	linkTag    func(string) string
}

// wrapLanguage adapts an upstream tree-sitter binding's raw `Language() unsafe.Pointer`
// into a closure returning the runtime's wrapped *Language. We wrap per-call
// because tree_sitter.NewLanguage is cheap and avoids package-init ordering
// issues with the bound grammars.
func wrapLanguage(raw func() unsafe.Pointer) func() *tree_sitter.Language {
	return func() *tree_sitter.Language { return tree_sitter.NewLanguage(raw()) }
}

var languages = map[string]*languageSpec{
	"booklit": {
		name:       "booklit",
		language:   wrapLanguage(booklitgrammar.Language),
		highlights: booklitHighlightsQuery,
		injections: booklitInjectionsQuery,
		links:      booklitLinksQuery,
		linkTag:    booklitLinkTag,
	},
	"go": {
		name:       "go",
		language:   wrapLanguage(tsgo.Language),
		highlights: goHighlightsQuery,
	},
	"javascript": {
		name:       "javascript",
		language:   wrapLanguage(tsjavascript.Language),
		highlights: javascriptHighlightsQuery,
	},
	"html": {
		name:       "html",
		language:   wrapLanguage(tshtml.Language),
		highlights: htmlHighlightsQuery,
	},
	"bash": {
		name:       "bash",
		language:   wrapLanguage(tsbash.Language),
		highlights: bashHighlightsQuery,
	},
	"dang": {
		// danglang.Language already returns a *tree_sitter.Language; no
		// wrapper needed (and adding one would pointlessly call
		// tree_sitter.NewLanguage on an already-wrapped value).
		name:       "dang",
		language:   danglang.Language,
		highlights: dangHighlightsQuery,
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
	"dang":             "dang",
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
//
// When the parent language has an injections query, captured byte ranges hand
// off to a child language's highlighter via the resolver. We build child
// configurations lazily inside a configRegistry so each one is compiled at
// most once per Chunks call and all of them are released together.
func Chunks(language, source string, opts Options) ([]Chunk, error) {
	spec := lookupLanguage(language)
	if spec == nil {
		return []Chunk{{HTML: html.EscapeString(source)}}, nil
	}

	cfgs := newConfigRegistry()
	defer cfgs.close()

	cfg, err := cfgs.compile(spec)
	if err != nil {
		return nil, err
	}

	var links []linkRange
	if opts.LinkReferences && spec.links != "" {
		links, err = collectLinks(spec, cfg.Language, []byte(source))
		if err != nil {
			return nil, err
		}
	}

	highlighter := tshighlight.New()
	defer highlighter.Parser.Close()
	events := highlighter.Highlight(context.Background(), *cfg, []byte(source), cfgs.resolve)
	return renderChunks(events, []byte(source), links)
}

// configRegistry compiles and caches tshighlight.Configurations for the
// duration of a single Chunks call. The injection-resolver callback uses it
// to find or build a child language's configuration on demand.
type configRegistry struct {
	configs map[string]*tshighlight.Configuration
}

func newConfigRegistry() *configRegistry {
	return &configRegistry{configs: map[string]*tshighlight.Configuration{}}
}

func (r *configRegistry) compile(spec *languageSpec) (*tshighlight.Configuration, error) {
	if existing, ok := r.configs[spec.name]; ok {
		return existing, nil
	}
	lang := spec.language()
	cfg, err := tshighlight.NewConfiguration(lang, spec.name, []byte(spec.highlights), []byte(spec.injections), nil)
	if err != nil {
		return nil, fmt.Errorf("configure %s highlighter: %w", spec.name, err)
	}
	cfg.Configure(captureNames)
	r.configs[spec.name] = cfg
	return cfg, nil
}

// resolve is the InjectionCallback passed to tshighlight. It returns nil for
// unknown languages so the highlighter falls back to escaped source for that
// range — better than crashing on a language we can't recognize.
func (r *configRegistry) resolve(name string) *tshighlight.Configuration {
	if name == "" {
		return nil
	}
	if existing, ok := r.configs[name]; ok {
		return existing
	}
	spec, ok := languages[name]
	if !ok {
		return nil
	}
	cfg, err := r.compile(spec)
	if err != nil {
		return nil
	}
	return cfg
}

func (r *configRegistry) close() {
	for _, c := range r.configs {
		c.Query.Close()
		if c.CombinedInjectionsQuery != nil {
			c.CombinedInjectionsQuery.Close()
		}
	}
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

// booklitLinkTag converts a tree-sitter capture (typically a JSX-style
// component name) into the kebab-case tag that Booklit uses to name targets,
// matching the convention applied elsewhere when a Go method like
// IncludeSection is exposed as the `\include-section` invoke.
//
// Examples:
//
//	"title"           → "title"
//	"include-section" → "include-section"
//	"IncludeSection"  → "include-section"
//	"OutputFrame"     → "output-frame"
//	"TableOfContents" → "table-of-contents"
//	"HTMLRenderer"    → "html-renderer"
//
// Underscores, spaces, and existing dashes all normalize to a single dash;
// runs of separators collapse and leading/trailing dashes are trimmed.
func booklitLinkTag(s string) string {
	runes := []rune(s)
	if len(runes) == 0 {
		return ""
	}

	var out strings.Builder
	out.Grow(len(s) + 4)

	// lastWasDash starts true so a leading separator is suppressed, mirroring
	// the trailing Trim below.
	lastWasDash := true
	writeDash := func() {
		if !lastWasDash {
			out.WriteByte('-')
			lastWasDash = true
		}
	}
	writeRune := func(r rune) {
		out.WriteRune(r)
		lastWasDash = false
	}

	for i, r := range runes {
		switch {
		case r == '_' || r == ' ' || r == '-':
			writeDash()
		case unicode.IsUpper(r):
			if startsNewWord(runes, i) {
				writeDash()
			}
			writeRune(unicode.ToLower(r))
		default:
			writeRune(unicode.ToLower(r))
		}
	}
	return strings.Trim(out.String(), "-")
}

// startsNewWord reports whether the uppercase rune at index i should be
// preceded by a word boundary. Two cases qualify:
//
//   - camelCase: the previous rune is lowercase or a digit (e.g. "fooBar"
//     splits between 'o' and 'B').
//   - End of an acronym: the previous rune is also uppercase but the next
//     rune is lowercase (e.g. "HTMLRenderer" splits between 'L' and 'R'
//     because 'e' follows).
func startsNewWord(runes []rune, i int) bool {
	if i == 0 {
		return false
	}
	prev := runes[i-1]
	if unicode.IsLower(prev) || unicode.IsDigit(prev) {
		return true
	}
	return i+1 < len(runes) && unicode.IsLower(runes[i+1])
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
(tag name: (identifier) @reference)
`

// booklitInjectionsQuery routes the contents of JSX-style {expr} regions
// (the booklit grammar's `tag_expr` node) to the Dang highlighter. The
// captured range includes the surrounding `{` and `}` because tshighlight's
// query layer here does not support the `#offset!` directive that would
// trim them — Dang parses `{...}` as a block expression, so the extra
// braces are tolerated by the child parser even if it adds an outer
// punctuation.bracket span.
//
// The trailing dummy pattern is a workaround for tshighlight's
// injectionForMatch, which short-circuits when the configuration lacks an
// @injection.language capture *index*, even for patterns that supply the
// name via #set!. The dummy gives the index something to bind to without
// actually injecting (no @injection.content, so the caller's `contentNode
// != nil` guard rejects it).
const booklitInjectionsQuery = `
((tag_expr) @injection.content
  (#set! injection.language "dang"))

((delimiter) @injection.language)
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

// dangHighlightsQuery mirrors ~/src/dang/treesitter/queries/highlights.scm.
// Dot-suffixed capture names (e.g. keyword.control, constant.numeric) fall
// back to their root capture (keyword, constant) via tshighlight's
// Configure() dot-fallback, so we do not need to add new entries to
// captureStyles for the more specific Dang captures.
const dangHighlightsQuery = `
;; Keywords
[
  (let_token)
  (pub_token)
] @keyword
[
  (type_token)
  (interface_token)
  (union_token)
  (implements_token)
  (enum_token)
  (scalar_token)
  (if_token)
  (else_token)
  (for_token)
  (break_token)
  (continue_token)
  (case_token)
  (directive_token)
  (on_token)
  (import_token)
  (new_token)
  (try_token)
  (catch_token)
  (raise_token)
  (return_token)
  (and_token)
  (or_token)
] @keyword.control

(self_keyword) @variable.special

;; Literals
(string) @string
(doc_string) @string
(triple_quote_string) @string
(int) @constant.numeric
(boolean) @constant.builtin.boolean
(null) @constant.builtin

;; Comments
(comment_token) @comment.line
(upper_token) @type

;; Directives
(directive_name) @function.macro
(directive_application
  (id) @function.macro)
(directive_location
  (upper_id) @constant.builtin)

;; Operators and punctuation
[
  (equal_token)
  (plus_equal_token)
  (double_interro_token)
  (bang_token)
  (arrow_token)
  (ampersand_token)
] @operator
["{{" "}}" "{" "}" "[" "]" "(" ")"] @punctuation.bracket
[
  (comma_token)
  (dot_token)
] @punctuation.delimiter
["@" "|"] @punctuation.special

;; Identifiers
(symbol) @variable
(call (symbol) @function.method)

;; Key-value pairs
(key_value
  (word_token) @property)

;; Field selections
(select_or_call
  (field_id) @function.method)
(field_selection
  (id) @property)

;; Type-bearing slots
(type_and_block_slot
  (symbol) @function.method)
(type_and_args_and_block_slot
  (symbol) @function.method)
(type_and_value_slot
  (symbol) @function.method)
(value_only_slot
  (symbol) @function.method)
(type_only_slot
  (symbol) @function.method)
(type_only_fun_slot
  (symbol) @function.method)

(arg_with_block_default
  (symbol) @variable.parameter)
(arg_with_type
  (symbol) @variable.parameter)
(arg_with_default
  (symbol) @variable.parameter)

;; Class definitions
(class (symbol) @type)
(implements (symbol) @type)
(interface (symbol) @type)
(enum (symbol) @type)
(enum (caps_symbol) @property)
(scalar (symbol) @type)
`
