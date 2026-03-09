package marklit

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// invokePlaceholder is used to replace block-level @invoke{...} sequences
// before goldmark parsing. After parsing, these are converted back.
// We use \x00 + "BLI" + decimal + \x00 to avoid goldmark splitting on
// punctuation/emphasis characters.
const invokePlaceholderPrefix = "\x00BLI"

// extractedInvoke holds a block-level invocation extracted during preprocessing.
type extractedInvoke struct {
	Function string
	RawArgs  [][]byte
}

// preprocess scans source for @name{...} patterns where the braces span
// multiple lines (block-level invocations). It extracts them and replaces
// them with unique inline markers that survive goldmark's block parsing.
//
// Single-line @name{...} invocations are left alone for the inline parser
// to handle.
func preprocess(source []byte) ([]byte, []extractedInvoke) {
	var extractions []extractedInvoke
	var result []byte

	i := 0
	for i < len(source) {
		// Look for @ followed by a valid function name
		if source[i] != '@' {
			result = append(result, source[i])
			i++
			continue
		}

		// Try to parse @name
		start := i
		i++ // skip @

		// @@ escape
		if i < len(source) && source[i] == '@' {
			result = append(result, '@', '@')
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
		args, endPos, multiline := parseAllBracedArgs(source, i)
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
		})

		marker := fmt.Sprintf("%s%d\x00", invokePlaceholderPrefix, idx)
		result = append(result, marker...)
		i = endPos
	}

	return result, extractions
}

// parseAllBracedArgs parses consecutive {...} groups starting at pos.
// Returns the raw content of each arg, the position after the last },
// and whether any arg was multiline.
func parseAllBracedArgs(source []byte, pos int) (args [][]byte, endPos int, multiline bool) {
	i := pos
	for i < len(source) && source[i] == '{' {
		content, end, ml := parseBracedContent(source, i)
		if end < 0 {
			// Unbalanced — return what we have
			return args, i, multiline
		}
		args = append(args, content)
		if ml {
			multiline = true
		}
		i = end
	}
	return args, i, multiline
}

// parseBracedContent reads a single {...} group starting at pos (which must
// point to '{'). Returns the inner content, the position after '}', and
// whether it spans multiple lines.
func parseBracedContent(source []byte, pos int) (content []byte, endPos int, multiline bool) {
	if pos >= len(source) || source[pos] != '{' {
		return nil, -1, false
	}

	depth := 0
	start := pos + 1 // skip opening {

	for i := pos; i < len(source); i++ {
		switch source[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				content := source[start:i]
				ml := bytes.ContainsAny(content, "\n\r")
				return content, i + 1, ml
			}
		case '\n', '\r':
			multiline = true
		}
	}

	return nil, -1, multiline
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
	idx, err := strconv.Atoi(rest[:nullIdx])
	if err != nil {
		return -1, false
	}
	return idx, true
}
