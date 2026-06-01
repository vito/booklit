// Package marklit parses MarkDangJSX (Markdown + JSX + Dang expressions)
// source documents into Booklit AST nodes.
package marklit

import (
	"bytes"

	"github.com/vito/booklit/ast"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Parse parses a MarkDangJSX source document into a Booklit AST node.
func Parse(source []byte) ast.Node {
	return parseArg(source, false)
}

// ParseInlineArg parses inline argument content (single-line, no block
// elements) into a Booklit AST node. Used for parsing the content inside
// inline JSX elements.
//
// Unlike Parse, this unwraps single-paragraph results so that inline
// arguments produce flat content rather than block-wrapped content.
func ParseInlineArg(source []byte) ast.Node {
	if len(bytes.TrimSpace(source)) == 0 {
		return ast.String(source)
	}
	node := Parse(source)
	node = unwrapInlineResult(node)

	// Goldmark trims leading/trailing whitespace from paragraph text.
	// In Booklit inline args, whitespace is significant (e.g.
	// \aux{The } needs the trailing space). Restore any stripped
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

// stripIndent removes a leading newline, detects the indentation of the first
// content line, strips that indentation from all lines, and removes a trailing
// whitespace-only line (the line before the closing brace). It lets block JSX
// children keep their natural indentation in source without being interpreted
// as Markdown code blocks.
func stripIndent(content []byte) []byte {
	if len(content) > 0 && content[0] == '\n' {
		content = content[1:]
	} else if len(content) > 1 && content[0] == '\r' && content[1] == '\n' {
		content = content[2:]
	}

	if len(content) == 0 {
		return content
	}

	var prefix []byte
	for _, ch := range content {
		if ch == ' ' || ch == '\t' {
			prefix = append(prefix, ch)
		} else {
			break
		}
	}

	lines := bytes.Split(content, []byte("\n"))

	if len(lines) > 0 && len(bytes.TrimSpace(lines[len(lines)-1])) == 0 {
		lines = lines[:len(lines)-1]
	}

	if len(prefix) > 0 {
		for i, line := range lines {
			if bytes.HasPrefix(line, prefix) {
				lines[i] = line[len(prefix):]
			} else {
				lines[i] = bytes.TrimLeft(line, " \t")
			}
		}
	}

	return bytes.Join(lines, []byte("\n"))
}

// ParseArg parses a full argument (may contain block elements, paragraphs,
// etc.) into a Booklit AST node. Used for block-level JSX children.
//
// Leading common indentation is stripped before goldmark parsing so that
// indented block children don't get reinterpreted as code blocks.
func ParseArg(source []byte) ast.Node {
	return parseArg(source, true)
}

func parseArg(source []byte, doStripIndent bool) ast.Node {
	if doStripIndent {
		source = stripIndent(source)
	}

	processed := preprocess(source)

	p := newParser()
	reader := text.NewReader(processed)
	doc := p.Parse(reader)

	c := &converter{source: processed}
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

// newParser builds a goldmark parser with the Booklit JSX + reference
// extensions and GFM table support registered.
//
// The default HTML block parser (CommonMark §4.6 types 6/7) is stripped:
// our JSX block parser now claims both PascalCase and lowercase `<tag>`
// openings, so leaving HTMLBlockParser enabled would just race against
// us and occasionally win, gobbling content as raw HTML before we get a
// chance to parse interleaved JSX or `{expr}`.
func newParser() parser.Parser {
	return parser.NewParser(
		parser.WithBlockParsers(
			append(
				blockParsersWithoutHTMLBlock(),
				util.Prioritized(NewJSXBlockParser(), 100),
			)...,
		),
		parser.WithInlineParsers(
			append(
				parser.DefaultInlineParsers(),
				util.Prioritized(NewReferenceInlineParser(), 99),
				util.Prioritized(NewJSXInlineParser(), 98),
				util.Prioritized(NewJSXExpressionInlineParser(), 97),
			)...,
		),
		parser.WithParagraphTransformers(
			append(
				parser.DefaultParagraphTransformers(),
				util.Prioritized(extension.NewTableParagraphTransformer(), 200),
			)...,
		),
		parser.WithASTTransformers(
			util.Prioritized(extension.NewTableASTTransformer(), 0),
		),
		parser.WithHeadingAttribute(),
	)
}

// blockParsersWithoutHTMLBlock returns goldmark's default block parsers
// minus the HTMLBlockParser (priority 900). See newParser for the
// rationale.
func blockParsersWithoutHTMLBlock() []util.PrioritizedValue {
	defaults := parser.DefaultBlockParsers()
	filtered := make([]util.PrioritizedValue, 0, len(defaults))
	for _, p := range defaults {
		if p.Priority == 900 {
			continue
		}
		filtered = append(filtered, p)
	}
	return filtered
}

// Extension is a goldmark.Extender that adds Booklit JSX + reference syntax
// support.
type Extension struct{}

// Extend implements goldmark.Extender.
func (e *Extension) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(NewJSXBlockParser(), 100),
		),
		parser.WithInlineParsers(
			util.Prioritized(NewReferenceInlineParser(), 99),
			util.Prioritized(NewJSXInlineParser(), 98),
			util.Prioritized(NewJSXExpressionInlineParser(), 97),
		),
	)
}
