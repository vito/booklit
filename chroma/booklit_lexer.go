package chroma

import (
	. "github.com/alecthomas/chroma/v2" // nolint
	"github.com/alecthomas/chroma/v2/lexers"
)

// Booklit is a lexer for Booklit (Markdown + @invoke) syntax.
var Booklit = lexers.Register(MustNewLexer(
	&Config{
		Name:      "Booklit",
		Aliases:   []string{"booklit", "lit"},
		Filenames: []string{"*.lit"},
		MimeTypes: []string{"text/x-booklit"},
	},
	func() Rules {
		return Rules{
			"root": {
				{Pattern: `[^@{}`+"``"+`]+`, Type: Text, Mutator: nil},
				{Pattern: `\{\{\{`, Type: StringDouble, Mutator: Push("verbatim")},
				{Pattern: `@@`, Type: Text, Mutator: nil},
				{Pattern: `@([a-z][a-z0-9-]*)`, Type: Keyword, Mutator: nil},
				{Pattern: "[`]+", Type: StringBacktick, Mutator: nil},
				{Pattern: `[{}]`, Type: NameBuiltin, Mutator: nil},
			},
			"verbatim": {
				{Pattern: `\}\}\}`, Type: StringDouble, Mutator: Pop(1)},
				{Pattern: `[^}]+`, Type: StringDouble, Mutator: nil},
				{Pattern: `}[^\}]`, Type: StringDouble, Mutator: nil},
			},
		}
	},
))
