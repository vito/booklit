package marklit

import "bytes"

// preprocess normalizes line endings and strips `{- ... -}` comments
// before goldmark sees the source. Backtick code spans are preserved
// verbatim so authors can write `{-` literally inside code.
func preprocess(source []byte) []byte {
	source = bytes.ReplaceAll(source, []byte("\r\n"), []byte("\n"))
	return stripComments(source)
}

// stripComments removes `{- ... -}` comments from source. Comments may be
// nested: `{- outer {- inner -} still comment -}`. Unmatched `{-` is left
// as-is. Backtick code spans are skipped so comments inside code are literal.
func stripComments(source []byte) []byte {
	if !bytes.Contains(source, []byte("{-")) {
		return source
	}

	var out []byte
	i := 0
	for i < len(source) {
		if source[i] == '`' {
			end := skipBacktickSpan(source, i)
			out = append(out, source[i:end+1]...)
			i = end + 1
			continue
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

// skipBacktickSpan skips over a backtick code span starting at pos.
// Returns the index of the closing backtick(s). Handles runs of multiple
// backticks (e.g. `` `` ... `` ``).
func skipBacktickSpan(source []byte, pos int) int {
	n := 0
	for pos+n < len(source) && source[pos+n] == '`' {
		n++
	}
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
	return pos
}

// findCommentEnd finds the end of a (possibly nested) `{- ... -}` comment
// starting at pos. Returns the index after the closing `-}`, or -1 if
// unmatched.
func findCommentEnd(source []byte, pos int) int {
	depth := 1
	i := pos + 2
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
	return -1
}
