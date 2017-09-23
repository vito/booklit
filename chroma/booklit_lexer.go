package chroma

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers"
)

// Booklit lexer.
var Booklit = lexers.Register(MustNewLexer(
	&Config{
		Name:      "Booklit",
		Aliases:   []string{"booklit"},
		Filenames: []string{"*.lit"},
		MimeTypes: []string{"text/x-booklit"},
	},
	Rules{
		"root": {
			{`[^\\{}]+`, Text, nil},
			{`\{\{\{`, StringDouble, Push("verbatim")},
			{`\{-`, CommentMultiline, Push("comment")},
			{`[{}]`, NameBuiltin, nil},
			{`\\([a-z-]+)`, Keyword, nil},
			{`\\[\\{}]+`, Text, nil},
		},
		"verbatim": {
			{`\}\}\}`, StringDouble, Pop(1)},
			{`[^}]+`, StringDouble, nil},
			{`}[^\}]`, StringDouble, nil},
		},
		"comment": {
			{`[^-{}]+`, CommentMultiline, nil},
			{`\{-`, CommentMultiline, Push()},
			{`-\}`, CommentMultiline, Pop(1)},
			{`[-{}]`, CommentMultiline, nil},
		},
	},
))
