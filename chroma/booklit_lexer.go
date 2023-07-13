package chroma

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers"
)

// Booklit is a lexer for Booklit syntax.
var Booklit = lexers.Register(MustNewLexer(
	&Config{
		Name:      "Booklit",
		Aliases:   []string{"booklit"},
		Filenames: []string{"*.lit"},
		MimeTypes: []string{"text/x-booklit"},
	},
	Rules{
		"root": {
			{Pattern: `[^\\{}]+`, Type: Text, Mutator: nil},
			{Pattern: `\{\{\{`, Type: StringDouble, Mutator: Push("verbatim")},
			{Pattern: `\{-`, Type: CommentMultiline, Mutator: Push("comment")},
			{Pattern: `[{}]`, Type: NameBuiltin, Mutator: nil},
			{Pattern: `\\([a-z-]+)`, Type: Keyword, Mutator: nil},
			{Pattern: `\\[\\{}]+`, Type: Text, Mutator: nil},
		},
		"verbatim": {
			{Pattern: `\}\}\}`, Type: StringDouble, Mutator: Pop(1)},
			{Pattern: `[^}]+`, Type: StringDouble, Mutator: nil},
			{Pattern: `}[^\}]`, Type: StringDouble, Mutator: nil},
		},
		"comment": {
			{Pattern: `[^-{}]+`, Type: CommentMultiline, Mutator: nil},
			{Pattern: `\{-`, Type: CommentMultiline, Mutator: Push()},
			{Pattern: `-\}`, Type: CommentMultiline, Mutator: Pop(1)},
			{Pattern: `[-{}]`, Type: CommentMultiline, Mutator: nil},
		},
	},
))
