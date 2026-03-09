// Package marklit parses Markdown documents with Booklit @invoke extensions
// and produces Booklit AST nodes.
package marklit

import (
	"bytes"

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
//
// Unlike Parse, this unwraps single-paragraph results so that inline
// arguments produce flat content rather than block-wrapped content.
// This matches the behavior of the old PEG parser where inline args
// like \title{Hello} produced a Sequence, not a Paragraph.
func ParseInlineArg(source []byte) ast.Node {
	if len(bytes.TrimSpace(source)) == 0 {
		return ast.String(source)
	}
	node := Parse(source)
	return unwrapInlineResult(node)
}

// unwrapInlineResult strips unnecessary Paragraph/Sequence wrapping from
// a parsed inline argument. If the result is a single Paragraph with one
// line, we return just the line's contents.
func unwrapInlineResult(node ast.Node) ast.Node {
	// Single paragraph with one line → return the line content
	if para, ok := node.(ast.Paragraph); ok && len(para) == 1 {
		return para[0]
	}
	// Sequence containing a single paragraph → unwrap that too
	if seq, ok := node.(ast.Sequence); ok && len(seq) == 1 {
		if para, ok := seq[0].(ast.Paragraph); ok && len(para) == 1 {
			return para[0]
		}
	}
	return node
}

// ParseArg parses a full argument (may contain block elements, paragraphs,
// etc.) into a Booklit AST node. Used for block-level @invoke arguments.
//
// Leading common indentation is stripped before parsing to prevent goldmark
// from interpreting indented content as code blocks. This is necessary
// because args nested inside multiple invocations accumulate indentation
// (e.g. @section{@table{  @table-row{...}}} has 4+ spaces).
func ParseArg(source []byte) ast.Node {
	return Parse(stripIndent(source))
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
