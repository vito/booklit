//go:build !cgo

// Package treehighlight renders code HTML, using tree-sitter when cgo is available.
package treehighlight

import "html"

// Chunk is a rendered fragment of highlighted source. HTML chunks should be
// emitted as raw HTML. Link chunks should be rendered as Booklit references.
type Chunk struct {
	HTML     string
	LinkTag  string
	LinkText string
}

// Options controls highlighting behavior.
type Options struct {
	// LinkReferences turns language-specific link captures into Link chunks when
	// tree-sitter support is available.
	LinkReferences bool
}

// HTML returns escaped source in the renderer's code wrapper. Tree-sitter
// highlighting is unavailable when Booklit is built with CGO_ENABLED=0.
func HTML(language, source string, inline bool) (string, error) {
	return PlainHTML(source, inline), nil
}

// Chunks returns escaped, unhighlighted source. Tree-sitter highlighting and
// link captures are unavailable when Booklit is built with CGO_ENABLED=0.
func Chunks(language, source string, opts Options) ([]Chunk, error) {
	return []Chunk{{HTML: html.EscapeString(source)}}, nil
}

// PlainHTML returns source escaped into the renderer's code wrapper without
// syntax spans.
func PlainHTML(source string, inline bool) string {
	if inline {
		return `<code style=";-webkit-text-size-adjust:none;">` + html.EscapeString(source) + `</code>`
	}
	return `<pre style=";-webkit-text-size-adjust:none;"><code>` + html.EscapeString(source) + `</code></pre>`
}
