// Package contentjson is the wire format for booklit.Content. It lets the
// serializable subset of a content tree cross a process boundary — most
// importantly, be produced by a Dagger module and decoded back into native
// booklit.Content by Booklit itself.
//
// The on-the-wire schema lives in the dependency-free subpackage wire, so a
// producer (e.g. a Dagger module) can build content without importing booklit.
// This package converts between wire.Node and native booklit.Content.
//
// Only "pure data" content is representable. Content that is bound to live
// in-process state — Section, TableOfContents, Lazy — cannot be serialized and
// returns an error from Marshal. Content that *refers* to in-process state by
// name — Reference and Target — crosses the wire carrying only the tag name;
// Unmarshal re-attaches the caller's *Section, so cross-references resolve as
// usual.
package contentjson

import (
	"fmt"

	"github.com/vito/booklit"
	"github.com/vito/booklit/contentjson/wire"
)

// Marshal serializes content into the wire format. It returns an error if the
// tree contains content that is bound to in-process state (Section,
// TableOfContents, Lazy).
func Marshal(content booklit.Content) ([]byte, error) {
	n, err := encode(content)
	if err != nil {
		return nil, err
	}
	return wire.Marshal(n)
}

// Unmarshal decodes the wire format back into native booklit.Content. sec is
// the section that Reference/Target nodes resolve against; it may be nil only
// if the tree contains no such nodes.
func Unmarshal(data []byte, sec *booklit.Section) (booklit.Content, error) {
	n, err := wire.Unmarshal(data)
	if err != nil {
		return nil, err
	}
	return decode(n, sec)
}

func encode(content booklit.Content) (*wire.Node, error) {
	switch v := content.(type) {
	case nil:
		return nil, nil
	case booklit.String:
		return wire.String(string(v)), nil
	case booklit.Sequence:
		items, err := encodeAll(v)
		if err != nil {
			return nil, err
		}
		return wire.Seq(items...), nil
	case booklit.Paragraph:
		items, err := encodeAll(v)
		if err != nil {
			return nil, err
		}
		return wire.Para(items...), nil
	case booklit.Preformatted:
		items, err := encodeAll(v)
		if err != nil {
			return nil, err
		}
		return wire.Pre(items...), nil
	case booklit.RawElement:
		inner, err := encode(v.Content)
		if err != nil {
			return nil, err
		}
		return wire.Element(v.Tag, v.Attrs, inner), nil
	case booklit.RawFragment:
		return wire.Fragment(v.HTML), nil
	case booklit.Link:
		inner, err := encode(v.Content)
		if err != nil {
			return nil, err
		}
		return wire.Link(v.Target, inner), nil
	case booklit.Image:
		return wire.Image(v.Path, v.Description), nil
	case booklit.List:
		items, err := encodeAll(v.Items)
		if err != nil {
			return nil, err
		}
		return wire.List(v.Ordered, items...), nil
	case booklit.Table:
		rows := make([][]*wire.Node, len(v.Rows))
		for i, row := range v.Rows {
			cols, err := encodeAll(row)
			if err != nil {
				return nil, err
			}
			rows[i] = cols
		}
		return wire.Table(rows...), nil
	case booklit.Definitions:
		defs := make([]wire.Def, len(v))
		for i, def := range v {
			subject, err := encode(def.Subject)
			if err != nil {
				return nil, err
			}
			body, err := encode(def.Definition)
			if err != nil {
				return nil, err
			}
			defs[i] = wire.Def{Subject: subject, Def: body}
		}
		return wire.Definitions(defs...), nil
	case booklit.Aux:
		inner, err := encode(v.Content)
		if err != nil {
			return nil, err
		}
		return wire.Aux(inner), nil
	case *booklit.Reference:
		inner, err := encode(v.Content)
		if err != nil {
			return nil, err
		}
		n := wire.Ref(v.TagName, inner)
		n.Optional = v.Optional
		return n, nil
	case booklit.Target:
		title, err := encode(v.Title)
		if err != nil {
			return nil, err
		}
		inner, err := encode(v.Content)
		if err != nil {
			return nil, err
		}
		return wire.Target(v.TagName, title, inner), nil
	default:
		return nil, fmt.Errorf("contentjson: cannot serialize %T: it is bound to in-process state", content)
	}
}

func encodeAll(cs []booklit.Content) ([]*wire.Node, error) {
	if len(cs) == 0 {
		return nil, nil
	}
	out := make([]*wire.Node, len(cs))
	for i, c := range cs {
		n, err := encode(c)
		if err != nil {
			return nil, err
		}
		out[i] = n
	}
	return out, nil
}

func decode(n *wire.Node, sec *booklit.Section) (booklit.Content, error) {
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
	case "element":
		inner, err := decode(n.Content, sec)
		if err != nil {
			return nil, err
		}
		return booklit.RawElement{Tag: n.HTMLTag, Attrs: n.Attrs, Content: inner}, nil
	case "fragment":
		return booklit.RawFragment{HTML: n.S}, nil
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

func decodeAll(ns []*wire.Node, sec *booklit.Section) ([]booklit.Content, error) {
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

