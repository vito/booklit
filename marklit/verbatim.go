package marklit

import (
	"bytes"

	"github.com/vito/booklit/ast"
)

// ArgType indicates how a braced argument should be parsed.
type ArgType int

const (
	// ArgNormal is a standard {…} argument parsed as Markdown with \invoke.
	ArgNormal ArgType = iota
	// ArgPreformatted is a {{…}} argument with preserved whitespace
	// and \invoke parsing but no Markdown formatting.
	ArgPreformatted
	// ArgVerbatim is a {{{…}}} argument with no parsing at all.
	ArgVerbatim
)

// stripIndent removes a leading newline, detects the indentation of the first
// content line, strips that indentation from all lines, and removes a trailing
// whitespace-only line (the line before the closing braces).
func stripIndent(content []byte) []byte {
	// Strip leading newline
	if len(content) > 0 && content[0] == '\n' {
		content = content[1:]
	} else if len(content) > 1 && content[0] == '\r' && content[1] == '\n' {
		content = content[2:]
	}

	if len(content) == 0 {
		return content
	}

	// Detect indentation prefix from first line
	var prefix []byte
	for _, ch := range content {
		if ch == ' ' || ch == '\t' {
			prefix = append(prefix, ch)
		} else {
			break
		}
	}

	lines := bytes.Split(content, []byte("\n"))

	// Strip trailing whitespace-only line
	if len(lines) > 0 && len(bytes.TrimSpace(lines[len(lines)-1])) == 0 {
		lines = lines[:len(lines)-1]
	}

	// Strip common indentation prefix
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

// verbatimToNode converts raw verbatim content ({{{…}}}) into a Booklit AST
// node. Block-form verbatim (where the raw content starts with a newline,
// i.e. {{{<newline>...}}}) always produces ast.Preformatted — even if
// stripping indent reduces it to a single line. This ensures block-level
// rendering (e.g. <pre> for syntax-highlighted code). Inline-form verbatim
// ({{{content}}}) produces ast.String for single-line content.
func verbatimToNode(raw []byte) ast.Node {
	content := stripIndent(raw)
	rawIsBlock := len(raw) > 0 && (raw[0] == '\n' || raw[0] == '\r')

	if !rawIsBlock && !bytes.Contains(content, []byte("\n")) {
		return ast.String(content)
	}

	lines := bytes.Split(content, []byte("\n"))
	seqs := make([]ast.Sequence, len(lines))
	for i, line := range lines {
		seqs[i] = ast.Sequence{ast.String(line)}
	}

	return ast.Preformatted(seqs)
}

// ParsePreformattedArg parses preformatted argument content ({{…}}).
// Whitespace structure is preserved and \invoke syntax is recognized, but
// Markdown formatting (* / ** / []() etc.) is not applied.
//
// Always produces ast.Preformatted (block content), matching the old parser's
// {{…}} behavior which required a newline after {{ and always produced a
// preformatted block.
func ParsePreformattedArg(source []byte) ast.Node {
	content := stripIndent(stripComments(source))

	if len(content) == 0 {
		return ast.Preformatted{ast.Sequence{ast.String("")}}
	}

	lines := bytes.Split(content, []byte("\n"))
	seqs := make([]ast.Sequence, len(lines))
	for i, line := range lines {
		seqs[i] = parsePreformattedLine(line)
	}

	return ast.Preformatted(seqs)
}

// parsePreformattedLine parses a single line of preformatted content,
// recognizing \invoke{arg} patterns but treating everything else as literal
// text.
func parsePreformattedLine(line []byte) ast.Sequence {
	var nodes []ast.Node
	i := 0
	textStart := 0

	for i < len(line) {
		if line[i] != '\\' {
			i++
			continue
		}

		// \\ escape — produces literal backslash
		if i+1 < len(line) && line[i+1] == '\\' {
			if textStart < i {
				nodes = append(nodes, ast.String(line[textStart:i]))
			}
			nodes = append(nodes, ast.String("\\"))
			i += 2
			textStart = i
			continue
		}

		// Try \name
		j := i + 1
		if j >= len(line) || !isLowerAlpha(line[j]) {
			i++
			continue
		}

		nameStart := j
		for j < len(line) && isNameChar(line[j]) {
			j++
		}
		funcName := string(line[nameStart:j])

		// Parse args
		var args []ast.Node
		k := j
		for k < len(line) && line[k] == '{' {
			argNode, end := parsePreformattedBracedArg(line, k)
			if end < 0 {
				break
			}
			args = append(args, argNode)
			k = end
		}

		// Emit text before invoke
		if textStart < i {
			nodes = append(nodes, ast.String(line[textStart:i]))
		}

		nodes = append(nodes, ast.Invoke{
			Function:  funcName,
			Arguments: args,
		})
		i = k
		textStart = i
	}

	if textStart < len(line) {
		nodes = append(nodes, ast.String(line[textStart:]))
	}

	if len(nodes) == 0 {
		return ast.Sequence{ast.String("")}
	}

	return ast.Sequence(nodes)
}

// parsePreformattedBracedArg parses a braced argument within a preformatted
// line. Supports {…} (parsed as inline Markdown) and {{{…}}} (verbatim).
func parsePreformattedBracedArg(data []byte, pos int) (ast.Node, int) {
	if pos >= len(data) || data[pos] != '{' {
		return nil, -1
	}

	// Check for {{{…}}}
	if pos+2 < len(data) && data[pos+1] == '{' && data[pos+2] == '{' {
		start := pos + 3
		for i := start; i+2 < len(data); i++ {
			if data[i] == '}' && data[i+1] == '}' && data[i+2] == '}' {
				return verbatimToNode(data[start:i]), i + 3
			}
		}
		return nil, -1
	}

	// Normal {…} with brace depth
	depth := 0
	start := pos + 1
	for i := pos; i < len(data); i++ {
		switch data[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return ParseInlineArg(data[start:i]), i + 1
			}
		}
	}
	return nil, -1
}
