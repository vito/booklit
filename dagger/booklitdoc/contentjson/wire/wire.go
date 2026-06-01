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

	S        string           `json:"s,omitempty"`        // string
	Items    []*Node          `json:"items,omitempty"`    // seq, para, pre, list
	Content  *Node            `json:"content,omitempty"`  // styled, link, aux, ref, target
	Style    string           `json:"style,omitempty"`    // styled
	Block    bool             `json:"block,omitempty"`    // styled
	Partials map[string]*Node `json:"partials,omitempty"` // styled
	Target   string           `json:"target,omitempty"`   // link
	Path     string           `json:"path,omitempty"`     // image
	Desc     string           `json:"desc,omitempty"`     // image
	Ordered  bool             `json:"ordered,omitempty"`  // list
	Rows     [][]*Node        `json:"rows,omitempty"`     // table
	Defs     []Def            `json:"defs,omitempty"`     // definitions
	Tag      string           `json:"tag,omitempty"`      // ref, target
	Title    *Node            `json:"title,omitempty"`    // target
	Optional bool             `json:"optional,omitempty"` // ref
}

// Def is one entry in a definitions node.
type Def struct {
	Subject *Node `json:"subject"`
	Def     *Node `json:"def"`
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

// Pre is preformatted content (block), e.g. a code block.
func Pre(items ...*Node) *Node { return &Node{Kind: "pre", Items: items} }

// Styled renders content with the named template.
func Styled(style string, content *Node) *Node {
	return &Node{Kind: "styled", Style: style, Content: content}
}

// StyledBlock is Styled forced to block layout.
func StyledBlock(style string, content *Node) *Node {
	return &Node{Kind: "styled", Style: style, Block: true, Content: content}
}

// RawHTML passes an HTML string through the renderer unescaped.
func RawHTML(html string) *Node { return Styled("raw-html", String(html)) }

// Link is a hyperlink.
func Link(target string, content *Node) *Node {
	return &Node{Kind: "link", Target: target, Content: content}
}

// Image embeds an image.
func Image(path, desc string) *Node { return &Node{Kind: "image", Path: path, Desc: desc} }

// List is an ordered or unordered list.
func List(ordered bool, items ...*Node) *Node {
	return &Node{Kind: "list", Ordered: ordered, Items: items}
}

// Table is a grid of rows and columns.
func Table(rows ...[]*Node) *Node { return &Node{Kind: "table", Rows: rows} }

// Definitions is a list of subject/definition pairs.
func Definitions(defs ...Def) *Node { return &Node{Kind: "defs", Defs: defs} }

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
