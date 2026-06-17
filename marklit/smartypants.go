package marklit

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Smart-quote replacement runes.
const (
	leftDoubleQuote  = '“' // “
	rightDoubleQuote = '”' // ”
	leftSingleQuote  = '‘' // ‘
	rightSingleQuote = '’' // ’ (also the apostrophe)
)

// smartQuotes strips Markdown backslash escapes and applies SmartyPants-style
// typographic substitution to plain flow text: straight double and single
// quotes become their curly equivalents based on the surrounding context.
//
// This only ever runs on plain paragraph/flow text. Verbatim and preformatted
// content — code spans, code blocks, and verbatim/preformatted \invoke
// arguments — reach the AST through other converters (convertCodeSpan,
// convertCodeBlock, verbatimToNode, ParsePreformattedArg) and never pass
// through convertText, so code is left exactly as written.
//
// prev is the last rune of the preceding text in the same flow block (0 at a
// block boundary), which lets a quote opening one text segment and closing a
// later one — e.g. "*bold*" — resolve correctly across emphasis. The returned
// rune is the last rune processed, to thread as prev into the next segment.
//
// A backslash escape is the explicit opt-out: `\"` and `\'` (like any other
// `\<punct>`) emit the literal straight character, matching CommonMark
// escaping.
func smartQuotes(prev rune, s string) (string, rune) {
	var out strings.Builder
	out.Grow(len(s))

	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])

		// Backslash escapes: emit the escaped punctuation literally (straight),
		// stripping the backslash. This is the escape hatch for literal quotes.
		if r == '\\' && i+size < len(s) && isASCIIPunct(s[i+size]) {
			next := rune(s[i+size])
			out.WriteRune(next)
			prev = next
			i += size + 1
			continue
		}

		switch r {
		case '"':
			if opensQuote(prev) {
				out.WriteRune(leftDoubleQuote)
			} else {
				out.WriteRune(rightDoubleQuote)
			}
		case '\'':
			// A single quote right after a letter or digit is an apostrophe or
			// a closing quote (don't, it's, the dogs', 'hi'). Otherwise it
			// opens at a boundary and closes anywhere else.
			switch {
			case isAlphaNumeric(prev):
				out.WriteRune(rightSingleQuote)
			case opensQuote(prev):
				out.WriteRune(leftSingleQuote)
			default:
				out.WriteRune(rightSingleQuote)
			}
		default:
			out.WriteRune(r)
		}

		prev = r
		i += size
	}

	return out.String(), prev
}

// opensQuote reports whether a quote following prev should open rather than
// close — i.e. prev is a word boundary: start-of-text, whitespace, or opening
// punctuation.
func opensQuote(prev rune) bool {
	switch prev {
	case 0, '(', '[', '{', '-', '–', '—', leftDoubleQuote, leftSingleQuote:
		return true
	}
	return unicode.IsSpace(prev)
}

func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}
