// Package wire is the dependency-free representation of Booklit's content JSON
// wire format. It deliberately imports nothing but encoding/json so it can be
// used from an out-of-process producer — in particular a Dagger module — that
// must not depend on the whole booklit package.
//
// A producer builds a *Node with the constructors below and serializes it with
// Marshal; Booklit's contentjson package converts between *Node and native
// booklit.Content. Keep this schema in lockstep with contentjson.
package wire

import "encoding/json"

// Node is one content node. It is a flat union: Kind discriminates, and only
// the fields relevant to that kind are populated.
type Node struct {
	Kind string `json:"k"`

	S        string  `json:"s,omitempty"`        // string, fragment
	Items    []*Node `json:"items,omitempty"`    // seq, para
	Content  *Node   `json:"content,omitempty"`  // element, aux, ref, target
	HTMLTag  string  `json:"htag,omitempty"`     // element
	Attrs    string  `json:"attrs,omitempty"`    // element
	Tag      string  `json:"tag,omitempty"`      // ref, target
	Title    *Node   `json:"title,omitempty"`    // target
	Optional bool    `json:"optional,omitempty"` // ref
}

// Marshal serializes a node tree to the wire format.
func Marshal(n *Node) ([]byte, error) { return json.Marshal(n) }

// Unmarshal parses the wire format into a node tree. A JSON null parses to a
// nil *Node.
func Unmarshal(data []byte) (*Node, error) {
	var n *Node
	if err := json.Unmarshal(data, &n); err != nil {
		return nil, err
	}
	return n, nil
}

// Constructors — the higher-level builder API a producer composes with.

// String is flow text.
func String(s string) *Node { return &Node{Kind: "string", S: s} }

// Seq concatenates content.
func Seq(items ...*Node) *Node { return &Node{Kind: "seq", Items: items} }

// Para is a paragraph (block).
func Para(items ...*Node) *Node { return &Node{Kind: "para", Items: items} }

// Element is a raw HTML element with a tag, pre-rendered attribute string
// (leading space, attribute values HTML-escaped), and body content. A nil
// content renders self-closing.
func Element(tag, attrs string, content *Node) *Node {
	return &Node{Kind: "element", HTMLTag: tag, Attrs: attrs, Content: content}
}

// Fragment passes a pre-rendered HTML string through the renderer
// verbatim. Use only for markup whose structure Booklit doesn't need to
// see — the typical case is a syntax highlighter's per-token span
// wrappers.
func Fragment(html string) *Node { return &Node{Kind: "fragment", S: html} }

// Aux is auxiliary content that can be stripped in some contexts.
func Aux(content *Node) *Node { return &Node{Kind: "aux", Content: content} }

// Ref links to a tag defined elsewhere; Booklit resolves it against the
// section the content is decoded into. content is optional display content.
func Ref(tag string, content *Node) *Node {
	return &Node{Kind: "ref", Tag: tag, Content: content}
}

// OptionalRef is a Ref that displays its content instead of erroring when the
// tag is missing.
func OptionalRef(tag string, content *Node) *Node {
	return &Node{Kind: "ref", Tag: tag, Content: content, Optional: true}
}

// Target creates an anchor for tag with the given title and body content.
func Target(tag string, title, content *Node) *Node {
	return &Node{Kind: "target", Tag: tag, Title: title, Content: content}
}
