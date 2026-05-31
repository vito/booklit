// Package contentjson is the wire format for booklit.Content. It lets the
// serializable subset of a content tree cross a process boundary — most
// importantly, be produced by a Dagger module and decoded back into native
// booklit.Content by Booklit itself.
//
// Both ends import this package, so there is a single source of truth for the
// schema: a Dagger module Marshals, Booklit Unmarshals.
//
// Only "pure data" content is representable. Content that is bound to live
// in-process state — Section, TableOfContents, Lazy — cannot be serialized and
// returns an error from Marshal. Content that *refers* to in-process state by
// name — Reference and Target — crosses the wire carrying only the tag name;
// Unmarshal re-attaches the caller's *Section, so cross-references resolve as
// usual.
package contentjson

import (
	"encoding/json"
	"fmt"

	"github.com/vito/booklit"
)

// node is the on-the-wire shape of a single content node. It is a flat union:
// "k" discriminates, and only the fields relevant to that kind are populated.
type node struct {
	Kind string `json:"k"`

	S        string            `json:"s,omitempty"`        // string
	Items    []*node           `json:"items,omitempty"`    // seq, para, pre, list
	Content  *node             `json:"content,omitempty"`  // styled, link, aux, ref, target
	Style    string            `json:"style,omitempty"`    // styled
	Block    bool              `json:"block,omitempty"`    // styled
	Partials map[string]*node  `json:"partials,omitempty"` // styled
	Target   string            `json:"target,omitempty"`   // link
	Path     string            `json:"path,omitempty"`     // image
	Desc     string            `json:"desc,omitempty"`     // image
	Ordered  bool              `json:"ordered,omitempty"`  // list
	Rows     [][]*node         `json:"rows,omitempty"`     // table
	Defs     []defNode         `json:"defs,omitempty"`     // definitions
	Tag      string            `json:"tag,omitempty"`      // ref, target
	Title    *node             `json:"title,omitempty"`    // target
	Optional bool              `json:"optional,omitempty"` // ref
}

type defNode struct {
	Subject *node `json:"subject"`
	Def     *node `json:"def"`
}

// Marshal serializes content into the wire format. It returns an error if the
// tree contains content that is bound to in-process state (Section,
// TableOfContents, Lazy).
func Marshal(content booklit.Content) ([]byte, error) {
	n, err := encode(content)
	if err != nil {
		return nil, err
	}
	return json.Marshal(n)
}

// Unmarshal decodes the wire format back into native booklit.Content. sec is
// the section that Reference/Target nodes resolve against; it may be nil only
// if the tree contains no such nodes.
func Unmarshal(data []byte, sec *booklit.Section) (booklit.Content, error) {
	var n *node
	if err := json.Unmarshal(data, &n); err != nil {
		return nil, err
	}
	return decode(n, sec)
}

func encode(content booklit.Content) (*node, error) {
	switch v := content.(type) {
	case nil:
		return nil, nil
	case booklit.String:
		return &node{Kind: "string", S: string(v)}, nil
	case booklit.Sequence:
		items, err := encodeAll(v)
		if err != nil {
			return nil, err
		}
		return &node{Kind: "seq", Items: items}, nil
	case booklit.Paragraph:
		items, err := encodeAll(v)
		if err != nil {
			return nil, err
		}
		return &node{Kind: "para", Items: items}, nil
	case booklit.Preformatted:
		items, err := encodeAll(v)
		if err != nil {
			return nil, err
		}
		return &node{Kind: "pre", Items: items}, nil
	case booklit.Styled:
		inner, err := encode(v.Content)
		if err != nil {
			return nil, err
		}
		partials, err := encodePartials(v.Partials)
		if err != nil {
			return nil, err
		}
		return &node{Kind: "styled", Style: string(v.Style), Block: v.Block, Content: inner, Partials: partials}, nil
	case booklit.Link:
		inner, err := encode(v.Content)
		if err != nil {
			return nil, err
		}
		return &node{Kind: "link", Target: v.Target, Content: inner}, nil
	case booklit.Image:
		return &node{Kind: "image", Path: v.Path, Desc: v.Description}, nil
	case booklit.List:
		items, err := encodeAll(v.Items)
		if err != nil {
			return nil, err
		}
		return &node{Kind: "list", Ordered: v.Ordered, Items: items}, nil
	case booklit.Table:
		rows := make([][]*node, len(v.Rows))
		for i, row := range v.Rows {
			cols, err := encodeAll(row)
			if err != nil {
				return nil, err
			}
			rows[i] = cols
		}
		return &node{Kind: "table", Rows: rows}, nil
	case booklit.Definitions:
		defs := make([]defNode, len(v))
		for i, def := range v {
			subject, err := encode(def.Subject)
			if err != nil {
				return nil, err
			}
			body, err := encode(def.Definition)
			if err != nil {
				return nil, err
			}
			defs[i] = defNode{Subject: subject, Def: body}
		}
		return &node{Kind: "defs", Defs: defs}, nil
	case booklit.Aux:
		inner, err := encode(v.Content)
		if err != nil {
			return nil, err
		}
		return &node{Kind: "aux", Content: inner}, nil
	case *booklit.Reference:
		inner, err := encode(v.Content)
		if err != nil {
			return nil, err
		}
		return &node{Kind: "ref", Tag: v.TagName, Optional: v.Optional, Content: inner}, nil
	case booklit.Target:
		title, err := encode(v.Title)
		if err != nil {
			return nil, err
		}
		inner, err := encode(v.Content)
		if err != nil {
			return nil, err
		}
		return &node{Kind: "target", Tag: v.TagName, Title: title, Content: inner}, nil
	default:
		return nil, fmt.Errorf("contentjson: cannot serialize %T: it is bound to in-process state", content)
	}
}

func encodeAll(cs []booklit.Content) ([]*node, error) {
	if len(cs) == 0 {
		return nil, nil
	}
	out := make([]*node, len(cs))
	for i, c := range cs {
		n, err := encode(c)
		if err != nil {
			return nil, err
		}
		out[i] = n
	}
	return out, nil
}

func encodePartials(partials booklit.Partials) (map[string]*node, error) {
	if len(partials) == 0 {
		return nil, nil
	}
	out := make(map[string]*node, len(partials))
	for k, v := range partials {
		n, err := encode(v)
		if err != nil {
			return nil, err
		}
		out[k] = n
	}
	return out, nil
}

func decode(n *node, sec *booklit.Section) (booklit.Content, error) {
	if n == nil {
		return booklit.Empty, nil
	}
	switch n.Kind {
	case "string":
		return booklit.String(n.S), nil
	case "seq":
		items, err := decodeAll(n.Items, sec)
		if err != nil {
			return nil, err
		}
		return booklit.Sequence(items), nil
	case "para":
		items, err := decodeAll(n.Items, sec)
		if err != nil {
			return nil, err
		}
		return booklit.Paragraph(items), nil
	case "pre":
		items, err := decodeAll(n.Items, sec)
		if err != nil {
			return nil, err
		}
		return booklit.Preformatted(items), nil
	case "styled":
		inner, err := decode(n.Content, sec)
		if err != nil {
			return nil, err
		}
		partials, err := decodePartials(n.Partials, sec)
		if err != nil {
			return nil, err
		}
		return booklit.Styled{Style: booklit.Style(n.Style), Block: n.Block, Content: inner, Partials: partials}, nil
	case "link":
		inner, err := decode(n.Content, sec)
		if err != nil {
			return nil, err
		}
		return booklit.Link{Content: inner, Target: n.Target}, nil
	case "image":
		return booklit.Image{Path: n.Path, Description: n.Desc}, nil
	case "list":
		items, err := decodeAll(n.Items, sec)
		if err != nil {
			return nil, err
		}
		return booklit.List{Items: items, Ordered: n.Ordered}, nil
	case "table":
		rows := make([][]booklit.Content, len(n.Rows))
		for i, row := range n.Rows {
			cols, err := decodeAll(row, sec)
			if err != nil {
				return nil, err
			}
			rows[i] = cols
		}
		return booklit.Table{Rows: rows}, nil
	case "defs":
		defs := make(booklit.Definitions, len(n.Defs))
		for i, def := range n.Defs {
			subject, err := decode(def.Subject, sec)
			if err != nil {
				return nil, err
			}
			body, err := decode(def.Def, sec)
			if err != nil {
				return nil, err
			}
			defs[i] = booklit.Definition{Subject: subject, Definition: body}
		}
		return defs, nil
	case "aux":
		inner, err := decode(n.Content, sec)
		if err != nil {
			return nil, err
		}
		return booklit.Aux{Content: inner}, nil
	case "ref":
		inner, err := decode(n.Content, sec)
		if err != nil {
			return nil, err
		}
		return &booklit.Reference{Section: sec, TagName: n.Tag, Content: inner, Optional: n.Optional}, nil
	case "target":
		title, err := decode(n.Title, sec)
		if err != nil {
			return nil, err
		}
		inner, err := decode(n.Content, sec)
		if err != nil {
			return nil, err
		}
		return booklit.Target{TagName: n.Tag, Title: title, Content: inner}, nil
	default:
		return nil, fmt.Errorf("contentjson: unknown content kind %q", n.Kind)
	}
}

func decodeAll(ns []*node, sec *booklit.Section) ([]booklit.Content, error) {
	if len(ns) == 0 {
		return nil, nil
	}
	out := make([]booklit.Content, len(ns))
	for i, n := range ns {
		c, err := decode(n, sec)
		if err != nil {
			return nil, err
		}
		out[i] = c
	}
	return out, nil
}

func decodePartials(ns map[string]*node, sec *booklit.Section) (booklit.Partials, error) {
	if len(ns) == 0 {
		return nil, nil
	}
	out := make(booklit.Partials, len(ns))
	for k, n := range ns {
		c, err := decode(n, sec)
		if err != nil {
			return nil, err
		}
		out[k] = c
	}
	return out, nil
}
