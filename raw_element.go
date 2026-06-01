package booklit

import "github.com/vito/booklit/internal/htmltags"

// RawElement is structured raw HTML: a tag name, a pre-rendered
// attribute string, and the element's body content. It is emitted for
// every lowercase JSX element (`<div>`, `<span>`, …), for the
// `<RawHTML>` builtin's container, and for the `<pre>` / `<code>`
// wrappers around fenced code blocks.
//
// The previous design lifted these through `Styled{Style: "raw-html"}`
// (open and close fragments wrapping the body, with a `Block` flag
// flagging the wrapper as block-level). That carried a bytes-as-text
// hazard (raw markup leaking into search indexes via Styled.Content.
// String()) and required string-parsing the open fragment at later
// stages. RawElement keeps the tag, attrs, and content separately so
// the renderer composes the HTML directly and consumers like
// stringifyEverything see only the inner content's text.
type RawElement struct {
	// Tag is the lowercase HTML element name (e.g. "div", "span", "pre").
	Tag string

	// Attrs is the already-rendered attribute substring, ready to splice
	// between the tag name and the closing `>` — leading space included
	// when non-empty, attribute values already HTML-escaped. The
	// evaluator builds this from JSX props; the code builtins build it
	// inline.
	Attrs string

	// Content is the element's body. A nil Content means the element
	// has no body and the renderer emits it self-closing
	// (`<tag attrs/>`), preserving today's behavior for childless
	// lowercase JSX like `<br/>` or `<hr/>`.
	Content Content
}

// IsFlow returns true when Tag is not a block-level HTML element. The
// classification comes from internal/htmltags, which mirrors
// CommonMark/HTML5's block-tag list. An unknown tag (e.g. a custom
// element like `<my-widget>`) is treated as flow — the same default
// HTML applies when rendering an unrecognized inline element.
func (con RawElement) IsFlow() bool {
	return !htmltags.Block[con.Tag]
}

// String returns the body content's text, or "" when the element has
// no body. Tag and Attrs are markup bytes, not user-visible text, so
// they are deliberately suppressed — plain-text consumers (search
// index, stringifyEverything) see only what the user actually wrote
// inside the element.
func (con RawElement) String() string {
	if con.Content == nil {
		return ""
	}
	return con.Content.String()
}

// Visit calls VisitRawElement.
func (con RawElement) Visit(visitor Visitor) error {
	return visitor.VisitRawElement(con)
}
