package marklit

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// invokePlaceholder is used to replace block-level \invoke{...} sequences
// before goldmark parsing. After parsing, these are converted back.
// We use \x00 + "BLI" + decimal + \x00 to avoid goldmark splitting on
// punctuation/emphasis characters.
const invokePlaceholderPrefix = "\x00BLI"

// extractedInvoke holds a block-level invocation extracted during preprocessing.
type extractedInvoke struct {
	Function string
	RawArgs  [][]byte
	ArgTypes []ArgType // parallel to RawArgs
}

// preprocess scans source for \name{...} patterns where the braces span
// multiple lines (block-level invocations). It extracts them and replaces
// them with unique inline markers that survive goldmark's block parsing.
//
// Single-line \name{...} invocations are left alone for the inline parser
// to handle.
func preprocess(source []byte) ([]byte, []extractedInvoke) {
	// Normalize CRLF to LF so all downstream processing only deals with \n.
	source = bytes.ReplaceAll(source, []byte("\r\n"), []byte("\n"))

	// Strip {- ... -} comments (possibly nested) before goldmark sees the source.
	source = stripComments(source)

	var extractions []extractedInvoke
	var result []byte

	i := 0
	for i < len(source) {
		// Look for \ followed by a valid function name
		if source[i] != '\\' {
			result = append(result, source[i])
			i++
			continue
		}

		// Try to parse \name
		start := i
		i++ // skip \

		// \\ escape — pass through for goldmark to handle
		if i < len(source) && source[i] == '\\' {
			result = append(result, '\\', '\\')
			i++
			continue
		}

		// Parse function name
		if i >= len(source) || !isLowerAlpha(source[i]) {
			result = append(result, source[start:i]...)
			continue
		}

		nameStart := i
		for i < len(source) && isNameChar(source[i]) {
			i++
		}
		funcName := string(source[nameStart:i])

		// Check if there's a { following
		if i >= len(source) || source[i] != '{' {
			// No args — leave as-is for the inline parser
			result = append(result, source[start:i]...)
			continue
		}

		// Check if this brace-group spans multiple lines
		// (i.e. contains a newline before the matching close brace)
		args, argTypes, endPos, multiline := parseAllBracedArgs(source, i)
		if !multiline {
			// Single-line invocation — leave for the inline parser
			result = append(result, source[start:i]...)
			continue
		}

		// Multi-line (block) invocation — extract and replace with marker
		idx := len(extractions)
		extractions = append(extractions, extractedInvoke{
			Function: funcName,
			RawArgs:  args,
			ArgTypes: argTypes,
		})

		marker := fmt.Sprintf("%s%d\x00", invokePlaceholderPrefix, idx)
		result = append(result, marker...)
		i = endPos
	}

	return result, extractions
}

// parseAllBracedArgs parses consecutive {...} / {{...}} / {{{...}}} groups
// starting at pos. Returns the raw content of each arg, the type of each arg,
// the position after the last closing brace(s), and whether any arg was
// multiline.
func parseAllBracedArgs(source []byte, pos int) (args [][]byte, argTypes []ArgType, endPos int, multiline bool) {
	i := pos
	for i < len(source) && source[i] == '{' {
		content, end, ml, at := parseBracedContent(source, i)
		if end < 0 {
			// Unbalanced — return what we have
			return args, argTypes, i, multiline
		}
		args = append(args, content)
		argTypes = append(argTypes, at)
		if ml {
			multiline = true
		}
		i = end
	}
	return args, argTypes, i, multiline
}

// parseBracedContent reads a single braced group starting at pos (which must
// point to '{'). Detects {{{…}}} (verbatim), {{…}} (preformatted, requires
// newline after {{), and {…} (normal). Returns the inner content, the
// position after the closing brace(s), whether it spans multiple lines, and
// the argument type.
func parseBracedContent(source []byte, pos int) (content []byte, endPos int, multiline bool, argType ArgType) {
	if pos >= len(source) || source[pos] != '{' {
		return nil, -1, false, ArgNormal
	}

	// Check for {{{…}}} verbatim
	if pos+2 < len(source) && source[pos+1] == '{' && source[pos+2] == '{' {
		start := pos + 3
		for i := start; i+2 < len(source); i++ {
			if source[i] == '}' && source[i+1] == '}' && source[i+2] == '}' {
				content := source[start:i]
				ml := bytes.ContainsAny(content, "\n\r")
				return content, i + 3, ml, ArgVerbatim
			}
		}
		return nil, -1, true, ArgVerbatim
	}

	// Check for {{…}} preformatted (must be followed by newline)
	if pos+1 < len(source) && source[pos+1] == '{' {
		afterOpen := pos + 2
		if afterOpen < len(source) && (source[afterOpen] == '\n' || source[afterOpen] == '\r') {
			start := afterOpen
			depth := 0
			for i := start; i < len(source); i++ {
				switch source[i] {
				case '`':
					i = skipBacktickSpan(source, i)
				case '{':
					depth++
				case '}':
					if depth > 0 {
						depth--
					} else if i+1 < len(source) && source[i+1] == '}' {
						content := source[start:i]
						ml := bytes.ContainsAny(content, "\n\r")
						return content, i + 2, ml, ArgPreformatted
					}
				}
			}
			return nil, -1, true, ArgPreformatted
		}
		// {{ not followed by newline — fall through to normal {…} parsing
	}

	// Normal {…} with brace depth.
	// Backtick code spans are skipped so that braces inside them
	// (e.g. `{` or `}`) don't affect depth tracking.
	depth := 0
	start := pos + 1 // skip opening {

	for i := pos; i < len(source); i++ {
		switch source[i] {
		case '`':
			// Skip backtick code span
			i = skipBacktickSpan(source, i)
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				content := source[start:i]
				ml := bytes.ContainsAny(content, "\n\r")
				return content, i + 1, ml, ArgNormal
			}
		}
	}

	return nil, -1, true, ArgNormal
}

// skipBacktickSpan skips over a backtick code span starting at pos.
// Returns the index of the closing backtick(s). Handles runs of multiple
// backticks (e.g. `` `` ... `` ``).
func skipBacktickSpan(source []byte, pos int) int {
	// Count opening backticks
	n := 0
	for pos+n < len(source) && source[pos+n] == '`' {
		n++
	}
	// Find matching run of n backticks
	for i := pos + n; i+n <= len(source); i++ {
		match := true
		for j := 0; j < n; j++ {
			if source[i+j] != '`' {
				match = false
				break
			}
		}
		if match {
			return i + n - 1
		}
	}
	// No closing backticks found — return pos (just skip the one char)
	return pos
}

// resolvePlaceholders walks a string looking for placeholder markers and
// returns true if any are found. This is used during the Convert phase.
func isPlaceholder(text string) (int, bool) {
	if !strings.HasPrefix(text, invokePlaceholderPrefix) {
		return -1, false
	}
	rest := text[len(invokePlaceholderPrefix):]
	// Find the null terminator
	nullIdx := strings.IndexByte(rest, 0)
	if nullIdx < 0 {
		return -1, false
	}
	// The null terminator must be the last character (exact placeholder match)
	if nullIdx+1 != len(rest) {
		return -1, false
	}
	idx, err := strconv.Atoi(rest[:nullIdx])
	if err != nil {
		return -1, false
	}
	return idx, true
}

// stripComments removes {- ... -} comments from source. Comments may be
// nested: {- outer {- inner -} still comment -}. Unmatched {- is left as-is.
func stripComments(source []byte) []byte {
	// Quick check: if there's no {- at all, return as-is.
	if !bytes.Contains(source, []byte("{-")) {
		return source
	}

	var out []byte
	i := 0
	for i < len(source) {
		// Skip backtick code spans — comments inside are literal
		if source[i] == '`' {
			end := skipBacktickSpan(source, i)
			out = append(out, source[i:end+1]...)
			i = end + 1
			continue
		}
		// Skip over {{{...}}} verbatim blocks — comments inside are literal
		if i+2 < len(source) && source[i] == '{' && source[i+1] == '{' && source[i+2] == '{' {
			end := bytes.Index(source[i+3:], []byte("}}}"))
			if end >= 0 {
				end += i + 3 + 3 // past the closing }}}
				out = append(out, source[i:end]...)
				i = end
				continue
			}
		}
		if i+1 < len(source) && source[i] == '{' && source[i+1] == '-' {
			end := findCommentEnd(source, i)
			if end >= 0 {
				i = end
				continue
			}
		}
		out = append(out, source[i])
		i++
	}
	return out
}

// findCommentEnd finds the end of a (possibly nested) {- ... -} comment
// starting at pos. Returns the index after the closing -}, or -1 if
// unmatched.
func findCommentEnd(source []byte, pos int) int {
	depth := 1
	i := pos + 2 // skip opening {-
	for i < len(source) {
		if i+1 < len(source) && source[i] == '{' && source[i+1] == '-' {
			depth++
			i += 2
			continue
		}
		if i+1 < len(source) && source[i] == '-' && source[i+1] == '}' {
			depth--
			i += 2
			if depth == 0 {
				return i
			}
			continue
		}
		i++
	}
	return -1 // unmatched
}
