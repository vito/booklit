// Package marklit parses Markdown documents with Booklit @invoke extensions
// and produces Booklit AST nodes.
package marklit

import (
	"github.com/vito/booklit/ast"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Parse parses a Markdown+Booklit source document into a Booklit AST node.
func Parse(source []byte) ast.Node {
	// Pre-process to extract block-level @invoke{...} sequences that span
	// multiple lines/paragraphs. These get replaced with markers.
	processed, extractions := preprocess(source)

	p := newParser()
	reader := text.NewReader(processed)
	doc := p.Parse(reader)

	c := &converter{
		source:      processed,
		extractions: extractions,
	}
	return c.convertChildren(doc)
}

// ParseInlineArg parses inline argument content (single-line, no block
// elements) into a Booklit AST node. Used for parsing the content inside
// @invoke{...} braces.
func ParseInlineArg(source []byte) ast.Node {
	return Parse(source)
}

// ParseArg parses a full argument (may contain block elements, paragraphs,
// etc.) into a Booklit AST node. Used for block-level @invoke arguments.
func ParseArg(source []byte) ast.Node {
	return Parse(source)
}

// newParser builds a goldmark parser with the Booklit @invoke extension
// registered.
func newParser() parser.Parser {
	return parser.NewParser(
		parser.WithBlockParsers(parser.DefaultBlockParsers()...),
		parser.WithInlineParsers(
			append(
				parser.DefaultInlineParsers(),
				util.Prioritized(NewInvokeInlineParser(), 100),
			)...,
		),
		parser.WithParagraphTransformers(parser.DefaultParagraphTransformers()...),
	)
}

// Extension is a goldmark.Extender that adds Booklit @invoke syntax support.
type Extension struct{}

// Extend implements goldmark.Extender.
func (e *Extension) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(NewInvokeInlineParser(), 100),
		),
	)
}
