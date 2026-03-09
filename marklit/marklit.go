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
	return parseArg(source, false)
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
	node = unwrapInlineResult(node)

	// Goldmark trims leading/trailing whitespace from paragraph text.
	// In Booklit inline args, whitespace is significant (e.g.
	// @aux{The } needs the trailing space). Restore any stripped
	// whitespace by comparing against the original source.
	leading := leadingWhitespace(source)
	trailing := trailingWhitespace(source)
	if leading == "" && trailing == "" {
		return node
	}

	var nodes []ast.Node
	if leading != "" {
		nodes = append(nodes, ast.String(leading))
	}
	if seq, ok := node.(ast.Sequence); ok {
		nodes = append(nodes, seq...)
	} else {
		nodes = append(nodes, node)
	}
	if trailing != "" {
		nodes = append(nodes, ast.String(trailing))
	}
	return ast.Sequence(nodes)
}

func leadingWhitespace(s []byte) string {
	for i, ch := range s {
		if ch != ' ' && ch != '\t' {
			return string(s[:i])
		}
	}
	return string(s)
}

func trailingWhitespace(s []byte) string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] != ' ' && s[i] != '\t' {
			return string(s[i+1:])
		}
	}
	return string(s)
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
// Leading common indentation is stripped after preprocessing (which extracts
// {{{...}}} verbatim blocks) but before goldmark parsing. This prevents
// goldmark from interpreting indented content as code blocks without
// corrupting verbatim content that has different indentation.
func ParseArg(source []byte) ast.Node {
	return parseArg(source, true)
}

func parseArg(source []byte, doStripIndent bool) ast.Node {
	if doStripIndent {
		source = stripIndent(source)
	}

	processed, extractions := preprocess(source)

	p := newParser()
	reader := text.NewReader(processed)
	doc := p.Parse(reader)

	c := &converter{
		source:      processed,
		extractions: extractions,
	}
	result := c.convertChildren(doc)
	// An empty or whitespace-only arg should produce an empty String,
	// not an empty Sequence (which evaluates to nil and panics).
	if result == nil {
		return ast.String("")
	}
	if seq, ok := result.(ast.Sequence); ok && len(seq) == 0 {
		return ast.String("")
	}
	return result
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
		parser.WithHeadingAttribute(),
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
